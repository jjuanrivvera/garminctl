package commands

import (
	"bytes"
	"encoding/json"

	"github.com/spf13/cobra"

	"github.com/llehouerou/go-garmin/endpoint"
	"github.com/llehouerou/go-garmin/endpoint/definitions"

	"github.com/jjuanrivvera/garminctl/internal/output"
)

// curatedShadow lists registry command names already provided by a friendlier curated resource,
// so promoting the registry to the top level doesn't collide with them. Only `sleep` overlaps.
var curatedShadow = map[string]bool{"sleep": true}

func init() {
	registerCommand(func(root *cobra.Command) {
		for _, c := range newRegistryCommands() {
			root.AddCommand(c)
		}
	})
}

// newRegistryCommands promotes go-garmin's full endpoint registry (every documented Garmin
// Connect operation — metrics, activities, workouts, devices, exercises, calendar, biometric,
// …) to top-level commands, matching go-garmin's own `garmin` CLI. Each group resolves the
// active profile's client per invocation, re-renders the endpoint's JSON through garminctl's
// formatter (so `-o table|yaml|csv` works — go-garmin emits JSON only), and persists any OAuth2
// refresh back to the keyring. The curated resources (sleep, stress, body-composition, …) remain
// as friendlier shortcuts for the common reads.
func newRegistryCommands() []*cobra.Command {
	reg := endpoint.NewRegistry()
	definitions.RegisterAll(reg)
	gen := endpoint.NewCLIGenerator(reg)

	// Shared across the generated groups; only one command runs per process, so a single buffer
	// and save-func are safe (set in PreRun, consumed in PostRun).
	var buf bytes.Buffer
	var save func() error

	var cmds []*cobra.Command
	for _, c := range gen.GenerateCommands() {
		if curatedShadow[c.Name()] {
			continue
		}
		// The nearest PersistentPreRunE runs before any endpoint's RunE: wire the per-profile
		// client and redirect the generator's JSON into our buffer.
		c.PersistentPreRunE = func(cmd *cobra.Command, _ []string) error {
			cl, s, _, err := getClient(cmd.Context())
			if err != nil {
				return err
			}
			gen.SetClient(cl)
			save = s
			buf.Reset()
			gen.SetOutput(&buf)
			return nil
		}
		c.PersistentPostRunE = func(cmd *cobra.Command, _ []string) error {
			if buf.Len() > 0 {
				var v any
				if json.Unmarshal(buf.Bytes(), &v) == nil {
					if err := output.Render(cmd.OutOrStdout(), gf.output, v); err != nil {
						return err
					}
				} else { // non-JSON (rare) — pass through verbatim
					_, _ = cmd.OutOrStdout().Write(buf.Bytes())
				}
			}
			if save != nil {
				return save()
			}
			return nil
		}
		cmds = append(cmds, c)
	}
	return cmds
}
