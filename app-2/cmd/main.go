package main

import (
	"log"

	"github.com/BitCoinOffical/geo-announcements/app-2/internal/interfaces/http/handlers"
	"github.com/gin-gonic/gin"
)

// сервер заглушка куда отправляются webhooks
func main() {
	r := gin.Default()

	r.POST("/webhook", handlers.WebhookHandler)

	if err := r.Run(":9090"); err != nil {
		log.Fatal(err)
	}
}
