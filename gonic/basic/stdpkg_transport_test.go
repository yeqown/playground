package basic_test

import (
	"fmt"
	"net/http"
	"os"
	"testing"
	"time"
)

var httpClient *http.Client

const (
	MaxIdleConnections int = 20
	RequestTimeout     int = 5
)

func init() {
	httpClient = createHTTPClient()
}

// createHTTPClient for connection re-use
func createHTTPClient() *http.Client {
	client := &http.Client{
		Transport: &http.Transport{
			MaxIdleConnsPerHost: MaxIdleConnections,
			MaxConnsPerHost:     1,
		},
		Timeout: time.Duration(RequestTimeout) * time.Second,
	}

	return client
}

// HOW TO WATCH TCP CONNECTIONS of this application
func Test_transport(t *testing.T) {
	fmt.Println("PID=", os.Getpid())

	// wait for observer ready
	time.Sleep(10 * time.Second)

	for i := 0; i < 5; i++ {
		go func(idx int) {
			defer fmt.Printf("%d finished\n", idx)
			resp, err := httpClient.Get("http://www.baidu.com")
			if err != nil {
				panic(err)
			}
			defer resp.Body.Close()
			time.Sleep(1 * time.Second)
		}(i)
	}

	select {}
}
