// Package store is a local SQLite cache of daily Garmin metrics, so health data you've fetched
// stays queryable offline — Garmin Connect is a pull-only REST API with no history export, so
// garminctl accumulates its own. The driver is modernc.org/sqlite (pure Go, no cgo), matching
// the rest of the fleet. One row per (profile, metric, date); re-syncing upserts.
package store

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"

	_ "modernc.org/sqlite" // registers the "sqlite" database/sql driver
)

// Sample is one day's record of one metric for one profile. Data is the raw JSON payload as
// returned by the curated resource / API.
type Sample struct {
	Profile    string          `json:"profile"`
	Metric     string          `json:"metric"`
	Date       string          `json:"date"` // YYYY-MM-DD
	Data       json.RawMessage `json:"data"`
	RecordedAt time.Time       `json:"recorded_at"`
}

// Store is a handle to the on-disk SQLite database.
type Store struct{ db *sql.DB }

const schema = `
CREATE TABLE IF NOT EXISTS samples (
	profile     TEXT NOT NULL,
	metric      TEXT NOT NULL,
	date        TEXT NOT NULL,
	data        TEXT NOT NULL,
	recorded_at TEXT NOT NULL,
	PRIMARY KEY (profile, metric, date)
);
CREATE INDEX IF NOT EXISTS idx_samples_range ON samples(profile, metric, date);
`

// Open opens (creating if needed) the SQLite store and initializes its schema idempotently. The
// parent dir is created 0700 and the file chmod'd 0600 — it holds personal health data.
func Open(dbPath string) (*Store, error) {
	if err := os.MkdirAll(filepath.Dir(dbPath), 0o700); err != nil {
		return nil, fmt.Errorf("create store dir: %w", err)
	}
	db, err := sql.Open("sqlite", dbPath)
	if err != nil {
		return nil, fmt.Errorf("open store: %w", err)
	}
	// One connection: modernc.org/sqlite serializes writers internally, and a single garminctl
	// invocation never needs concurrent connections. busy_timeout smooths two processes.
	db.SetMaxOpenConns(1)
	if _, err := db.Exec(`PRAGMA busy_timeout = 5000;`); err != nil {
		_ = db.Close()
		return nil, fmt.Errorf("configure store: %w", err)
	}
	if _, err := db.Exec(schema); err != nil {
		_ = db.Close()
		return nil, fmt.Errorf("migrate store: %w", err)
	}
	if err := os.Chmod(dbPath, 0o600); err != nil {
		_ = db.Close()
		return nil, fmt.Errorf("chmod store: %w", err)
	}
	return &Store{db: db}, nil
}

// Close releases the database handle.
func (s *Store) Close() error { return s.db.Close() }

// Put upserts one day's metric — re-syncing a date overwrites it with fresh data.
func (s *Store) Put(profile, metric, date string, data json.RawMessage) error {
	_, err := s.db.Exec(
		`INSERT INTO samples (profile, metric, date, data, recorded_at) VALUES (?,?,?,?,?)
		 ON CONFLICT(profile, metric, date) DO UPDATE SET data=excluded.data, recorded_at=excluded.recorded_at`,
		profile, metric, date, string(data), time.Now().UTC().Format(time.RFC3339))
	if err != nil {
		return fmt.Errorf("store %s/%s/%s: %w", profile, metric, date, err)
	}
	return nil
}

// Get returns one day's stored metric. ok is false when the date isn't cached.
func (s *Store) Get(profile, metric, date string) (data json.RawMessage, ok bool, err error) {
	var raw string
	row := s.db.QueryRow(`SELECT data FROM samples WHERE profile=? AND metric=? AND date=?`, profile, metric, date)
	switch err := row.Scan(&raw); err {
	case nil:
		return json.RawMessage(raw), true, nil
	case sql.ErrNoRows:
		return nil, false, nil
	default:
		return nil, false, err
	}
}

// Range returns stored samples for a metric across [from, to] inclusive, ordered by date. Dates
// are compared as YYYY-MM-DD strings, which sort chronologically.
func (s *Store) Range(profile, metric, from, to string) ([]Sample, error) {
	rows, err := s.db.Query(
		`SELECT date, data, recorded_at FROM samples
		 WHERE profile=? AND metric=? AND date>=? AND date<=? ORDER BY date`,
		profile, metric, from, to)
	if err != nil {
		return nil, err
	}
	defer func() { _ = rows.Close() }()

	var out []Sample
	for rows.Next() {
		var date, data, recorded string
		if err := rows.Scan(&date, &data, &recorded); err != nil {
			return nil, err
		}
		ts, _ := time.Parse(time.RFC3339, recorded)
		out = append(out, Sample{Profile: profile, Metric: metric, Date: date, Data: json.RawMessage(data), RecordedAt: ts})
	}
	return out, rows.Err()
}
