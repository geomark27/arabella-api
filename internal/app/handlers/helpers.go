package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// getUserIDFromContext extracts the authenticated user's ID from the Gin context.
// The "user_id" key is set by AuthMiddleware after successful JWT validation.
// Returns the userID and true if found and valid, or 0 and false otherwise.
func getUserIDFromContext(c *gin.Context) (uint, bool) {
	value, exists := c.Get("user_id")
	if !exists {
		return 0, false
	}
	userID, ok := value.(uint)
	return userID, ok
}

// respondUnauthorized sends a standardized 401 Unauthorized JSON response and aborts the request.
func respondUnauthorized(c *gin.Context) {
	c.JSON(http.StatusUnauthorized, gin.H{
		"error": "unauthorized",
	})
	c.Abort()
}

// parseIntParam parses an integer query parameter with a fallback default value.
func parseIntParam(c *gin.Context, key string, defaultValue int) int {
	if value := c.Query(key); value != "" {
		if parsed, err := strconv.Atoi(value); err == nil {
			return parsed
		}
	}
	return defaultValue
}
