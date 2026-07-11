package config

import (
	"path/filepath"
	"testing"
)

func TestResolvePrecedence(t *testing.T) {
	t.Setenv("XDG_CONFIG_HOME", t.TempDir())
	t.Setenv("GARMINCTL_PROFILE", "")
	if got := Resolve("flagval"); got != "flagval" {
		t.Errorf("flag should win: %s", got)
	}
	t.Setenv("GARMINCTL_PROFILE", "envval")
	if got := Resolve(""); got != "envval" {
		t.Errorf("env should win over default: %s", got)
	}
	t.Setenv("GARMINCTL_PROFILE", "")
	if got := Resolve(""); got != "default" {
		t.Errorf("fallback should be default: %s", got)
	}
}

func TestSaveLoadRoundTrip(t *testing.T) {
	t.Setenv("XDG_CONFIG_HOME", t.TempDir())
	c := &Config{}
	c.AddProfile("juan")
	c.AddProfile("vane")
	c.AddProfile("juan") // duplicate ignored
	if len(c.Profiles) != 2 {
		t.Fatalf("duplicate not ignored: %v", c.Profiles)
	}
	if c.DefaultProfile != "juan" {
		t.Errorf("first profile should become default: %s", c.DefaultProfile)
	}
	if err := Save(c); err != nil {
		t.Fatal(err)
	}
	got, err := Load()
	if err != nil {
		t.Fatal(err)
	}
	if got.DefaultProfile != "juan" || len(got.Profiles) != 2 {
		t.Errorf("round-trip lost data: %+v", got)
	}
}

func TestLoadMissingIsEmpty(t *testing.T) {
	t.Setenv("XDG_CONFIG_HOME", filepath.Join(t.TempDir(), "does-not-exist"))
	c, err := Load()
	if err != nil {
		t.Fatalf("missing config should not error: %v", err)
	}
	if c.DefaultProfile != "" || len(c.Profiles) != 0 {
		t.Errorf("expected empty config, got %+v", c)
	}
}

func TestDirAndPathXDG(t *testing.T) {
	t.Setenv("XDG_CONFIG_HOME", "/x")
	d, err := Dir()
	if err != nil || d != "/x/garminctl" {
		t.Errorf("Dir() = %q, %v", d, err)
	}
	p, err := Path()
	if err != nil || p != "/x/garminctl/config.yaml" {
		t.Errorf("Path() = %q, %v", p, err)
	}
}

func TestDirHomeFallback(t *testing.T) {
	t.Setenv("XDG_CONFIG_HOME", "")
	t.Setenv("HOME", "/tmp/fakehome")
	d, err := Dir()
	if err != nil || d != "/tmp/fakehome/.garminctl-cli" {
		t.Errorf("Dir() HOME fallback = %q, %v", d, err)
	}
}
