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

	"github.com/BitCoinOffical/geo-announcements/app-1/internal/api/middleware"
	"github.com/BitCoinOffical/geo-announcements/app-1/internal/domain/rules"
	"github.com/BitCoinOffical/geo-announcements/app-1/internal/interfaces/http/dto"
	"github.com/BitCoinOffical/geo-announcements/app-1/internal/interfaces/http/models"
	"github.com/go-playground/assert/v2"
	"github.com/go-playground/validator/v10"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"go.uber.org/mock/gomock"
	"go.uber.org/zap"
)

const (
	testApiKey                 = "SreUeZhc6VEaSl7Vlg9kB3dbQ9xUr0EB"
	testStatsTimeWindowMinutes = 5
)

// набор тестов для NewIncidentHandler
func TestGetIncidentsHandler_request_ok(t *testing.T) {
	gin.SetMode(gin.TestMode)
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	incSvc := mock_handlers.NewMockincidentService(ctrl)
	logger := zap.NewNop()
	cfg := config.Config{
		App: config.AppConfig{
			ApiKey: testApiKey,
		},
	}

	h := handlers.NewIncidentHandler(incSvc, logger, &cfg.App)
	incSvc.EXPECT().GetIncidents(gomock.Any(), 1, 10).Return([]models.Incident{}, nil)

	r := gin.New()
	r.GET("/api/v1/incidents/", middleware.CheckApiKey(cfg.App.ApiKey), h.GetIncidentsHandler)

	req, _ := http.NewRequest(http.MethodGet, "/api/v1/incidents/?page=1&limit=10", nil)
	req.Header.Set("X-API-KEY", cfg.App.ApiKey)

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)
}
func TestGetIncidentsHandler_bad_request(t *testing.T) {
	gin.SetMode(gin.TestMode)
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	incSvc := mock_handlers.NewMockincidentService(ctrl)
	logger := zap.NewNop()
	cfg := config.Config{
		App: config.AppConfig{
			ApiKey: testApiKey,
		},
	}

	h := handlers.NewIncidentHandler(incSvc, logger, &cfg.App)

	r := gin.New()
	r.GET("/api/v1/incidents/", middleware.CheckApiKey(cfg.App.ApiKey), h.GetIncidentsHandler)

	req, _ := http.NewRequest(http.MethodGet, "/api/v1/incidents/?page=w&limit=10", nil) //передаем не правильный ?page=w&limit=10
	req.Header.Set("X-API-KEY", cfg.App.ApiKey)

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusBadRequest, w.Code)
}
func TestGetIncidentsHandler_server_error(t *testing.T) {
	gin.SetMode(gin.TestMode)
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	incSvc := mock_handlers.NewMockincidentService(ctrl)
	logger := zap.NewNop()
	cfg := config.Config{
		App: config.AppConfig{
			ApiKey: testApiKey,
		},
	}

	dbErr := errors.New("db is down") //имитируем падения бд

	h := handlers.NewIncidentHandler(incSvc, logger, &cfg.App)
	incSvc.EXPECT().GetIncidents(gomock.Any(), 1, 10).Return([]models.Incident{}, dbErr)
	r := gin.New()
	r.GET("/api/v1/incidents/", middleware.CheckApiKey(cfg.App.ApiKey), h.GetIncidentsHandler)

	req, _ := http.NewRequest(http.MethodGet, "/api/v1/incidents/?page=1&limit=10", nil) //передаем не валидный ?page=w&limit=10
	req.Header.Set("X-API-KEY", cfg.App.ApiKey)

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusInternalServerError, w.Code)
}

