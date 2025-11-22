
package http

import (
	"context"
	"github.com/unitechio/einfra-be/internal/constants"
	"github.com/unitechio/einfra-be/internal/domain"
	"net/http"
	"os"
	"strings"

	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
)

// AuthMiddlewareConfig holds the dependencies for the auth middleware.
type AuthMiddlewareConfig struct {
	UserRepo domain.UserRepository
}

// AuthMiddleware creates a new middleware for JWT authentication.
func (amc *AuthMiddlewareConfig) AuthMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		// 1. Get the Authorization header
		authHeader := c.Request().Header.Get("Authorization")
		if authHeader == "" {
			return echo.NewHTTPError(http.StatusUnauthorized, "missing Authorization header")
		}

		// 2. Validate the header format (Bearer <token>)
		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
			return echo.NewHTTPError(http.StatusUnauthorized, "invalid Authorization header format")
		}
		tokenString := parts[1]

		// 3. Parse and validate the token
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, echo.NewHTTPError(http.StatusUnauthorized, "unexpected signing method")
			}
			return []byte(os.Getenv("JWT_SECRET")), nil
		})

		if err != nil || !token.Valid {
			return echo.NewHTTPError(http.StatusUnauthorized, "invalid or expired token")
		}

		// 4. Extract claims and get user ID
		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			return echo.NewHTTPError(http.StatusUnauthorized, "invalid token claims")
		}

		userID, ok := claims["user_id"].(string)
		if !ok {
			return echo.NewHTTPError(http.StatusUnauthorized, "invalid user ID in token")
		}

		// 5. Fetch the user from the repository
		user, err := amc.UserRepo.GetByID(c.Request().Context(), userID)
		if err != nil {
			return echo.NewHTTPError(http.StatusUnauthorized, "user not found")
		}

		// 6. Store user in context for later handlers
		ctx := context.WithValue(c.Request().Context(), constants.UserContextKey, user)
		c.SetRequest(c.Request().WithContext(ctx))

		return next(c)
	}
}

// RoleMiddleware creates a new middleware for role-based authorization.
func RoleMiddleware(allowedRoles ...domain.Role) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			// 1. Get user from context (set by AuthMiddleware)
			user, ok := c.Request().Context().Value(constants.UserContextKey).(*domain.User)
			if !ok || user == nil {
				return echo.NewHTTPError(http.StatusForbidden, "user not found in context")
			}

			// 2. Check if the user's role is in the list of allowed roles
			allowed := false
			for _, role := range allowedRoles {
				if user.Role == role {
					allowed = true
					break
				}
			}

			if !allowed {
				return echo.NewHTTPError(http.StatusForbidden, "insufficient permissions")
			}

			return next(c)
		}
	}
}
