package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type SystemHandler struct {
	logger *zap.Logger
	DB     DBPinger
	RDB    RedisPinger
}

func NewSystemHandler(db DBPinger, rdb RedisPinger, logger *zap.Logger) *SystemHandler {
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

	c.JSON(status, gin.H{
		"postgres": postg,
		"redis":    redis,
		"status":   status,
	})
}
