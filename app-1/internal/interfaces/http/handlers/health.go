package handlers

import (
	"database/sql"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
)

type SystemHandler struct {
	logger *zap.Logger
	DB     *sql.DB
	RDB    *redis.Client
}

func NewSystemHandler(db *sql.DB, rdb *redis.Client, logger *zap.Logger) *SystemHandler {
	return &SystemHandler{DB: db, RDB: rdb, logger: logger}
}

func (h *SystemHandler) GetSystemHealth(c *gin.Context) {
	status := http.StatusOK
	postg := "ok"
	redis := "ok"
	if err := h.DB.Ping(); err != nil {
		postg = err.Error()
		status = http.StatusServiceUnavailable
	}

	if err := h.RDB.Ping(c).Err(); err != nil {
		redis = err.Error()
		status = http.StatusServiceUnavailable
	}

	h.logger.Info("status",
		zap.String("postgres", postg),
		zap.String("redis", redis),
		zap.Int("status", status),
	)

	c.JSON(http.StatusServiceUnavailable, gin.H{
		"postgres": postg,
		"redis":    redis,
		"status":   status,
	})
}
