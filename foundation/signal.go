package foundation

import (
	"os"
	"os/signal"
)

var Exit = os.Exit
var SignalCh = make(chan os.Signal, 1)

type SignalHandler interface {
	SignalTearDown()
}

// Shutdown when notice interrupt signal
func SignalHandling(sh SignalHandler) {
	signal.Notify(SignalCh, os.Interrupt)
	go func(sh SignalHandler, s chan os.Signal) {
		sig := <-s
		SLogger.Infof("Catch signal %s: Shutting down.\n", sig)
		sh.SignalTearDown()
		Exit(0)
	}(sh, SignalCh)
}
