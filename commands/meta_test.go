package commands

import (
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/zalando/go-keyring"

	"github.com/jjuanrivvera/garminctl/internal/config"
)

func TestVersionCommand(t *testing.T) {
	out, _, err := execRoot(t, "version")
	if err != nil || out == "" {
		t.Fatalf("version: err=%v out=%q", err, out)
	}
	jsonOut, _, err := execRoot(t, "version", "--json")
	if err != nil || !strings.Contains(jsonOut, "version") {
		t.Errorf("version --json: err=%v out=%q", err, jsonOut)
	}
}

func TestCompletionCommand(t *testing.T) {
	for _, sh := range []string{"bash", "zsh", "fish", "powershell"} {
		out, _, err := execRoot(t, "completion", sh)
		if err != nil || out == "" {
			t.Errorf("completion %s: err=%v empty=%v", sh, err, out == "")
		}
	}
	if _, _, err := execRoot(t, "completion", "nonsense"); err == nil {
		t.Error("invalid shell should error")
	}
}

func TestConfigListUsePath(t *testing.T) {
	t.Setenv("XDG_CONFIG_HOME", t.TempDir())
	t.Setenv("GARMINCTL_PROFILE", "")
	c := &config.Config{}
	c.AddProfile("me")
	c.AddProfile("alt")
	if err := config.Save(c); err != nil {
		t.Fatal(err)
	}

	out, _, err := execRoot(t, "config", "list")
	if err != nil || !strings.Contains(out, "me") || !strings.Contains(out, "alt") {
		t.Fatalf("config list: err=%v out=%q", err, out)
	}
	if !strings.Contains(out, "* me") {
		t.Errorf("first profile should be default: %q", out)
	}

	if _, _, err := execRoot(t, "config", "use", "alt"); err != nil {
		t.Fatalf("config use: %v", err)
	}
	out, _, _ = execRoot(t, "config", "list")
	if !strings.Contains(out, "* alt") {
		t.Errorf("default should switch to alt: %q", out)
	}

	pathOut, _, err := execRoot(t, "config", "path")
	if err != nil || !strings.Contains(pathOut, "config.yaml") {
		t.Errorf("config path: err=%v out=%q", err, pathOut)
	}

	if _, _, err := execRoot(t, "config", "use", "nobody"); err == nil {
		t.Error("config use of unknown profile should error")
	}
}

func TestDoctorHealthy(t *testing.T) {
	keyring.MockInit()
	t.Setenv("XDG_CONFIG_HOME", t.TempDir())
	t.Setenv("GARMINCTL_PROFILE", "")
	gdir := t.TempDir()
	writeGarthTokens(t, gdir, time.Now().Add(time.Hour).Unix())
	if _, _, err := execRoot(t, "--profile", "me", "auth", "import", "--from", gdir); err != nil {
		t.Fatal(err)
	}
	out, _, err := execRoot(t, "doctor")
	if err != nil {
		t.Fatalf("doctor should be healthy: %v (%s)", err, out)
	}
	if !strings.Contains(out, "me") || !strings.Contains(out, "token valid") {
		t.Errorf("doctor output: %q", out)
	}
}

func TestDoctorNoProfiles(t *testing.T) {
	keyring.MockInit()
	t.Setenv("XDG_CONFIG_HOME", t.TempDir())
	out, _, err := execRoot(t, "doctor")
	if err != nil {
		t.Fatalf("doctor with no profiles should not error: %v", err)
	}
	if !strings.Contains(out, "no profiles") {
		t.Errorf("doctor should note no profiles: %q", out)
	}
}

func TestInitImportsAndReportsMissing(t *testing.T) {
	keyring.MockInit()
	t.Setenv("XDG_CONFIG_HOME", t.TempDir())
	t.Setenv("GARMINCTL_PROFILE", "")

	// Missing dir → guidance, no error.
	out, _, err := execRoot(t, "init", "--from", filepath.Join(t.TempDir(), "nope"))
	if err != nil || !strings.Contains(out, "No tokens found") {
		t.Fatalf("init missing: err=%v out=%q", err, out)
	}

	// Present dir → import.
	gdir := t.TempDir()
	writeGarthTokens(t, gdir, time.Now().Add(time.Hour).Unix())
	out, _, err = execRoot(t, "--profile", "me", "init", "--from", gdir)
	if err != nil || !strings.Contains(out, "imported") {
		t.Fatalf("init import: err=%v out=%q", err, out)
	}
	// The imported profile is usable.
	if status, _, err := execRoot(t, "--profile", "me", "auth", "status"); err != nil || !strings.Contains(status, "true") {
		t.Errorf("status after init: err=%v out=%q", err, status)
	}
}

