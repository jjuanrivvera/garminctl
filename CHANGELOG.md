# Changelog

All notable changes to garminctl are documented here. The format follows
[Keep a Changelog](https://keepachangelog.com/en/1.1.0/), and this project adheres to
[Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [0.1.0] - 2026-07-11

Initial release.

### Added

- **Auth that survives.** Wraps `llehouerou/go-garmin` for OAuth1 → OAuth2 exchange with
  automatic refresh before every request, persisting the refreshed token back to the keyring —
  fixing the recurring `GarminConnectAuthenticationError` on long-running setups.
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