// набор тестов для GetIncidentByIDHandler
func TestGetIncidentByIDHandler_request_ok(t *testing.T) {
	gin.SetMode(gin.TestMode)
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	incSvc := mock_handlers.NewMockincidentService(ctrl)
	logger := zap.NewNop()
	cfg := config.Config{
		App: config.AppConfig{
			ApiKey: testApiKey,
		},
	}

	h := handlers.NewIncidentHandler(incSvc, logger, &cfg.App)
	incSvc.EXPECT().GetIncidentByID(gomock.Any(), 1).Return(&models.Incident{}, nil)

	r := gin.New()
	r.GET("/api/v1/incidents", middleware.CheckApiKey(cfg.App.ApiKey), h.GetIncidentByIDHandler)

	req, _ := http.NewRequest(http.MethodGet, "/api/v1/incidents?id=1", nil)
	req.Header.Set("X-API-KEY", cfg.App.ApiKey)

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)
}
func TestGetIncidentByIDHandler_bad_request(t *testing.T) {
	gin.SetMode(gin.TestMode)
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	incSvc := mock_handlers.NewMockincidentService(ctrl)
	logger := zap.NewNop()
	cfg := config.Config{
		App: config.AppConfig{
			ApiKey: testApiKey,
		},
	}

	h := handlers.NewIncidentHandler(incSvc, logger, &cfg.App)

	r := gin.New()
	r.GET("/api/v1/incidents", middleware.CheckApiKey(cfg.App.ApiKey), h.GetIncidentByIDHandler)

	req, _ := http.NewRequest(http.MethodGet, "/api/v1/incidents?id=w", nil) //передаем не валидный id
	req.Header.Set("X-API-KEY", cfg.App.ApiKey)

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusBadRequest, w.Code)
}
func TestGetIncidentByIDHandler_server_error(t *testing.T) {
	gin.SetMode(gin.TestMode)
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	incSvc := mock_handlers.NewMockincidentService(ctrl)
	logger := zap.NewNop()
	cfg := config.Config{
		App: config.AppConfig{
			ApiKey: testApiKey,
		},
	}

	h := handlers.NewIncidentHandler(incSvc, logger, &cfg.App)
	dbErr := errors.New("db is down") //имитируем падения бд
	incSvc.EXPECT().GetIncidentByID(gomock.Any(), 1).Return(&models.Incident{}, dbErr)

	r := gin.New()
	r.GET("/api/v1/incidents", middleware.CheckApiKey(cfg.App.ApiKey), h.GetIncidentByIDHandler)

	req, _ := http.NewRequest(http.MethodGet, "/api/v1/incidents?id=1", nil)
	req.Header.Set("X-API-KEY", cfg.App.ApiKey)

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusInternalServerError, w.Code)
}

// набор тестов для GetIncidentStatHandler
func TestGetIncidentStatHandler_request_ok(t *testing.T) {
	gin.SetMode(gin.TestMode)
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	incSvc := mock_handlers.NewMockincidentService(ctrl)
	logger := zap.NewNop()
	cfg := config.Config{
		App: config.AppConfig{
			ApiKey:                 testApiKey,
			StatsTimeWindowMinutes: testStatsTimeWindowMinutes,
		},
	}

	h := handlers.NewIncidentHandler(incSvc, logger, &cfg.App)
	incSvc.EXPECT().GetIncidentStat(gomock.Any(), testStatsTimeWindowMinutes).Return(&models.UsersInDangerousZones{}, nil)

	r := gin.New()
	r.GET("/api/v1/incidents/stats", middleware.CheckApiKey(cfg.App.ApiKey), h.GetIncidentStatHandler)

	req, _ := http.NewRequest(http.MethodGet, "/api/v1/incidents/stats", nil)
	req.Header.Set("X-API-KEY", cfg.App.ApiKey)

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)
}
func TestGetIncidentStatHandler_server_error(t *testing.T) {
	gin.SetMode(gin.TestMode)
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	incSvc := mock_handlers.NewMockincidentService(ctrl)
	logger := zap.NewNop()
	cfg := config.Config{
		App: config.AppConfig{
			ApiKey:                 testApiKey,
			StatsTimeWindowMinutes: testStatsTimeWindowMinutes,
		},
	}

	h := handlers.NewIncidentHandler(incSvc, logger, &cfg.App)
	dbErr := errors.New("db is down") //имитируем падения бд
	incSvc.EXPECT().GetIncidentStat(gomock.Any(), testStatsTimeWindowMinutes).Return(&models.UsersInDangerousZones{}, dbErr)

	r := gin.New()
	r.GET("/api/v1/incidents/stats", middleware.CheckApiKey(cfg.App.ApiKey), h.GetIncidentStatHandler)

	req, _ := http.NewRequest(http.MethodGet, "/api/v1/incidents/stats", nil)
	req.Header.Set("X-API-KEY", cfg.App.ApiKey)

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusInternalServerError, w.Code)
}

