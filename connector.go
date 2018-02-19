package main

import (
	"fmt"
	"net"
	"os"
	"os/signal"

	"go.uber.org/zap"
)

const (
	SocketFilePath    = "/tmp/iceflake-worker-%d.sock"
	ListenNetworkType = "unix"
)

func socketConnect(workerId int64) {
	f := fmt.Sprintf(SocketFilePath, workerId)
	defer os.Remove(f)

	listener, err := net.Listen(ListenNetworkType, f)
	if err != nil {
		logger.Fatal("Error: ", zap.Error(err))
	}

	// Shutdown when notice interrupt signal
	signalCh := make(chan os.Signal, 1)
	signal.Notify(signalCh, os.Interrupt)
	go signalHandler(listener, signalCh)

	// Listen socket
	for {
		conn, err := listener.Accept()
		if err != nil {
			os.Remove(f)
			logger.Fatal("Error: ", zap.Error(err))
		}

		go sendUUID(uint64(workerId), conn)
	}
}

func signalHandler(listener net.Listener, s chan os.Signal) {
	sig := <-s
	sLogger.Infof("Catch signal %s: Shutting down.\n", sig)
	listener.Close()
	os.Exit(0)
}

func sendUUID(workerId uint64, conn net.Conn) {
	g := NewGenerator(workerId)
	uuid := g.Generate()
	conn.Write([]byte(uuid))
	conn.Close()
}
