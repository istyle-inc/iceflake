package main

import (
	"testing"
	"time"
)

var (
	workerId uint64 = 1
)

func TestNewGenerator(t *testing.T) {
	generator := NewGenerator(workerId)
	expectGenerator := Generator{
		WorkerId:         1,
		SequentialNumber: 1,
	}

	if *generator != expectGenerator {
		t.Error("Not equal generator")
	}
}

func TestGenerator_Generate(t *testing.T) {
	generator := NewGenerator(workerId)
	uuid, _ := generator.Generate()
	if uuid == "" {
		t.Error("Generate result is empty")
	}
}

func TestGenerator_Generate_RollbackTimestampError(t *testing.T) {
	generator := NewGenerator(workerId)
	var futureEpochTime = time.Date(2017, 1, 1, 0, 0, 0, 0, time.UTC)
	generator.LastTimestamp = generator.GetTimestamp(futureEpochTime)

	_, err := generator.Generate()
	if err == nil {
		t.Error("Expect rollback timestamp error")
	}
}

func BenchmarkGenerateUUID(b *testing.B) {
	for i := 0; i < b.N; i++ {
		g := NewGenerator(uint64(i))
		g.Generate()
	}
}
