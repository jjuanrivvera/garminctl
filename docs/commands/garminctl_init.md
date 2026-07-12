## garminctl init

Guided first-time setup (import existing tokens)

### Synopsis

First-run setup. Looks for an existing garth / python-garminconnect token directory
(default ~/.garminconnect) and imports it into the keyring as a profile, so a working Python
setup carries straight over. If none is found, prints the two ways to authenticate.

The profile name comes from the global --profile flag (default "default").

```
garminctl init [flags]
```

### Options

```
      --from string   garth token dir to import (default ~/.garminconnect)
  -h, --help          help for init
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

