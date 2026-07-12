## garminctl fitnessstats get

Get fitness statistics

### Synopsis

Get aggregated fitness statistics for activities including calories, distance, and duration over a date range

```
garminctl fitnessstats get [aggregation] [metrics] [group_by_activity_type] [standardized_units] [flags]
```

### Options

```
      --aggregation string       Aggregation period: daily, weekly, monthly, yearly (default: weekly)
      --end string               End date (YYYY-MM-DD)
      --group_by_activity_type   Group stats by activity type (e.g., running, hiking)
  -h, --help                     help for get
      --metrics string           Comma-separated metrics: calories, distance, duration, avgSpeed, maxHr, avgHr, elevationGain, avgRunCadence, avgGroundContactBalance, avgStrideLength, avgVerticalOscillation, avgVerticalRatio, avgGroundContactTime, aerobicTrainingEffect, anaerobicTrainingEffect (default: calories,distance,duration)
      --standardized_units       Use standardized units in response
      --start string             Start date (YYYY-MM-DD)
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

* [garminctl fitnessstats](garminctl_fitnessstats.md)	 - fitnessstats commands