// набор тестов для CreateIncidents
func TestCreateIncidents_request_ok(t *testing.T) {
	gin.SetMode(gin.TestMode)
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	incSvc := mock_handlers.NewMockincidentService(ctrl)

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

	h := handlers.NewIncidentHandler(incSvc, logger, &cfg.App)

	incSvc.EXPECT().CreateIncidents(gomock.Any(), gomock.AssignableToTypeOf(&dto.IncidentDTO{})).Return(nil)
	incSvc.EXPECT().UpdateZones(gomock.Any(), gomock.AssignableToTypeOf(&dto.IncidentDTO{})).Return(nil)

	r := gin.New()
	r.POST("/api/v1/incidents", middleware.CheckApiKey(cfg.App.ApiKey), h.CreateIncidentsHandler)
	body, _ := json.Marshal(map[string]any{
		"lat": 55.755,
		"lon": 37.617,
	})
	req, _ := http.NewRequest(http.MethodPost, "/api/v1/incidents", bytes.NewReader(body))
	req.Header.Set("X-API-KEY", cfg.App.ApiKey)
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)
}
func TestCreateIncidents_bad_request(t *testing.T) {
	gin.SetMode(gin.TestMode)
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	incSvc := mock_handlers.NewMockincidentService(ctrl)

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

	h := handlers.NewIncidentHandler(incSvc, logger, &cfg.App)

	router := gin.New()
	router.POST("/api/v1/incidents", middleware.CheckApiKey(cfg.App.ApiKey), h.CreateIncidentsHandler)
	body, _ := json.Marshal(map[string]any{
		"lat": 55.755,
		"lon": -307.617, //данные выходят за диапазон значений для долготы (От 0° до 180°)
	})
	req, _ := http.NewRequest(http.MethodPost, "/api/v1/incidents", bytes.NewReader(body))
	req.Header.Set("X-API-KEY", cfg.App.ApiKey)
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestCreateIncidents_server_error(t *testing.T) {
	gin.SetMode(gin.TestMode)
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	incSvc := mock_handlers.NewMockincidentService(ctrl)

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

	h := handlers.NewIncidentHandler(incSvc, logger, &cfg.App)
	dbErr := errors.New("db is down") //имитируем падения бд
	incSvc.EXPECT().CreateIncidents(gomock.Any(), gomock.Any()).Return(dbErr)
	router := gin.New()
	router.POST("/api/v1/incidents", middleware.CheckApiKey(cfg.App.ApiKey), h.CreateIncidentsHandler)
	body, _ := json.Marshal(map[string]any{
		"lat": 55.755,
		"lon": 37.617,
	})
	req, _ := http.NewRequest(http.MethodPost, "/api/v1/incidents", bytes.NewReader(body))
	req.Header.Set("X-API-KEY", cfg.App.ApiKey)
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusInternalServerError, w.Code)
}

// набор тестов для UpdateIncidentsByIDHandler
func TestUpdateIncidentsByIDHandler_request_ok(t *testing.T) {
	gin.SetMode(gin.TestMode)
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	incSvc := mock_handlers.NewMockincidentService(ctrl)

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

	h := handlers.NewIncidentHandler(incSvc, logger, &cfg.App)
	incSvc.EXPECT().UpdateIncidentsByID(gomock.Any(), gomock.AssignableToTypeOf(&dto.IncidentDTO{}), 1).Return(nil)
	router := gin.New()
	router.PUT("/api/v1/incidents", middleware.CheckApiKey(cfg.App.ApiKey), h.UpdateIncidentsByIDHandler)
	body, _ := json.Marshal(map[string]any{
		"lat": 55.755,
		"lon": 37.617,
	})
	req, _ := http.NewRequest(http.MethodPut, "/api/v1/incidents?id=1", bytes.NewReader(body))
	req.Header.Set("X-API-KEY", cfg.App.ApiKey)
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)
}
func TestUpdateIncidentsByIDHandler_bad_request(t *testing.T) {
	gin.SetMode(gin.TestMode)
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	incSvc := mock_handlers.NewMockincidentService(ctrl)

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

	h := handlers.NewIncidentHandler(incSvc, logger, &cfg.App)
	router := gin.New()
	router.PUT("/api/v1/incidents", middleware.CheckApiKey(cfg.App.ApiKey), h.UpdateIncidentsByIDHandler)
	body, _ := json.Marshal(map[string]any{
		"lat": 505.755, //данные выходят за диапазон значений для долготы (От 0° до 180°)
		"lon": 37.617,
	})
	req, _ := http.NewRequest(http.MethodPut, "/api/v1/incidents?id=1", bytes.NewReader(body)) //передаем не валидный id
	req.Header.Set("X-API-KEY", cfg.App.ApiKey)
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusBadRequest, w.Code)
}
func TestUpdateIncidentsByIDHandler_server_error(t *testing.T) {
	gin.SetMode(gin.TestMode)
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	incSvc := mock_handlers.NewMockincidentService(ctrl)

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

	h := handlers.NewIncidentHandler(incSvc, logger, &cfg.App)
	dbErr := errors.New("db is down") //имитируем падения бд
	incSvc.EXPECT().UpdateIncidentsByID(gomock.Any(), gomock.AssignableToTypeOf(&dto.IncidentDTO{}), 1).Return(dbErr)

	router := gin.New()
	router.PUT("/api/v1/incidents", middleware.CheckApiKey(cfg.App.ApiKey), h.UpdateIncidentsByIDHandler)
	body, _ := json.Marshal(map[string]any{
		"lat": 55.755,
		"lon": 37.617,
	})

	req, _ := http.NewRequest(http.MethodPut, "/api/v1/incidents?id=1", bytes.NewReader(body))
	req.Header.Set("X-API-KEY", cfg.App.ApiKey)
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusInternalServerError, w.Code)
}

