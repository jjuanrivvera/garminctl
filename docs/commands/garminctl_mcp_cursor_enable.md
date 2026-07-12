## garminctl mcp cursor enable

Add server to Cursor config

### Synopsis

Add this application as an MCP server in Cursor

```
garminctl mcp cursor enable [flags]
```

### Options

```
      --config-path string   Path to Cursor config file
  -e, --env stringToString   Environment variables (e.g., --env KEY1=value1 --env KEY2=value2) (default [])
  -h, --help                 help for enable
      --log-level string     Log level (debug, info, warn, error)
      --server-name string   Name for the MCP server (default: derived from executable name)
      --workspace            Add to workspace settings (.cursor/mcp.json) instead of user settings
```

### Options inherited from parent commands

```
      --dry-run          print the equivalent request instead of sending it
      --no-color         disable colored output
  -o, --output string    output format: table|json|yaml|csv (default "table")
      --profile string   profile (Garmin account) to use; env GARMINCTL_PROFILE
```

### SEE ALSO

* [garminctl mcp cursor](garminctl_mcp_cursor.md)	 - Manage Cursor MCP servers

