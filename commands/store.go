package commands

import (
	"encoding/json"
	"path/filepath"

	"github.com/jjuanrivvera/garminctl/internal/config"
	"github.com/jjuanrivvera/garminctl/internal/store"
)

// openStore opens the per-user offline SQLite store (config dir / store.db). The caller closes it.
func openStore() (*store.Store, error) {
	dir, err := config.Dir()
	if err != nil {
		return nil, err
	}
	return store.Open(filepath.Join(dir, "store.db"))
}

// cacheSample best-effort records a freshly fetched result to the offline store. A store error
// never fails the read — the data was already fetched and rendered successfully.
func cacheSample(profile, metric, date string, v any) {
	b, err := json.Marshal(v)
	if err != nil {
		return
	}
	st, err := openStore()
	if err != nil {
		return
	}
	defer func() { _ = st.Close() }()
	_ = st.Put(profile, metric, date, b)
}

// offlineSample reads a cached sample as a generic value for rendering. ok is false when the date
// isn't in the store.
func offlineSample(profile, metric, date string) (v any, ok bool, err error) {
	st, err := openStore()
	if err != nil {
		return nil, false, err
	}
	defer func() { _ = st.Close() }()
	raw, ok, err := st.Get(profile, metric, date)
	if err != nil || !ok {
		return nil, ok, err
	}
	if err := json.Unmarshal(raw, &v); err != nil {
		return nil, false, err
	}
	return v, true, nil
}
