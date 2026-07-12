## garminctl auth login

Log in with Garmin credentials and store the session in the keyring

```
garminctl auth login [flags]
```

### Examples

```
  garminctl --profile juan auth login
```

### Options

```
      --email string   Garmin account email (omit to be prompted)
  -h, --help           help for login
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

