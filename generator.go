package iceflake

import (
	"errors"
	"sync"
	"time"

	"github.com/istyle-inc/iceflake/foundation"
)

// Generator interface of generator generates each Unique ID
type Generator interface {
	Generate() (uint64, error)
}

const (
	workerIDBits            = 10
	sequenceBits            = 12
	initialSequentialNumber = 1
)

// GeneratorService local implement of Generator
type GeneratorService struct {
	w        uint64
	baseTime time.Time
	lastTS   uint64
	seq      uint64
	gate     sync.Mutex
}

// New return new Generator instance
func New(workerID uint64, baseTime time.Time) Generator {
	return &GeneratorService{
		w:        workerID,
		baseTime: baseTime,
		lastTS:   0,
		seq:      initialSequentialNumber,
	}
}

// Generate generate unique id
func (g *GeneratorService) Generate() (uint64, error) {
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
func (g *GeneratorService) GetTimeInt() uint64 {
	return uint64(foundation.InternalTimer.Now().Sub(g.baseTime).Round(time.Millisecond)) / uint64(time.Millisecond)
}
