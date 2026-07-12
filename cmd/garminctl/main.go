// Command garminctl is a command-line tool for the Garmin Connect API.
package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/jjuanrivvera/garminctl/commands"
)

func main() {
	// signal.NotifyContext makes Ctrl-C (SIGINT/SIGTERM) cancel in-flight work: token refresh,
	// retry backoff, and rate-limit waits all observe this context.
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()
	os.Exit(commands.Main(ctx, os.Args[1:]))
}
