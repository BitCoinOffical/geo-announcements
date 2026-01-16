package handlers

import (
	"context"
	"net/http"
	"time"

	"github.com/BitCoinOffical/geo-announcements/app-1/config"
	"github.com/BitCoinOffical/geo-announcements/app-1/internal/api/response"
	"github.com/BitCoinOffical/geo-announcements/app-1/internal/interfaces/http/dto"
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
	retry   webhookRetry
}

func NewLocationHandler(service locationService, logger *zap.Logger, cfg *config.AppConfig, retry webhookRetry) *LocationHandler {
	return &LocationHandler{service: service, logger: logger, cfg: cfg, retry: retry}
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

	err := h.service.CreateLocation(c.Request.Context(), &dt, userID)
	if err != nil {
		h.logger.Error("create location error", zap.Error(err))
		return
	}
	zones, err := h.service.GetDangerZones(c.Request.Context(), &dt)
	if err != nil {
		h.logger.Error("get danger zones error", zap.Error(err))
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
		h.retry.Retry(ctx, &webhook)

	}()

}
