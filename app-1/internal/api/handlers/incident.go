package handlers

import (
	"net/http"
	"strconv"

	"github.com/BitCoinOffical/geo-announcements/app-1/config"
	"github.com/BitCoinOffical/geo-announcements/app-1/internal/api/response"
	"github.com/BitCoinOffical/geo-announcements/app-1/internal/interfaces/http/dto"
	"github.com/BitCoinOffical/geo-announcements/app-1/internal/interfaces/http/services"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
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
	if err != nil {
		response.BadRequest(c, err, "invalid request body", h.logger)
		return
	}

	limit, err := strconv.Atoi(c.Request.URL.Query().Get("limit"))
	if err != nil {
		response.BadRequest(c, err, "invalid request body", h.logger)
		return
	}

	rows, err := h.service.GetIncidents(c.Request.Context(), page, limit)
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
	row, err := h.service.GetIncidentByID(c.Request.Context(), id)
	if err != nil {
		response.InternalServerError(c, err, "failed to get incident by id", h.logger)
		return
	}
	c.JSON(http.StatusOK, row)
}

func (h *IncidentHandler) GetIncidentStatHandler(c *gin.Context) {

	users, err := h.service.GetIncidentStat(c.Request.Context(), h.cfg)
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
	if err := h.service.CreateIncidents(c.Request.Context(), &dto); err != nil {
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

	if err := h.service.UpdateIncidentsByID(c.Request.Context(), &dto, id); err != nil {
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

	if err := h.service.DeleteIncidentsByID(c.Request.Context(), id); err != nil {
		response.InternalServerError(c, err, "failed to delete incident by id", h.logger)
		return
	}

	c.Status(http.StatusOK)
}
