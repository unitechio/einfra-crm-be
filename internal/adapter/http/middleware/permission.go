package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"mymodule/internal/domain" // Corrected project path
)

const (
	// UserRoleKey is the key used to store the user's role in the Gin context.
	// This key MUST be set by a preceding authentication middleware (e.g., AuthMiddleware).
	UserRoleKey = "userRole"
)

// RequireRoles creates a middleware handler that checks if the user's role,
// retrieved from the Gin context, is one of the allowed roles.
// This middleware MUST run AFTER an authentication middleware has successfully run.
func RequireRoles(allowedRoles ...domain.Role) gin.HandlerFunc {
	return func(c *gin.Context) {
		// 1. Get the user's role from the Gin context.
		roleFromContext, exists := c.Get(UserRoleKey)
		if !exists {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "User role not found in context. Is authentication middleware missing?"})
			return
		}

		// 2. Assert that the role is in the expected string format.
		userRole, ok := roleFromContext.(string)
		if !ok {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Invalid role format in context"})
			return
		}

		// 3. Check if the user's role is in the list of allowed roles for this endpoint.
		for _, allowedRole := range allowedRoles {
			// We cast the `domain.Role` constant to a string for comparison.
			if userRole == string(allowedRole) {
				c.Next() // Role is allowed, continue to the next handler in the chain.
				return
			}
		}

		// 4. If the loop completes, the user's role is not in the allowed list. Access denied.
		c.AbortWithStatusJSON(http.StatusForbidden, gin.H{
			"error": "Forbidden: You do not have the required role to access this resource",
		})
	}
}
