## garminctl api

Make a raw authenticated request to the Garmin Connect API

### Synopsis

Escape hatch for Connect endpoints the typed commands don't wrap.

<path> is everything after https://connectapi.garmin.com — for example
/usersummary-service/usersummary/daily?calendarDate=2026-07-10

The request is signed with the active profile's session token. The typed commands refresh that
token through go-garmin and persist it; if a raw call returns HTTP 401, run any resource
command once (e.g. `garminctl sleep`) to refresh, then retry. With --dry-run, prints the
equivalent curl (token redacted) instead of sending.

```
garminctl api <path> [flags]
```

### Examples

```
  garminctl api /userprofile-service/userprofile
  garminctl api -X GET /usersummary-service/usersummary/daily
  garminctl --dry-run api /userprofile-service/userprofile
```

### Options

```
      --data string     request body (JSON) for POST/PUT
  -h, --help            help for api
  -X, --method string   HTTP method (GET|POST|PUT|DELETE) (default "GET")
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

