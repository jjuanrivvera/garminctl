## garminctl history

Query the local store for a metric across a date range (offline)

### Synopsis

history reads the offline store (populated by `garminctl sync` or by earlier reads) for
one metric over a date range and renders one row per day — pair it with -o csv for a trend you can
open in a spreadsheet. It never hits the network. Defaults to the last 30 days.

```
garminctl history <metric> [flags]
```

### Examples

```
  garminctl history body-composition --from 2026-01-01 -o csv
  garminctl history sleep --from 2026-06-01 --to 2026-07-10
```

### Options

```
      --from string   start date YYYY-MM-DD (default: 30 days ago)
  -h, --help          help for history
      --to string     end date YYYY-MM-DD (default: today)
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

