package commands

import (
	"encoding/json"
	"fmt"

	"github.com/spf13/cobra"

	"github.com/jjuanrivvera/garminctl/internal/api"
	"github.com/jjuanrivvera/garminctl/internal/garmin"
	"github.com/jjuanrivvera/garminctl/internal/output"
)

func init() {
	registerCommand(func(root *cobra.Command) {
		var method, data string
		cmd := &cobra.Command{
			Use:   "api <path>",
			Short: "Make a raw authenticated request to the Garmin Connect API",
			Long: `Escape hatch for Connect endpoints the typed commands don't wrap.

<path> is everything after https://connectapi.garmin.com — for example
/usersummary-service/usersummary/daily?calendarDate=2026-07-10

The request is signed with the active profile's session token. The typed commands refresh that
token through go-garmin and persist it; if a raw call returns HTTP 401, run any resource
command once (e.g. ` + "`garminctl sleep`" + `) to refresh, then retry. With --dry-run, prints the
equivalent curl (token redacted) instead of sending.`,
			Args: cobra.ExactArgs(1),
			Example: `  garminctl api /userprofile-service/userprofile
  garminctl api -X GET /usersummary-service/usersummary/daily
  garminctl --dry-run api /userprofile-service/userprofile`,
			RunE: func(cmd *cobra.Command, args []string) error {
				ctx := cmd.Context()
				c, save, _, err := getClient(ctx)
				if err != nil {
					return err
				}
				defer func() { _ = save() }()

				sessionJSON, err := garmin.DumpSession(c)
				if err != nil {
					return err
				}
				token, baseURL, err := garmin.SessionToken(sessionJSON)
				if err != nil {
					return err
				}

				opts := []api.Option{}
				if gf.dryRun {
					opts = append(opts, api.WithDryRun(cmd.OutOrStdout()))
				}
				if testHTTPClient != nil { // test seam: mock the raw transport too
					opts = append(opts, api.WithHTTPClient(testHTTPClient))
				}
				client := api.New(token, baseURL, opts...)

				var body []byte
				if data != "" {
					body = []byte(data)
				}
				resp, err := client.Do(ctx, method, args[0], body)
				if err != nil {
					return err
				}
				if resp == nil { // --dry-run already wrote the curl
					return nil
				}

				// Render structured JSON through the shared formatter; fall back to the raw body
				// for non-JSON responses so nothing is lost.
				var v any
				if json.Unmarshal(resp, &v) != nil {
					fmt.Fprintln(cmd.OutOrStdout(), string(resp))
					return nil
				}
				return output.Render(cmd.OutOrStdout(), gf.output, v)
			},
		}
		cmd.Flags().StringVarP(&method, "method", "X", "GET", "HTTP method (GET|POST|PUT|DELETE)")
		cmd.Flags().StringVar(&data, "data", "", "request body (JSON) for POST/PUT")
		root.AddCommand(cmd)
	})
}
