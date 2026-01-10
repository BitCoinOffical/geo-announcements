package handlers

import (
	"context"
	"net/http"
	"time"

	"github.com/BitCoinOffical/geo-announcements/app-1/config"
	"github.com/BitCoinOffical/geo-announcements/app-1/internal/interfaces/http/dto"
	"github.com/BitCoinOffical/geo-announcements/app-1/internal/interfaces/http/queue"
	"github.com/BitCoinOffical/geo-announcements/app-1/internal/interfaces/http/response"
	"github.com/BitCoinOffical/geo-announcements/app-1/internal/interfaces/http/services"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

const (
	ClientID    = "X-Client-Id"
	timeOut     = 5
	maxAttempts = 5
)

type LocationHandler struct {
	cfg     *config.AppConfig
	logger  *zap.Logger
	service *services.LocationService
	queue   *queue.WebHookQueue
}

func NewLocationHandler(service *services.LocationService, queue *queue.WebHookQueue, logger *zap.Logger, cfg *config.AppConfig) *LocationHandler {
	return &LocationHandler{service: service, queue: queue, logger: logger, cfg: cfg}
}

func (h *LocationHandler) CreateLocationHandler(c *gin.Context) {
	userID := c.GetHeader(ClientID)

	if userID == "" {
		userID = uuid.NewString()
		c.Header(ClientID, userID)
	}

	if _, err := uuid.Parse(userID); err != nil {
		response.BadRequest(c, err, "invalid request body", h.logger)
		return
	}

	var dt dto.LocationDTO
	if err := c.ShouldBindJSON(&dt); err != nil {
		response.BadRequest(c, err, "invalid request body", h.logger)
		return
	}

	zones, err := h.service.CreateLocationService(c.Request.Context(), &dt, userID)
	if err != nil {
		h.logger.Error("create location error", zap.Error(err))
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "dunger zones": zones})

	go func() {
		ctx, canel := context.WithTimeout(context.Background(), timeOut*time.Second)
		defer canel()
		webhook := dto.WebHookDTO{
			URL:        h.cfg.WEBHOOK_URL,
			User_id:    userID,
			Payload:    dt,
			RetryCount: 0,
		}
		for i := range maxAttempts {
			if i == maxAttempts {
				h.logger.Error("exceeded the number of attempts")
				break
			}
			if err := h.queue.EnqueueWebHook(ctx, &webhook, h.cfg.QUEUE_KEY); err == nil {
				return
			}
			h.logger.Error("enqueue webhook error", zap.Error(err))
			select {
			case <-ctx.Done():
				return
			case <-time.After(time.Second):
				continue
			}
		}
	}()

}
