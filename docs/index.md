# garminctl

Drive [Garmin Connect](https://connect.garmin.com) from the terminal — body composition, sleep,
heart rate, stress, body battery, and every other Connect endpoint — with named profiles for
several accounts, OS-keyring token storage, and table/json/yaml/csv output.

garminctl wraps [`llehouerou/go-garmin`](https://github.com/llehouerou/go-garmin), which does the
reverse-engineered auth (OAuth1 → OAuth2 exchange with **automatic refresh**) and the typed
endpoint surface. garminctl adds keyring-backed sessions, named profiles, a one-command import
from an existing `garth` / `python-garminconnect` setup, an MCP server, and an agent safety guard.

!!! note "The bug it fixes"
    A long-lived cron that reads Garmin data eventually dies with
    `GarminConnectAuthenticationError` because the short-lived OAuth2 token expired and nothing
    refreshed it. garminctl refreshes before every request from the ~1-year OAuth1 token and
    persists the new token back to the keyring — so it keeps working unattended.

## Install

=== "Homebrew"

    ```bash
    brew install jjuanrivvera/garminctl/garminctl-cli
    ```

=== "Scoop"

    ```bash
    scoop install https://raw.githubusercontent.com/jjuanrivvera/scoop-garminctl/main/garminctl.json
    ```

=== "Script"

    ```bash
    curl -fsSL https://raw.githubusercontent.com/jjuanrivvera/garminctl/main/install.sh | sh
    ```

=== "Go"

    ```bash
    go install github.com/jjuanrivvera/garminctl/cmd/garminctl@latest
    ```

## Quick start

```bash
# Import an existing garth / python-garminconnect session (no MFA flow needed)
garminctl init                                     # auto-detects ~/.garminconnect
garminctl auth import --from ~/.garminconnect --profile me

# Read your data (each resource takes --date, default today)
garminctl sleep --date 2026-07-09
garminctl body-composition -o json
garminctl stress

# Anything the typed resources don't wrap: the full endpoint registry, or a raw request
garminctl connect --help
garminctl api /usersummary-service/usersummary/daily
```

## Learn more

- **[Command reference](commands/garminctl.md)** — every command, flag, and example, generated
  from the live CLI.
- **[Source & releases](https://github.com/jjuanrivvera/garminctl)** — GitHub repo, changelog,
  and signed release binaries.

Built with [cliwright](https://cliwright.jjuanrivvera.com).
