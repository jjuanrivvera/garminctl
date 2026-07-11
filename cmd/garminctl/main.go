// Command garminctl is a command-line tool for the Garmin Connect API.
package main

import (
	"os"

	"github.com/jjuanrivvera/garminctl/commands"
)

func main() { os.Exit(commands.Main(os.Args[1:])) }
