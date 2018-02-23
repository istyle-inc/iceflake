package main

import (
	"strconv"
	"sync"
	"time"

	"github.com/pkg/errors"
)

const (
	WorkerIdBits                   = 10
	SequenceBits                   = 12
	InitialSequentialNumber uint64 = 1
	DecimalRadix10                 = 10
)

var (
	sequentialNumber uint64 = 1
	GeneratorLock    sync.Mutex
	ServiceEpochTime = time.Date(2018, 1, 1, 0, 0, 0, 0, time.UTC)
)

type Generator struct {
	Timestamp        uint64
	LastTimestamp    uint64
	WorkerId         uint64
	SequentialNumber uint64
}

func NewGenerator(workerId uint64) *Generator {
	return &Generator{
		WorkerId:         workerId,
		SequentialNumber: InitialSequentialNumber,
	}
}

// Generate UUID
func (g *Generator) Generate() (string, error) {
	GeneratorLock.Lock()
	defer GeneratorLock.Unlock()

	g.Timestamp = g.GetTimestamp(ServiceEpochTime)
	if g.Timestamp < g.LastTimestamp {
		return "", errors.New("system clock was rolled back")
	}
	if g.Timestamp == g.LastTimestamp {
		sequentialNumber++
		g.SequentialNumber = sequentialNumber
	}
	if g.Timestamp > g.LastTimestamp {
		sequentialNumber = InitialSequentialNumber
	}
	g.LastTimestamp = g.Timestamp

	uuid := strconv.FormatUint(
		(g.Timestamp<<(WorkerIdBits+SequenceBits))|g.WorkerId<<SequenceBits|g.SequentialNumber,
		DecimalRadix10,
	)
	return uuid, nil
}

func (g *Generator) GetTimestamp(epochTime time.Time) uint64 {
	return uint64(time.Now().Sub(epochTime).Round(time.Millisecond)) / uint64(time.Millisecond)
}
