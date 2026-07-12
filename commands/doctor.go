package commands

import (
	"fmt"
	"time"

	"github.com/spf13/cobra"

	"github.com/jjuanrivvera/garminctl/internal/config"
	"github.com/jjuanrivvera/garminctl/internal/garmin"
)

func init() {
	registerCommand(func(root *cobra.Command) {
		cmd := &cobra.Command{
			Use:   "doctor",
			Short: "Diagnose garminctl setup (config, keyring, sessions)",
			Long: `Offline health check: verifies the config file is readable and inspects each profile's
stored session — whether both tokens are present and whether the OAuth2 token is still valid.
An expired OAuth2 token is not a failure; the next call refreshes it from the OAuth1 token.`,
			Args: cobra.NoArgs,
			RunE: func(cmd *cobra.Command, _ []string) error {
				out := cmd.OutOrStdout()
				healthy := true

				if p, err := config.Path(); err == nil {
					fmt.Fprintf(out, "✓ config path: %s\n", p)
				} else {
					healthy = false
					fmt.Fprintf(out, "✗ config path: %v\n", err)
				}

				c, err := config.Load()
				if err != nil {
					fmt.Fprintf(out, "✗ load config: %v\n", err)
					return fmt.Errorf("doctor found problems")
				}
				if len(c.Profiles) == 0 {
					fmt.Fprintln(out, "! no profiles configured — run `garminctl init`")
					return nil
				}

				for _, p := range c.Profiles {
					sess, err := keyringStore().Get(p)
					if err != nil || sess == "" {
						healthy = false
						fmt.Fprintf(out, "✗ profile %q: no session in keyring\n", p)
						continue
					}
					expiry, authed, err := garmin.SessionInfo(sess)
					switch {
					case err != nil:
						healthy = false
						fmt.Fprintf(out, "✗ profile %q: unreadable session: %v\n", p, err)
					case !authed:
						healthy = false
						fmt.Fprintf(out, "✗ profile %q: session missing OAuth1/OAuth2 tokens\n", p)
					case time.Now().After(expiry):
						fmt.Fprintf(out, "✓ profile %q: OAuth2 expired — refreshes on next call\n", p)
					default:
						fmt.Fprintf(out, "✓ profile %q: token valid until %s\n", p, expiry.Format(time.RFC3339))
					}
				}

				if !healthy {
					return fmt.Errorf("doctor found problems")
				}
				fmt.Fprintln(out, "\nall checks passed")
				return nil
			},
		}
		root.AddCommand(cmd)
	})
}
