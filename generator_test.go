package main

import (
	"testing"
)

func BenchmarkGenerateUUID(b *testing.B) {
	for i := 0; i < b.N; i++ {
		g := NewGenerator(uint64(i))
		g.Generate()
	}
}
