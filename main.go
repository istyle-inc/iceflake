package main

import (
	"flag"

	"github.com/istyle-inc/iceflake/constantvalues"
	"github.com/istyle-inc/iceflake/foundation"

	"syscall"

	"github.com/syossan27/tebata"
	"go.uber.org/zap"
)

var (
	// Flags
	workerIDOption   = flag.Int64("w", constantvalues.DefaultWorkerID, "Setting worker id of iceflake")
	socketPathOption = flag.String("s", constantvalues.DefaultSocketFilePath, "Setting socket path")
)

func main() {
	defer foundation.Logger.Sync()
	flag.Parse()

	connector := NewConnector(
		*socketPathOption,
		constantvalues.ListenType,
		NewGenerator(uint64(*workerIDOption)),
	)

	err := connector.Listen()
	if err != nil {
		foundation.Logger.Fatal("Error: ", zap.Error(err))
	}

	// Catch interrupt signal for stop listen
	t := tebata.New(syscall.SIGINT, syscall.SIGKILL, syscall.SIGHUP, syscall.SIGTERM)
	t.Reserve(connector.SignalTearDown)

	err = connector.AcceptListener()
	if err != nil {
		foundation.Logger.Fatal("Error: ", zap.Error(err))
	}
}
