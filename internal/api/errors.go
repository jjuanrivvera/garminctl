package api

import (
	"fmt"
	"net/http"
	"strings"
)

// APIError is a failed Garmin Connect request carrying an actionable, status-keyed hint — so a
// 401 tells you to refresh instead of just printing "request failed". Callers can errors.As it
// to branch on Status; the Error() string is what the user sees.
type APIError struct {
	Status int    // HTTP status code
	Path   string // request path, for context
	Body   string // response body (may be empty), already status-checked
}

func (e *APIError) Error() string {
	msg := fmt.Sprintf("garmin api %s: HTTP %d", e.Path, e.Status)
	if h := hintForStatus(e.Status); h != "" {
		msg += " — " + h
	}
	if b := strings.TrimSpace(e.Body); b != "" {
		msg += ": " + b
	}
	return msg
}

// hintForStatus maps an HTTP status to a next step the user can actually take. Garmin Connect
// is a reverse-engineered API with no error envelope, so the status is the only signal.
func hintForStatus(status int) string {
	switch status {
	case http.StatusUnauthorized: // 401
		return "session token rejected — run a resource command (e.g. `garminctl sleep`) to refresh, or re-import with `garminctl auth import`"
	case http.StatusForbidden: // 403
		return "Garmin denied this request — the endpoint may require a different scope or your account lacks access"
	case http.StatusNotFound: // 404
		return "no such endpoint, or no data for that date — check the path against connectapi.garmin.com"
	case http.StatusTooManyRequests: // 429
		return "rate-limited by Garmin — wait and retry; garminctl already backs off on idempotent calls"
	}
	if status >= 500 {
		return "Garmin server error — usually transient; retry shortly"
	}
	return ""
}
