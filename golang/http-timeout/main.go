// 编写一个 http 服务，设置 readTimeout 和 writeTimeout 均为 1s
// 其中有一个 handler（/timeout） 的执行时间超过 5s，通过 go http client 来测试服务的超时情况

package main

import (
	"fmt"

	"net/http"
	"time"
)

// curl -i http://localhost:8080/timeout

func main() {
	http.HandleFunc("/timeout", func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("start", time.Now().Format("2006-01-02 15:04:05"))
		time.Sleep(5 * time.Second)
		fmt.Fprintln(w, "timeout")
		fmt.Println("end, ", time.Now().Format("2006-01-02 15:04:05"))
	})

	server := &http.Server{
		Addr:         ":8080",
		ReadTimeout:  1 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	if err := server.ListenAndServe(); err != nil {
		fmt.Println(err)
	}
}
