package main

import (
	"context"
	"fmt"
	"github.com/cloudflare/tableflip"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	// Create a TableFlip object. This is the main entrypoint to the library.
	upg, err := tableflip.New(tableflip.Options{})
	if err != nil {
		panic(err)
	}
	defer upg.Stop()

	// Create a listener that will be automatically rebind on SIGHUP.
	// This is the same as net.Listen("tcp", "localhost:8080").
	ln, err := upg.Listen("tcp", "localhost:8080")
	if err != nil {
		panic(err)
	}

	srv := newServer()

	go func() {
		fmt.Printf("Listening on %s, called\n", "localhost:8080")

		err := srv.Serve(ln)
		if err != nil {
			panic(err)
		}
	}()

	go func() {
		sig := make(chan os.Signal, 1)
		signal.Notify(sig, syscall.SIGHUP)
		for range sig {
			// 核心的 Upgrade 调用
			err := upg.Upgrade()
			if err != nil {
				log.Println("Upgrade failed:", err)
			}
		}
	}()

	// Tell the parent process that we are ready to accept connections.
	_ = upg.Ready()
	// Wait for a SIGHUP or SIGTERM.
	<-upg.Exit()

	_ = srv.Shutdown(context.TODO())
}

var version = "1.0.4"

func newServer() *http.Server {

	mux := &http.ServeMux{}
	mux.HandleFunc("/version", func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte(version + "\n"))
	})

	return &http.Server{
		Handler: mux,
	}
}
