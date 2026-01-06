package handlers

import (
	"net/http"
	"strconv"

	"github.com/BitCoinOffical/geo-announcements/internal/interfaces/http/dto"
	"github.com/BitCoinOffical/geo-announcements/internal/interfaces/http/services"
	"github.com/gin-gonic/gin"
)

type IncidentHandler struct {
	service *services.IncidentService
}

func NewIncidentHandler(service *services.IncidentService) *IncidentHandler {
	return &IncidentHandler{service: service}
}

func (h *IncidentHandler) GetIncidentsHandler(c *gin.Context) {

	page, err := strconv.Atoi(c.Request.URL.Query().Get("page"))
	if err != nil || page < 1 {
		page = 1
	}

	limit, err := strconv.Atoi(c.Request.URL.Query().Get("limit"))
	if err != nil || limit < 1 || limit > 100 {
		limit = 10
	}

	rows, err := h.service.GetIncidentsService(c.Request.Context(), page, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, rows)
}

func (h *IncidentHandler) GetIncidentByIDHandler(c *gin.Context) {

	id, err := strconv.Atoi(c.Request.URL.Query().Get("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}
	row, err := h.service.GetIncidentByIDService(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, row)
}

func (h *IncidentHandler) GetIncidentStatHandler(c *gin.Context) {

	count, err := h.service.GetIncidentStatService(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, count)
}

func (h *IncidentHandler) CreateIncidentsHandler(c *gin.Context) {

	var dto dto.IncidentDTO
	if err := c.ShouldBindJSON(&dto); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}
	if err := h.service.CreateIncidentsService(c.Request.Context(), &dto); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}
	c.Status(http.StatusOK)
}

func (h *IncidentHandler) UpdateIncidentsByIDHandler(c *gin.Context) {

	id, err := strconv.Atoi(c.Request.URL.Query().Get("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	var dto dto.IncidentDTO
	if err := c.ShouldBindJSON(&dto); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	if err := h.service.UpdateIncidentsByIDService(c.Request.Context(), &dto, id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.Status(http.StatusOK)

}

func (h *IncidentHandler) DeleteIncidentsByIDHandler(c *gin.Context) {

	id, err := strconv.Atoi(c.Request.URL.Query().Get("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	if err := h.service.DeleteIncidentsByIDService(c.Request.Context(), id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.Status(http.StatusOK)
}
