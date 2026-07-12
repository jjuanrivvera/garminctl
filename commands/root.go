// Package commands holds the garminctl command tree.
package commands

import (
	"github.com/spf13/cobra"
)

// globalFlags holds the persistent, cross-command flags resolved once per invocation.
type globalFlags struct {
	profile string
	output  string
	noColor bool
	dryRun  bool
}

var gf globalFlags

// commandRegistrars is populated by each command file's init(); NewRootCmd applies them so
// adding a command is a single init() with zero edits to shared code.
var commandRegistrars []func(*cobra.Command)

func registerCommand(fn func(*cobra.Command)) { commandRegistrars = append(commandRegistrars, fn) }

// NewRootCmd builds the root command tree.
func NewRootCmd() *cobra.Command {
	root := &cobra.Command{
		Use:   "garminctl",
		Short: "Drive Garmin Connect from the terminal",
		Long: `garminctl reads your Garmin Connect health data — body composition, sleep, heart
rate, stress, body battery, respiration, and intensity minutes — plus the full Connect endpoint
surface via ` + "`connect`" + `, with named profiles for several accounts, OS-keyring token storage,
and table/json/yaml/csv output.`,
		Example: `  garminctl auth import --from ~/.garminconnect --profile juan
  garminctl sleep --date 2026-07-09 -o json
  garminctl stress
  garminctl --profile vane body-composition
  garminctl doctor`,
		SilenceUsage:  true,
		SilenceErrors: true,
	}
	p := root.PersistentFlags()
	p.StringVar(&gf.profile, "profile", "", "profile (Garmin account) to use; env GARMINCTL_PROFILE")
	p.StringVarP(&gf.output, "output", "o", "table", "output format: table|json|yaml|csv")
	p.BoolVar(&gf.noColor, "no-color", false, "disable colored output")
	p.BoolVar(&gf.dryRun, "dry-run", false, "print the equivalent request instead of sending it")
	for _, fn := range commandRegistrars {
		fn(root)
	}
	return root
}
