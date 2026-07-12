package commands

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/spf13/cobra"

	gm "github.com/llehouerou/go-garmin"

	"github.com/jjuanrivvera/garminctl/internal/config"
	"github.com/jjuanrivvera/garminctl/internal/garmin"
)

func init() {
	registerCommand(func(root *cobra.Command) { root.AddCommand(newAuthCmd()) })
}

func newAuthCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "auth",
		Short: "Manage Garmin Connect authentication",
		Long: `Store and verify the Garmin session for a profile. A session (OAuth1 + OAuth2 tokens)
lives in your OS keyring, keyed by profile. Bring one in with 'auth import' (from an existing
garth / python-garminconnect token dir) or 'auth login' (email + password).`,
	}
	cmd.AddCommand(newAuthImportCmd(), newAuthLoginCmd(), newAuthStatusCmd(), newAuthLogoutCmd())
	return cmd
}

func newAuthImportCmd() *cobra.Command {
	var from string
	cmd := &cobra.Command{
		Use:   "import",
		Short: "Import an existing garth / python-garminconnect token dir into the keyring",
		Long: `Migrate an existing garth session (oauth1_token.json + oauth2_token.json, e.g.
~/.garminconnect) into garminctl's keyring under the active profile. No login required — the
OAuth1 token (valid ~1 year) drives OAuth2 refresh from here, so this fixes the recurring
"username and password are required" failure of a cron that never refreshed its cached tokens.`,
		Example: `  garminctl --profile me auth import --from ~/.garminconnect
  garminctl --profile alt auth import --from ~/.garminconnect-alt`,
		Args: cobra.NoArgs,
		RunE: func(cmd *cobra.Command, _ []string) error {
			if from == "" {
				home, _ := os.UserHomeDir()
				from = filepath.Join(home, ".garminconnect")
			}
			sessionJSON, err := garmin.ImportGarth(from)
			if err != nil {
				return err
			}
			profile := config.Resolve(gf.profile)
			if err := store().Set(profile, sessionJSON); err != nil {
				return fmt.Errorf("store session: %w", err)
			}
			cfg, _ := config.Load()
			cfg.AddProfile(profile)
			_ = config.Save(cfg)
			exp, _, _ := garmin.SessionInfo(sessionJSON)
			fmt.Fprintf(cmd.ErrOrStderr(), "imported garth session for profile %q from %s (OAuth2 valid until %s)\n",
				profile, from, exp.Local().Format(time.RFC822))
			return nil
		},
	}
	cmd.Flags().StringVar(&from, "from", "", "garth token dir (default: ~/.garminconnect)")
	return cmd
}

func newAuthLoginCmd() *cobra.Command {
	var email string
	cmd := &cobra.Command{
		Use:     "login",
		Short:   "Log in with Garmin credentials and store the session in the keyring",
		Example: `  garminctl --profile me auth login`,
		Args:    cobra.NoArgs,
		RunE: func(cmd *cobra.Command, _ []string) error {
			var err error
			if email == "" {
				if email, err = promptLine(cmd, "Garmin email: "); err != nil {
					return err
				}
			}
			password, err := promptSecret(cmd, "Garmin password (hidden): ")
			if err != nil {
				return err
			}
			c := gm.New(gm.Options{})
			if err := c.Login(cmd.Context(), email, password); err != nil {
				return fmt.Errorf("login failed: %w", err)
			}
			sessionJSON, err := garmin.DumpSession(c)
			if err != nil {
				return err
			}
			profile := config.Resolve(gf.profile)
			if err := store().Set(profile, sessionJSON); err != nil {
				return err
			}
			cfg, _ := config.Load()
			cfg.AddProfile(profile)
			_ = config.Save(cfg)
			fmt.Fprintf(cmd.ErrOrStderr(), "logged in and stored session for profile %q\n", profile)
			return nil
		},
	}
	cmd.Flags().StringVar(&email, "email", "", "Garmin account email (omit to be prompted)")
	return cmd
}

func newAuthStatusCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "status",
		Short: "Show whether the active profile has a stored, valid session",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, _ []string) error {
			profile := config.Resolve(gf.profile)
			sessionJSON, err := store().Get(profile)
			if err != nil || sessionJSON == "" {
				return fmt.Errorf("no session for profile %q — run `garminctl auth import` or `garminctl auth login`", profile)
			}
			exp, authed, err := garmin.SessionInfo(sessionJSON)
			if err != nil {
				return fmt.Errorf("stored session is unreadable: %w", err)
			}
			w := cmd.OutOrStdout()
			fmt.Fprintf(w, "profile:        %s\n", profile)
			fmt.Fprintf(w, "authenticated:  %v\n", authed)
			fmt.Fprintf(w, "OAuth2 expiry:  %s (auto-refreshed on use)\n", exp.Local().Format(time.RFC822))
			return nil
		},
	}
}

func newAuthLogoutCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "logout",
		Short: "Remove the stored session for the active profile",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, _ []string) error {
			profile := config.Resolve(gf.profile)
			if err := store().Delete(profile); err != nil {
				return err
			}
			fmt.Fprintf(cmd.ErrOrStderr(), "removed session for profile %q\n", profile)
			return nil
		},
	}
}
