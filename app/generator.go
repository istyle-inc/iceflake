package app

import (
	"errors"
	"sync"
	"time"

	"github.com/istyle-inc/iceflake/foundation"
)

// IDGenerator interface of generator generates each Unique ID
type IDGenerator interface {
	Generate() (uint64, error)
}

const (
	workerIDBits            = 10
	sequenceBits            = 12
	initialSequentialNumber = 1
)

var (
	sequentialNumber uint64 = 1
)

// IceFlakeGenerator local implement of IDGenerator
type IceFlakeGenerator struct {
	w        uint64
	baseTime time.Time
	lastTS   uint64
	seq      uint64
	gate     sync.Mutex
}

// NewIDGenerator return new IDGenerator instance
func NewIDGenerator(workerID uint64, baseTime time.Time) IDGenerator {
	return &IceFlakeGenerator{
		w:        workerID,
		baseTime: baseTime,
		lastTS:   0,
		seq:      initialSequentialNumber,
	}
}

// Generate generate unique id
func (g *IceFlakeGenerator) Generate() (uint64, error) {
	g.gate.Lock()
	defer g.gate.Unlock()

	ts := g.GetTimeInt()
	switch {
	case ts < g.lastTS:
		return 0, errors.New("system clock was rolled back")
	case ts == g.lastTS:
		g.seq = g.seq + 1
	case ts > g.lastTS:
		g.seq = initialSequentialNumber
	}
	g.lastTS = ts
	return ts<<(workerIDBits+sequenceBits) | g.w<<sequenceBits | g.seq, nil
}

// GetTimeInt get uint value differ between now and epochtime
func (g *IceFlakeGenerator) GetTimeInt() uint64 {
	return uint64(foundation.InternalTimer.Now().Sub(g.baseTime).Round(time.Millisecond)) / uint64(time.Millisecond)
}
