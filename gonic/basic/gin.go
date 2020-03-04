package main

import (
	"github.com/gin-gonic/gin"
)

//post form-data a
func main() {
	r := gin.Default()
	r.POST("/test", TestShoudbind)
	r.Run(":8080")
}

func TestShoudbind(c *gin.Context) {
	a := A{}
	c.ShouldBind(&a)
	c.JSON(200, a)
}

type A struct {
	B string `form:"b" json:"b"`
	L A      `json:"l"`
}
