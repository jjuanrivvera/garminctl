# Changelog

All notable changes to garminctl are documented here. The format follows
[Keep a Changelog](https://keepachangelog.com/en/1.1.0/), and this project adheres to
[Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [0.3.0] - 2026-07-12

### Added

- **`steps` — daily step count, goal, and distance.** The last advertised-but-missing curated
  metric. go-garmin declares a `StepsService` but implements no methods and its registry omits the
  endpoint, so `steps` fetches Garmin's `/usersummary-service/stats/steps/daily/{date}/{date}`
  through the raw client. It slots into the curated-metric machinery like the others: `--date`,
  `-o table/json/yaml/csv`, `--offline`, `history`, and `sync`. Because a raw request doesn't
  refresh OAuth2 the way a typed go-garmin call does, `steps` first nudges a refresh when the
  session is near expiry, so it never 401s on a stale token.

### Security

- Bump the Go toolchain to **go1.25.12**, clearing the reachable standard-library advisories
  (crypto/tls GO-2026-5856, crypto/x509, net/http, net/textproto) that govulncheck flagged.

## [0.2.1] - 2026-07-12

### Fixed

- The hidden secret prompt (`auth login` / `auth import`) now reads in **raw mode** instead of
  `term.ReadPassword`'s canonical mode (capped at MAX_CANON, 1024 bytes on macOS), so a long
  pasted token no longer hangs the prompt until Ctrl-C. Bracketed-paste markers are still stripped
  as a defensive guard.

## [0.2.0] - 2026-07-11

### Changed

- **Full surface at the top level.** go-garmin's complete endpoint registry (metrics, activities,
  workouts, devices, exercises, calendar, biometric, hrv, weight, wellness, …) is now promoted to
  top-level commands, matching go-garmin's own `garmin` CLI — instead of being nested under
  `connect` (removed). The 7 curated shortcuts (sleep, body-composition, stress, …) remain. The
  promoted commands re-render go-garmin's JSON through garminctl's formatter, so `-o table/yaml/csv`
  works on them too (go-garmin emits JSON only).

### Added

- **Offline store.** A local SQLite database (pure-Go `modernc.org/sqlite`) keeps the health data
  you fetch. Every read caches its day; `garminctl sync [--from --to] [--metrics …]` backfills a
  range; `garminctl --offline <metric>` serves a day with no network; and `garminctl history
  <metric> --from --to` renders one row per day (`-o csv` for a spreadsheet-ready trend). The
  renderer now turns any array of objects into a real table/CSV, not one wrapped cell.
- **Workout-write guardrails.** Promoting the registry brings in `workouts`
  create/update/delete/schedule/unschedule (the only typed writes). The agent guard now blocks
  them on the Bash surface, and `workouts` is excluded from the MCP tool surface.

## [0.1.0] - 2026-07-11

Initial release.

### Added

- **Keyring-backed auth.** Wraps `llehouerou/go-garmin` (which does the OAuth1 → OAuth2 exchange
  and automatic refresh) and persists the session in the OS keyring instead of a plaintext file,
  so refreshed tokens survive across runs.
- **garth import.** `garminctl auth import --from <dir>` translates an existing
  `garth` / `python-garminconnect` session (`~/.garminconnect`) into a keyring profile; the
  token files are read-only. `garminctl init` auto-detects the default directory.
- **Login / status / logout.** Hidden password entry; offline `auth status`.
- **Named profiles** for multiple accounts, with `config list` / `config use`, the global
  `--profile` flag, and `GARMINCTL_PROFILE`.
- **Typed read surface:** `body-composition`, `sleep`, `heart-rate`, `stress`, `body-battery`,
  `respiration`, `intensity-minutes` — each with `--date` and json/yaml/csv/table.
- **`connect` bridge** exposing go-garmin's full 68-endpoint registry.
- **`api` escape hatch** — raw authenticated request with `--dry-run` curl (token redacted)
  and idempotent-only retry.
- **Agent surface:** `mcp` server (ophis) and `agent guard` with a PreToolUse hook that blocks
  the mutation vectors (`auth logout`, `alias set`, `api` writes) and fails safe without `jq`.
- **Meta commands:** `doctor`, `config`, `init`, `version --check`, `completion`, `alias`.
- **Distribution:** Homebrew, Scoop, deb/rpm/apk, `install.sh`; cosign-signed releases + SBOM.

[0.2.0]: https://github.com/jjuanrivvera/garminctl/releases/tag/v0.2.0
[0.1.0]: https://github.com/jjuanrivvera/garminctl/releases/tag/v0.1.0
