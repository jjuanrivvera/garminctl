## garminctl exercises list

List/search exercises

### Synopsis

List exercises with optional filters. All filters are combined with AND logic.

```
garminctl exercises list [category] [muscle] [equipment] [search] [flags]
```

### Options

```
      --category string    Filter by category (e.g., BENCH_PRESS)
      --equipment string   Filter by equipment (e.g., DUMBBELL)
  -h, --help               help for list
      --muscle string      Filter by muscle group (e.g., CHEST)
      --search string      Search exercise names
```

### Options inherited from parent commands

```
      --dry-run          print the equivalent request instead of sending it
      --no-color         disable colored output
  -o, --output string    output format: table|json|yaml|csv (default "table")
      --profile string   profile (Garmin account) to use; env GARMINCTL_PROFILE
```

### SEE ALSO

* [garminctl exercises](garminctl_exercises.md)	 - exercises commands

