## garminctl workouts create

Create a new workout

### Synopsis

Create a new workout with segments and steps. Use --file to read from a file, --json to pass inline JSON, or pipe JSON to stdin.

```
garminctl workouts create [flags]
```

### Options

```
  -f, --file string   Read JSON body from file
  -h, --help          help for create
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

* [garminctl workouts](garminctl_workouts.md)	 - workouts commands

