package main

import (
	"crypto/tls"
	"log"
	"net"
	"net/http"
	"time"
)

// std http2

func main() {
	if _, err := tls.LoadX509KeyPair("../http2.pem", "../http2.key"); err != nil {
		log.Fatalf("ERR_CERT_INVALID: %v", err)
	}

	mux := http.NewServeMux()
	mux.HandleFunc("/http2/std", func(w http.ResponseWriter, req *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("codegen http2 standard"))
	})

	http2Srv := &http.Server{
		ReadTimeout:  5000 * time.Millisecond,
		WriteTimeout: 5000 * time.Millisecond,
		Handler:      mux,
	}

	var (
		listener net.Listener
		err      error
	)

	if listener, err = net.Listen("tcp", ":8080"); err != nil {
		log.Fatal(err)
	}

	// 除了使用ServeTLS来启动支持HTTPS/2特性的服务端之外，
	// 还可以通过http2.ConfigureServer来为http.Server启动HTTPS/2特性并直接使用Serve来启动服务。
	log.Println("/http2/std is up on :8080")
	if err := http2Srv.ServeTLS(listener, "../http2.pem", "../http2.key"); err != nil {
		log.Fatal(err)
	}
}
