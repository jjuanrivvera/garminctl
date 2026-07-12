## garminctl completion

Generate a shell completion script

### Synopsis

Output a completion script for your shell. See `garminctl completion <shell> --help` for install instructions.

```
garminctl completion [bash|zsh|fish|powershell]
```

### Examples

```
  source <(garminctl completion bash)
  garminctl completion zsh > "${fpath[1]}/_garminctl"
  garminctl completion fish > ~/.config/fish/completions/garminctl.fish
```

### Options

```
  -h, --help   help for completion
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

