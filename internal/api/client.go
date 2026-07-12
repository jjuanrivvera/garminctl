// Package api is a thin authenticated HTTP client for raw Garmin Connect requests — the engine
// behind `garminctl api`, the escape hatch for endpoints the typed surface doesn't wrap.
//
// The typed commands go through go-garmin, which owns the reverse-engineered auth (OAuth1 →
// OAuth2 exchange, lazy refresh) and rate limiting. This client is deliberately minimal: it
// signs a one-off request with the session's current bearer token, renders an equivalent curl
// under --dry-run, and retries only idempotent methods. It never touches the keyring or the
// refresh flow — the caller supplies a token the typed path has already kept fresh.
package api

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"time"
)

const (
	userAgent      = "GCM-iOS-5.19.1.2" // must match go-garmin so Garmin accepts the request
	maxAttempts    = 3
	requestTimeout = 30 * time.Second
)

// Client issues authenticated raw requests against the Garmin Connect API.
type Client struct {
	http    *http.Client
	baseURL string
	token   string
	dryRun  bool
	dryRunW io.Writer // where the --dry-run curl is written (default os.Stderr)
}

// Option configures a Client.
type Option func(*Client)

// WithBaseURL overrides the API base (default derived from the session domain). Used by tests
// to point at an httptest server.
func WithBaseURL(u string) Option { return func(c *Client) { c.baseURL = strings.TrimRight(u, "/") } }

// WithHTTPClient injects a transport (test seam).
func WithHTTPClient(h *http.Client) Option { return func(c *Client) { c.http = h } }

// WithDryRun makes Do print the equivalent curl to w instead of sending the request.
func WithDryRun(w io.Writer) Option {
	return func(c *Client) { c.dryRun, c.dryRunW = true, w }
}

// New builds a Client that signs requests with the given OAuth2 bearer token.
func New(token, baseURL string, opts ...Option) *Client {
	c := &Client{
		http:    &http.Client{Timeout: requestTimeout},
		baseURL: strings.TrimRight(baseURL, "/"),
		token:   token,
		dryRunW: os.Stderr,
	}
	for _, o := range opts {
		o(c)
	}
	return c
}

// idempotent reports whether method is safe to retry: repeating it has no additional effect,
// so a transient 5xx/429 can be retried without risking a duplicate write.
func idempotent(method string) bool {
	switch method {
	case http.MethodGet, http.MethodHead, http.MethodPut, http.MethodDelete, http.MethodOptions:
		return true
	}
	return false
}

// Do sends method+path (path is everything after the base URL, query string included) and
// returns the response body. Under --dry-run it writes the equivalent curl (token redacted)
// and returns nil, nil. Idempotent methods are retried on 429/5xx with capped backoff; a
// non-idempotent method (POST) is never retried, so a write can't be silently duplicated.
func (c *Client) Do(ctx context.Context, method, path string, body []byte) ([]byte, error) {
	method = strings.ToUpper(method)
	url := c.baseURL + path
	if c.dryRun {
		fmt.Fprintln(c.dryRunW, c.curl(method, url, body))
		return nil, nil
	}

	var lastErr error
	for attempt := range maxAttempts {
		if attempt > 0 { // only reached for idempotent methods
			select {
			case <-ctx.Done():
				return nil, ctx.Err()
			case <-time.After(backoff(attempt)):
			}
		}

		var reqBody io.Reader
		if body != nil {
			reqBody = bytes.NewReader(body)
		}
		req, err := http.NewRequestWithContext(ctx, method, url, reqBody)
		if err != nil {
			return nil, err
		}
		req.Header.Set("Authorization", "Bearer "+c.token)
		req.Header.Set("User-Agent", userAgent)
		req.Header.Set("nk", "NT")
		if body != nil {
			req.Header.Set("Content-Type", "application/json")
		}

		resp, err := c.http.Do(req)
		if err != nil {
			lastErr = err
			if idempotent(method) {
				continue
			}
			return nil, err
		}
		data, _ := io.ReadAll(resp.Body)
		_ = resp.Body.Close()

		if resp.StatusCode == http.StatusTooManyRequests || resp.StatusCode >= 500 {
			lastErr = &APIError{Status: resp.StatusCode, Path: path, Body: string(data)}
			if idempotent(method) {
				continue
			}
			return data, lastErr
		}
		if resp.StatusCode >= 400 {
			return data, &APIError{Status: resp.StatusCode, Path: path, Body: string(data)}
		}
		return data, nil
	}
	return nil, lastErr
}

// backoff returns the delay before retry attempt n (1-based): 200ms, 400ms, …, capped.
func backoff(attempt int) time.Duration {
	d := time.Duration(200*(1<<(attempt-1))) * time.Millisecond
	return min(d, 2*time.Second)
}

// curl renders the equivalent curl invocation with the bearer token redacted — the --dry-run
// output. Values are shell-quoted so the printed command is safe to copy-paste.
func (c *Client) curl(method, url string, body []byte) string {
	var b strings.Builder
	b.WriteString("curl -sS -X " + method + " " + shellQuote(url))
	b.WriteString(" -H " + shellQuote("Authorization: Bearer "+redact(c.token)))
	b.WriteString(" -H " + shellQuote("User-Agent: "+userAgent))
	b.WriteString(" -H " + shellQuote("nk: NT"))
	if body != nil {
		b.WriteString(" -H " + shellQuote("Content-Type: application/json"))
		b.WriteString(" --data " + shellQuote(string(body)))
	}
	return b.String()
}

// redact masks all but the last 4 characters of a token, so a --dry-run curl reveals the shape
// without leaking the credential.
func redact(tok string) string {
	if len(tok) <= 4 {
		return "****"
	}
	return "****" + tok[len(tok)-4:]
}

// shellQuote wraps s in single quotes, escaping embedded single quotes, so the value survives
// a POSIX shell verbatim.
func shellQuote(s string) string {
	return "'" + strings.ReplaceAll(s, "'", `'\''`) + "'"
}
