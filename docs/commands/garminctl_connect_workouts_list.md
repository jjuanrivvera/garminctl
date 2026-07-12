## garminctl connect workouts list

List workouts

### Synopsis

List workouts with pagination including name, sport type, and estimated duration

```
garminctl connect workouts list [start] [limit] [flags]
```

### Options

```
  -h, --help        help for list
      --limit int   Maximum number of workouts to return (defaults to 20)
      --start int   Starting index (0-based, defaults to 0)
```

### Options inherited from parent commands

```
      --dry-run          print the equivalent request instead of sending it
      --no-color         disable colored output
  -o, --output string    output format: table|json|yaml|csv (default "table")
      --profile string   profile (Garmin account) to use; env GARMINCTL_PROFILE
```

### SEE ALSO

* [garminctl connect workouts](garminctl_connect_workouts.md)	 - workouts commands

