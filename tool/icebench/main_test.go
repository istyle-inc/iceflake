package main

import (
	"os"
	"testing"
)

func BenchmarkConnectIceFlake(b *testing.B) {
	// get path from env. then check stat
	path := os.Getenv("ICEFLAKE_SOCKETFILE_PATH")
	_, err := os.Stat(path)
	// if sock file doesn't exist or path is empty, skip this bench
	if len(path) == 0 || err != nil {
		b.Error("could not find whether iceflake is running or not, failed")
	}
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			_, err := connectIceFlake(path)
			if err != nil {
				b.Error("got failed: ", err)
			}
		}
	})
}
