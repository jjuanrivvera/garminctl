package garmin

import (
	"context"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"testing"
	"time"
)

func writeGarth(t *testing.T, dir string, oauth1, oauth2 string, expiresAt int64) {
	t.Helper()
	if oauth1 == "" {
		oauth1 = `{"oauth_token":"OA1","oauth_token_secret":"OA1S","mfa_token":"","domain":"garmin.com"}`
	}
	if oauth2 == "" {
		oauth2 = `{"access_token":"OA2","refresh_token":"OA2R","scope":"CONNECT_READ","expires_at":` + strconv.FormatInt(expiresAt, 10) + `}`
	}
	if err := os.WriteFile(filepath.Join(dir, "oauth1_token.json"), []byte(oauth1), 0o600); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(dir, "oauth2_token.json"), []byte(oauth2), 0o600); err != nil {
		t.Fatal(err)
	}
}

func TestImportGarthRoundTrip(t *testing.T) {
	dir := t.TempDir()
	exp := time.Now().Add(time.Hour).Unix()
	writeGarth(t, dir, "", "", exp)
	js, err := ImportGarth(dir)
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(js, "OA1") || !strings.Contains(js, "OA2") {
		t.Errorf("tokens missing from session json: %s", js)
	}
	gotExp, authed, err := SessionInfo(js)
	if err != nil {
		t.Fatal(err)
	}
	if !authed {
		t.Error("expected authenticated session")
	}
	if gotExp.Unix() != exp {
		t.Errorf("expiry mismatch: got %d want %d", gotExp.Unix(), exp)
	}
}

func TestImportGarthMissingDir(t *testing.T) {
	if _, err := ImportGarth(t.TempDir()); err == nil {
		t.Error("missing token files should error")
	}
}

func TestImportGarthIncomplete(t *testing.T) {
	dir := t.TempDir()
	writeGarth(t, dir, `{"oauth_token":""}`, `{"access_token":""}`, 0)
	if _, err := ImportGarth(dir); err == nil {
		t.Error("empty tokens should error")
	}
}

func TestNewClientAndDump(t *testing.T) {
	dir := t.TempDir()
	writeGarth(t, dir, "", "", time.Now().Add(time.Hour).Unix())
	js, _ := ImportGarth(dir)
	c, err := NewClient(context.Background(), js, nil)
	if err != nil {
		t.Fatal(err)
	}
	dumped, err := DumpSession(c)
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(dumped, "OA1") {
		t.Errorf("dumped session missing token: %s", dumped)
	}
}

func TestNewClientBadJSON(t *testing.T) {
	if _, err := NewClient(context.Background(), "not json", nil); err == nil {
		t.Error("bad session json should error")
	}
}

func TestSessionInfoBadJSON(t *testing.T) {
	if _, _, err := SessionInfo("{bad"); err == nil {
		t.Error("bad json should error")
	}
}
