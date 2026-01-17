package middleware_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/BitCoinOffical/geo-announcements/app-1/internal/api/middleware"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/assert/v2"
)

const (
	testApiKey = "SreUeZhc6VEaSl7Vlg9kB3dbQ9xUr0EB"
)

// набор тестов для CheckApiKey
func TestCheckApiKey_request_ok(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.OPTIONS("/", middleware.CheckApiKey(testApiKey), func(ctx *gin.Context) {
		ctx.Status(http.StatusOK)
	})

	req, _ := http.NewRequest(http.MethodOptions, "/", nil)
	req.Header.Set("X-API-KEY", testApiKey)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)
}

func TestCheckApiKey_server_error(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.OPTIONS("/", middleware.CheckApiKey(testApiKey), func(ctx *gin.Context) {
		ctx.Status(http.StatusOK)
	})

	req, _ := http.NewRequest(http.MethodOptions, "/", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusInternalServerError, w.Code)
}
