## garminctl

Drive Garmin Connect from the terminal

### Synopsis

garminctl reads your Garmin Connect health data — body composition, sleep, heart
rate, stress, body battery, respiration, and intensity minutes — plus the full Connect endpoint
surface via `connect`, with named profiles for several accounts, OS-keyring token storage,
and table/json/yaml/csv output.

### Examples

```
  garminctl auth import --from ~/.garminconnect --profile me
  garminctl sleep --date 2026-07-09 -o json
  garminctl stress
  garminctl --profile alt body-composition
  garminctl doctor
```

### Options

```
      --dry-run                  print the equivalent request instead of sending it
  -h, --help                     help for garminctl
      --no-color                 disable colored output
      --offline garminctl sync   read from the local store instead of the Garmin API (see garminctl sync)
  -o, --output string            output format: table|json|yaml|csv (default "table")
      --profile string           profile (Garmin account) to use; env GARMINCTL_PROFILE
```

### SEE ALSO

* [garminctl activities](garminctl_activities.md)	 - activities commands
* [garminctl agent](garminctl_agent.md)	 - AI-agent integration helpers
* [garminctl api](garminctl_api.md)	 - Make a raw authenticated request to the Garmin Connect API
* [garminctl auth](garminctl_auth.md)	 - Manage Garmin Connect authentication
* [garminctl biometric](garminctl_biometric.md)	 - biometric commands
* [garminctl body-battery](garminctl_body-battery.md)	 - Body Battery events for a day
* [garminctl body-composition](garminctl_body-composition.md)	 - Weight, BMI, and body fat for a day
* [garminctl calendar](garminctl_calendar.md)	 - calendar commands
* [garminctl completion](garminctl_completion.md)	 - Generate a shell completion script
* [garminctl config](garminctl_config.md)	 - Manage garminctl configuration and profiles
* [garminctl courses](garminctl_courses.md)	 - courses commands
* [garminctl devices](garminctl_devices.md)	 - devices commands
* [garminctl doctor](garminctl_doctor.md)	 - Diagnose garminctl setup (config, keyring, sessions)
* [garminctl exercises](garminctl_exercises.md)	 - exercises commands
* [garminctl fitnessage](garminctl_fitnessage.md)	 - fitnessage commands
* [garminctl fitnessstats](garminctl_fitnessstats.md)	 - fitnessstats commands
* [garminctl heart-rate](garminctl_heart-rate.md)	 - Daily and resting heart rate for a day
* [garminctl history](garminctl_history.md)	 - Query the local store for a metric across a date range (offline)
* [garminctl hrv](garminctl_hrv.md)	 - hrv commands
* [garminctl init](garminctl_init.md)	 - Guided first-time setup (import existing tokens)
* [garminctl intensity-minutes](garminctl_intensity-minutes.md)	 - Intensity minutes for a day
* [garminctl mcp](garminctl_mcp.md)	 - MCP server management
* [garminctl metrics](garminctl_metrics.md)	 - metrics commands
* [garminctl profile](garminctl_profile.md)	 - profile commands
* [garminctl respiration](garminctl_respiration.md)	 - All-day respiration for a day
* [garminctl sleep](garminctl_sleep.md)	 - Sleep stages, duration, and score for a day
* [garminctl stress](garminctl_stress.md)	 - All-day stress for a day
* [garminctl sync](garminctl_sync.md)	 - Backfill daily metrics into the local store for offline use
* [garminctl version](garminctl_version.md)	 - Print version, commit, and build date
* [garminctl weight](garminctl_weight.md)	 - weight commands
* [garminctl wellness](garminctl_wellness.md)	 - wellness commands
* [garminctl workouts](garminctl_workouts.md)	 - workouts commands

