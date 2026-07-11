// Package garmin wraps llehouerou/go-garmin with keyring-backed sessions and garth token
// import. go-garmin does the reverse-engineered auth (OAuth1 + OAuth2 exchange, refresh) and
// the endpoint surface; garminctl adds keyring storage, named profiles, and the import path
// from an existing python-garminconnect / garth session.
package garmin

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	gm "github.com/llehouerou/go-garmin"
)

// session mirrors go-garmin's on-disk authState JSON — the exact shape LoadSession/SaveSession
// serialize, so we can build one from imported tokens and hand it to the client.
type session struct {
	OAuth1Token        string    `json:"oauth1_token"`
	OAuth1Secret       string    `json:"oauth1_secret"`
	MFAToken           string    `json:"mfa_token,omitempty"`
	OAuth2AccessToken  string    `json:"oauth2_access_token"`
	OAuth2RefreshToken string    `json:"oauth2_refresh_token"`
	OAuth2Expiry       time.Time `json:"oauth2_expiry"`
	OAuth2Scope        string    `json:"oauth2_scope,omitempty"`
	Domain             string    `json:"domain"`
}

// garth token file shapes (python-garminconnect / garth on-disk format).
type garthOAuth1 struct {
	OAuthToken       string `json:"oauth_token"`
	OAuthTokenSecret string `json:"oauth_token_secret"`
	MFAToken         string `json:"mfa_token"`
	Domain           string `json:"domain"`
}

type garthOAuth2 struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	Scope        string `json:"scope"`
	ExpiresAt    int64  `json:"expires_at"` // unix seconds
}

// ImportGarth reads a garth token directory (oauth1_token.json + oauth2_token.json) and returns
// a go-garmin session JSON string ready to store in the keyring — how an existing
// python-garminconnect / garth session (e.g. ~/.garminconnect) migrates into garminctl.
func ImportGarth(dir string) (string, error) {
	var o1 garthOAuth1
	if err := readJSON(filepath.Join(dir, "oauth1_token.json"), &o1); err != nil {
		return "", fmt.Errorf("read oauth1_token.json: %w", err)
	}
	var o2 garthOAuth2
	if err := readJSON(filepath.Join(dir, "oauth2_token.json"), &o2); err != nil {
		return "", fmt.Errorf("read oauth2_token.json: %w", err)
	}
	if o1.OAuthToken == "" || o2.AccessToken == "" {
		return "", fmt.Errorf("incomplete tokens in %s (need oauth_token + access_token)", dir)
	}
	s := session{
		OAuth1Token:        o1.OAuthToken,
		OAuth1Secret:       o1.OAuthTokenSecret,
		MFAToken:           o1.MFAToken,
		Domain:             o1.Domain,
		OAuth2AccessToken:  o2.AccessToken,
		OAuth2RefreshToken: o2.RefreshToken,
		OAuth2Scope:        o2.Scope,
		OAuth2Expiry:       time.Unix(o2.ExpiresAt, 0),
	}
	b, err := json.Marshal(s)
	if err != nil {
		return "", err
	}
	return string(b), nil
}

func readJSON(path string, v any) error {
	b, err := os.ReadFile(path) // #nosec G304 -- path is a user-supplied token dir, by design
	if err != nil {
		return err
	}
	return json.Unmarshal(b, v)
}

// NewClient builds an authenticated go-garmin client from a stored session JSON. The client
// auto-refreshes the OAuth2 token (via the long-lived OAuth1 token) before each call when it's
// near expiry — the reliability that the Python cron setup lacked.
func NewClient(_ context.Context, sessionJSON string) (*gm.Client, error) {
	c := gm.New(gm.Options{})
	if err := c.LoadSession(strings.NewReader(sessionJSON)); err != nil {
		return nil, fmt.Errorf("load session: %w", err)
	}
	return c, nil
}

// SessionInfo reports the OAuth2 expiry and whether both tokens are present, parsed from a
// stored session JSON — a status check with no live API call.
func SessionInfo(sessionJSON string) (expiry time.Time, authenticated bool, err error) {
	var s session
	if err = json.Unmarshal([]byte(sessionJSON), &s); err != nil {
		return time.Time{}, false, err
	}
	return s.OAuth2Expiry, s.OAuth1Token != "" && s.OAuth2AccessToken != "", nil
}

// DumpSession serializes the client's current session (with any refreshed tokens) for
// persistence back to the keyring after a run.
func DumpSession(c *gm.Client) (string, error) {
	var buf bytes.Buffer
	if err := c.SaveSession(&buf); err != nil {
		return "", err
	}
	return buf.String(), nil
}
