package handlers

import (
	"context"
	"net/http"
	"time"

	"github.com/BitCoinOffical/geo-announcements/app-1/config"
	"github.com/BitCoinOffical/geo-announcements/app-1/internal/api/response"
	"github.com/BitCoinOffical/geo-announcements/app-1/internal/interfaces/http/dto"
	"github.com/BitCoinOffical/geo-announcements/app-1/internal/interfaces/http/queue"
	"github.com/BitCoinOffical/geo-announcements/app-1/internal/retry"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

const (
	ClientID = "X-Client-Id"
	timeOut  = 5
)

type LocationHandler struct {
	cfg     *config.AppConfig
	logger  *zap.Logger
	service locationService
	queue   *queue.WebHookQueue
}

func NewLocationHandler(service locationService, queue *queue.WebHookQueue, logger *zap.Logger, cfg *config.AppConfig) *LocationHandler {
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

	zones, err := h.service.CreateLocation(c.Request.Context(), &dt, userID)
	if err != nil {
		h.logger.Error("create location error", zap.Error(err))
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "dunger zones": zones})

	go func() {
		ctx, canel := context.WithTimeout(context.Background(), timeOut*time.Second)
		defer canel()
		webhook := dto.WebHookDTO{
			URL:        h.cfg.WebhookUrl,
			User_id:    userID,
			Payload:    dt,
			RetryCount: 0,
		}
		retry.Retry(ctx, &webhook, h.logger, h.queue, h.cfg)

	}()

}
