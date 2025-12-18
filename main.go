package main

import (
	"log"
	"net/http"

	"github.com/AryanAg08/Simplified-Splitwise/controllers"
	"github.com/AryanAg08/Simplified-Splitwise/workers/cache"
	"github.com/AryanAg08/Simplified-Splitwise/workers/db"
	"github.com/gin-gonic/gin"
)

func main() {

	db.ConnectMongo()
	cache.ConnectRedis()

	router := gin.Default()

	router.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status":  "ok",
			"message": "pong",
		})
	})

	groupController := &controllers.GroupControllers{}
	groupController.InitGroupController(router)

	if err := router.Run(":8000"); err != nil {
		log.Fatalf("Failed to start the server: %v", err)
	}
}
