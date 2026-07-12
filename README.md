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
| Output | table / json / yaml / csv | json |
| Install | **brew / scoop / deb / rpm / apk** + install.sh | `go install` |
| Agent guard | **yes** (PreToolUse hook) | — |
| MCP server | yes | yes |
| Typed data surface | sleep, body-composition, stress, body-battery, heart-rate, respiration, intensity-minutes (+ everything else via `connect`) | broader: also metrics, workouts, exercises, calendar, biometric |

In short: garminctl trades some of go-garmin's command breadth for keyring security, multiple
accounts, token import, and prebuilt packages, in the [jjuanrivvera fleet](https://cliwright.jjuanrivvera.com) style.

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

Each resource takes an optional `--date` (default today) and honors the global `-o` format:

```bash
garminctl sleep --date 2026-07-09
garminctl body-composition -o json           # weight, BMI, body-fat %
garminctl stress
garminctl body-battery
garminctl heart-rate
garminctl respiration
garminctl intensity-minutes
```

### Everything else: the `connect` bridge

`connect` exposes go-garmin's **full endpoint registry** (68 endpoints) for anything the typed
resources above don't wrap — metrics, activities, devices, and more:

```bash
garminctl connect --help                     # list every endpoint group
garminctl connect activities list
```

### Raw escape hatch: `api`

For a Connect endpoint neither the resources nor `connect` wrap:

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

garminctl's typed surface is read-only, so the guard blocks only the mutation vectors: `auth
logout` (deletes the session), `alias set` (mints indirections), and `api` with a write method.
See [AGENTS.md](AGENTS.md).

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
