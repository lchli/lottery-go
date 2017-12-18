// lottery-server project main.go
package main

import (
	"lottery-server/api"

	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.Default()
	router.Use(func(c *gin.Context) {

		//c.AbortWithStatusJSON(500, "Forbidden")

	})
	router.POST("/register", api.Register)
	router.GET("/login", api.Login)
	router.GET("/queryTopic", api.QueryTopic)
	router.Run(":8080")
}
