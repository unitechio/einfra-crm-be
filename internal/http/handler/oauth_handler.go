
package http

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/unitechio/einfra-be/internal/domain"
)

// OAuthHandler handles the HTTP requests for OAuth2.
type OAuthHandler struct {
	OAuthUsecase domain.OAuthUsecase
	OAuthProviders map[string]domain.OAuthProvider
}

// NewOAuthHandler creates a new OAuthHandler.
func NewOAuthHandler(e *echo.Echo, oauthUsecase domain.OAuthUsecase, providers map[string]domain.OAuthProvider) {
	h := &OAuthHandler{
		OAuthUsecase: oauthUsecase,
		OAuthProviders: providers,
	}

	e.GET("/oauth/:provider/login", h.Login)
	e.GET("/oauth/:provider/callback", h.Callback)
}

// Login redirects the user to the provider's login page.
func (h *OAuthHandler) Login(c echo.Context) error {
	providerName := c.Param("provider")
	provider, ok := h.OAuthProviders[providerName]
	if !ok {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Provider not supported"})
	}

	// TODO: Generate and store a random state string in the session for CSRF protection.
	state := "random-state-string"
	url := provider.GetAuthURL(state)
	return c.Redirect(http.StatusTemporaryRedirect, url)
}

// Callback handles the callback from the provider.
func (h *OAuthHandler) Callback(c echo.Context) error {
	providerName := c.Param("provider")
	state := c.QueryParam("state")
	code := c.QueryParam("code")

	// TODO: Validate the state string.

	user, token, err := h.OAuthUsecase.HandleCallback(c.Request().Context(), providerName, state, code)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	// On success, you might set a cookie or return the token in the response.
	return c.JSON(http.StatusOK, map[string]interface{}{
		"message": "Login successful",
		"token":   token,
		"user":    user,
	})
}
