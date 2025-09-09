
package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"mymodule/internal/adapter/http/session"
	"mymodule/internal/domain"
	"mymodule/internal/errorx"
)

// AuthHandler handles authentication-related requests.
type AuthHandler struct {
	authUsecase domain.AuthUsecase
}

// NewAuthHandler creates a new AuthHandler.
func NewAuthHandler(authUsecase domain.AuthUsecase) *AuthHandler {
	return &AuthHandler{authUsecase: authUsecase}
}

// Login handles the login request, creates a session, and returns a JWT token.
func (h *AuthHandler) Login(c *gin.Context) {
	var req domain.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.Error(errorx.New(http.StatusBadRequest, "Invalid request body"))
		return
	}

	// Usecase returns the user object and a JWT token string.
	user, token, err := h.authUsecase.Login(req.Username, req.Password)
	if err != nil {
		c.Error(errorx.New(http.StatusUnauthorized, "Invalid username or password"))
		return
	}

	// --- Session Creation ---
	sess, err := session.Store.Get(c.Request, session.SessionName)
	if err != nil {
		c.Error(errorx.New(http.StatusInternalServerError, "Failed to get session").WithStack())
		return
	}

	// Set user details in the session.
	sess.Values[session.UserIDKey] = user.ID
	sess.Values[session.UserRoleKey] = user.Role
	sess.Options.HttpOnly = true
	sess.Options.Secure = c.Request.TLS != nil // Set to true if using HTTPS
	sess.Options.MaxAge = 86400 * 7             // 7 days

	if err := sess.Save(c.Request, c.Writer); err != nil {
		c.Error(errorx.New(http.StatusInternalServerError, "Failed to save session").WithStack())
		return
	}

	// --- JWT Token Response ---
	c.JSON(http.StatusOK, gin.H{
		"message": "Login successful",
		"token":   token, // Return the JWT token for non-browser clients
	})
}

// Logout handles the logout request and destroys the session.
func (h *AuthHandler) Logout(c *gin.Context) {
	sess, err := session.Store.Get(c.Request, session.SessionName)
	if err != nil {
		c.Error(errorx.New(http.StatusInternalServerError, "Failed to get session").WithStack())
		return
	}

	// Invalidate the session by setting MaxAge to -1.
	sess.Options.MaxAge = -1

	if err := sess.Save(c.Request, c.Writer); err != nil {
		c.Error(errorx.New(http.StatusInternalServerError, "Failed to invalidate session").WithStack())
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Logout successful"})
}
