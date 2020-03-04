package main

import (
	"log"
	"net/http"
)

// std net/http

func main() {
	http.HandleFunc("/http1x/std", func(w http.ResponseWriter, req *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("hello http1x standard"))
	})

	log.Println("/http1x/std is up on :8080")
	if err := http.ListenAndServeTLS(":8080", "../https.pem", "../https.key", nil); err != nil {
		log.Fatal(err)
	}
}
