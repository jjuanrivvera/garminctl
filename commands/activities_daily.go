package commands

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	gm "github.com/llehouerou/go-garmin"

	"github.com/jjuanrivvera/garminctl/internal/api"
	"github.com/jjuanrivvera/garminctl/internal/garmin"
)

// fetchActivitiesDaily returns the activities recorded on one day. go-garmin's
// ActivityService.List paginates by index only (start/limit, no date filter), so — like steps —
// this goes through the raw client, using the search endpoint's own date parameters. The
// response unmarshals into go-garmin's Activity type so rendered fields match `activities list`.
func fetchActivitiesDaily(ctx context.Context, c *gm.Client, d time.Time) (any, error) {
	ensureFreshToken(ctx, c)

	sessionJSON, err := garmin.DumpSession(c)
	if err != nil {
		return nil, err
	}
	token, baseURL, err := garmin.SessionToken(sessionJSON)
	if err != nil {
		return nil, err
	}
	var opts []api.Option
	if testHTTPClient != nil { // test seam: mock the raw transport, like the api command
		opts = append(opts, api.WithHTTPClient(testHTTPClient))
	}
	raw := api.New(token, baseURL, opts...)

	day := d.Format("2006-01-02")
	path := fmt.Sprintf(
		"/activitylist-service/activities/search/activities?startDate=%s&endDate=%s&limit=50",
		day, day)
	body, err := raw.Do(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, err
	}

	var acts []gm.Activity
	if err := json.Unmarshal(body, &acts); err != nil {
		return nil, fmt.Errorf("decode activities for %s: %w", day, err)
	}
	return acts, nil
}
