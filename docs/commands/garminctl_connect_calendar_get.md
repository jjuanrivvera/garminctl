## garminctl connect calendar get

Get calendar

### Synopsis

Get calendar items including activities, workouts, and weight entries. Parameters are hierarchical: month requires year, day requires both month and start.

```
garminctl connect calendar get <year> [month] [day] [start] [flags]
```

### Options

```
      --day int     Day of month (requires month and start)
  -h, --help        help for get
      --month int   Month (0-11, January=0)
      --start int   Week start day, 1=Monday (required when day is provided)
```

### Options inherited from parent commands

```
      --dry-run          print the equivalent request instead of sending it
      --no-color         disable colored output
  -o, --output string    output format: table|json|yaml|csv (default "table")
      --profile string   profile (Garmin account) to use; env GARMINCTL_PROFILE
```

### SEE ALSO

* [garminctl connect calendar](garminctl_connect_calendar.md)	 - calendar commands

