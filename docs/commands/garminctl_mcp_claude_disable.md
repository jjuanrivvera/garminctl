## garminctl mcp claude disable

Remove server from Claude config

### Synopsis

Remove this application from Claude Desktop MCP servers

```
garminctl mcp claude disable [flags]
```

### Options

```
      --config-path string   Path to Claude config file
  -h, --help                 help for disable
      --server-name string   Name of the MCP server to remove (default: derived from executable name)
```

### Options inherited from parent commands

```
      --dry-run                  print the equivalent request instead of sending it
      --no-color                 disable colored output
      --offline garminctl sync   read from the local store instead of the Garmin API (see garminctl sync)
  -o, --output string            output format: table|json|yaml|csv (default "table")
      --profile string           profile (Garmin account) to use; env GARMINCTL_PROFILE
```

### SEE ALSO

* [garminctl mcp claude](garminctl_mcp_claude.md)	 - Manage Claude Desktop MCP servers

