## garminctl connect workouts update

Update an existing workout

### Synopsis

Update an existing workout. Use --file to read from a file, --json to pass inline JSON, or pipe JSON to stdin.

```
garminctl connect workouts update <workout_id> [flags]
```

### Options

```
  -f, --file string   Read JSON body from file
  -h, --help          help for update
      --json string   JSON body as string
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

