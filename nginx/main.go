package main

import (
	"encoding/json"
	"fmt"
	"net/http"
)

func main() {
	// server 1
	mux1 := http.NewServeMux()
	mux1.HandleFunc("/h", genHandleFunc("I'm mux1"))

	s1 := http.Server{
		Handler: mux1,
		Addr:    ":8081",
	}

	go func() {
		fmt.Println(s1.ListenAndServe())
	}()

	// server 2
	mux2 := http.NewServeMux()
	mux2.HandleFunc("/h", genHandleFunc("I'm mux2"))

	s2 := http.Server{
		Handler: mux2,
		Addr:    ":8082",
	}

	go func() {
		fmt.Println(s2.ListenAndServe())
	}()

	// default
	mux0 := http.NewServeMux()
	mux0.HandleFunc("/h", genHandleFunc("I'm mux0"))

	s0 := http.Server{
		Handler: mux0,
		Addr:    ":8080",
	}

	fmt.Println(s0.ListenAndServe())
}

func genHandleFunc(prefix string) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		byt, _ := json.Marshal(map[string]string{
			"prefix": prefix,
		})

		fmt.Fprintln(w, string(byt))
	}
}
