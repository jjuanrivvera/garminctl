package commands

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/jjuanrivvera/garminctl/internal/output"
	"github.com/jjuanrivvera/garminctl/internal/version"
)

// Main builds and runs the garminctl root command for the given args, returning a process exit
// code. It lives here (not in package main) so the entry-point logic is testable.
func Main(args []string) int {
	// signal.NotifyContext makes Ctrl-C (SIGINT/SIGTERM) cancel in-flight work: token refresh,
	// retry backoff, and rate-limit waits all observe this context.
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

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
