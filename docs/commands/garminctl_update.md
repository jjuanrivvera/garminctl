## garminctl update

Update garminctl to the latest release

### Synopsis

Check GitHub for a newer release and, if one exists, download it, verify it
against the release checksums, and replace the running binary in place.

```
garminctl update [flags]
```

### Options

```
  -h, --help   help for update
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
* [garminctl update check](garminctl_update_check.md)	 - Check for a newer release without installing it

