package handlers

import (
	"context"
	"log"
	"net/http"
	"os"

	"github.com/BitCoinOffical/geo-announcements/app-1/internal/interfaces/http/dto"
	"github.com/BitCoinOffical/geo-announcements/app-1/internal/interfaces/http/queue"
	"github.com/BitCoinOffical/geo-announcements/app-1/internal/interfaces/http/services"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type LocationHandler struct {
	service *services.LocationService
	queue   *queue.WebHookQueue
}

func NewLocationHandler(service *services.LocationService, queue *queue.WebHookQueue) *LocationHandler {
	return &LocationHandler{service: service, queue: queue}
}

func (h *LocationHandler) CreateLocationHandler(c *gin.Context) {
	userID := c.GetHeader("X-Client-Id")

	if userID == "" {
		userID = uuid.NewString()
		c.Header("X-Client-Id", userID)
	}

	if _, err := uuid.Parse(userID); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	var dt dto.LocationDTO
	if err := c.ShouldBindJSON(&dt); err != nil {
		log.Println(err)
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
	})

	go func() {
		ctx := context.Background()
		if err := h.service.CreateLocationService(ctx, &dt, userID); err != nil {
			log.Println("create location error:", err)
			return
		}

		webhook := dto.WebHookDTO{
			URL:        os.Getenv("WEBHOOK_URL"),
			User_id:    userID,
			Payload:    dt,
			RetryCount: 0,
		}
		if err := h.queue.EnqueueWebHook(ctx, &webhook); err != nil {
			log.Println("enqueue webhook error:", err)

		}
	}()

}
