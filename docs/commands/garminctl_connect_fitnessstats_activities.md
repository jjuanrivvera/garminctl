## garminctl connect fitnessstats activities

Get individual activity data

### Synopsis

Get individual activity data without aggregation, including activity names, types, and training effects

```
garminctl connect fitnessstats activities [activity_type] [metrics] [flags]
```

### Options

```
      --activity_type string   Filter by activity type (e.g., running, hiking, cycling)
      --end string             End date (YYYY-MM-DD)
  -h, --help                   help for activities
      --metrics string         Comma-separated metrics: calories, distance, duration, avgSpeed, maxHr, avgHr, elevationGain, avgRunCadence, avgGroundContactBalance, avgStrideLength, avgVerticalOscillation, avgVerticalRatio, avgGroundContactTime, startLocal, activityType, activitySubType, name, aerobicTrainingEffect, anaerobicTrainingEffect (default: name,startLocal,activityType,duration,distance,calories)
      --start string           Start date (YYYY-MM-DD)
```

### Options inherited from parent commands

```
      --dry-run          print the equivalent request instead of sending it
      --no-color         disable colored output
  -o, --output string    output format: table|json|yaml|csv (default "table")
      --profile string   profile (Garmin account) to use; env GARMINCTL_PROFILE
```

### SEE ALSO

* [garminctl connect fitnessstats](garminctl_connect_fitnessstats.md)	 - fitnessstats commands

