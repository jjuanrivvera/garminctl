## garminctl courses download-fit

Download course as FIT

### Synopsis

Download a course/route as a FIT file. Output goes to stdout by default, use -o to write to a file.

```
garminctl courses download-fit <course_id> [flags]
```

### Options

```
  -h, --help            help for download-fit
  -o, --output string   Output file path
```

### Options inherited from parent commands

```
      --dry-run                  print the equivalent request instead of sending it
      --no-color                 disable colored output
      --offline garminctl sync   read from the local store instead of the Garmin API (see garminctl sync)
      --profile string           profile (Garmin account) to use; env GARMINCTL_PROFILE
```

### SEE ALSO

* [garminctl courses](garminctl_courses.md)	 - courses commands

