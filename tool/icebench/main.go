package main

import (
	"errors"
	"flag"
	"log"
	"net"
)

var (
	// Flags
	socketPathOption = flag.String("s", DefaultSocketFilePath, "Setting socket path")
)

const (
	// DefaultSocketFilePath socket file path using dial default
	DefaultSocketFilePath = "/var/run/iceflake/iceflake-worker-1.sock"
	// ListenType dial type
	ListenType = "unix"
)

func main() {
	flag.Parse()
	r, err := connectIceFlake(*socketPathOption)
	log.Println(r, err)
}

func connectIceFlake(socketPath string) (result string, err error) {
	var c net.Conn
	c, err = net.Dial(ListenType, socketPath)
	if err != nil {
		return
	}
	defer func() {
		_ = c.Close()
	}()
	buf := make([]byte, 1024)
	var n int
	n, err = c.Read(buf)
	if n == 0 {
		err = errors.New("iceflake returned nothing")
	}
	if err != nil {
		return
	}
	return string(buf[:n]), nil
}
