## garminctl auth import

Import an existing garth / python-garminconnect token dir into the keyring

### Synopsis

Migrate an existing garth session (oauth1_token.json + oauth2_token.json, e.g.
~/.garminconnect) into garminctl's keyring under the active profile. No login required — the
OAuth1 token (valid ~1 year) drives OAuth2 refresh from here, so this fixes the recurring
"username and password are required" failure of a cron that never refreshed its cached tokens.

```
garminctl auth import [flags]
```

### Examples

```
  garminctl --profile me auth import --from ~/.garminconnect
  garminctl --profile alt auth import --from ~/.garminconnect-alt
```

### Options

```
      --from string   garth token dir (default: ~/.garminconnect)
  -h, --help          help for import
```

### Options inherited from parent commands

```
      --dry-run          print the equivalent request instead of sending it
      --no-color         disable colored output
  -o, --output string    output format: table|json|yaml|csv (default "table")
      --profile string   profile (Garmin account) to use; env GARMINCTL_PROFILE
```

### SEE ALSO

* [garminctl auth](garminctl_auth.md)	 - Manage Garmin Connect authentication

