package api

import (
	"errors"
	"net/http"
	"strings"
	"testing"
)

func TestAPIErrorHints(t *testing.T) {
	cases := map[int]string{
		http.StatusUnauthorized:    "refresh",
		http.StatusForbidden:       "denied",
		http.StatusNotFound:        "no such endpoint",
		http.StatusTooManyRequests: "rate-limited",
		http.StatusBadGateway:      "server error",
	}
	for status, want := range cases {
		e := &APIError{Status: status, Path: "/x", Body: "oops"}
		got := e.Error()
		if !strings.Contains(got, want) {
			t.Errorf("status %d: hint missing %q in %q", status, want, got)
		}
		if !strings.Contains(got, "/x") {
			t.Errorf("status %d: path missing in %q", status, got)
		}
	}

	// A plain 400 has no canned hint but still reports status + body.
	e := &APIError{Status: http.StatusBadRequest, Path: "/y", Body: "bad"}
	if !strings.Contains(e.Error(), "HTTP 400") || !strings.Contains(e.Error(), "bad") {
		t.Errorf("400: %q", e.Error())
	}

	// errors.As lets callers branch on status.
	var target *APIError
	if !errors.As(error(e), &target) || target.Status != http.StatusBadRequest {
		t.Error("errors.As should recover APIError")
	}
}
