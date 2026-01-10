package middleware

import (
	"log"
	"net/http"

	"github.com/BitCoinOffical/geo-announcements/app-1/config"
	"github.com/gin-gonic/gin"
)

const (
	apiKey  = "X-API-KEY"
	message = "Invalid API Key"
)

func CheckApiKey(cfg *config.AppConfig) gin.HandlerFunc {
	return func(c *gin.Context) {
		apiKey := c.Request.Header.Get(apiKey)
		log.Println(cfg.API_KEY)
		if apiKey != cfg.API_KEY {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
				"message": message,
			})
			return
		}
		c.Next()
	}
}
