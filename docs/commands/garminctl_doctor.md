## garminctl doctor

Diagnose garminctl setup (config, keyring, sessions)

### Synopsis

Offline health check: verifies the config file is readable and inspects each profile's
stored session — whether both tokens are present and whether the OAuth2 token is still valid.
An expired OAuth2 token is not a failure; the next call refreshes it from the OAuth1 token.

```
garminctl doctor [flags]
```

### Options

```
  -h, --help   help for doctor
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

