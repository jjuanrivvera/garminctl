package commands

import (
	"context"
	"fmt"
	"net/http"

	gm "github.com/llehouerou/go-garmin"

	"github.com/jjuanrivvera/garminctl/internal/auth"
	"github.com/jjuanrivvera/garminctl/internal/config"
	"github.com/jjuanrivvera/garminctl/internal/garmin"
)

// testHTTPClient is a test seam: when non-nil it is injected into the go-garmin client so tests
// can mock the transport. It is nil in production (the default HTTP client is used).
var testHTTPClient *http.Client

// store returns the keyring-backed token store (encrypted-file fallback rooted at the config dir).
func store() auth.Store {
	dir, _ := config.Dir()
	return auth.New(dir)
}

// getClient resolves the active profile, loads its session from the keyring, and returns an
// authenticated go-garmin client plus a save func that persists any refreshed tokens back to
// the keyring. Defer save() after using the client so a token refreshed mid-call is not lost.
func getClient(ctx context.Context) (client *gm.Client, save func() error, profile string, err error) {
	profile = config.Resolve(gf.profile)
	sessionJSON, err := store().Get(profile)
	if err != nil || sessionJSON == "" {
		return nil, nil, profile, fmt.Errorf(
			"no session for profile %q — run `garminctl auth import` or `garminctl auth login`", profile)
	}
	c, err := garmin.NewClient(ctx, sessionJSON, testHTTPClient)
	if err != nil {
		return nil, nil, profile, err
	}
	save = func() error {
		dumped, derr := garmin.DumpSession(c)
		if derr != nil {
			return derr
		}
		if dumped != sessionJSON { // only write when the OAuth2 token actually refreshed
			return store().Set(profile, dumped)
		}
		return nil
	}
	return c, save, profile, nil
}
