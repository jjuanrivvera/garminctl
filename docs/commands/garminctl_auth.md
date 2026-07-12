## garminctl auth

Manage Garmin Connect authentication

### Synopsis

Store and verify the Garmin session for a profile. A session (OAuth1 + OAuth2 tokens)
lives in your OS keyring, keyed by profile. Bring one in with 'auth import' (from an existing
garth / python-garminconnect token dir) or 'auth login' (email + password).

### Options

```
  -h, --help   help for auth
```

### Options inherited from parent commands

```
      --dry-run          print the equivalent request instead of sending it
      --no-color         disable colored output
  -o, --output string    output format: table|json|yaml|csv (default "table")
      --profile string   profile (Garmin account) to use; env GARMINCTL_PROFILE
```

### SEE ALSO

* [garminctl](garminctl.md)	 - Drive Garmin Connect from the terminal
* [garminctl auth import](garminctl_auth_import.md)	 - Import an existing garth / python-garminconnect token dir into the keyring
* [garminctl auth login](garminctl_auth_login.md)	 - Log in with Garmin credentials and store the session in the keyring
* [garminctl auth logout](garminctl_auth_logout.md)	 - Remove the stored session for the active profile
* [garminctl auth status](garminctl_auth_status.md)	 - Show whether the active profile has a stored, valid session

