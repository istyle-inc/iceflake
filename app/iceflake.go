package app

import (
	"bufio"
	"context"
	"fmt"
	"net"
	"time"

	"github.com/apex/log"
)

// Option struct hold option value
type Option struct {
	// ListenerType listen tcp, unix or else
	ListenerType string
	// WorkerID given worker-id to this app instance
	WorkerID int
	// Addr tcp address or unix socket path
	Addr string
	// BaseTime
	BaseTime time.Time
}

// IceFlake app itself
type IceFlake struct {
	// Listener instance of net.Listener implementation
	net.Listener
	// IDGenerator Implant IDGenerator
	IDGenerator
	preparing chan interface{}
	baseTime  time.Time
	option    *Option
}

func listener(o *Option) (net.Listener, error) {
	return net.Listen(o.ListenerType, o.Addr)
}

// New return new IceFlake instance
func New(o *Option) (*IceFlake, error) {
	l, err := listener(o)
	if err != nil {
		return nil, err
	}
	return &IceFlake{
		Listener:    l,
		IDGenerator: NewIDGenerator(o.WorkerID),
		baseTime:    o.BaseTime,
		option:      o,
		preparing:   make(chan interface{}),
	}, nil
}

// Preparing returns channel to notify listener's listened up
func (i *IceFlake) Preparing() chan interface{} {
	return i.preparing
}

// Listen process start with listener
func (i *IceFlake) Listen(ctx context.Context) error {
	go func() {
		<-ctx.Done()
		if err := i.Close(); err != nil {
			log.Errorf("error occurred when closing: %s", err)
		}
	}()
	close(i.preparing)
	for {
		conn, err := i.Accept()
		if err != nil {
			select {
			case <-ctx.Done():
				log.Info("Shutting down listener")
				return nil
			default:
				log.Errorf("error occurred while processing: %s", err)
				return err
			}
		}
		go i.handle(ctx, conn)
	}
}

func (i *IceFlake) handle(ctx context.Context, conn net.Conn) {
	innerCtx, cancel := context.WithCancel(ctx)
	defer cancel()
	go func() {
		<-innerCtx.Done()
		_ = conn.Close()
		log.Info("connection closed")
	}()
	id, err := i.Generate()
	if err != nil {
		log.Errorf("error with generation id: %s", err)
		return
	}
	writer := bufio.NewWriter(conn)
	_, err = fmt.Fprint(writer, id)
	if err != nil {
		log.Errorf("error with writing to stream: %s", err)
		return
	}
	if err = writer.Flush(); err != nil {
		log.Errorf("error with flushing stream: %s", err)
		return
	}
}
