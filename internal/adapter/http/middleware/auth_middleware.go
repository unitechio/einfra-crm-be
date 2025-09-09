
package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"mymodule/internal/domain"
	"mymodule/internal/errorx"
)

// AuthMiddleware creates a middleware handler for JWT authentication.
func AuthMiddleware(tokenRepo domain.TokenRepository) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.Error(errorx.New(http.StatusUnauthorized, "Authorization header is missing"))
			c.Abort()
			return
		}

		// The header should be in the format "Bearer <token>"
		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
			c.Error(errorx.New(http.StatusUnauthorized, "Authorization header format must be Bearer {token}"))
			c.Abort()
			return
		}

		tokenString := parts[1]

		claims, err := tokenRepo.ValidateToken(tokenString)
		if err != nil {
			c.Error(errorx.New(http.StatusUnauthorized, "Invalid or expired token").WithCause(err))
			c.Abort()
			return
		}

		// Set user information in the context for downstream handlers.
		c.Set("userID", claims.UserID)
		c.Set("userRole", claims.Role)

		c.Next()
	}
}
