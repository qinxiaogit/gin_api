package main

import (
	"github.com/gin-gonic/gin"
	"go_gin_api/router"
)

func main() {
	r := gin.Default()
	//r.Use(middleware.LoggerToFile())
	//r.Use(middleware.LoggerToMongo())
	//r.Use(middleware.LoggerToES())
	//r.Use(middleware.LoggerToMQ())
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})
	router.InitRouter(r)

	//go middleware.MqConsumer()
	r.Run()
}
