package handlers_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/BitCoinOffical/geo-announcements/app-1/config"
	"github.com/BitCoinOffical/geo-announcements/app-1/internal/api/handlers"
	"github.com/BitCoinOffical/geo-announcements/app-1/internal/api/handlers/mocks"
	"github.com/BitCoinOffical/geo-announcements/app-1/internal/api/middleware"

	"github.com/gin-gonic/gin"
	"go.uber.org/mock/gomock"
	"go.uber.org/zap"
)

func TestHandler_request_ok(t *testing.T) {
	gin.SetMode(gin.TestMode)

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	incSvc := mocks.NewMockincidentService(ctrl)

	logger := zap.NewNop()
	cfg, _ := config.NewLoadConfig()

	h := handlers.NewIncidentHandler(incSvc, logger, &cfg.App)

	incSvc.EXPECT().CreateIncidents(gomock.Any(), gomock.Any()).Return(nil)

	router := gin.New()
	router.POST("/api/v1/incidents", middleware.CheckApiKey(cfg.App.ApiKey), h.CreateIncidentsHandler)

	reqBody, _ := json.Marshal(map[string]any{
		"lat": 55.755,
		"lon": 37.617,
	})
	req, _ := http.NewRequest(http.MethodPost, "/api/v1/incidents", bytes.NewReader(reqBody))

	req.Header.Set("X-API-KEY", cfg.App.ApiKey)
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("unexpected status: %d, body=%s", w.Code, w.Body.String())
	}
}
