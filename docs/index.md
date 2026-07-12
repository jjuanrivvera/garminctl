# garminctl

A [Garmin Connect](https://connect.garmin.com) CLI: read your health and activity data — sleep,
body composition, stress, heart rate, activities, training metrics, and the full Garmin Connect
endpoint surface — with OS-keyring token storage, multiple named accounts, a local offline store,
an agent safety guard, and prebuilt binaries.

Built on [`llehouerou/go-garmin`](https://github.com/llehouerou/go-garmin), the Go library (and
CLI) that does the Garmin Connect auth and endpoint work. garminctl generates the same command
surface and adds keyring storage, named profiles, token import, an offline SQLite store,
multi-format output (`table`/`json`/`yaml`/`csv`), an agent guard, and packaging.

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

# Keep your data offline: backfill a range, then read/trend it with no network
garminctl sync --from 2026-01-01
garminctl --offline sleep --date 2026-07-09
garminctl history body-composition --from 2026-01-01 -o csv

# Or a raw request for anything else
garminctl api /usersummary-service/usersummary/daily
```

## Offline data

Garmin Connect is pull-only with no history export, so garminctl keeps a local SQLite store.
Reads cache the day they fetch; `garminctl sync [--from --to] [--metrics …]` backfills a range;
`garminctl --offline <metric>` then serves a day with no network; and `garminctl history <metric>
--from --to` renders one row per day — pair it with `-o csv` for a spreadsheet-ready trend. The
store lives at `<config dir>/store.db` (chmod 0600, per profile).

## For AI agents

`garminctl mcp` exposes the read surface as MCP tools, and `garminctl agent guard --host
claude-code` emits a PreToolUse safety hook. garminctl is read-focused, so the guard blocks the
mutation vectors — `workouts` writes, `api` with a write method, `auth logout`, and `alias set`;
`workouts` and `sync` stay out of the MCP surface.

## Learn more

- **[Command reference](commands/garminctl.md)** — every command, flag, and example, generated
  from the live CLI.
- **[Source & releases](https://github.com/jjuanrivvera/garminctl)** — GitHub repo, changelog,
  and signed release binaries.
- **[go-garmin](https://github.com/llehouerou/go-garmin)** — the upstream library and CLI that
  does the Garmin Connect work.

Built with [cliwright](https://cliwright.jjuanrivvera.com).
