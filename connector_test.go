package main

import (
	"bytes"
	"encoding/json"
	"io"
	"os"
	"syscall"
	"testing"

	"github.com/istyle-inc/iceflake/foundation"

	"github.com/koron/go-dproxy"
	"github.com/syossan27/tebata"
	"go.uber.org/zap"
)

var (
	testSocketFilePath        = "/run/iceflake-test.sock"
	testListenType            = "unix"
	testWorkerId       uint64 = 1
	connector                 = newConnector()
)

func newConnector() *Connector {
	generator := NewGenerator(testWorkerId)
	return NewConnector(
		testSocketFilePath,
		testListenType,
		generator,
	)
}

func TestNewConnector(t *testing.T) {
	defer os.Remove(testSocketFilePath)

	generator := NewGenerator(testWorkerId)

	connector := &Connector{
		SocketFilePath: testSocketFilePath,
		ListenType:     testListenType,
		Generator:      generator,
	}

	newConnector := NewConnector(
		testSocketFilePath,
		testListenType,
		generator,
	)

	if *connector != *newConnector {
		t.Error("Create invalid connector")
	}
}

func TestConnector_Listen(t *testing.T) {
	defer os.Remove(testSocketFilePath)

	err := connector.Listen()
	if err != nil {
		t.Error("Fail connector listen: ", err)
	}
	connector.Listener.Close()
}

func TestConnector_AcceptListener(t *testing.T) {
	defer os.Remove(testSocketFilePath)

	// Override exit function used for signalHandler
	done := make(chan int, 1)

	// Change zap output to none
	stderr := os.Stderr
	r, w, _ := os.Pipe()
	os.Stderr = w
	foundation.Logger, _ = zap.NewProduction()
	foundation.SLogger = foundation.Logger.Sugar()

	connector := newConnector()
	err := connector.Listen()
	if err != nil {
		t.Error("Fail connector listen: ", err)
	}

	tbt := tebata.New(syscall.SIGINT, syscall.SIGKILL)
	tbt.Reserve(
		func() {
			foundation.SLogger.Infof("Shutting down.\n")
			done <- 1
		},
	)
	go connector.AcceptListener()
	p, err := os.FindProcess(os.Getpid())
	if err != nil {
		t.Error("Fail find process: ", err)
	}
	p.Signal(os.Interrupt)

	// Waiting done signalHandler
	<-done

	// Tear down handling zap output
	w.Close()
	var buf bytes.Buffer
	io.Copy(&buf, r)
	os.Stderr = stderr

	if buf.Len() == 0 {
		t.Error("Output empty")
	}

	var stdOutJSON interface{}
	json.Unmarshal(buf.Bytes(), &stdOutJSON)
	msg, err := dproxy.New(stdOutJSON).M("msg").String()
	if msg != "Shutting down.\n" {
		t.Error("Invalid output")
	}
}
