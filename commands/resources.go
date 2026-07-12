package commands

import (
	"context"
	"fmt"
	"time"

	gm "github.com/llehouerou/go-garmin"
	"github.com/spf13/cobra"

	"github.com/jjuanrivvera/garminctl/internal/output"
)

// dateFetch fetches a resource for a given date from an authenticated client.
type dateFetch func(ctx context.Context, c *gm.Client, date time.Time) (any, error)

func init() {
	registerCommand(func(root *cobra.Command) {
		root.AddCommand(
			newDateResource("body-composition", "Weight, BMI, and body fat for a day",
				func(ctx context.Context, c *gm.Client, d time.Time) (any, error) { return c.Weight.GetDaily(ctx, d) }),
			newDateResource("sleep", "Sleep stages, duration, and score for a day",
				func(ctx context.Context, c *gm.Client, d time.Time) (any, error) { return c.Sleep.GetDaily(ctx, d) }),
			newDateResource("stress", "All-day stress for a day",
				func(ctx context.Context, c *gm.Client, d time.Time) (any, error) {
					return c.Wellness.GetDailyStress(ctx, d)
				}),
			newDateResource("body-battery", "Body Battery events for a day",
				func(ctx context.Context, c *gm.Client, d time.Time) (any, error) {
					return c.Wellness.GetBodyBatteryEvents(ctx, d)
				}),
			newDateResource("heart-rate", "Daily and resting heart rate for a day",
				func(ctx context.Context, c *gm.Client, d time.Time) (any, error) {
					return c.Wellness.GetDailyHeartRate(ctx, d)
				}),
			newDateResource("respiration", "All-day respiration for a day",
				func(ctx context.Context, c *gm.Client, d time.Time) (any, error) {
					return c.Wellness.GetDailyRespiration(ctx, d)
				}),
			newDateResource("intensity-minutes", "Intensity minutes for a day",
				func(ctx context.Context, c *gm.Client, d time.Time) (any, error) {
					return c.Wellness.GetDailyIntensityMinutes(ctx, d)
				}),
		)
	})
}

// newDateResource builds a `garminctl <resource> [--date YYYY-MM-DD]` command that fetches one
// day's record for the active profile and renders it. Any OAuth2 refresh during the call is
// persisted back to the keyring via the deferred save().
func newDateResource(use, short string, fetch dateFetch) *cobra.Command {
	var dateStr string
	cmd := &cobra.Command{
		Use:     use,
		Short:   short,
		Args:    cobra.NoArgs,
		Example: fmt.Sprintf("  garminctl --profile me %s --date 2026-07-10 -o json", use),
		RunE: func(cmd *cobra.Command, _ []string) error {
			date := time.Now()
			if dateStr != "" {
				d, err := time.Parse("2006-01-02", dateStr)
				if err != nil {
					return fmt.Errorf("invalid --date %q (want YYYY-MM-DD): %w", dateStr, err)
				}
				date = d
			}
			c, save, _, err := getClient(cmd.Context())
			if err != nil {
				return err
			}
			defer func() { _ = save() }()
			result, err := fetch(cmd.Context(), c, date)
			if err != nil {
				return err
			}
			return output.Render(cmd.OutOrStdout(), gf.output, result)
		},
	}
	cmd.Flags().StringVar(&dateStr, "date", "", "date YYYY-MM-DD (default: today)")
	return cmd
}
