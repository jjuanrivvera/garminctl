package commands

import (
	"context"
	"fmt"
	"os"

	"github.com/jjuanrivvera/garminctl/internal/output"
	"github.com/jjuanrivvera/garminctl/internal/version"
)

// Main builds and runs the garminctl root command for the given args, returning a process exit
// code. It lives here (not in package main) so the entry-point logic is testable, and takes the
// signal-cancelled context from main() so Ctrl-C propagates into every request.
func Main(ctx context.Context, args []string) int {
	root := NewRootCmd()
	root.Version = version.Get().Version
	root.SetVersionTemplate(version.String() + "\n")

	// Expand user-defined aliases before cobra parses, so an alias can map to any command path
	// without shadowing a built-in.
	root.SetArgs(ExpandAliases(args))

	if err := root.ExecuteContext(ctx); err != nil {
		// Error text can carry an API-returned body; strip terminal escapes before printing so a
		// crafted value can't hijack the terminal.
		fmt.Fprintln(os.Stderr, "Error:", output.SanitizeTerminal(err.Error()))
		return 1
	}
	return 0
}
