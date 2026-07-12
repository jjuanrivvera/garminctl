# Changelog

All notable changes to garminctl are documented here. The format follows
[Keep a Changelog](https://keepachangelog.com/en/1.1.0/), and this project adheres to
[Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

### Changed

- **Full surface at the top level.** go-garmin's complete endpoint registry (metrics, activities,
  workouts, devices, exercises, calendar, biometric, hrv, weight, wellness, …) is now promoted to
  top-level commands, matching go-garmin's own `garmin` CLI — instead of being nested under
  `connect` (removed). The 7 curated shortcuts (sleep, body-composition, stress, …) remain. The
  promoted commands re-render go-garmin's JSON through garminctl's formatter, so `-o table/yaml/csv`
  works on them too (go-garmin emits JSON only).

### Added

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

[0.1.0]: https://github.com/jjuanrivvera/garminctl/releases/tag/v0.1.0
