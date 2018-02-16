package main

import (
	"strconv"
	"sync"
	"time"
)

const (
	WorkerIdBits                   = 10
	SequenceBits                   = 12
	InitialSequentialNumber uint64 = 1
	DecimalNumberType              = 10
)

var (
	lastTimestamp    uint64
	sequentialNumber uint64 = 1
	GeneratorLock    sync.Mutex
	ServiceEpochTime = time.Date(2018, 1, 1, 0, 0, 0, 0, time.UTC)
)

type Generator struct {
	Timestamp        uint64
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
func (g *Generator) Generate() string {
	GeneratorLock.Lock()
	defer GeneratorLock.Unlock()

	g.Timestamp = g.GetTimestamp()
	if g.Timestamp < lastTimestamp {
		logger.Fatal("system clock was rolled back")
	}
	if g.Timestamp == lastTimestamp {
		sequentialNumber++
		g.SequentialNumber = sequentialNumber
	}
	if g.Timestamp > lastTimestamp {
		sequentialNumber = InitialSequentialNumber
	}
	lastTimestamp = g.Timestamp

	uuid := strconv.FormatUint(
		(g.Timestamp<<(WorkerIdBits+SequenceBits))|g.WorkerId<<SequenceBits|g.SequentialNumber,
		DecimalNumberType,
	)
	return uuid
}

func (g *Generator) GetTimestamp() uint64 {
	return uint64(time.Now().Sub(ServiceEpochTime).Round(time.Millisecond)) / uint64(time.Millisecond)
}
