package main

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

// gin net/http

func main() {
	r := gin.Default()

	r.GET("/http1x/gin", func(c *gin.Context) {
		c.String(http.StatusOK, "codegen http1x gin framework")
	})

	log.Println("/http1x/gin is up on :8080")
	if err := r.Run(":8080"); err != nil {
		log.Fatal(err)
	}
}
