# AGENTS.md — working in the garminctl repo

`garminctl` is a command-line tool for **Garmin Connect**, built to the cliwright standard
(Go + Cobra + GoReleaser). It wraps [`llehouerou/go-garmin`](https://github.com/llehouerou/go-garmin)
for the reverse-engineered auth and endpoint surface, and adds keyring sessions, named
profiles, garth token import, an MCP server, and an agent guard. This file orients an AI agent
(or human) contributing.

## The one rule that matters

**`make verify` is the gate.** A change is done only when `make verify` exits `0`. It runs
`make check` (fmt, vet, golangci-lint, tests) + `spec-check` (the built surface matches
`api-manifest.json`) + `spec-completeness` (the manifest wraps the enumerated go-garmin
registry) + `cover-check` (≥80% coverage) + `dod-check.sh`. Run the full `make verify` for any
change to the command surface or a documented behavior — not just `make check`.

## Architecture (where things live)

- `internal/garmin/` — the wrapper around go-garmin: garth token import (`ImportGarth`),
  keyring session load/dump, the `NewClient` that auto-refreshes OAuth2, and `SessionToken`
  for the raw `api` hatch. **This is the crux** — go-garmin owns the auth; garminctl owns
  persistence and profiles.
- `internal/api/` — a thin authenticated raw client (bearer token, idempotent-only retry,
  `--dry-run` curl) behind `garminctl api`. Deliberately minimal; it does not touch refresh.
- `commands/` — thin, declarative command files, each registered from its own `init()` via
  `registerCommand` (zero edits to shared code to add one). `resources.go` is the typed
  read surface; `connect.go` bridges go-garmin's full 68-endpoint registry.
- `internal/{config,auth,output,version}` — profiles + manual precedence (no Viper), keyring
  token storage (+ encrypted-file fallback), the table/json/yaml/csv renderer, build metadata.
- `cmd/garminctl/main.go` — entry point: `signal.NotifyContext` (Ctrl-C cancels in-flight
  work: token refresh, retry backoff) then `commands.Main`.

## Agent safety

`garminctl agent guard --host claude-code` emits a PreToolUse hook. garminctl's surface is
read-only, so the guard blocks only the mutation vectors:

- `auth logout` — deletes the stored session from the keyring;
- `alias set` — could mint a shorthand that expands to a blocked command;
- `api` with a write HTTP method (`-X POST|PUT|DELETE|PATCH`) — the only path that mutates
  Garmin data. `api` GET (a read) passes.

The hook de-obfuscates quote/backslash tricks, catches path-invoked binaries
(`./bin/garminctl`), tolerates the known global flags between the binary and the subcommand,
and fails safe when `jq` is absent. `commands/agent_hook_test.go` runs it through real bash.
MCP-only operation is the hard guarantee; the Bash rails are best-effort.

## House rules

- Comments explain **WHY**, not WHAT.
- Thread `cmd.Context()` everywhere; never `context.Background()` (it breaks Ctrl-C). Tests
  use `t.Context()`.
- Secrets live in the OS keyring — never in config-in-repo, code, or commit messages. The
  garth import path is read-only on the token files.
- Never cross account boundaries: profiles are separate by design.
- The refresh guarantee is the whole point — the typed surface refreshes through go-garmin and
  `getClient`'s `save()` persists the new token. Don't add a code path that reads data without
  going through that save-back.
