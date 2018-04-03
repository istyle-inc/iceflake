package app

import (
	"context"
	"net"
	"time"

	"github.com/golang/protobuf/proto"
	"github.com/istyle-inc/iceflake/foundation"
	"github.com/istyle-inc/iceflake/pbdef"
	"go.uber.org/zap"
)

// Option struct hold option value
type Option struct {
	// ListenerType listen tcp, unix or else
	ListenerType string
	// WorkerID given worker-id to this app instance
	WorkerID uint64
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
		IDGenerator: NewIDGenerator(o.WorkerID, o.BaseTime),
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
			foundation.Logger.Error("error occurred when closing", zap.Error(err))
		}
	}()
	close(i.preparing)
	for {
		conn, err := i.Accept()
		if err != nil {
			select {
			case <-ctx.Done():
				foundation.Logger.Info("Shutting down listener")
				return nil
			default:
				foundation.Logger.Error("error occurred while processing", zap.Error(err))
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
		foundation.Logger.Info("connection closed")
	}()
	id, err := i.Generate()
	if err != nil {
		foundation.Logger.Error("error with generation id", zap.Error(err))
		return
	}
	flake := pbdef.IceFlake{Id: id}
	data, err := proto.Marshal(&flake)
	_, err = conn.Write(data)
	if err != nil {
		foundation.Logger.Error("error with writing to strea", zap.Error(err))
		return
	}
}
