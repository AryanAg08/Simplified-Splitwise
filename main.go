package main

import (
	"log"

	"github.com/AryanAg08/Simplified-Splitwise/workers/cache"
	"github.com/AryanAg08/Simplified-Splitwise/workers/db"
	"github.com/gin-gonic/gin"
)

func main() {

	db.ConnectMongo()
	cache.ConnectRedis()

	router := gin.Default()

	if err := router.Run(":8000"); err != nil {
		log.Fatalf("Failed to start the server: %v", err)
	}
}
