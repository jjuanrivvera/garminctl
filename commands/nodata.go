package commands

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/spf13/cobra"
)

// A day the watch never synced does not come back as an error: Garmin answers 200 with a
// fully-null body, and Go's zero values turn those nulls into 0. "No data" then looks exactly
// like a measurement — and a resting heart rate of 0 is physiologically impossible, but a
// script reading the JSON has no way to know that (issue #16: an automated daily summary
// reported "resting heart rate: 0" as if the watch had measured it).
//
// The output itself is left alone: changing the shape of a payload would break every consumer
// to help the few that check. Instead an empty day says so on stderr and sets a distinct exit
// status, so a script can branch on "nothing to report" without parsing anything.

// ExitNoData is returned when a read succeeded but the day holds no measurement. It is
// deliberately not 1: "the request failed" and "there is nothing there" are different answers
// and a caller usually wants to treat them differently.
const ExitNoData = 4

// requestEcho are fields the API returns because they were asked for, not because anything was
// measured. A payload carrying only these is empty however full it looks — the sleep response
// for an unsynced day still echoes its calendarDate.
var requestEcho = map[string]bool{
	"calendardate":  true,
	"date":          true,
	"userprofileid": true,
	"userprofilepk": true,
	"uuid":          true,
}

// hasData reports whether a payload carries a single real value. It works on the JSON form
// rather than the Go struct because the structs come from a third-party client (go-garmin)
// whose numeric fields are plain ints: by the time the value reaches here, an absent field and
// a measured zero are the same 0. The JSON keeps nulls, which is the distinction that matters.
//
// It is a heuristic, and it errs toward saying there IS data: anything it cannot parse counts
// as present, so a false positive costs an exit status, never the output.
func hasData(v any) bool {
	b, err := json.Marshal(v)
	if err != nil {
		return true
	}
	var decoded any
	if err := json.Unmarshal(b, &decoded); err != nil {
		return true
	}
	return anyMeasured(decoded)
}

func anyMeasured(v any) bool {
	switch t := v.(type) {
	case nil:
		return false
	case bool:
		return t
	case float64:
		return t != 0
	case string:
		return strings.TrimSpace(t) != ""
	case []any:
		for _, e := range t {
			if anyMeasured(e) {
				return true
			}
		}
		return false
	case map[string]any:
		for k, e := range t {
			if requestEcho[strings.ToLower(k)] {
				continue
			}
			if anyMeasured(e) {
				return true
			}
		}
		return false
	default:
		return true
	}
}

// noteIfEmpty says on stderr that a day holds nothing and records it for the exit status.
// It returns whether the day was empty so a caller can skip work that assumes data.
func noteIfEmpty(cmd *cobra.Command, resource, date string, v any) bool {
	if hasData(v) {
		return false
	}
	gf.noData = true
	fmt.Fprintf(cmd.ErrOrStderr(),
		"note: no %s data for %s — the watch may not have synced that day (exit %d)\n",
		resource, date, ExitNoData)
	return true
}
