package handlers

import (
	"database/sql"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
)

type SystemHandler struct {
	DB  *sql.DB
	RDB *redis.Client
}

func NewSystemHandler(db *sql.DB, rdb *redis.Client) *SystemHandler {
	return &SystemHandler{DB: db, RDB: rdb}
}

func (h *SystemHandler) GetSystemHealth(c *gin.Context) {
	status := http.StatusOK
	postg := "ok"
	if err := h.DB.Ping(); err != nil {
		postg = err.Error()
		status = http.StatusServiceUnavailable
	}

	redis := "ok"
	if err := h.RDB.Ping(c).Err(); err != nil {
		redis = err.Error()
		status = http.StatusServiceUnavailable
	}
	c.JSON(http.StatusServiceUnavailable, gin.H{
		"postgres": postg,
		"redis":    redis,
		"status":   status,
	})
}
