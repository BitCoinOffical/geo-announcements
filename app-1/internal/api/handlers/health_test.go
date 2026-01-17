package handlers_test

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/BitCoinOffical/geo-announcements/app-1/config"
	"github.com/BitCoinOffical/geo-announcements/app-1/internal/api/handlers"
	mock_handlers "github.com/BitCoinOffical/geo-announcements/app-1/internal/api/handlers/mocks"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/assert/v2"
	"github.com/redis/go-redis/v9"
	"go.uber.org/mock/gomock"
	"go.uber.org/zap"
)

// набор тестов для TestGetSystemHealth
func TestGetSystemHealth_request_ok(t *testing.T) {
	gin.SetMode(gin.TestMode)
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	db := mock_handlers.NewMockDBPinger(ctrl)
	rdb := mock_handlers.NewMockRedisPinger(ctrl)

	logger := zap.NewNop()
	cfg := config.Config{
		App: config.AppConfig{
			ApiKey: testApiKey,
		},
	}

	db.EXPECT().Ping().Return(nil)
	rdb.EXPECT().Ping(gomock.Any()).Return(redis.NewStatusCmd(context.Background()))

	h := handlers.NewSystemHandler(db, rdb, logger)
	r := gin.New()
	r.GET("/api/v1/system/health", h.GetSystemHealth)

	req, _ := http.NewRequest(http.MethodGet, "/api/v1/system/health", nil)
	req.Header.Set("X-API-KEY", cfg.App.ApiKey)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)
}

func TestGetSystemHealth_server_error(t *testing.T) {
	gin.SetMode(gin.TestMode)
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	db := mock_handlers.NewMockDBPinger(ctrl)
	rdb := mock_handlers.NewMockRedisPinger(ctrl)

	logger := zap.NewNop()
	cfg := config.Config{
		App: config.AppConfig{
			ApiKey: testApiKey,
		},
	}

	db.EXPECT().Ping().Return(errors.New("postgres is down"))
	rdb.EXPECT().Ping(gomock.Any()).Return(redis.NewStatusCmd(context.Background()))

	h := handlers.NewSystemHandler(db, rdb, logger)
	r := gin.New()
	r.GET("/api/v1/system/health", h.GetSystemHealth)

	req, _ := http.NewRequest(http.MethodGet, "/api/v1/system/health", nil)
	req.Header.Set("X-API-KEY", cfg.App.ApiKey)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusServiceUnavailable, w.Code)
}
