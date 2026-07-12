package commands

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

// blockedBashPaths are garminctl subcommand paths an agent must never run on the Bash surface.
// garminctl is read-focused, but promoting go-garmin's full registry brings in workout
// management (create/update/delete/schedule/unschedule) — the only typed writes — plus the two
// indirect vectors (deleting the stored session, minting an alias that expands to a blocked
// command). The raw `api` hatch is gated separately (by HTTP method) inside the hook.
var blockedBashPaths = []string{
	"alias set",
	"auth logout",
	"workouts create",
	"workouts update",
	"workouts delete",
	"workouts schedule",
	"workouts unschedule",
}

// bashPattern is the Bash permission pattern a host gates for a garminctl subcommand path.
func bashPattern(path string) string { return "Bash(garminctl " + path + ":*)" }

func init() {
	registerCommand(func(root *cobra.Command) {
		agentCmd := &cobra.Command{
			Use:   "agent",
			Short: "AI-agent integration helpers",
			Long:  "Generate safety configuration for AI agents that drive garminctl.",
		}

		var host, out string
		guard := &cobra.Command{
			Use:   "guard --host <claude-code|codex|opencode>",
			Short: "Generate agent-safety config that blocks mutating garminctl operations",
			Long: `garminctl's surface is read-focused Garmin Connect data, so the guard blocks the
handful of mutation vectors rather than a rich destructive taxonomy:

  • "workouts create|update|delete|schedule|unschedule" — the only typed writes;
  • the raw "api" escape hatch with a write method (-X POST|PUT|DELETE|PATCH);
  • "auth logout" — deletes the stored session from the keyring;
  • "alias set" — could mint a shorthand that expands to a blocked command before cobra parses.

Reads — every resource, every promoted registry endpoint, and "api" GET — are allowed.

For claude-code the output includes a PreToolUse hook (.claude/hooks/garminctl-guard.sh) that
strips quote/backslash obfuscation, matches the binary even when path-invoked
(./bin/garminctl, /usr/local/bin/garminctl), and gates the "api" hatch by HTTP method.

MCP-only operation is the hard guarantee; the Bash rails are best-effort — the hook defeats
quoting and path prefixes, but not variable indirection (m=DELETE; garminctl api x -X $m) or
shell aliases.`,
			Example: `  garminctl agent guard --host claude-code
  garminctl agent guard --host codex --out ~/.codex/config.toml
  garminctl agent guard --host opencode`,
			RunE: func(cmd *cobra.Command, _ []string) error {
				var content string
				var err error
				switch host {
				case "claude-code", "claude":
					content, err = renderClaudeCode()
				case "codex":
					content, err = renderCodex()
				case "opencode":
					content, err = renderOpenCode()
				default:
					return fmt.Errorf("unknown --host %q (want claude-code|codex|opencode)", host)
				}
				if err != nil {
					return err
				}
				if out != "" {
					if err := os.WriteFile(out, []byte(content), 0o600); err != nil {
						return err
					}
					fmt.Fprintf(cmd.ErrOrStderr(), "wrote %s safety config to %s\n", host, out)
					return nil
				}
				fmt.Fprint(cmd.OutOrStdout(), content)
				return nil
			},
		}
		guard.Flags().StringVar(&host, "host", "", "target agent host: claude-code|codex|opencode (required)")
		guard.Flags().StringVar(&out, "out", "", "write to this file instead of stdout")
		_ = guard.MarkFlagRequired("host")

		agentCmd.AddCommand(guard)
		root.AddCommand(agentCmd)
	})
}
