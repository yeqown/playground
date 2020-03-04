/*
 Why this program will blocked?
*/
package main

import (
	"fmt"
	"runtime"
	"sync/atomic"
)

var (
	end = make(chan struct{})
	i   int32
)

func main() {
	// what hanppens when num less than 4 OR num bigger than 4?
	// why 4?
	runtime.GOMAXPROCS(4)

	// 4 goroutine + 1 main.main
	names := []string{"A", "B", "C", "D", "E"}
	for idx, v := range names {
		go foo(int32(idx), v)
	}

	for range names {
		<-end
	}
}

func foo(threadNum int32, threadName string) {
	for {
		v := atomic.LoadInt32((*int32)(&i))
		if v >= 30 {
			break
		}

		if v%5 == threadNum {
			fmt.Printf("%d: %s\n", i, threadName)
			atomic.AddInt32((*int32)(&i), int32(1))
		} else {
			// time.Sleep(0)
			continue
		}
	}

	end <- struct{}{}
}
