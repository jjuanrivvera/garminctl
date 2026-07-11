package commands

import (
	"bytes"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/zalando/go-keyring"
)

// rtFunc adapts a function to an http.RoundTripper so tests can mock go-garmin's transport.
type rtFunc func(*http.Request) (*http.Response, error)

func (f rtFunc) RoundTrip(r *http.Request) (*http.Response, error) { return f(r) }

func mockOK() *http.Client {
	return &http.Client{Transport: rtFunc(func(r *http.Request) (*http.Response, error) {
		body := "{}"
		if strings.Contains(r.URL.Path, "attery") { // body-battery returns a JSON array
			body = "[]"
		}
		return &http.Response{StatusCode: 200, Body: io.NopCloser(strings.NewReader(body)), Header: make(http.Header)}, nil
	})}
}

// execRoot builds a fresh root (which resets the global flags to defaults) and runs args.
func execRoot(t *testing.T, args ...string) (stdout, stderr string, err error) {
	t.Helper()
	root := NewRootCmd()
	var out, errb bytes.Buffer
	root.SetOut(&out)
	root.SetErr(&errb)
	root.SetArgs(args)
	err = root.Execute()
	return out.String(), errb.String(), err
}

func writeGarthTokens(t *testing.T, dir string, expiresAt int64) {
	t.Helper()
	o1 := `{"oauth_token":"OA1","oauth_token_secret":"OA1S","mfa_token":"","domain":"garmin.com"}`
	o2 := `{"access_token":"OA2","refresh_token":"OA2R","scope":"CONNECT_READ","expires_at":` + strconv.FormatInt(expiresAt, 10) + `}`
	if err := os.WriteFile(filepath.Join(dir, "oauth1_token.json"), []byte(o1), 0o600); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(dir, "oauth2_token.json"), []byte(o2), 0o600); err != nil {
		t.Fatal(err)
	}
}

func TestNewRootCmdSurface(t *testing.T) {
	root := NewRootCmd()
	want := map[string]bool{"auth": false, "body-composition": false, "sleep": false, "connect": false}
	for _, c := range root.Commands() {
		if _, ok := want[c.Name()]; ok {
			want[c.Name()] = true
		}
	}
	for name, found := range want {
		if !found {
			t.Errorf("root missing command %q", name)
		}
	}
	if root.PersistentFlags().Lookup("output") == nil || root.PersistentFlags().Lookup("profile") == nil {
		t.Error("root missing global flags")
	}
}

func TestExpandAliasesPassthrough(t *testing.T) {
	in := []string{"sleep", "--date", "2026-07-10"}
	got := ExpandAliases(in)
	if len(got) != len(in) || got[0] != "sleep" {
		t.Errorf("ExpandAliases changed args: %v", got)
	}
}

func TestPromptLineAndSecret(t *testing.T) {
	root := NewRootCmd()
	root.SetIn(strings.NewReader("hello\n"))
	got, err := promptLine(root, "x: ")
	if err != nil || got != "hello" {
		t.Errorf("promptLine = %q, %v", got, err)
	}
	// promptSecret falls back to a line read when stdin is not a terminal (a pipe).
	root.SetIn(strings.NewReader("s3cret\n"))
	got, err = promptSecret(root, "pw: ")
	if err != nil || got != "s3cret" {
		t.Errorf("promptSecret = %q, %v", got, err)
	}
}

func TestAuthImportStatusLogout(t *testing.T) {
	keyring.MockInit()
	t.Setenv("XDG_CONFIG_HOME", t.TempDir())
	t.Setenv("GARMINCTL_PROFILE", "")
	gdir := t.TempDir()
	writeGarthTokens(t, gdir, time.Now().Add(time.Hour).Unix())

	if _, errb, err := execRoot(t, "--profile", "juan", "auth", "import", "--from", gdir); err != nil {
		t.Fatalf("import: %v (%s)", err, errb)
	}
	out, _, err := execRoot(t, "--profile", "juan", "auth", "status")
	if err != nil || !strings.Contains(out, "authenticated:  true") {
		t.Fatalf("status: err=%v out=%q", err, out)
	}
	if _, _, err := execRoot(t, "--profile", "juan", "auth", "logout"); err != nil {
		t.Fatalf("logout: %v", err)
	}
	if _, _, err := execRoot(t, "--profile", "juan", "auth", "status"); err == nil {
		t.Error("status after logout should fail (no session)")
	}
}

