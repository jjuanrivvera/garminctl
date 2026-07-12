## garminctl mcp cursor list

Show Cursor MCP servers

### Synopsis

Show all MCP servers configured in Cursor

```
garminctl mcp cursor list [flags]
```

### Options

```
      --config-path string   Path to Cursor config file
  -h, --help                 help for list
      --workspace            List from workspace settings (.cursor/mcp.json) instead of user settings
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

