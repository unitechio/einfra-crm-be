package middleware

import (
	"context"
	"net/http"
	"strings"

	"github.com/unitechio/einfra-be/internal/domain"
	"github.com/gin-gonic/gin"
)

// AuthMiddleware validates JWT tokens and sets user context
func AuthMiddleware(authRepo domain.AuthRepository) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get token from Authorization header
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "Authorization header required",
				"code":  "AUTH_HEADER_MISSING",
			})
			c.Abort()
			return
		}

		// Extract token from "Bearer <token>"
		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || parts[0] != "Bearer" {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "Invalid authorization header format",
				"code":  "INVALID_AUTH_HEADER",
			})
			c.Abort()
			return
		}

		token := parts[1]

		// Validate token
		claims, err := authRepo.ValidateToken(c.Request.Context(), token, domain.TokenTypeAccess)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "Invalid or expired token",
				"code":  "INVALID_TOKEN",
			})
			c.Abort()
			return
		}

		// Set user information in context
		c.Set("user_id", claims.UserID)
		c.Set("username", claims.Username)
		c.Set("email", claims.Email)
		c.Set("role_id", claims.RoleID)
		c.Set("role_name", claims.RoleName)
		c.Set("permissions", claims.Permissions)

		c.Next()
	}
}

// OptionalAuthMiddleware validates JWT tokens but doesn't require them
func OptionalAuthMiddleware(authRepo domain.AuthRepository) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.Next()
			return
		}

		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || parts[0] != "Bearer" {
			c.Next()
			return
		}

		token := parts[1]
		claims, err := authRepo.ValidateToken(c.Request.Context(), token, domain.TokenTypeAccess)
		if err == nil {
			c.Set("user_id", claims.UserID)
			c.Set("username", claims.Username)
			c.Set("email", claims.Email)
			c.Set("role_id", claims.RoleID)
			c.Set("role_name", claims.RoleName)
			c.Set("permissions", claims.Permissions)
		}

		c.Next()
	}
}

// PermissionMiddleware checks if user has required permission
func PermissionMiddleware(requiredPermission string) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get permissions from context (set by AuthMiddleware)
		permissionsInterface, exists := c.Get("permissions")
		if !exists {
			c.JSON(http.StatusForbidden, gin.H{
				"error": "No permissions found",
				"code":  "NO_PERMISSIONS",
			})
			c.Abort()
			return
		}

		permissions, ok := permissionsInterface.([]string)
		if !ok {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "Invalid permissions format",
				"code":  "INVALID_PERMISSIONS",
			})
			c.Abort()
			return
		}

		// Check if user has required permission
		hasPermission := false
		for _, perm := range permissions {
			if perm == requiredPermission || perm == "*" {
				hasPermission = true
				break
			}
		}

		if !hasPermission {
			c.JSON(http.StatusForbidden, gin.H{
				"error":               "Insufficient permissions",
				"code":                "INSUFFICIENT_PERMISSIONS",
				"required_permission": requiredPermission,
			})
			c.Abort()
			return
		}

		c.Next()
	}
}

// RoleMiddleware checks if user has required role
func RoleMiddleware(requiredRole string) gin.HandlerFunc {
	return func(c *gin.Context) {
		roleNameInterface, exists := c.Get("role_name")
		if !exists {
			c.JSON(http.StatusForbidden, gin.H{
				"error": "No role found",
				"code":  "NO_ROLE",
			})
			c.Abort()
			return
		}

		roleName, ok := roleNameInterface.(string)
		if !ok {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "Invalid role format",
				"code":  "INVALID_ROLE",
			})
			c.Abort()
			return
		}

		if roleName != requiredRole && roleName != "admin" {
			c.JSON(http.StatusForbidden, gin.H{
				"error":         "Insufficient role",
				"code":          "INSUFFICIENT_ROLE",
				"required_role": requiredRole,
			})
			c.Abort()
			return
		}

		c.Next()
	}
}

// GetUserID gets user ID from context
func GetUserID(c *gin.Context) (string, bool) {
	userID, exists := c.Get("user_id")
	if !exists {
		return "", false
	}
	id, ok := userID.(string)
	return id, ok
}

// GetUsername gets username from context
func GetUsername(c *gin.Context) (string, bool) {
	username, exists := c.Get("username")
	if !exists {
		return "", false
	}
	name, ok := username.(string)
	return name, ok
}

// GetPermissions gets permissions from context
func GetPermissions(c *gin.Context) ([]string, bool) {
	permissions, exists := c.Get("permissions")
	if !exists {
		return nil, false
	}
	perms, ok := permissions.([]string)
	return perms, ok
}

// HasPermission checks if user has a specific permission
func HasPermission(c *gin.Context, permission string) bool {
	permissions, ok := GetPermissions(c)
	if !ok {
		return false
	}

	for _, perm := range permissions {
		if perm == permission || perm == "*" {
			return true
		}
	}
	return false
}

// SetUserContext sets user context for downstream services
func SetUserContext(c *gin.Context, ctx context.Context) context.Context {
	if userID, ok := GetUserID(c); ok {
		ctx = context.WithValue(ctx, "user_id", userID)
	}
	if username, ok := GetUsername(c); ok {
		ctx = context.WithValue(ctx, "username", username)
	}
	return ctx
}
