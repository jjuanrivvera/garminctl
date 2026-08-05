package commands

import (
	"strconv"

	"github.com/njayp/ophis"
	"github.com/spf13/cobra"
)

// MCP tool annotations (cliwright GOAL.md §3b).
//
// Hosts running a read-only MCP session allow a tool only when readOnlyHint is strictly
// true, and treat a missing annotation as a write rather than as unknown. With nothing
// annotated the host finds no allowed tool and drops the entire server, so every command
// has to declare itself.
//
// garminctl is read-dominant, but the classification is still fail-safe: only the explicit
// allowlist below is advertised read-only, so a verb this file has not seen costs an
// approval prompt instead of being auto-approved. That matters more here than the list
// length suggests — the surface tracks go-garmin's registry, so new commands arrive without
// anyone editing this file.

// readLeaves are leaf command names that only read Garmin Connect (or the local store).
// Most of the surface is bare metric nouns rather than verbs, so they are enumerated.
var readLeaves = map[string]bool{
	"acclimation": true, "activities": true, "activities-daily": true, "body-battery": true,
	"body-composition": true, "categories": true, "check": true, "daily": true,
	"details": true, "display": true, "download-fit": true, "download-gpx": true,
	"endurance": true, "equipment": true, "exercise-sets": true, "ftp": true,
	"ftp-range": true, "gear": true, "get": true, "heart-rate": true,
	"hill": true, "history": true, "hr-zones": true, "hrv": true,
	"intensity-minutes": true, "lactate-threshold": true, "list": true, "load-balance": true,
	"lt-hr-range": true, "lt-speed-range": true, "muscles": true, "power-weight": true,
	"power-zones": true, "race-predictions": true, "range": true, "readiness": true,
	"respiration": true, "settings": true, "sleep": true, "social": true,
	"split-summaries": true, "splits": true, "spo2": true, "stats": true,
	"steps": true, "stress": true, "training-readiness": true, "training-status": true,
	"training-status-daily": true, "typed-splits": true, "types": true, "vo2max": true,
	"vo2max-range": true, "weather": true,
}

// destructiveLeaves cannot be undone. "courses delete" permanently removes a course from
// the Garmin Connect account.
var destructiveLeaves = map[string]bool{
	"delete":     true,
	"unschedule": true,
}

// applyMCPAnnotations stamps hints on every command in the tree. Call it once, after the
// tree is fully assembled — ophis reads cmd.Annotations at run time, in registerTools.
func applyMCPAnnotations(root *cobra.Command) {
	walkCommands(root, func(cmd *cobra.Command) {
		if !cmd.Runnable() || cmd.Hidden || cmd.Name() == "help" {
			return
		}
		switch {
		case destructiveLeaves[cmd.Name()]:
			setMCPAnnotations(cmd, false, true)
		case readLeaves[cmd.Name()]:
			setMCPAnnotations(cmd, true, false)
		default:
			// Fail-safe: "import", "update", the workouts writes, and anything new.
			setMCPAnnotations(cmd, false, false)
		}
	})
}

// walkCommands visits every command below root, leaves first.
func walkCommands(root *cobra.Command, visit func(*cobra.Command)) {
	var walk func(*cobra.Command)
	walk = func(c *cobra.Command) {
		for _, sub := range c.Commands() {
			walk(sub)
			visit(sub)
		}
	}
	walk(root)
}

// setMCPAnnotations merges MCP hints into cmd.Annotations, preserving unrelated keys.
//
// destructiveHint is written only when true; leaving it off lets the host apply the spec
// default (true when readOnlyHint is false), which fails safe.
func setMCPAnnotations(cmd *cobra.Command, readOnly, destructive bool) {
	if cmd.Annotations == nil {
		cmd.Annotations = map[string]string{}
	}
	cmd.Annotations[ophis.AnnotationReadOnly] = strconv.FormatBool(readOnly)
	cmd.Annotations[ophis.AnnotationOpenWorld] = "true"
	if destructive {
		cmd.Annotations[ophis.AnnotationDestructive] = "true"
	}
}
