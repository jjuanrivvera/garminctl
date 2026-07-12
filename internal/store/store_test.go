package store

import (
	"encoding/json"
	"path/filepath"
	"testing"
)

func open(t *testing.T) *Store {
	t.Helper()
	s, err := Open(filepath.Join(t.TempDir(), "store.db"))
	if err != nil {
		t.Fatalf("open: %v", err)
	}
	t.Cleanup(func() { _ = s.Close() })
	return s
}

func TestPutGetUpsert(t *testing.T) {
	s := open(t)

	// Missing → ok=false.
	if _, ok, err := s.Get("me", "sleep", "2026-07-09"); err != nil || ok {
		t.Fatalf("missing Get: ok=%v err=%v", ok, err)
	}

	if err := s.Put("me", "sleep", "2026-07-09", json.RawMessage(`{"score":72}`)); err != nil {
		t.Fatal(err)
	}
	data, ok, err := s.Get("me", "sleep", "2026-07-09")
	if err != nil || !ok || string(data) != `{"score":72}` {
		t.Fatalf("Get after Put: %q ok=%v err=%v", data, ok, err)
	}

	// Re-Put upserts (no duplicate row, fresh data).
	if err := s.Put("me", "sleep", "2026-07-09", json.RawMessage(`{"score":80}`)); err != nil {
		t.Fatal(err)
	}
	data, _, _ = s.Get("me", "sleep", "2026-07-09")
	if string(data) != `{"score":80}` {
		t.Errorf("upsert did not replace: %q", data)
	}

	// Different profile is isolated.
	if _, ok, _ := s.Get("alt", "sleep", "2026-07-09"); ok {
		t.Error("profiles must be isolated")
	}
}

func TestRange(t *testing.T) {
	s := open(t)
	for _, d := range []string{"2026-07-07", "2026-07-08", "2026-07-10"} {
		if err := s.Put("me", "stress", d, json.RawMessage(`{"avg":30}`)); err != nil {
			t.Fatal(err)
		}
	}
	// Out-of-range and other-metric rows must not appear.
	_ = s.Put("me", "stress", "2026-07-20", json.RawMessage(`{}`))
	_ = s.Put("me", "sleep", "2026-07-08", json.RawMessage(`{}`))

	got, err := s.Range("me", "stress", "2026-07-07", "2026-07-10")
	if err != nil {
		t.Fatal(err)
	}
	if len(got) != 3 {
		t.Fatalf("range len = %d, want 3", len(got))
	}
	// Ordered by date.
	if got[0].Date != "2026-07-07" || got[2].Date != "2026-07-10" {
		t.Errorf("range not ordered: %v", got)
	}
}
