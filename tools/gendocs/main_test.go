package main

import (
	"os"
	"path/filepath"
	"testing"
)

// TestGeneratesReference runs the generator in a temp working dir and checks it emits the
// command reference — so a broken cobra tree or doc API change fails the build, not the deploy.
func TestGeneratesReference(t *testing.T) {
	wd, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = os.Chdir(wd) })
	if err := os.Chdir(t.TempDir()); err != nil {
		t.Fatal(err)
	}

	main()

	for _, want := range []string{"garminctl.md", "garminctl_auth.md", "garminctl_sleep.md"} {
		if _, err := os.Stat(filepath.Join("docs", "commands", want)); err != nil {
			t.Errorf("expected generated %s: %v", want, err)
		}
	}
}
