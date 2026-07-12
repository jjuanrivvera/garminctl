## garminctl fitnessage stats

Get fitness age statistics

### Synopsis

Get daily fitness age statistics including fitness age, achievable fitness age, RHR, BMI, and vigorous activity days. Note: date range must be 28 days or less.

```
garminctl fitnessage stats [flags]
```

### Options

```
      --end string     End date (YYYY-MM-DD)
  -h, --help           help for stats
      --start string   Start date (YYYY-MM-DD)
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

* [garminctl fitnessage](garminctl_fitnessage.md)	 - fitnessage commands

