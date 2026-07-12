## garminctl body-battery

Body Battery events for a day

```
garminctl body-battery [flags]
```

### Examples

```
  garminctl --profile me body-battery --date 2026-07-10 -o json
```

### Options

```
      --date string   date YYYY-MM-DD (default: today)
  -h, --help          help for body-battery
```

### Options inherited from parent commands

```
      --dry-run          print the equivalent request instead of sending it
      --no-color         disable colored output
  -o, --output string    output format: table|json|yaml|csv (default "table")
      --profile string   profile (Garmin account) to use; env GARMINCTL_PROFILE
```

### SEE ALSO

* [garminctl](garminctl.md)	 - Drive Garmin Connect from the terminal

