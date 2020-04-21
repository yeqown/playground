package main

import (
	"time"
)

// reference doc:
// https://draveness.me/golang/docs/part3-runtime/ch07-memory/golang-garbage-collector/#%E6%A0%87%E8%AE%B0%E6%B8%85%E9%99%A4

// open GC Flag and pprof to analyze GC, run command:
// GODEBUG=gctrace=1 go run main.go -> gc.out
//
// build an app always malloc memory and release them
func main() {
	var cnt int

	for {
		// 1KB * 100
		apply := make([]byte, 1024)
		time.Sleep(1 * time.Microsecond)
		_ = apply

		cnt++
		if cnt > 500000 {
			// 500,000 = 500K * 100KB = 50000KB = 50MB
			break
		}
	}

	println("i'm done")
}
