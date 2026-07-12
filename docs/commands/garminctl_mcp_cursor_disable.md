## garminctl mcp cursor disable

Remove server from Cursor config

### Synopsis

Remove this application from Cursor MCP servers

```
garminctl mcp cursor disable [flags]
```

### Options

```
      --config-path string   Path to Cursor config file
  -h, --help                 help for disable
      --server-name string   Name of the MCP server to remove (default: derived from executable name)
      --workspace            Remove from workspace settings (.cursor/mcp.json) instead of user settings
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

* [garminctl mcp cursor](garminctl_mcp_cursor.md)	 - Manage Cursor MCP servers

