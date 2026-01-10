package handlers

import (
	"net/http"
	"strconv"

	"github.com/BitCoinOffical/geo-announcements/app-1/config"
	"github.com/BitCoinOffical/geo-announcements/app-1/internal/interfaces/http/dto"
	"github.com/BitCoinOffical/geo-announcements/app-1/internal/interfaces/http/response"
	"github.com/BitCoinOffical/geo-announcements/app-1/internal/interfaces/http/services"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

const (
	defaultPage     = 1
	defaultMaxLimit = 10
	defaultMinLimit = 1
)

type IncidentHandler struct {
	cfg     *config.AppConfig
	logger  *zap.Logger
	service *services.IncidentService
}

func NewIncidentHandler(service *services.IncidentService, logger *zap.Logger, cfg *config.AppConfig) *IncidentHandler {
	return &IncidentHandler{service: service, logger: logger, cfg: cfg}
}

func (h *IncidentHandler) GetIncidentsHandler(c *gin.Context) {
	page, err := strconv.Atoi(c.Request.URL.Query().Get("page"))
	if err != nil || page < defaultPage {
		h.logger.Warn("invalid page query param, fallback to default",
			zap.Int("page", page),
			zap.Int("default_page", defaultPage),
			zap.String("path", c.FullPath()),
			zap.Error(err),
		)
		page = defaultPage
	}

	limit, err := strconv.Atoi(c.Request.URL.Query().Get("limit"))
	if err != nil || limit < defaultMinLimit || limit > defaultMaxLimit {
		h.logger.Warn("invalid limit query param, fallback to default",
			zap.Int("limit", limit),
			zap.Int("default_limit", defaultMaxLimit),
			zap.String("path", c.FullPath()),
			zap.Error(err),
		)
		limit = defaultMaxLimit
	}

	rows, err := h.service.GetIncidentsService(c.Request.Context(), page, limit)
	if err != nil {
		response.InternalServerError(c, err, "failed to get incidents", h.logger)
		return
	}

	c.JSON(http.StatusOK, rows)
}

func (h *IncidentHandler) GetIncidentByIDHandler(c *gin.Context) {
	id, err := strconv.Atoi(c.Request.URL.Query().Get("id"))
	if err != nil {
		response.BadRequest(c, err, "invalid request body", h.logger)
		return
	}
	row, err := h.service.GetIncidentByIDService(c.Request.Context(), id)
	if err != nil {
		response.InternalServerError(c, err, "failed to get incident by id", h.logger)
		return
	}
	c.JSON(http.StatusOK, row)
}

func (h *IncidentHandler) GetIncidentStatHandler(c *gin.Context) {

	users, err := h.service.GetIncidentStatService(c.Request.Context(), h.cfg)
	if err != nil {
		response.InternalServerError(c, err, "failed to get incidents stat", h.logger)
		return
	}
	c.JSON(http.StatusOK, users)
}

func (h *IncidentHandler) CreateIncidentsHandler(c *gin.Context) {

	var dto dto.IncidentDTO
	if err := c.ShouldBindJSON(&dto); err != nil {
		response.BadRequest(c, err, "invalid request body", h.logger)
		return
	}
	if err := h.service.CreateIncidentsService(c.Request.Context(), &dto); err != nil {
		response.InternalServerError(c, err, "failed to create incident", h.logger)
		return
	}
	c.Status(http.StatusOK)
}

func (h *IncidentHandler) UpdateIncidentsByIDHandler(c *gin.Context) {

	id, err := strconv.Atoi(c.Request.URL.Query().Get("id"))
	if err != nil {
		response.BadRequest(c, err, "invalid request body", h.logger)
		return
	}

	var dto dto.IncidentDTO
	if err := c.ShouldBindJSON(&dto); err != nil {
		response.BadRequest(c, err, "invalid request body", h.logger)
		return
	}

	if err := h.service.UpdateIncidentsByIDService(c.Request.Context(), &dto, id); err != nil {
		response.InternalServerError(c, err, "failed to update incident by id", h.logger)
		return
	}

	c.Status(http.StatusOK)

}

func (h *IncidentHandler) DeleteIncidentsByIDHandler(c *gin.Context) {

	id, err := strconv.Atoi(c.Request.URL.Query().Get("id"))
	if err != nil {
		response.BadRequest(c, err, "invalid request body", h.logger)
		return
	}

	if err := h.service.DeleteIncidentsByIDService(c.Request.Context(), id); err != nil {
		response.InternalServerError(c, err, "failed to delete incident by id", h.logger)
		return
	}

	c.Status(http.StatusOK)
}
