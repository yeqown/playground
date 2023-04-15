package main

import (
	"log"

	"swagger-demo/controllers"

	"github.com/gin-gonic/gin"
)

func main() {
	loadConf()

	eng := gin.New()
	eng.GET("/base/get", controllers.Get)
	eng.POST("/base/post", controllers.Post)

	log.Fatal(eng.Run(":8088"))
}
