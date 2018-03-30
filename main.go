package main

import (
	"flag"

	"github.com/istyle-inc/iceflake/foundation"

	"syscall"

	"github.com/syossan27/tebata"
	"go.uber.org/zap"
)

var (
	// Flags
	workerIdOption   = flag.Int64("w", 1, "Setting worker id of iceflake")
	socketPathOption = flag.String("s", "/run/iceflake-worker-1.sock", "Setting socket path")
)

const (
	ListenType = "unix"
)

func main() {
	defer foundation.Logger.Sync()
	flag.Parse()

	connector := NewConnector(
		*socketPathOption,
		ListenType,
		NewGenerator(uint64(*workerIdOption)),
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
