package commands

// ExpandAliases expands user-defined command aliases before cobra parses, so an alias can
// map to any command path without shadowing a built-in. The `alias` meta-command manages the
// stored set; this passthrough keeps the entry point stable until it is wired in.
func ExpandAliases(args []string) []string {
	return args
}
