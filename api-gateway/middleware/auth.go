package middleware

import (
	"context"
	"net/http"
	"strings"

	pbUser "be-golang-chapter-57/user-service/proto"

	"github.com/gin-gonic/gin"
)

// AuthMiddleware checks the JWT token and validates it with the user service.
func AuthMiddleware(userClient pbUser.UserServiceClient) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header required"})
			c.Abort()
			return
		}

		// Extract the token from the header
		token := strings.TrimSpace(strings.TrimPrefix(authHeader, "Bearer"))
		if token == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Bearer token required"})
			c.Abort()
			return
		}

		// Validate the token with user service
		req := &pbUser.ValidateTokenRequest{Token: token}
		res, err := userClient.ValidateToken(context.Background(), req)
		if err != nil || !res.Valid {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			c.Abort()
			return
		}

		// Store the user ID in context
		c.Set("user_id", res.UserId)
		c.Next()
	}
}
