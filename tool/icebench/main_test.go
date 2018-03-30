package main

import "testing"
import "os"

func BenchmarkConnectIceFlake(b *testing.B) {
	// get path from env. then check stat
	path := os.Getenv("ICEFLAKE_SOCKETFILE_PATH")
	_, err := os.Stat(path)
	// if sock file doesn't exist or path is empty, skip this bench
	if len(path) == 0 || err != nil {
		b.Skip("could not find whether iceflake is running or not, skipped")
	}

	// Start bench
	for i := 0; i < b.N; i++ {
		_, err := connectIceFlake(path)
		if err != nil {
			b.Error(err)
		}
	}
}