func TestAgentGuardAllHosts(t *testing.T) {
	claude, _, err := execRoot(t, "agent", "guard", "--host", "claude-code")
	if err != nil || !strings.Contains(claude, "PreToolUse") || !strings.Contains(claude, "garminctl-guard.sh") {
		t.Fatalf("claude-code guard: err=%v out=%q", err, claude)
	}
	codex, _, err := execRoot(t, "agent", "guard", "--host", "codex")
	if err != nil || !strings.Contains(codex, "sandbox_mode") {
		t.Fatalf("codex guard: err=%v out=%q", err, codex)
	}
	opencode, _, err := execRoot(t, "agent", "guard", "--host", "opencode")
	if err != nil || !strings.Contains(opencode, "opencode.ai") {
		t.Fatalf("opencode guard: err=%v out=%q", err, opencode)
	}
	if _, _, err := execRoot(t, "agent", "guard", "--host", "bogus"); err == nil {
		t.Error("unknown host should error")
	}
	if _, _, err := execRoot(t, "agent", "guard"); err == nil {
		t.Error("missing required --host should error")
	}
}

func TestAPICommandDryRunAndMocked(t *testing.T) {
	keyring.MockInit()
	t.Setenv("XDG_CONFIG_HOME", t.TempDir())
	t.Setenv("GARMINCTL_PROFILE", "")
	gdir := t.TempDir()
	writeGarthTokens(t, gdir, time.Now().Add(time.Hour).Unix())
	if _, _, err := execRoot(t, "--profile", "me", "auth", "import", "--from", gdir); err != nil {
		t.Fatal(err)
	}

	// --dry-run prints a curl with the token redacted, no network.
	out, _, err := execRoot(t, "--profile", "me", "--dry-run", "api", "/userprofile-service/userprofile")
	if err != nil || !strings.Contains(out, "curl") {
		t.Fatalf("dry-run api: err=%v out=%q", err, out)
	}
	if strings.Contains(out, "OA2") { // the imported access token must not leak
		t.Errorf("token leaked in dry-run: %q", out)
	}

	// Mocked transport returns {} — the command renders it.
	testHTTPClient = mockOK()
	t.Cleanup(func() { testHTTPClient = nil })
	if _, _, err := execRoot(t, "--profile", "me", "api", "/userprofile-service/userprofile", "-o", "json"); err != nil {
		t.Errorf("mocked api: %v", err)
	}
}

func TestDoctorUnhealthy(t *testing.T) {
	keyring.MockInit()
	t.Setenv("XDG_CONFIG_HOME", t.TempDir())
	// A profile listed in config but with no keyring session → doctor fails.
	c := &config.Config{}
	c.AddProfile("ghost")
	if err := config.Save(c); err != nil {
		t.Fatal(err)
	}
	if _, _, err := execRoot(t, "doctor"); err == nil {
		t.Error("doctor should fail when a profile has no session")
	}
	// A corrupt session exercises the unreadable-session branch.
	if err := store().Set("ghost", "not json"); err != nil {
		t.Fatal(err)
	}
	if _, _, err := execRoot(t, "doctor"); err == nil {
		t.Error("doctor should fail on a corrupt session")
	}
}

func TestDoctorExpiredTokenStillHealthy(t *testing.T) {
	keyring.MockInit()
	t.Setenv("XDG_CONFIG_HOME", t.TempDir())
	t.Setenv("GARMINCTL_PROFILE", "")
	gdir := t.TempDir()
	writeGarthTokens(t, gdir, time.Now().Add(-time.Hour).Unix()) // already expired
	if _, _, err := execRoot(t, "--profile", "me", "auth", "import", "--from", gdir); err != nil {
		t.Fatal(err)
	}
	out, _, err := execRoot(t, "doctor")
	if err != nil {
		t.Fatalf("expired token is not a failure: %v (%s)", err, out)
	}
	if !strings.Contains(out, "refreshes on next call") {
		t.Errorf("doctor should note expiry is auto-refreshed: %q", out)
	}
}

func TestAPIDryRunWithBody(t *testing.T) {
	keyring.MockInit()
	t.Setenv("XDG_CONFIG_HOME", t.TempDir())
	t.Setenv("GARMINCTL_PROFILE", "")
	gdir := t.TempDir()
	writeGarthTokens(t, gdir, time.Now().Add(time.Hour).Unix())
	if _, _, err := execRoot(t, "--profile", "me", "auth", "import", "--from", gdir); err != nil {
		t.Fatal(err)
	}
	out, _, err := execRoot(t, "--profile", "me", "--dry-run", "api", "/weight-service/weight", "-X", "POST", "--data", `{"value":72.5}`)
	if err != nil {
		t.Fatalf("dry-run api with body: %v", err)
	}
	if !strings.Contains(out, "POST") || !strings.Contains(out, "--data") || !strings.Contains(out, "value") {
		t.Errorf("dry-run curl missing method/body: %q", out)
	}
}
