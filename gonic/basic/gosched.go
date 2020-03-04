package main

import (
	"fmt"
	"sync/atomic"
	"time"
	"runtime"
)

var (
	end = make(chan struct{})
	i   int32
)

func threadPrint(threadNum int32, threadName string) {

	for {
		v := atomic.LoadInt32((*int32)(&i))
		if v >= 50 {
			break
		}

		if v%5 == threadNum {
			fmt.Printf("%d: %s\n", i, threadName)
			atomic.AddInt32((*int32)(&i), int32(1))
		} else {
			// why called time.Sleep(0) ?
			// https://www.bwangel.me/2019/04/10/go-scheduler-pitfall/
			time.Sleep(0)
			// runtime.Gosched()
			continue
		}
	}
	end <- struct{}{}
}

func main() {
	fmt.Printf("num of cpu: %d\n", runtime.NumCPU())
	names := []string{"A", "B", "C", "D", "E"}

	for idx, name := range names {
		go threadPrint(int32(idx), name)
	}

	for _ = range names {
		<-end
	}
}
