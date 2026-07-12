package commands

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/spf13/cobra"
)

func init() {
	registerCommand(func(root *cobra.Command) {
		var fromStr, toStr, metricsCSV string
		cmd := &cobra.Command{
			Use:   "sync",
			Short: "Backfill daily metrics into the local store for offline use",
			Long: `sync fetches a date range of daily metrics for the active profile and records them in the
local SQLite store, so ` + "`garminctl --offline <metric>`" + ` and ` + "`garminctl history`" + ` work without the
network. Re-syncing a date overwrites it with fresh data. Defaults to the last 7 days and all
curated metrics.`,
			Args: cobra.NoArgs,
			Example: `  garminctl sync
  garminctl sync --from 2026-01-01 --to 2026-07-10
  garminctl sync --metrics sleep,body-composition --from 2026-06-01`,
			RunE: func(cmd *cobra.Command, _ []string) error {
				to := time.Now()
				if toStr != "" {
					t, err := time.Parse("2006-01-02", toStr)
					if err != nil {
						return fmt.Errorf("invalid --to %q (want YYYY-MM-DD): %w", toStr, err)
					}
					to = t
				}
				from := to.AddDate(0, 0, -6) // last 7 days inclusive
				if fromStr != "" {
					f, err := time.Parse("2006-01-02", fromStr)
					if err != nil {
						return fmt.Errorf("invalid --from %q (want YYYY-MM-DD): %w", fromStr, err)
					}
					from = f
				}
				if from.After(to) {
					return fmt.Errorf("--from %s is after --to %s", from.Format("2006-01-02"), to.Format("2006-01-02"))
				}
				metrics, err := selectMetrics(metricsCSV)
				if err != nil {
					return err
				}

				c, save, profile, err := getClient(cmd.Context())
				if err != nil {
					return err
				}
				defer func() { _ = save() }()
				st, err := openStore()
				if err != nil {
					return err
				}
				defer func() { _ = st.Close() }()

				var days, stored, failed int
				for d := from; !d.After(to); d = d.AddDate(0, 0, 1) {
					days++
					key := d.Format("2006-01-02")
					for _, r := range metrics {
						result, err := r.fetch(cmd.Context(), c, d)
						if err != nil {
							failed++
							fmt.Fprintf(cmd.ErrOrStderr(), "  ! %s %s: %v\n", key, r.name, err)
							continue
						}
						b, mErr := json.Marshal(result)
						if mErr != nil {
							failed++
							continue
						}
						if err := st.Put(profile, r.name, key, b); err != nil {
							failed++
							continue
						}
						stored++
					}
				}
				fmt.Fprintf(cmd.OutOrStdout(),
					"synced %d day(s) × %d metric(s) → %d stored, %d failed (profile %q)\n",
					days, len(metrics), stored, failed, profile)
				return nil
			},
		}
		cmd.Flags().StringVar(&fromStr, "from", "", "start date YYYY-MM-DD (default: 7 days ago)")
		cmd.Flags().StringVar(&toStr, "to", "", "end date YYYY-MM-DD (default: today)")
		cmd.Flags().StringVar(&metricsCSV, "metrics", "", "comma-separated metrics to sync (default: all)")
		root.AddCommand(cmd)
	})
}

// selectMetrics filters curatedResources by a comma-separated name list (empty selects all).
func selectMetrics(csv string) ([]curatedResource, error) {
	if strings.TrimSpace(csv) == "" {
		return curatedResources, nil
	}
	byName := map[string]curatedResource{}
	var all []string
	for _, r := range curatedResources {
		byName[r.name] = r
		all = append(all, r.name)
	}
	var out []curatedResource
	for _, n := range strings.Split(csv, ",") {
		n = strings.TrimSpace(n)
		r, ok := byName[n]
		if !ok {
			return nil, fmt.Errorf("unknown metric %q (have: %s)", n, strings.Join(all, ", "))
		}
		out = append(out, r)
	}
	return out, nil
}
