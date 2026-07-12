# Security Policy

## Reporting a vulnerability

Email **jjuanrivvera@gmail.com** with details and reproduction steps. Please do not open a
public issue for a security report. Expect an acknowledgement within a few days.

## How garminctl handles credentials

- **Tokens live in the OS keyring** (macOS Keychain, GNOME Keyring / libsecret, Windows
  Credential Manager) via `zalando/go-keyring`. When no keyring is available (headless Linux,
  CI), garminctl falls back to an **AES-256-GCM encrypted file** whose key comes from
  `GARMINCTL_KEYRING_PASSWORD`. Tokens are never written to the config file or the repo.
- **Import is read-only.** `garminctl auth import` reads `oauth1_token.json` and
  `oauth2_token.json` from the directory you point it at and never writes back to them.
- **Password entry is hidden.** `garminctl auth login` reads the password with
  `term.ReadPassword` (no echo, no shell history); garminctl never takes a password on the
  command line or an environment variable.
- **`--dry-run` redacts.** The equivalent curl printed for `garminctl api` masks all but the
  last four characters of the bearer token.
- **Error output is sanitized.** API error bodies are stripped of terminal escape sequences
  before printing, so a crafted response cannot manipulate your terminal.

## Agent guard

`garminctl agent guard` generates host safety config that blocks the mutation vectors (`auth
logout`, `alias set`, and `api` writes) on the Bash and MCP surfaces. The Bash rails are
best-effort (they defeat quoting and path prefixes, not variable indirection); **MCP-only
operation, or a read-only sandbox, is the hard guarantee.** See [AGENTS.md](AGENTS.md).

## Supply chain

Releases are built by GoReleaser in CI, checksummed, signed with cosign (keyless / OIDC), and
ship a CycloneDX SBOM. Verify a downloaded binary against the published checksum and signature.

## Supported versions

Fixes land on the latest minor release. Pre-1.0, only the most recent tag is supported.
