package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"net/http"
)

var (
	addr   string
	prefix string
)

func main() {
	flag.StringVar(&addr, "addr", ":8080", "specify the address")
	flag.StringVar(&prefix, "prefix", "mux0", "specify the prefix string")
	flag.Parse()

	// default
	mux0 := http.NewServeMux()
	mux0.HandleFunc("/h", genHandleFunc(prefix))

	s0 := http.Server{
		Handler: mux0,
		Addr:    addr,
	}

	fmt.Printf("starting server on: %s\n", addr)
	fmt.Println(s0.ListenAndServe())
}

func genHandleFunc(prefix string) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		fmt.Println("recv an request")

		byt, _ := json.Marshal(map[string]string{
			"prefix": prefix,
		})

		fmt.Fprintln(w, string(byt))
	}
}
