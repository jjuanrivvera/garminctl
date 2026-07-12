# garminctl

A [Garmin Connect](https://connect.garmin.com) CLI: read your health and activity data — sleep,
body composition, stress, heart rate, activities, training metrics, and the full Garmin Connect
endpoint surface — with OS-keyring token storage, multiple named accounts, a **local offline
store**, an agent safety guard, and prebuilt binaries (Homebrew / Scoop / deb / rpm / apk).

Built on [`llehouerou/go-garmin`](https://github.com/llehouerou/go-garmin), the Go library (and
CLI) that does the Garmin Connect auth and endpoint work. garminctl generates the same command
surface and adds keyring storage, named profiles, token import, an offline SQLite store,
multi-format output (`table`/`json`/`yaml`/`csv`), an agent guard, and packaging.

## Install

```bash
# Homebrew (macOS/Linux)
brew install jjuanrivvera/garminctl/garminctl-cli

# Scoop (Windows)
scoop install https://raw.githubusercontent.com/jjuanrivvera/scoop-garminctl/main/garminctl.json

# install script (Linux/macOS)
curl -fsSL https://raw.githubusercontent.com/jjuanrivvera/garminctl/main/install.sh | sh

# from source
go install github.com/jjuanrivvera/garminctl/cmd/garminctl@latest
```

Debian/RPM/Alpine packages are attached to each [release](https://github.com/jjuanrivvera/garminctl/releases).

## Authenticate

No MFA flow — reuse tokens you already have, or log in with email + password.

**Import an existing session (recommended).** If you already have working Garmin tokens under
`~/.garminconnect` (`oauth1_token.json` + `oauth2_token.json`), import them into the keyring:

```bash
garminctl init                          # auto-detects ~/.garminconnect
garminctl auth import --from ~/.garminconnect            --profile me
garminctl auth import --from ~/.garminconnect-alt        --profile alt
```

Nothing is written back to the token directory.

**Or log in fresh:**

```bash
garminctl auth login --email you@example.com     # prompts for the password (hidden)
```

**Check status (offline — no API call):**

```bash
garminctl auth status
# profile:        me
# authenticated:  true
# oauth2 expiry:  2026-06-14T09:22:10Z  (expired — refreshes on next call)
```

## Read your data

Curated shortcuts cover the common daily reads — each takes an optional `--date` (default today)
and honors the global `-o` format:

```bash
garminctl sleep --date 2026-07-09
garminctl body-composition -o json           # weight, BMI, body-fat %
garminctl stress
garminctl body-battery
garminctl heart-rate
garminctl respiration
garminctl intensity-minutes
garminctl steps                              # daily step count, goal, and distance
```

### The full surface

Every Garmin Connect endpoint (68 in all) is available as a top-level command, grouped by service,
each honoring `-o table/yaml/csv`:

```bash
garminctl metrics vo2max
garminctl activities list
garminctl devices list
garminctl hrv daily
garminctl weight range --start=2026-07-01 --end=2026-07-09
garminctl workouts list
garminctl --help                             # every group at the top level
```

Workout **writes** (`create`/`update`/`delete`/`schedule`/`unschedule`) are the only mutations; the
[agent guard](#for-ai-agents) blocks them by default.

### Raw escape hatch: `api`

For a Connect endpoint the typed commands don't wrap:

```bash
garminctl api /usersummary-service/usersummary/daily
garminctl --dry-run api /userprofile-service/userprofile   # prints the equivalent curl
```

`api` signs the request with the active profile's token (redacted under `--dry-run`).

## Offline data

Garmin Connect is pull-only with no history export, so garminctl keeps a local SQLite store. Every
read caches the day it fetched, and `sync` backfills a range so your data is available with no
network afterward:

```bash
garminctl sync                                   # last 7 days, all metrics
garminctl sync --from 2026-01-01 --metrics sleep,body-composition

garminctl --offline sleep --date 2026-07-09      # served from the store, no API call
garminctl history body-composition --from 2026-01-01 -o csv   # a trend: one row per day
```

`history` renders one row per day — with `-o csv` you get a spreadsheet-ready trend. The store
lives at `<config dir>/store.db` (chmod 0600, per profile).

## Profiles

```bash
garminctl config list                # * marks the default
garminctl config use alt
garminctl --profile me sleep         # one-off override; env GARMINCTL_PROFILE also works
```

## Output

`-o table` (default), `json`, `yaml`, `csv`. A single record flattens one level (one row per
field); a **list** (e.g. `activities list`, `workouts list`, `history`) renders as a real table —
a header plus one row per record — so `-o csv` is spreadsheet-ready.

## For AI agents

```bash
garminctl mcp                                   # expose the read surface as MCP tools
garminctl agent guard --host claude-code        # emit a PreToolUse safety hook
```

garminctl is read-focused, so the guard blocks the mutation vectors: `workouts`
create/update/delete/schedule/unschedule, `api` with a write method, `auth logout`, and
`alias set`. `workouts` is also excluded from the MCP tool surface. See [AGENTS.md](AGENTS.md).

## Diagnostics

```bash
garminctl doctor        # offline: config + keyring + each profile's token state
garminctl version --check
```

## Credits

The Garmin Connect client, auth, and endpoint registry are
[`llehouerou/go-garmin`](https://github.com/llehouerou/go-garmin). garminctl is the packaging and
UX layer around it. Built with [cliwright](https://cliwright.jjuanrivvera.com).

## License

MIT — see [LICENSE](LICENSE).
