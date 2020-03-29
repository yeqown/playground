package main

import (
	"log"

	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()
	r.GET("/test", TestShoudbind)
	r.Run(":8080")
}

func TestShoudbind(c *gin.Context) {
	header := c.GetHeader("X-Not")
	log.Println("header=", header)

	// a := A{}
	// c.ShouldBind(&a)
	c.JSON(200, gin.H{"messgae": "ok"})
}

// type A struct {
// 	B string `form:"b" json:"b"`
// 	L A      `json:"l"`
// }
