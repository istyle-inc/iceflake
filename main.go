package main

import (
	"flag"
	"fmt"

	"iceflake/foundation"

	"go.uber.org/zap"
)

var (
	// Flags
	workerIdOption = flag.Int64("w", 1, "Setting worker id of iceflake")
)

const (
	SocketFilePath = "/tmp/iceflake-worker-%d.sock"
	ListenType     = "unix"
)

func main() {
	defer foundation.Logger.Sync()
	flag.Parse()

	connector := NewConnector(
		fmt.Sprintf(SocketFilePath, *workerIdOption),
		ListenType,
		NewGenerator(uint64(*workerIdOption)),
	)

	err := connector.Listen()
	if err != nil {
		foundation.Logger.Fatal("Error: ", zap.Error(err))
	}
	foundation.SignalHandling(connector) // Catch interrupt signal for stop listen
	err = connector.AcceptListener()
	if err != nil {
		foundation.Logger.Fatal("Error: ", zap.Error(err))
	}
}
