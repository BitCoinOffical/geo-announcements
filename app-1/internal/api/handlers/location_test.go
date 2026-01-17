package handlers_test

import (
	"bytes"
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/BitCoinOffical/geo-announcements/app-1/config"
	"github.com/BitCoinOffical/geo-announcements/app-1/internal/api/handlers"
	mock_handlers "github.com/BitCoinOffical/geo-announcements/app-1/internal/api/handlers/mocks"

	"github.com/BitCoinOffical/geo-announcements/app-1/internal/domain/rules"
	"github.com/BitCoinOffical/geo-announcements/app-1/internal/interfaces/http/dto"
	"github.com/BitCoinOffical/geo-announcements/app-1/internal/interfaces/http/models"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/assert/v2"
	"github.com/go-playground/validator/v10"
	"go.uber.org/mock/gomock"
	"go.uber.org/zap"
)

func TestCreateLocationHandler_request_ok(t *testing.T) {
	gin.SetMode(gin.TestMode)
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	locSvc := mock_handlers.NewMocklocationService(ctrl)
	retry := mock_handlers.NewMockwebhookRetry(ctrl)
	logger := zap.NewNop()
	cfg := config.Config{
		App: config.AppConfig{
			ApiKey: testApiKey,
		},
	}
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		if err := v.RegisterValidation("lat", rules.ValidateLat); err != nil {
			log.Printf("error lat validate: %v", err)
		}
	}

	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		if err := v.RegisterValidation("lon", rules.ValidateLon); err != nil {
			log.Printf("error lat validate: %v", err)
		}
	}

	h := handlers.NewLocationHandler(locSvc, logger, &cfg.App, retry)
	locSvc.EXPECT().CreateLocation(gomock.Any(), gomock.AssignableToTypeOf(&dto.LocationDTO{}), gomock.Any()).Return(nil)
	locSvc.EXPECT().GetDangerZones(gomock.Any(), gomock.AssignableToTypeOf(&dto.LocationDTO{})).Return([]models.DangerousZones{}, nil)
	body, _ := json.Marshal(map[string]any{
		"lat": 55.755,
		"lon": 37.617,
	})

	r := gin.New()
	r.POST("/api/v1/location/check", h.CreateLocationHandler)

	req, _ := http.NewRequest(http.MethodPost, "/api/v1/location/check", bytes.NewReader(body))

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)
}
func TestCreateLocationHandler_bad_request(t *testing.T) {
	gin.SetMode(gin.TestMode)
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	locSvc := mock_handlers.NewMocklocationService(ctrl)
	retry := mock_handlers.NewMockwebhookRetry(ctrl)
	logger := zap.NewNop()
	cfg := config.Config{
		App: config.AppConfig{
			ApiKey: testApiKey,
		},
	}
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		if err := v.RegisterValidation("lat", rules.ValidateLat); err != nil {
			log.Printf("error lat validate: %v", err)
		}
	}

	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		if err := v.RegisterValidation("lon", rules.ValidateLon); err != nil {
			log.Printf("error lat validate: %v", err)
		}
	}

	h := handlers.NewLocationHandler(locSvc, logger, &cfg.App, retry)
	body, _ := json.Marshal(map[string]any{ //данные выходят за диапазон значений для долготы (От 0° до 180°)
		"lat": 505.755,
		"lon": 307.617,
	})

	r := gin.New()
	r.POST("/api/v1/location/check", h.CreateLocationHandler)

	req, _ := http.NewRequest(http.MethodPost, "/api/v1/location/check", bytes.NewReader(body))

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusBadRequest, w.Code)
}
func TestCreateLocationHandler_server_error(t *testing.T) {
	gin.SetMode(gin.TestMode)
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	locSvc := mock_handlers.NewMocklocationService(ctrl)
	retry := mock_handlers.NewMockwebhookRetry(ctrl)
	logger := zap.NewNop()
	cfg := config.Config{
		App: config.AppConfig{
			ApiKey: testApiKey,
		},
	}
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		if err := v.RegisterValidation("lat", rules.ValidateLat); err != nil {
			log.Printf("error lat validate: %v", err)
		}
	}

	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		if err := v.RegisterValidation("lon", rules.ValidateLon); err != nil {
			log.Printf("error lat validate: %v", err)
		}
	}

	h := handlers.NewLocationHandler(locSvc, logger, &cfg.App, retry)
	dbErr := errors.New("db is down") //имитируем падения бд
	locSvc.EXPECT().CreateLocation(gomock.Any(), gomock.AssignableToTypeOf(&dto.LocationDTO{}), gomock.Any()).Return(dbErr)
	body, _ := json.Marshal(map[string]any{
		"lat": 55.755,
		"lon": 37.617,
	})

	r := gin.New()
	r.POST("/api/v1/location/check", h.CreateLocationHandler)

	req, _ := http.NewRequest(http.MethodPost, "/api/v1/location/check", bytes.NewReader(body))

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)
}
