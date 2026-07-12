package commands

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/jjuanrivvera/garminctl/internal/config"
)

func init() {
	registerCommand(func(root *cobra.Command) {
		cmd := &cobra.Command{
			Use:   "config",
			Short: "Manage garminctl configuration and profiles",
			Long:  "Inspect the config file and switch the default Garmin account (profile).",
		}
		cmd.AddCommand(configListCmd(), configUseCmd(), configPathCmd())
		root.AddCommand(cmd)
	})
}

func configListCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "list",
		Short: "List configured profiles",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, _ []string) error {
			c, err := config.Load()
			if err != nil {
				return err
			}
			if len(c.Profiles) == 0 {
				fmt.Fprintln(cmd.OutOrStdout(), "no profiles yet — run `garminctl auth import` or `garminctl init`")
				return nil
			}
			for _, p := range c.Profiles {
				marker := "  "
				if p == c.DefaultProfile {
					marker = "* " // active default
				}
				fmt.Fprintf(cmd.OutOrStdout(), "%s%s\n", marker, p)
			}
			return nil
		},
	}
}

func configUseCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "use <profile>",
		Short: "Set the default profile",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := config.Load()
			if err != nil {
				return err
			}
			if !containsProfile(c.Profiles, args[0]) {
				return fmt.Errorf("unknown profile %q — run `garminctl config list`", args[0])
			}
			c.DefaultProfile = args[0]
			if err := config.Save(c); err != nil {
				return err
			}
			fmt.Fprintf(cmd.OutOrStdout(), "default profile is now %q\n", args[0])
			return nil
		},
	}
}

func configPathCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "path",
		Short: "Print the config file path",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, _ []string) error {
			p, err := config.Path()
			if err != nil {
				return err
			}
			fmt.Fprintln(cmd.OutOrStdout(), p)
			return nil
		},
	}
}

func containsProfile(profiles []string, name string) bool {
	for _, p := range profiles {
		if p == name {
			return true
		}
	}
	return false
}
