# garminctl

A [Garmin Connect](https://connect.garmin.com) CLI: read your health and activity data — sleep,
body composition, stress, heart rate, activities, training metrics, and the full Garmin Connect
endpoint surface — with OS-keyring token storage, multiple named accounts, an agent safety guard,
and prebuilt binaries.

Built on [`llehouerou/go-garmin`](https://github.com/llehouerou/go-garmin), the Go library (and
CLI) that does the Garmin Connect auth and endpoint work. garminctl generates the same command
surface and adds keyring storage, named profiles, token import, multi-format output
(`table`/`json`/`yaml`/`csv`), an agent guard, and packaging.

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
# Import an existing Garmin session (~/.garminconnect); no MFA flow needed
garminctl init                                     # auto-detects ~/.garminconnect
garminctl auth import --from ~/.garminconnect --profile me

# Curated shortcuts for common daily reads (--date defaults to today)
garminctl sleep --date 2026-07-09
garminctl body-composition -o json
garminctl stress

# The full endpoint surface at the top level (honoring -o table/yaml/csv)
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
