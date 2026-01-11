package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

const (
	apiKey  = "X-API-KEY"
	message = "Invalid API Key"
)

func CheckApiKey(appApiKey string) gin.HandlerFunc {
	return func(c *gin.Context) {
		apiKey := c.Request.Header.Get(apiKey)
		if apiKey != appApiKey {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
				"message": message,
			})
			return
		}
		c.Next()
	}
}
