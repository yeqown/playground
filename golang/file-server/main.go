package main

import (
	"log"
	"net/http"
)

func main() {

	http.HandleFunc("/upload", UploadFileHandler)
	http.HandleFunc("/download", DownloadFileHandler)
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "index.html")
	})

	log.Println("Listening on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func renderError(w http.ResponseWriter, err error) {
	log.Println(err.Error())
	http.Error(w, err.Error(), http.StatusInternalServerError)
}

func renderOK(w http.ResponseWriter, status int, message string) {
	w.WriteHeader(status)
	_, _ = w.Write([]byte(message))
}
