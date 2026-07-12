## garminctl agent guard

Generate agent-safety config that blocks mutating garminctl operations

### Synopsis

garminctl's surface is read-only Garmin Connect health data, so the guard blocks the
few mutation vectors rather than a rich destructive taxonomy:

  • the raw "api" escape hatch with a write method (-X POST|PUT|DELETE|PATCH) — the only way
    to mutate Garmin data through garminctl;
  • "auth logout" — deletes the stored session from the keyring;
  • "alias set" — could mint a shorthand that expands to a blocked command before cobra parses.

Reads — every resource, every "connect" endpoint, and "api" GET — are allowed.

For claude-code the output includes a PreToolUse hook (.claude/hooks/garminctl-guard.sh) that
strips quote/backslash obfuscation, matches the binary even when path-invoked
(./bin/garminctl, /usr/local/bin/garminctl), and gates the "api" hatch by HTTP method.

MCP-only operation is the hard guarantee; the Bash rails are best-effort — the hook defeats
quoting and path prefixes, but not variable indirection (m=DELETE; garminctl api x -X $m) or
shell aliases.

```
garminctl agent guard --host <claude-code|codex|opencode> [flags]
```

### Examples

```
  garminctl agent guard --host claude-code
  garminctl agent guard --host codex --out ~/.codex/config.toml
  garminctl agent guard --host opencode
```

### Options

```
  -h, --help          help for guard
      --host string   target agent host: claude-code|codex|opencode (required)
      --out string    write to this file instead of stdout
```

### Options inherited from parent commands

```
      --dry-run          print the equivalent request instead of sending it
      --no-color         disable colored output
  -o, --output string    output format: table|json|yaml|csv (default "table")
      --profile string   profile (Garmin account) to use; env GARMINCTL_PROFILE
```

### SEE ALSO

* [garminctl agent](garminctl_agent.md)	 - AI-agent integration helpers

