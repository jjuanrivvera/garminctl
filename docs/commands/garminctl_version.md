## garminctl version

Print version, commit, and build date

### Synopsis

Print build metadata. With --check, compare against the latest GitHub release.

```
garminctl version [flags]
```

### Examples

```
  garminctl version
  garminctl version --json
  garminctl version --check
```

### Options

```
      --check   check for a newer release on GitHub
  -h, --help    help for version
      --json    output as JSON
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

