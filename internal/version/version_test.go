package version

import (
	"strings"
	"testing"
)

func TestGetMatchesVars(t *testing.T) {
	got := Get()
	if got.Version != Version || got.Commit != Commit || got.Date != Date {
		t.Errorf("Get() = %+v, want vars %s/%s/%s", got, Version, Commit, Date)
	}
}

func TestStringContainsName(t *testing.T) {
	s := String()
	if !strings.Contains(s, "garminctl") || !strings.Contains(s, Version) {
		t.Errorf("String() = %q", s)
	}
}
