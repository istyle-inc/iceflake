package iceflake

import (
	"context"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/istyle-inc/iceflake/app"
	"github.com/istyle-inc/iceflake/foundation"
	"github.com/istyle-inc/iceflake/tests/mocks"
)

func TestListen(t *testing.T) {
	// setup mock
	c := gomock.NewController(t)
	defer c.Finish()
	defer func() {
		foundation.InternalTimer = foundation.NewLocalTimer()
	}()
	mock := mocks.NewMockTimer(c)
	mock.EXPECT().Now().Times(1).Return(
		time.Date(2018, 4, 1, 0, 0, 0, 0, time.UTC))
	foundation.InternalTimer = mock

	tmpDir, _ := ioutil.TempDir("", "iceflake")
	defer os.RemoveAll(tmpDir)
	fp := filepath.Join(tmpDir, "iceflake.sock")
	o := &app.Option{
		ListenerType: "unix",
		WorkerID:     1,
		Addr:         fp,
		BaseTime:     time.Date(2018, 1, 1, 0, 0, 0, 0, time.UTC),
	}

	ice, err := app.New(o)
	if err != nil {
		t.Error("error: ", err)
	}
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	go ice.Listen(ctx)
	<-ice.Preparing()

	cli := New("unix", fp)
	result, err := cli.Get()
	if err != nil {
		t.Error("error: ", err)
	}
	expected := uint64(32614907904004097)
	if result.Id != expected {
		t.Error("result expected: ", expected, " but ", result.Id)
	}
}
