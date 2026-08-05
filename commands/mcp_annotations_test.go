package commands

import (
	"strings"
	"testing"

	"github.com/njayp/ophis"
	"github.com/spf13/cobra"
)

func findCmd(root *cobra.Command, path string) *cobra.Command {
	cur := root
	for _, part := range strings.Fields(path) {
		var next *cobra.Command
		for _, sub := range cur.Commands() {
			if sub.Name() == part {
				next = sub
				break
			}
		}
		if next == nil {
			return nil
		}
		cur = next
	}
	return cur
}

func assertHints(t *testing.T, root *cobra.Command, path, readOnly, destructive string) {
	t.Helper()
	cmd := findCmd(root, path)
	if cmd == nil {
		t.Fatalf("command %q not found", path)
	}
	if got := cmd.Annotations[ophis.AnnotationReadOnly]; got != readOnly {
		t.Errorf("%q: readOnlyHint = %q, want %q", path, got, readOnly)
	}
	if got := cmd.Annotations[ophis.AnnotationDestructive]; got != destructive {
		t.Errorf("%q: destructiveHint = %q, want %q", path, got, destructive)
	}
}

// TestApplyMCPAnnotations_Mutations is the one that matters. garminctl is read-dominant, so
// the risk is a mutation being advertised read-only and auto-approved by a host.
func TestApplyMCPAnnotations_Mutations(t *testing.T) {
	root := NewRootCmd()

	assertHints(t, root, "courses delete", "false", "true")
	assertHints(t, root, "courses import", "false", "")
	assertHints(t, root, "update", "false", "")

	for _, p := range []string{
		"workouts create", "workouts update", "workouts delete",
		"workouts schedule", "workouts unschedule",
	} {
		cmd := findCmd(root, p)
		if cmd == nil {
			continue // registry-dependent; covered by the blocked-path test below
		}
		if cmd.Annotations[ophis.AnnotationReadOnly] == "true" {
			t.Errorf("%q must not be advertised read-only", p)
		}
	}
}

func TestApplyMCPAnnotations_Reads(t *testing.T) {
	root := NewRootCmd()

	for _, p := range []string{
		"sleep", "steps", "stress", "hrv", "spo2", "respiration",
		"body-battery", "body-composition", "heart-rate", "training-readiness",
		"courses list", "courses get", "courses download-gpx",
		"activities list", "activities get", "devices list",
		"metrics vo2max", "profile settings", "update check",
	} {
		assertHints(t, root, p, "true", "")
	}
}

// TestApplyMCPAnnotations_EveryLeafAnnotated stops a new command landing undeclared. The
// surface tracks go-garmin's registry, so commands arrive without anyone editing the
// annotation file — this fails the build when that happens.
func TestApplyMCPAnnotations_EveryLeafAnnotated(t *testing.T) {
	root := NewRootCmd()

	var missing []string
	walkCommands(root, func(cmd *cobra.Command) {
		if !cmd.Runnable() || cmd.Hidden || cmd.Name() == "help" {
			return
		}
		if _, ok := cmd.Annotations[ophis.AnnotationReadOnly]; !ok {
			missing = append(missing, cmd.CommandPath())
		}
	})
	if len(missing) > 0 {
		t.Errorf("commands missing MCP annotations: %v", missing)
	}
}

// TestApplyMCPAnnotations_UnknownVerbIsWrite pins the fail-safe default.
func TestApplyMCPAnnotations_UnknownVerbIsWrite(t *testing.T) {
	root := &cobra.Command{Use: "garminctl"}
	grp := &cobra.Command{Use: "widgets"}
	grp.AddCommand(&cobra.Command{
		Use:  "frobnicate",
		RunE: func(_ *cobra.Command, _ []string) error { return nil },
	})
	root.AddCommand(grp)

	applyMCPAnnotations(root)

	assertHints(t, root, "widgets frobnicate", "false", "")
}

// TestBlockedBashPaths_CoverExposedDestructive checks the guard's hand-maintained list
// against the annotations: anything marked destructive must also be blocked on the Bash
// surface, or an agent can reach it by shelling out. "courses delete" was missing.
func TestBlockedBashPaths_CoverExposedDestructive(t *testing.T) {
	root := NewRootCmd()

	blocked := make(map[string]bool, len(blockedBashPaths))
	for _, p := range blockedBashPaths {
		blocked[p] = true
	}

	walkCommands(root, func(cmd *cobra.Command) {
		if cmd.Annotations[ophis.AnnotationDestructive] != "true" {
			return
		}
		path := strings.TrimPrefix(cmd.CommandPath(), root.Name()+" ")
		if !blocked[path] {
			t.Errorf("%q is destructive but not in blockedBashPaths", path)
		}
	})
}
