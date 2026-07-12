## garminctl courses import

Import a GPX course

### Synopsis

Import a course/route from a GPX file to Garmin Connect

```
garminctl courses import <file> [activity-type] [privacy] [flags]
```

### Options

```
      --activity-type int   Activity type ID (e.g. 1=running, 3=hiking, 5=cycling)
  -h, --help                help for import
      --privacy int         Privacy rule: 1=Public, 2=Private (default), 4=Group
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

* [garminctl courses](garminctl_courses.md)	 - courses commands

