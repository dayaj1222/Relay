package middleware

import (
	"context"
	"net/http"
	"relay/internals/users"
	"strings"

	"github.com/gin-gonic/gin"
)

func AuthRequired(userService *users.Service) gin.HandlerFunc {
	return func(c *gin.Context) {
		var apiToken string

		authHeader := c.GetHeader("Authorization")
		if authHeader != "" {
			parts := strings.Split(authHeader, " ")
			if len(parts) != 2 || parts[0] != "Bearer" {
				c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Authorization header format must be Bearer {token}"})
				return
			}
			apiToken = parts[1]
		} else {
			// WebSocket connections can't send headers — fall back to query param
			apiToken = c.Query("token")
			if apiToken == "" {
				c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Authorization required"})
				return
			}
		}

		userID, err := userService.FindUserIDByToken(c.Request.Context(), apiToken)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid or expired api token"})
			return
		}

		c.Set("user_id", userID)
		reqCtx := context.WithValue(c.Request.Context(), "user_id", userID)
		c.Request = c.Request.WithContext(reqCtx)
		c.Next()
	}
}
