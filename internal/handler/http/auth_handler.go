
package http

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"mymodule/internal/domain"
)

// AuthHandler handles HTTP requests for authentication.
type AuthHandler struct {
	authUsecase domain.AuthUsecase
}

// NewAuthHandler creates a new AuthHandler.
func NewAuthHandler(e *echo.Echo, au domain.AuthUsecase) {
	h := &AuthHandler{authUsecase: au}

	// Group for auth routes
	group := e.Group("/auth")

	group.POST("/login", h.handleLocalLogin)
}

// handleLocalLogin handles the local username/password login.
func (h *AuthHandler) handleLocalLogin(c echo.Context) error {
	var req domain.LoginRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request payload"})
	}

	user, token, err := h.authUsecase.Login(req.Username, req.Password)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"message": "Login successful",
		"user":    user,
		"token":   token,
	})
}
