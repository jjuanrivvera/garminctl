# garminctl

Drive [Garmin Connect](https://connect.garmin.com) from the terminal — body composition,
sleep, steps, heart rate, stress, body battery, and every other Connect endpoint — with named
profiles for several accounts, OS-keyring token storage, and table/json/yaml/csv output.

garminctl wraps [`llehouerou/go-garmin`](https://github.com/llehouerou/go-garmin), which does
the reverse-engineered auth (OAuth1 → OAuth2 exchange with **automatic refresh**) and the typed
endpoint surface. garminctl adds the parts you actually run: keyring-backed sessions, named
profiles, a one-command import from an existing `garth` / `python-garminconnect` setup, an MCP
server, and an agent safety guard.

> **The bug it fixes.** A long-lived cron that reads Garmin data eventually dies with
> `GarminConnectAuthenticationError` because the short-lived OAuth2 token expired and nothing
> refreshed it. garminctl refreshes before every request from the ~1-year OAuth1 token and
> persists the new token back to the keyring — so it keeps working unattended.

Built with [cliwright](https://cliwright.jjuanrivvera.com).

## Install

```bash
# Homebrew (macOS/Linux)
brew install jjuanrivvera/tap/garminctl

# Scoop (Windows)
scoop bucket add jjuanrivvera https://github.com/jjuanrivvera/scoop-bucket
scoop install garminctl

# install script (Linux/macOS)
curl -fsSL https://raw.githubusercontent.com/jjuanrivvera/garminctl/main/install.sh | sh

# from source
go install github.com/jjuanrivvera/garminctl/cmd/garminctl@latest
```

Debian/RPM/Alpine packages are attached to each [release](https://github.com/jjuanrivvera/garminctl/releases).

## Authenticate

garminctl has no MFA flow — it reuses tokens you already have or logs in with email + password.

**Import an existing garth / python-garminconnect session (recommended):**

```bash
garminctl init                          # auto-detects ~/.garminconnect
garminctl auth import --from ~/.garminconnect            --profile juan
garminctl auth import --from ~/.garminconnect-vane       --profile vane
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
# profile:        juan
# authenticated:  true
# oauth2 expiry:  2026-03-14T09:22:10Z  (expired — refreshes on next call)
```

## Read your data

Every resource takes an optional `--date` (default today) and honors the global `-o` format:

```bash
garminctl steps
garminctl sleep --date 2026-07-09
garminctl body-composition -o json           # weight, BMI, body-fat %
garminctl stress
garminctl body-battery
garminctl heart-rate
garminctl respiration
garminctl intensity-minutes
```

### Everything else: the `connect` bridge

The typed resources are the common reads; `connect` exposes go-garmin's **full endpoint
registry** (68 endpoints) for anything not surfaced directly:

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
garminctl config use vane
garminctl --profile juan steps       # one-off override; env GARMINCTL_PROFILE also works
```

## Output

`-o table` (default), `json`, `yaml`, `csv`. `table`/`csv` flatten one level so nested objects
stay one row per field.

## For AI agents

```bash
garminctl mcp                                   # expose the read surface as MCP tools
garminctl agent guard --host claude-code        # emit a PreToolUse safety hook
```

garminctl is read-only health data, so the guard blocks only the mutation vectors: `auth
logout` (deletes the session), `alias set` (mints indirections), and `api` with a write method.
See [AGENTS.md](AGENTS.md).

## Diagnostics

```bash
garminctl doctor        # offline: config + keyring + each profile's token state
garminctl version --check
```

## License

MIT — see [LICENSE](LICENSE).
