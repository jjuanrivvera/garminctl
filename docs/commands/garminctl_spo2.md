## garminctl spo2

Pulse ox (SpO2) summary for a day

```
garminctl spo2 [flags]
```

### Examples

```
  garminctl --profile me spo2 --date 2026-07-10 -o json
```

### Options

```
      --date string   date YYYY-MM-DD (default: today)
  -h, --help          help for spo2
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

