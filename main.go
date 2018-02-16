package main

import (
	"flag"
)

var (
	// Flags
	workerIdOption = flag.Int64("w", 1, "Setting worker id of iceflake")
)

func main() {
	defer logger.Sync()

	flag.Parse()
	socketConnect(*workerIdOption)
}
