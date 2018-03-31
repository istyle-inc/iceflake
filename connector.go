package main

import (
	"log"
	"net"
	"os"

	"github.com/istyle-inc/iceflake/foundation"
)

type Connector struct {
	SocketFilePath string
	ListenType     string
	Listener       net.Listener
	Generator      *Generator
}

func NewConnector(socketFilePath, listenType string, generator *Generator) *Connector {
	return &Connector{
		SocketFilePath: socketFilePath,
		ListenType:     listenType,
		Generator:      generator,
	}
}

func (c *Connector) Listen() error {
	l, err := net.Listen(c.ListenType, c.SocketFilePath)
	if err != nil {
		return err
	}

	c.Listener = l
	return nil
}

func (c *Connector) AcceptListener() error {
	// Listen socket
	for {
		conn, err := c.Listener.Accept()
		if err != nil {
			log.Println("Error: ", err)
			continue
		}

		// Process after accepted
		uuid, err := c.Generator.Generate()
		if err != nil {
			return err
		}

		// Send UUID
		go func(conn net.Conn, uuid []byte) {
			defer conn.Close()
			_, err := conn.Write([]byte(uuid))
			if err != nil {
				log.Println("Error: ", err)
			}
		}(conn, []byte(uuid))
	}
}

func (c *Connector) SignalTearDown() {
	foundation.SLogger.Infof("Shutting down.\n")
	c.Listener.Close()
	os.Remove(c.SocketFilePath)
	os.Exit(0)
}
