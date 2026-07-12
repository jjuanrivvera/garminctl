## garminctl

Drive Garmin Connect from the terminal

### Synopsis

garminctl reads your Garmin Connect health data — body composition, sleep, heart
rate, stress, body battery, respiration, and intensity minutes — plus the full Connect endpoint
surface via `connect`, with named profiles for several accounts, OS-keyring token storage,
and table/json/yaml/csv output.

### Examples

```
  garminctl auth import --from ~/.garminconnect --profile juan
  garminctl sleep --date 2026-07-09 -o json
  garminctl stress
  garminctl --profile vane body-composition
  garminctl doctor
```

### Options

```
      --dry-run          print the equivalent request instead of sending it
  -h, --help             help for garminctl
      --no-color         disable colored output
  -o, --output string    output format: table|json|yaml|csv (default "table")
      --profile string   profile (Garmin account) to use; env GARMINCTL_PROFILE
```

### SEE ALSO

* [garminctl agent](garminctl_agent.md)	 - AI-agent integration helpers
* [garminctl api](garminctl_api.md)	 - Make a raw authenticated request to the Garmin Connect API
* [garminctl auth](garminctl_auth.md)	 - Manage Garmin Connect authentication
* [garminctl body-battery](garminctl_body-battery.md)	 - Body Battery events for a day
* [garminctl body-composition](garminctl_body-composition.md)	 - Weight, BMI, and body fat for a day
* [garminctl completion](garminctl_completion.md)	 - Generate a shell completion script
* [garminctl config](garminctl_config.md)	 - Manage garminctl configuration and profiles
* [garminctl connect](garminctl_connect.md)	 - The full Garmin Connect endpoint surface (every documented operation)
* [garminctl doctor](garminctl_doctor.md)	 - Diagnose garminctl setup (config, keyring, sessions)
* [garminctl heart-rate](garminctl_heart-rate.md)	 - Daily and resting heart rate for a day
* [garminctl init](garminctl_init.md)	 - Guided first-time setup (import existing tokens)
* [garminctl intensity-minutes](garminctl_intensity-minutes.md)	 - Intensity minutes for a day
* [garminctl mcp](garminctl_mcp.md)	 - MCP server management
* [garminctl respiration](garminctl_respiration.md)	 - All-day respiration for a day
* [garminctl sleep](garminctl_sleep.md)	 - Sleep stages, duration, and score for a day
* [garminctl stress](garminctl_stress.md)	 - All-day stress for a day
* [garminctl version](garminctl_version.md)	 - Print version, commit, and build date

