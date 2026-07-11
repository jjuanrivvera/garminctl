package commands

import (
	"os"

	"github.com/spf13/cobra"

	"github.com/llehouerou/go-garmin/endpoint"
	"github.com/llehouerou/go-garmin/endpoint/definitions"
)

func init() {
	registerCommand(func(root *cobra.Command) { root.AddCommand(newConnectCmd()) })
}

// newConnectCmd exposes the full go-garmin endpoint registry (every documented Garmin Connect
// operation) under `garminctl connect …`, so the CLI covers the complete API surface, not just
// the curated top-level resources. The active profile's client is resolved per invocation and
// any OAuth2 refresh is persisted back to the keyring afterward. Output is JSON.
func newConnectCmd() *cobra.Command {
	reg := endpoint.NewRegistry()
	definitions.RegisterAll(reg)
	gen := endpoint.NewCLIGenerator(reg)
	gen.SetOutput(os.Stdout)

	var save func() error
	parent := &cobra.Command{
		Use:   "connect",
		Short: "The full Garmin Connect endpoint surface (every documented operation)",
		Long: `connect exposes the complete Garmin Connect endpoint registry — every documented
operation, grouped by service (sleep, wellness, activities, metrics, devices, …). These are the
raw endpoints with JSON output; the curated top-level commands (body-composition, sleep, …) are
friendlier for the common cases.`,
		// The nearest PersistentPreRunE runs before any endpoint's RunE, so we can set the
		// per-profile client here (it's read at call time by the generated handlers).
		PersistentPreRunE: func(cmd *cobra.Command, _ []string) error {
			c, s, _, err := getClient(cmd.Context())
			if err != nil {
				return err
			}
			gen.SetClient(c)
			save = s
			return nil
		},
		PersistentPostRunE: func(_ *cobra.Command, _ []string) error {
			if save != nil {
				return save()
			}
			return nil
		},
	}
	parent.AddCommand(gen.GenerateCommands()...)
	return parent
}
