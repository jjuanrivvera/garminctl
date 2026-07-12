package commands

import (
	"context"
	"fmt"
	"time"

	gm "github.com/llehouerou/go-garmin"
	"github.com/spf13/cobra"

	"github.com/jjuanrivvera/garminctl/internal/config"
	"github.com/jjuanrivvera/garminctl/internal/output"
)

// dateFetch fetches a resource for a given date from an authenticated client.
type dateFetch func(ctx context.Context, c *gm.Client, date time.Time) (any, error)

// curatedResource is one friendly daily-metric shortcut.
type curatedResource struct {
	name  string
	short string
	fetch dateFetch
}

// curatedResources are the daily health metrics exposed as friendly top-level shortcuts and
// backfilled to the offline store. Shared by the command registration and `garminctl sync`.
var curatedResources = []curatedResource{
	{"body-composition", "Weight, BMI, and body fat for a day",
		func(ctx context.Context, c *gm.Client, d time.Time) (any, error) { return c.Weight.GetDaily(ctx, d) }},
	{"sleep", "Sleep stages, duration, and score for a day",
		func(ctx context.Context, c *gm.Client, d time.Time) (any, error) { return c.Sleep.GetDaily(ctx, d) }},
	{"stress", "All-day stress for a day",
		func(ctx context.Context, c *gm.Client, d time.Time) (any, error) {
			return c.Wellness.GetDailyStress(ctx, d)
		}},
	{"body-battery", "Body Battery events for a day",
		func(ctx context.Context, c *gm.Client, d time.Time) (any, error) {
			return c.Wellness.GetBodyBatteryEvents(ctx, d)
		}},
	{"heart-rate", "Daily and resting heart rate for a day",
		func(ctx context.Context, c *gm.Client, d time.Time) (any, error) {
			return c.Wellness.GetDailyHeartRate(ctx, d)
		}},
	{"respiration", "All-day respiration for a day",
		func(ctx context.Context, c *gm.Client, d time.Time) (any, error) {
			return c.Wellness.GetDailyRespiration(ctx, d)
		}},
	{"intensity-minutes", "Intensity minutes for a day",
		func(ctx context.Context, c *gm.Client, d time.Time) (any, error) {
			return c.Wellness.GetDailyIntensityMinutes(ctx, d)
		}},
	// steps has no go-garmin service method (StepsService is an unimplemented stub) nor a registry
	// endpoint, so its fetch goes through the raw client — see fetchSteps in steps.go.
	{"steps", "Daily step count, goal, and distance for a day", fetchSteps},
}

func init() {
	registerCommand(func(root *cobra.Command) {
		for _, r := range curatedResources {
			root.AddCommand(newDateResource(r))
		}
	})
}

// parseDate parses a --date flag (YYYY-MM-DD), defaulting to today.
func parseDate(dateStr string) (time.Time, error) {
	if dateStr == "" {
		return time.Now(), nil
	}
	d, err := time.Parse("2006-01-02", dateStr)
	if err != nil {
		return time.Time{}, fmt.Errorf("invalid --date %q (want YYYY-MM-DD): %w", dateStr, err)
	}
	return d, nil
}

// newDateResource builds a `garminctl <resource> [--date YYYY-MM-DD]` command. Online, it fetches
// one day's record and also caches it to the offline store; with --offline it serves that day
// from the store instead. Any OAuth2 refresh during a live call is persisted via the deferred
// save().
func newDateResource(r curatedResource) *cobra.Command {
	var dateStr string
	cmd := &cobra.Command{
		Use:     r.name,
		Short:   r.short,
		Args:    cobra.NoArgs,
		Example: fmt.Sprintf("  garminctl --profile me %s --date 2026-07-10 -o json", r.name),
		RunE: func(cmd *cobra.Command, _ []string) error {
			date, err := parseDate(dateStr)
			if err != nil {
				return err
			}
			dateKey := date.Format("2006-01-02")
			profile := config.Resolve(gf.profile)

			if gf.offline { // serve the day from the local store — no network
				v, ok, err := offlineSample(profile, r.name, dateKey)
				if err != nil {
					return err
				}
				if !ok {
					return fmt.Errorf("no offline data for %s on %s — run `garminctl sync` while online", r.name, dateKey)
				}
				return output.Render(cmd.OutOrStdout(), gf.output, v)
			}

			c, save, _, err := getClient(cmd.Context())
			if err != nil {
				return err
			}
			defer func() { _ = save() }()
			result, err := r.fetch(cmd.Context(), c, date)
			if err != nil {
				return err
			}
			cacheSample(profile, r.name, dateKey, result) // grow the offline store as you read
			return output.Render(cmd.OutOrStdout(), gf.output, result)
		},
	}
	cmd.Flags().StringVar(&dateStr, "date", "", "date YYYY-MM-DD (default: today)")
	return cmd
}
