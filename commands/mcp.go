package commands

import (
	"github.com/njayp/ophis"
	"github.com/spf13/cobra"
)

// excludedFromMCP are command-name substrings kept out of the MCP tool surface: setup/meta
// commands an agent should not drive, the raw `api` escape hatch (which would bypass the typed
// surface), and `workouts` (create/update/delete/schedule). The `mcp` and `agent` subtrees are
// excluded too so an agent can neither re-enter the server nor disable its own guardrails.
//
// `courses` is NOT excluded — its reads (list/get/download-fit/download-gpx) are the point — so
// its two writes are handled by annotation instead: "courses import" falls to the write default
// and "courses delete" is destructive, and the agent guard blocks it on the Bash surface too.
var excludedFromMCP = []string{
	"agent", "auth", "config", "alias", "init", "doctor", "completion", "version", "api", "workouts",
	// `sync` bulk-fetches a date range — not something an agent should trigger. (`history`, a local
	// read, stays exposed.)
	"sync",
}

// secretFlags must never reach the MCP tool schema: an agent must not switch the account it
// runs as. The server uses whatever profile is active at startup.
var secretFlags = []string{"profile"}

func init() {
	registerCommand(func(root *cobra.Command) {
		// ophis walks the command tree and exposes each runnable leaf as an MCP tool, replaying
		// the cobra command on invocation so tools reuse the same client, keyring, and profile.
		root.AddCommand(ophis.Command(&ophis.Config{
			ToolNamePrefix: "garmin",
			Selectors: []ophis.Selector{{
				CmdSelector:           ophis.ExcludeCmdsContaining(excludedFromMCP...),
				InheritedFlagSelector: ophis.ExcludeFlags(secretFlags...),
			}},
		}))
	})
}
