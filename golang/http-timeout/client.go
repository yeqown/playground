package main

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"time"
)

func main() {
	client := newClient()
	resp, err := client.Get("http://localhost:8080/timeout")
	if err != nil {
		fmt.Println(err)
		return
	}
	defer resp.Body.Close()

	fmt.Println(resp.Status)
}

func newClient() *http.Client {
	httpTransport := &http.Transport{
		DialContext: func(ctx context.Context, network, addr string) (net.Conn, error) {
			c, err := net.DialTimeout(network, addr, 3*time.Second)
			if err != nil {
				return nil, err
			}
			return c, nil
		},
		DisableKeepAlives: false,

		MaxIdleConns:        50,
		MaxIdleConnsPerHost: 10,
		MaxConnsPerHost:     100,
		// keep same with http.DefaultTransport
		IdleConnTimeout:       90 * time.Second,
		TLSHandshakeTimeout:   10 * time.Second,
		ExpectContinueTimeout: 1 * time.Second,
	}
	return &http.Client{
		Timeout:   60 * time.Second,
		Transport: httpTransport,
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse // not support redirect
		},
	}
}
