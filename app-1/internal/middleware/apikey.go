package middleware

import (
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
)

func CheckApiKey() gin.HandlerFunc {
	return func(c *gin.Context) {
		apiKey := c.Request.Header.Get("X-API-KEY")
		if apiKey != os.Getenv("API_KEY") {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
				"message": "Invalid API Key",
			})
			return
		}
		c.Next()
	}
}
