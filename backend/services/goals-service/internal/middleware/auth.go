package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// AuthMiddleware ensures the X-User-ID header is present and valid
func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		userIDStr := c.GetHeader("X-User-ID")
		if userIDStr == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized: Missing X-User-ID header"})
			c.Abort()
			return
		}

		_, err := uuid.Parse(userIDStr)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized: Invalid X-User-ID header"})
			c.Abort()
			return
		}

		c.Next()
	}
}