// набор тестов для DeleteIncidentsByIDHandler
func TestDeleteIncidentsByIDHandler_request_ok(t *testing.T) {
	gin.SetMode(gin.TestMode)
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	incSvc := mock_handlers.NewMockincidentService(ctrl)

	logger := zap.NewNop()
	cfg := config.Config{
		App: config.AppConfig{
			ApiKey: testApiKey,
		},
	}

	h := handlers.NewIncidentHandler(incSvc, logger, &cfg.App)
	incSvc.EXPECT().DeleteIncidentsByID(gomock.Any(), 1).Return(nil)
	router := gin.New()
	router.DELETE("/api/v1/incidents", middleware.CheckApiKey(cfg.App.ApiKey), h.DeleteIncidentsByIDHandler)
	req, _ := http.NewRequest(http.MethodDelete, "/api/v1/incidents?id=1", nil)
	req.Header.Set("X-API-KEY", cfg.App.ApiKey)

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)
}
func TestDeleteIncidentsByIDHandler_bad_request(t *testing.T) {
	gin.SetMode(gin.TestMode)
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	incSvc := mock_handlers.NewMockincidentService(ctrl)

	logger := zap.NewNop()
	cfg := config.Config{
		App: config.AppConfig{
			ApiKey: testApiKey,
		},
	}

	h := handlers.NewIncidentHandler(incSvc, logger, &cfg.App)

	router := gin.New()
	router.DELETE("/api/v1/incidents", middleware.CheckApiKey(cfg.App.ApiKey), h.DeleteIncidentsByIDHandler)
	req, _ := http.NewRequest(http.MethodDelete, "/api/v1/incidents?id=w", nil) //передаем не валидный id
	req.Header.Set("X-API-KEY", cfg.App.ApiKey)

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusBadRequest, w.Code)
}
func TestDeleteIncidentsByIDHandler_server_error(t *testing.T) {
	gin.SetMode(gin.TestMode)
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	incSvc := mock_handlers.NewMockincidentService(ctrl)

	logger := zap.NewNop()
	cfg := config.Config{
		App: config.AppConfig{
			ApiKey: testApiKey,
		},
	}

	h := handlers.NewIncidentHandler(incSvc, logger, &cfg.App)
	dbErr := errors.New("db is down") //имитируем падение бд
	incSvc.EXPECT().DeleteIncidentsByID(gomock.Any(), 1).Return(dbErr)
	router := gin.New()
	router.DELETE("/api/v1/incidents", middleware.CheckApiKey(cfg.App.ApiKey), h.DeleteIncidentsByIDHandler)
	req, _ := http.NewRequest(http.MethodDelete, "/api/v1/incidents?id=1", nil)
	req.Header.Set("X-API-KEY", cfg.App.ApiKey)

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusInternalServerError, w.Code)
}

//набор тестов для
