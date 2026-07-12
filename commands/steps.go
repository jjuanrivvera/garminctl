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

// DailySteps is one day's step summary from Garmin's
// /usersummary-service/stats/steps/daily/{start}/{end} endpoint. go-garmin declares a
// StepsService but implements no methods, and the endpoint is absent from its registry, so —
// unlike the other curated metrics — steps is fetched through the raw client (the same path the
// `api` escape hatch uses). Unknown fields are ignored, so the struct need not be exhaustive.
type DailySteps struct {
	CalendarDate  string  `json:"calendarDate"`
	TotalSteps    int     `json:"totalSteps"`
	StepGoal      int     `json:"stepGoal"`
	TotalDistance float64 `json:"totalDistance"`
}

// fetchSteps returns the step summary for one day. It matches the curatedResource dateFetch
// signature so `steps` registers, syncs, and caches offline exactly like the go-garmin-backed
// metrics; the difference is the raw request underneath.
func fetchSteps(ctx context.Context, c *gm.Client, d time.Time) (any, error) {
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
	path := fmt.Sprintf("/usersummary-service/stats/steps/daily/%s/%s", day, day)
	body, err := raw.Do(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, err
	}

	// The endpoint returns an array (one element per day in the range); we ask for a single day.
	var days []DailySteps
	if err := json.Unmarshal(body, &days); err != nil {
		return nil, fmt.Errorf("decode steps for %s: %w", day, err)
	}
	if len(days) == 0 {
		return nil, fmt.Errorf("no step data for %s", day)
	}
	return days[0], nil
}

// ensureFreshToken nudges go-garmin to refresh its OAuth2 token when the session is at or near
// expiry. go-garmin only refreshes inside its own typed calls (doAPI), so a following RAW request
// would otherwise carry a stale token and 401 — the very reliability gap the other reads avoid by
// going through go-garmin. GetSocialProfile is a cheap authenticated GET whose side effect (the
// refresh) is what we want; the deferred save() in the caller then persists the new token. Errors
// are swallowed: if the refresh probe fails, the real steps request surfaces the actual error.
func ensureFreshToken(ctx context.Context, c *gm.Client) {
	sessionJSON, err := garmin.DumpSession(c)
	if err != nil {
		return
	}
	expiry, _, err := garmin.SessionInfo(sessionJSON)
	if err != nil {
		return
	}
	if time.Now().Before(expiry.Add(-2 * time.Minute)) {
		return // still valid with margin; skip the probe so a fresh token costs no extra call
	}
	_, _ = c.UserProfile.GetSocialProfile(ctx)
}
