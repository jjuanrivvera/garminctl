package commands

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"

	"github.com/jjuanrivvera/garminctl/internal/config"
	"github.com/jjuanrivvera/garminctl/internal/garmin"
)

func init() {
	registerCommand(func(root *cobra.Command) {
		var from string
		cmd := &cobra.Command{
			Use:   "init",
			Short: "Guided first-time setup (import existing tokens)",
			Long: `First-run setup. Looks for an existing garth / python-garminconnect token directory
(default ~/.garminconnect) and imports it into the keyring as a profile, so a working Python
setup carries straight over. If none is found, prints the two ways to authenticate.

The profile name comes from the global --profile flag (default "default").`,
			Args: cobra.NoArgs,
			RunE: func(cmd *cobra.Command, _ []string) error {
				out := cmd.OutOrStdout()
				dir := from
				if dir == "" {
					home, err := os.UserHomeDir()
					if err != nil {
						return err
					}
					dir = filepath.Join(home, ".garminconnect")
				}

				if !fileExists(filepath.Join(dir, "oauth1_token.json")) {
					fmt.Fprintf(out, `No tokens found at %s.

To get started, either:
  • import existing garth tokens:  garminctl auth import --from <dir>
  • or log in fresh:               garminctl auth login --email you@example.com
`, dir)
					return nil
				}

				profile := gf.profile
				if profile == "" {
					profile = "default"
				}
				sessionJSON, err := garmin.ImportGarth(dir)
				if err != nil {
					return err
				}
				if err := keyringStore().Set(profile, sessionJSON); err != nil {
					return err
				}
				c, err := config.Load()
				if err != nil {
					return err
				}
				c.AddProfile(profile)
				if err := config.Save(c); err != nil {
					return err
				}
				fmt.Fprintf(out, `✓ imported tokens from %s into profile %q

Try:
  garminctl sleep
  garminctl auth status
`, dir, profile)
				return nil
			},
		}
		cmd.Flags().StringVar(&from, "from", "", "garth token dir to import (default ~/.garminconnect)")
		root.AddCommand(cmd)
	})
}

// fileExists reports whether path is present, so the caller can branch without a captured
// error in scope (which the nilerr linter would flag on the informational return).
func fileExists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}
