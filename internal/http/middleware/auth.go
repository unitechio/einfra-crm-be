
package middleware

import (
	"context"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"mymodule/internal/util"
)

// contextKey is a private type to prevent collisions with other context keys.
// This avoids staticcheck error: SA1029
type contextKey string

const (
	// UserIDKey is the context key for the user's ID.
	UserIDKey contextKey = "userID"
	// UserRoleKey is the context key for the user's role.
	UserRoleKey contextKey = "userRole"
)

// AuthMiddleware validates the JWT token from the Authorization header.
func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Authorization header is missing"})
			return
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid Authorization header format"})
			return
		}

		tokenString := parts[1]
		claims, err := util.ValidateJWT(tokenString)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
			return
		}

		// Store user information in the request's context using a custom key type.
		ctx := context.WithValue(c.Request.Context(), UserIDKey, claims["user_id"])
		ctx = context.WithValue(ctx, UserRoleKey, claims["role"])
		c.Request = c.Request.WithContext(ctx)

		c.Next()
	}
}

// RBACMiddleware checks if the user has the required role.
func RBACMiddleware(requiredRole string) gin.HandlerFunc {
	return func(c *gin.Context) {
		userRole, ok := c.Request.Context().Value(UserRoleKey).(string)
		if !ok {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "User role not found or not a string in context"})
			return
		}

		if userRole != requiredRole {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "Insufficient permissions"})
			return
		}

		c.Next()
	}
}
