# garminctl

A [Garmin Connect](https://connect.garmin.com) CLI built on
[`llehouerou/go-garmin`](https://github.com/llehouerou/go-garmin), packaged the
[cliwright](https://cliwright.jjuanrivvera.com) way: **OS-keyring token storage, multiple named
accounts, one-command import from an existing `garth` / `python-garminconnect` setup, an agent
safety guard, and prebuilt binaries** (Homebrew / Scoop / deb / rpm / apk).

## Is garminctl the right tool for you?

go-garmin is **both a Go library and its own CLI** (`garmin`), and that CLI is broader than this
one — it also covers training metrics, workouts (including writes), the exercise library,
calendar, and biometric data as first-class commands, and it ships its own MCP server. All of
Garmin's auth (OAuth1 → OAuth2 exchange with automatic refresh) lives in go-garmin, and both
tools inherit it.

**Use go-garmin's `garmin` directly** if you want the most complete surface and you're happy with
`go install`, a single account, and a plaintext session file:

```bash
go install github.com/llehouerou/go-garmin/cmd/garmin@latest
```

**Use garminctl** if you want:

| | garminctl | go-garmin `garmin` |
|---|---|---|
| Token storage | **OS keyring** (+ encrypted-file fallback) | plaintext `~/.config/garmin/session.json` |
| Accounts | **multiple named profiles** (you + partner) | single session |
| Import existing tokens | **`garth` / `~/.garminconnect` import** | fresh `login` only |
| Output | **table / json / yaml / csv** (everywhere) | json |
| Install | **brew / scoop / deb / rpm / apk** + install.sh | `go install` |
| Agent guard | **yes** (PreToolUse hook) | — |
| MCP server | yes | yes |
| Command surface | the **same full registry** (sleep, wellness, activities, metrics, workouts, exercises, calendar, biometric, hrv, devices, …) **plus 7 curated shortcuts** | full registry |

Both tools generate their command surface from the same go-garmin endpoint registry, so the data
coverage is identical. garminctl adds keyring security, multiple accounts, token import,
multi-format output, an agent guard, and prebuilt packages — in the
[jjuanrivvera fleet](https://cliwright.jjuanrivvera.com) style.

> **Coming from a Python setup that died with `GarminConnectAuthenticationError`?** That was a
> problem in the `garth` / `python-garminconnect` stack, not in go-garmin. Either go-garmin's
> `garmin` or garminctl refreshes the OAuth2 token correctly (from the ~1-year OAuth1 token), so
> either one resolves it. garminctl additionally persists the refreshed token to the keyring.

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

**Import an existing garth / python-garminconnect session (recommended):**

```bash
garminctl init                          # auto-detects ~/.garminconnect
garminctl auth import --from ~/.garminconnect            --profile me
garminctl auth import --from ~/.garminconnect-alt        --profile alt
```

`import` reads `oauth1_token.json` + `oauth2_token.json` and stores the translated session in
your OS keyring. Nothing is written back to the token directory.

**Or log in fresh:**

```bash
garminctl auth login --email you@example.com     # prompts for the password (hidden)
```

**Check status (offline — no API call):**

```bash
garminctl auth status
# profile:        me
# authenticated:  true
# oauth2 expiry:  2026-03-14T09:22:10Z  (expired — refreshes on next call)
```

## Read your data

The **curated shortcuts** cover the common daily reads — each takes an optional `--date`
(default today) and honors the global `-o` format:

```bash
garminctl sleep --date 2026-07-09
garminctl body-composition -o json           # weight, BMI, body-fat %
garminctl stress
garminctl body-battery
garminctl heart-rate
garminctl respiration
garminctl intensity-minutes
```

### The full surface

go-garmin's complete endpoint registry (68 endpoints) is promoted to the **top level**, grouped
by service — the same commands as go-garmin's `garmin` CLI, but honoring `-o table/yaml/csv`:

```bash
garminctl metrics vo2max
garminctl activities list
garminctl devices list
garminctl hrv daily
garminctl weight range --start=2026-07-01 --end=2026-07-09
garminctl workouts list
garminctl --help                             # every group at the top level
```

Workout **writes** (`create`/`update`/`delete`/`schedule`/`unschedule`) are the only typed
mutations; the [agent guard](#for-ai-agents) blocks them by default.

### Raw escape hatch: `api`

For a Connect endpoint the typed commands don't wrap:

```bash
garminctl api /usersummary-service/usersummary/daily
garminctl --dry-run api /userprofile-service/userprofile   # prints the equivalent curl
```

`api` signs the request with the active profile's token (redacted under `--dry-run`). Writes
(`-X POST|PUT|DELETE`) are possible but unusual — the agent guard blocks them by default.

## Profiles

```bash
garminctl config list                # * marks the default
garminctl config use alt
garminctl --profile me sleep         # one-off override; env GARMINCTL_PROFILE also works
```

## Output

`-o table` (default), `json`, `yaml`, `csv`. `table`/`csv` flatten one level so nested objects
stay one row per field.

## For AI agents

```bash
garminctl mcp                                   # expose the read surface as MCP tools
garminctl agent guard --host claude-code        # emit a PreToolUse safety hook
```

garminctl is read-focused, so the guard blocks the few mutation vectors: `workouts`
create/update/delete/schedule/unschedule, `api` with a write method, `auth logout` (deletes the
session), and `alias set` (mints indirections). `workouts` is also kept out of the MCP tool
surface. See [AGENTS.md](AGENTS.md).

## Diagnostics

```bash
garminctl doctor        # offline: config + keyring + each profile's token state
garminctl version --check
```

## Credits

The Garmin Connect client, auth, and endpoint registry are
[`llehouerou/go-garmin`](https://github.com/llehouerou/go-garmin) — itself both a library and a
more complete CLI. garminctl is the packaging and UX layer around it. Built with
[cliwright](https://cliwright.jjuanrivvera.com).

## License

MIT — see [LICENSE](LICENSE).
