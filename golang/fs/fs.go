package main

import (
	"log"
	"net/http"

	"github.com/gin-contrib/static"
	"github.com/gin-gonic/gin"
)

func main() {
	engi := gin.New()

	engi.Use(authorize())
	engi.Use(static.Serve("/", static.LocalFile(".", true)))

	if err := engi.Run(":8000"); err != nil {
		log.Fatal(err)
	}
}

type authForm struct {
	AccessToken string `form:"access_token" binding:"required"`
}

var (
	tokenMap = map[string]bool{
		"O)*Y(BQJBOJX": true,
		"*)HBJH)(BSA":  false,
	}
)

func authorize() gin.HandlerFunc {
	return func(c *gin.Context) {
		var f = new(authForm)
		if err := c.ShouldBind(f); err != nil {
			log.Printf("Err=%v", err)
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": "参数非法"})
			return
		}

		v, ok := tokenMap[f.AccessToken]
		if !ok || !v {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": "Token过期"})
			return
		}

		c.Next()
	}
}
