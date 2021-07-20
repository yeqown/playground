package basic_test

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"runtime"
	"sync"
	"testing"
)

func Test_Client(t *testing.T) {
	num := 20

	wg := sync.WaitGroup{}
	wg.Add(num)
	for index := 0; index < num; index++ {
		go func() {
			defer wg.Done()
			resp, _ := http.Get("https://www.baidu.com")
			_, _ = ioutil.ReadAll(resp.Body)
		}()
	}

	println("break point for debug")
	wg.Wait()

	fmt.Printf("此时goroutine个数= %d\n", runtime.NumGoroutine())
	println("break point for debug")
}

func Test_Server(t *testing.T) {

}
