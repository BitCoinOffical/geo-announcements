package handlers

import (
	"log"

	"github.com/BitCoinOffical/geo-announcements/app-2/internal/interfaces/http/dto"
	"github.com/gin-gonic/gin"
)

func WebhookHandler(c *gin.Context) {
	var dto dto.WebhookPayload

	if err := c.ShouldBindJSON(&dto); err != nil {
		log.Println(err)
		return
	}

	log.Println(dto)
}
