package api

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestDoSuccessAndHeaders(t *testing.T) {
	var gotAuth, gotUA string
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		gotAuth = r.Header.Get("Authorization")
		gotUA = r.Header.Get("User-Agent")
		_, _ = w.Write([]byte(`{"ok":true}`))
	}))
	defer srv.Close()

	c := New("tok-abcd1234", srv.URL)
	body, err := c.Do(t.Context(), "GET", "/usersummary-service/usersummary/daily", nil)
	if err != nil {
		t.Fatalf("Do: %v", err)
	}
	if !strings.Contains(string(body), "ok") {
		t.Errorf("body = %q", body)
	}
	if gotAuth != "Bearer tok-abcd1234" {
		t.Errorf("auth header = %q", gotAuth)
	}
	if gotUA != userAgent {
		t.Errorf("user-agent = %q", gotUA)
	}
}

func TestDoRetriesIdempotentOn5xx(t *testing.T) {
	var calls int
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		calls++
		if calls < 2 {
			w.WriteHeader(http.StatusBadGateway)
			return
		}
		_, _ = w.Write([]byte(`[]`))
	}))
	defer srv.Close()

	c := New("t", srv.URL)
	if _, err := c.Do(t.Context(), "GET", "/x", nil); err != nil {
		t.Fatalf("Do: %v", err)
	}
	if calls != 2 {
		t.Errorf("expected a retry (2 calls), got %d", calls)
	}
}

func TestDoDoesNotRetryPost(t *testing.T) {
	var calls int
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		calls++
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer srv.Close()

	c := New("t", srv.URL)
	if _, err := c.Do(t.Context(), "POST", "/x", []byte(`{}`)); err == nil {
		t.Error("expected error on 500")
	}
	if calls != 1 {
		t.Errorf("POST must not be retried, got %d calls", calls)
	}
}

func TestDoDryRunPrintsCurlWithRedactedToken(t *testing.T) {
	var buf bytes.Buffer
	c := New("supersecret-token-9999", "https://connectapi.garmin.com", WithDryRun(&buf))
	body, err := c.Do(t.Context(), "get", "/usersummary-service/usersummary/daily?calendarDate=2026-07-10", nil)
	if err != nil || body != nil {
		t.Fatalf("dry-run should not send: body=%q err=%v", body, err)
	}
	curl := buf.String()
	if !strings.Contains(curl, "curl -sS -X GET") {
		t.Errorf("missing curl line: %q", curl)
	}
	if !strings.Contains(curl, "/usersummary-service/usersummary/daily") {
		t.Errorf("missing path: %q", curl)
	}
	if strings.Contains(curl, "supersecret-token-9999") {
		t.Errorf("token leaked into dry-run: %q", curl)
	}
	if !strings.Contains(curl, "****9999") {
		t.Errorf("token not redacted: %q", curl)
	}
}

func TestDo4xxReturnsError(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusForbidden)
		_, _ = w.Write([]byte(`forbidden`))
	}))
	defer srv.Close()
	c := New("t", srv.URL)
	if _, err := c.Do(t.Context(), "GET", "/x", nil); err == nil {
		t.Error("expected error on 403")
	}
}

func TestShellQuoteEscapes(t *testing.T) {
	if got := shellQuote("it's here"); got != `'it'\''s here'` {
		t.Errorf("shellQuote = %q", got)
	}
}

func TestOptionsAndRedactShort(t *testing.T) {
	var hit bool
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		hit = true
		_, _ = w.Write([]byte(`{}`))
	}))
	defer srv.Close()
	// WithBaseURL overrides the placeholder base; WithHTTPClient injects the server's client.
	c := New("tok", "https://placeholder.example", WithBaseURL(srv.URL), WithHTTPClient(srv.Client()))
	if _, err := c.Do(t.Context(), "GET", "/x", nil); err != nil {
		t.Fatal(err)
	}
	if !hit {
		t.Error("WithBaseURL/WithHTTPClient not applied")
	}
	if got := redact("ab"); got != "****" { // short tokens are fully masked
		t.Errorf("short redact = %q", got)
	}
}
