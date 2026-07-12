## garminctl sync

Backfill daily metrics into the local store for offline use

### Synopsis

sync fetches a date range of daily metrics for the active profile and records them in the
local SQLite store, so `garminctl --offline <metric>` and `garminctl history` work without the
network. Re-syncing a date overwrites it with fresh data. Defaults to the last 7 days and all
curated metrics.

```
garminctl sync [flags]
```

### Examples

```
  garminctl sync
  garminctl sync --from 2026-01-01 --to 2026-07-10
  garminctl sync --metrics sleep,body-composition --from 2026-06-01
```

### Options

```
      --from string      start date YYYY-MM-DD (default: 7 days ago)
  -h, --help             help for sync
      --metrics string   comma-separated metrics to sync (default: all)
      --to string        end date YYYY-MM-DD (default: today)
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

* [garminctl](garminctl.md)	 - Drive Garmin Connect from the terminal

