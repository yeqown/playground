package basic_test

import (
	"fmt"
	"runtime"
	"sync/atomic"
	"testing"
	"time"
)

var (
	end = make(chan struct{})
	i   int32
)

func Test_sched(t *testing.T) {
	fmt.Printf("num of cpu: %d\n", runtime.NumCPU())
	names := []string{"A", "B", "C", "D", "E"}

	threadPrint := func(threadNum int32, threadName string) {
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

	for idx, name := range names {
		go threadPrint(int32(idx), name)
	}

	for _ = range names {
		<-end
	}
}

// Why this program will blocked?
func Test_scheduler(t *testing.T) {
	// what hanppens when num less than 4 OR num bigger than 4?
	// why 4?
	runtime.GOMAXPROCS(4)

	foo := func(threadNum int32, threadName string) {
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

	// 4 goroutine + 1 main.main
	names := []string{"A", "B", "C", "D", "E"}
	for idx, v := range names {
		go foo(int32(idx), v)
	}

	for range names {
		<-end
	}
}
