package main

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()

	r.Any("/http2/gin", func(c *gin.Context) {
		c.String(http.StatusOK, "http2s gin framework")
	})

	if err := r.RunTLS(":8080", "../http2.pem", "../http2.key"); err != nil {
		log.Fatal(err)
	}
}
