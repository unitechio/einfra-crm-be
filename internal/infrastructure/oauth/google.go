
package oauth

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"

	"golang.org/x/oauth2"
	googleoauth "golang.org/x/oauth2/google"
	"mymodule/internal/config"
	"mymodule/internal/domain"
)

// GoogleProvider implements the OAuthProvider interface for Google.
type GoogleProvider struct {
	config *oauth2.Config
}

// NewGoogleProvider creates a new GoogleProvider.
func NewGoogleProvider(cfg *config.GoogleOAuthConfig) domain.OAuthProvider {
	return &GoogleProvider{
		config: &oauth2.Config{
			ClientID:     cfg.ClientID,
			ClientSecret: cfg.ClientSecret,
			RedirectURL:  cfg.RedirectURL,
			Scopes:       cfg.Scopes,
			Endpoint:     googleoauth.Endpoint,
		},
	}
}

// GetAuthURL generates the authentication URL for the user to visit.
func (g *GoogleProvider) GetAuthURL(state string) string {
	return g.config.AuthCodeURL(state)
}

// ExchangeCodeForToken exchanges the authorization code for an OAuth token.
// NOTE: This returns the provider's token, not our application's JWT.
func (g *GoogleProvider) ExchangeCodeForToken(ctx context.Context, code string) (string, error) {
	token, err := g.config.Exchange(ctx, code)
	if err != nil {
		return "", err
	}
	return token.AccessToken, nil
}

// GetUserInfo uses the access token to fetch user information from Google's user info endpoint.
type googleUserInfoResponse struct {
	ID    string `json:"id"`
	Email string `json:"email"`
	Name  string `json:"name"`
}

func (g *GoogleProvider) GetUserInfo(ctx context.Context, token string) (*domain.OAuthUserInfo, error) {
	const googleAPI = "https://www.googleapis.com/oauth2/v2/userinfo"

	req, err := http.NewRequestWithContext(ctx, "GET", googleAPI, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", "Bearer "+token)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, errors.New("failed to get user info from Google")
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var userInfo googleUserInfoResponse
	if err := json.Unmarshal(body, &userInfo); err != nil {
		return nil, err
	}

	return &domain.OAuthUserInfo{
		ID:    userInfo.ID,
		Email: userInfo.Email,
		Name:  userInfo.Name,
	}, nil
}