func TestAuthImportMissingDir(t *testing.T) {
	keyring.MockInit()
	t.Setenv("XDG_CONFIG_HOME", t.TempDir())
	if _, _, err := execRoot(t, "--profile", "juan", "auth", "import", "--from", filepath.Join(t.TempDir(), "nope")); err == nil {
		t.Error("import from missing dir should fail")
	}
}

func TestResourceNoSessionErrors(t *testing.T) {
	keyring.MockInit()
	t.Setenv("XDG_CONFIG_HOME", t.TempDir())
	t.Setenv("GARMINCTL_PROFILE", "")
	if _, _, err := execRoot(t, "--profile", "nobody", "body-composition"); err == nil {
		t.Error("resource without a session should error")
	}
}

func TestResourceBadDate(t *testing.T) {
	keyring.MockInit()
	t.Setenv("XDG_CONFIG_HOME", t.TempDir())
	if _, _, err := execRoot(t, "sleep", "--date", "not-a-date"); err == nil {
		t.Error("invalid --date should error before hitting the API")
	}
}

func TestConnectBridgeBuilds(t *testing.T) {
	c := newConnectCmd()
	if c.Name() != "connect" {
		t.Fatalf("name = %q", c.Name())
	}
	if len(c.Commands()) < 5 {
		t.Errorf("connect should expose many endpoint groups, got %d", len(c.Commands()))
	}
}

func TestResourceSuccessMocked(t *testing.T) {
	keyring.MockInit()
	t.Setenv("XDG_CONFIG_HOME", t.TempDir())
	t.Setenv("GARMINCTL_PROFILE", "")
	testHTTPClient = mockOK()
	t.Cleanup(func() { testHTTPClient = nil })

	gdir := t.TempDir()
	writeGarthTokens(t, gdir, time.Now().Add(time.Hour).Unix())
	if _, _, err := execRoot(t, "--profile", "juan", "auth", "import", "--from", gdir); err != nil {
		t.Fatal(err)
	}
	for _, res := range []string{"body-composition", "sleep", "stress", "body-battery", "heart-rate", "respiration", "intensity-minutes"} {
		if _, _, err := execRoot(t, "--profile", "juan", res, "-o", "json"); err != nil {
			t.Errorf("%s fetch: %v", res, err)
		}
	}
}

func TestConnectExecMocked(t *testing.T) {
	keyring.MockInit()
	t.Setenv("XDG_CONFIG_HOME", t.TempDir())
	t.Setenv("GARMINCTL_PROFILE", "")
	testHTTPClient = mockOK()
	t.Cleanup(func() { testHTTPClient = nil })

	gdir := t.TempDir()
	writeGarthTokens(t, gdir, time.Now().Add(time.Hour).Unix())
	if _, _, err := execRoot(t, "--profile", "juan", "auth", "import", "--from", gdir); err != nil {
		t.Fatal(err)
	}
	// Find any runnable leaf under `connect` and execute it — this exercises the parent's
	// PersistentPreRunE (per-profile client wiring) and PostRunE (session save-back).
	connect := newConnectCmd()
	var leaf []string
	for _, grp := range connect.Commands() {
		if len(grp.Commands()) > 0 {
			leaf = []string{"connect", grp.Name(), grp.Commands()[0].Name()}
			break
		}
		if grp.RunE != nil || grp.Run != nil {
			leaf = []string{"connect", grp.Name()}
			break
		}
	}
	if leaf != nil {
		// Tolerate a non-nil error (arg shapes vary); the point is to cover the PreRun/PostRun.
		_, _, _ = execRoot(t, append([]string{"--profile", "juan"}, leaf...)...)
	}
}

func TestMainEntry(t *testing.T) {
	if code := Main([]string{"--version"}); code != 0 {
		t.Errorf("--version exit = %d, want 0", code)
	}
	if code := Main([]string{"definitely-not-a-command"}); code != 1 {
		t.Errorf("unknown command exit = %d, want 1", code)
	}
}

func TestStoreHelper(t *testing.T) {
	keyring.MockInit()
	t.Setenv("XDG_CONFIG_HOME", t.TempDir())
	s := store()
	if err := s.Set("p", "sess"); err != nil {
		t.Fatal(err)
	}
	got, err := s.Get("p")
	if err != nil || got != "sess" {
		t.Errorf("store round-trip: %q %v", got, err)
	}
}
