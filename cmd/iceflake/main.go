package main

import (
	"context"
	"flag"
	"syscall"
	"time"

	"github.com/istyle-inc/iceflake/app"
	"github.com/istyle-inc/iceflake/constantvalues"
	"github.com/istyle-inc/iceflake/foundation"
	"github.com/syossan27/tebata"
	"go.uber.org/zap"
)

func main() {
	workerIDOption := flag.Int64("w", constantvalues.DefaultWorkerID, "Setting worker id of iceflake")
	socketPathOption := flag.String("s", constantvalues.DefaultSocketFilePath, "Setting socket path")
	defer foundation.Logger.Sync()
	flag.Parse()

	ctx, cancel := context.WithCancel(context.Background())
	app, err := app.New(&app.Option{
		BaseTime:     time.Date(2018, 1, 1, 0, 0, 0, 0, time.UTC),
		ListenerType: "unix",
		Addr:         *socketPathOption,
		WorkerID:     uint64(*workerIDOption),
	})
	t := tebata.New(syscall.SIGINT, syscall.SIGKILL, syscall.SIGHUP, syscall.SIGTERM)
	t.Reserve(cancel)
	if err != nil {
		return
	}
	if err := app.Listen(ctx); err != nil {
		foundation.Logger.Fatal("Error: ", zap.Error(err))
	}
}
