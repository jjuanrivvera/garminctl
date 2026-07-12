package commands

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/spf13/cobra"

	"github.com/jjuanrivvera/garminctl/internal/config"
	"github.com/jjuanrivvera/garminctl/internal/output"
)

func init() {
	registerCommand(func(root *cobra.Command) {
		var fromStr, toStr string
		cmd := &cobra.Command{
			Use:   "history <metric>",
			Short: "Query the local store for a metric across a date range (offline)",
			Long: `history reads the offline store (populated by ` + "`garminctl sync`" + ` or by earlier reads) for
one metric over a date range and renders one row per day — pair it with -o csv for a trend you can
open in a spreadsheet. It never hits the network. Defaults to the last 30 days.`,
			Args: cobra.ExactArgs(1),
			Example: `  garminctl history body-composition --from 2026-01-01 -o csv
  garminctl history sleep --from 2026-06-01 --to 2026-07-10`,
			RunE: func(cmd *cobra.Command, args []string) error {
				metric := args[0]
				to := time.Now()
				if toStr != "" {
					t, err := time.Parse("2006-01-02", toStr)
					if err != nil {
						return fmt.Errorf("invalid --to %q (want YYYY-MM-DD): %w", toStr, err)
					}
					to = t
				}
				from := to.AddDate(0, 0, -29) // last 30 days inclusive
				if fromStr != "" {
					f, err := time.Parse("2006-01-02", fromStr)
					if err != nil {
						return fmt.Errorf("invalid --from %q (want YYYY-MM-DD): %w", fromStr, err)
					}
					from = f
				}

				profile := config.Resolve(gf.profile)
				st, err := openStore()
				if err != nil {
					return err
				}
				defer func() { _ = st.Close() }()

				fromKey, toKey := from.Format("2006-01-02"), to.Format("2006-01-02")
				samples, err := st.Range(profile, metric, fromKey, toKey)
				if err != nil {
					return err
				}
				if len(samples) == 0 {
					return fmt.Errorf("no offline data for %s between %s and %s — run `garminctl sync`", metric, fromKey, toKey)
				}

				// One row per day: date + the metric's top-level fields flattened in, so -o csv/table
				// produces a real trend table. A non-object payload keeps a single "data" column.
				rows := make([]any, 0, len(samples))
				for _, s := range samples {
					row := map[string]any{"date": s.Date}
					var obj map[string]any
					if json.Unmarshal(s.Data, &obj) == nil {
						for k, v := range obj {
							if k != "date" {
								row[k] = v
							}
						}
					} else {
						var raw any
						_ = json.Unmarshal(s.Data, &raw)
						row["data"] = raw
					}
					rows = append(rows, row)
				}
				return output.Render(cmd.OutOrStdout(), gf.output, rows)
			},
		}
		cmd.Flags().StringVar(&fromStr, "from", "", "start date YYYY-MM-DD (default: 30 days ago)")
		cmd.Flags().StringVar(&toStr, "to", "", "end date YYYY-MM-DD (default: today)")
		root.AddCommand(cmd)
	})
}
