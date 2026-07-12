# garminctl

A [Garmin Connect](https://connect.garmin.com) CLI built on
[`llehouerou/go-garmin`](https://github.com/llehouerou/go-garmin), packaged the
[cliwright](https://cliwright.jjuanrivvera.com) way: OS-keyring token storage, multiple named
accounts, one-command import from an existing `garth` / `python-garminconnect` setup, an agent
safety guard, and prebuilt binaries.

## Is garminctl the right tool for you?

go-garmin is **both a Go library and its own CLI** (`garmin`), and that CLI is broader than this
one — it also covers training metrics, workouts (including writes), the exercise library,
calendar, and biometric data as first-class commands, and ships its own MCP server. Garmin's auth
(OAuth1 → OAuth2 exchange with automatic refresh) lives in go-garmin; both tools inherit it.

!!! tip "Prefer the official CLI when it fits"
    If you want the most complete surface and you're happy with `go install`, a single account,
    and a plaintext session file, use go-garmin's `garmin` directly:
    `go install github.com/llehouerou/go-garmin/cmd/garmin@latest`

garminctl is worth it when you want **keyring-encrypted tokens** (go-garmin writes a plaintext
`~/.config/garmin/session.json`), **multiple named accounts**, **import of existing
`~/.garminconnect` tokens**, an **agent guard**, table/json/yaml/csv output, and **prebuilt
packages** (brew / scoop / deb / rpm / apk). It trades some of go-garmin's command breadth for
those.

!!! note "Coming from a Python setup that died with `GarminConnectAuthenticationError`?"
    That was a problem in the `garth` / `python-garminconnect` stack, not in go-garmin. Either
    go-garmin's `garmin` or garminctl refreshes the OAuth2 token correctly from the ~1-year OAuth1
    token, so either one resolves it. garminctl additionally persists the refreshed token to the
    keyring.

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

# Curated shortcuts for common daily reads (--date defaults to today)
garminctl sleep --date 2026-07-09
garminctl body-composition -o json
garminctl stress

# The full go-garmin registry, promoted to the top level (honoring -o table/yaml/csv)
garminctl metrics vo2max
garminctl activities list
garminctl workouts list

# Or a raw request for anything else
garminctl api /usersummary-service/usersummary/daily
```

## Learn more

- **[Command reference](commands/garminctl.md)** — every command, flag, and example, generated
  from the live CLI.
- **[Source & releases](https://github.com/jjuanrivvera/garminctl)** — GitHub repo, changelog,
  and signed release binaries.
- **[go-garmin](https://github.com/llehouerou/go-garmin)** — the upstream library and CLI that
  does the Garmin Connect work.

Built with [cliwright](https://cliwright.jjuanrivvera.com).
