package main

import (
	"log"

	"github.com/BitCoinOffical/geo-announcements/internal/interfaces/http/handlers"
	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()

	r.POST("/webhook", handlers.WebhookHandler)

	if err := r.Run(":9090"); err != nil {
		log.Fatal(err)
	}
}
