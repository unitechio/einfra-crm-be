
package oauth

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/microsoft"
	"github.com/unitechio/einfra-be/internal/config"
	"github.com/unitechio/einfra-be/internal/domain"
)

// AzureProvider implements the OAuthProvider interface for Azure AD.
type AzureProvider struct {
	config *oauth2.Config
}

// NewAzureProvider creates a new AzureProvider.
func NewAzureProvider(cfg *config.AzureOAuthConfig) domain.OAuthProvider {
	return &AzureProvider{
		config: &oauth2.Config{
			ClientID:     cfg.ClientID,
			ClientSecret: cfg.ClientSecret,
			RedirectURL:  cfg.RedirectURL,
			Scopes:       cfg.Scopes,
			Endpoint:     microsoft.AzureADEndpoint(cfg.Tenant),
		},
	}
}

// GetAuthURL generates the authentication URL for the user to visit.
func (a *AzureProvider) GetAuthURL(state string) string {
	return a.config.AuthCodeURL(state)
}

// ExchangeCodeForToken exchanges the authorization code for an OAuth token.
func (a *AzureProvider) ExchangeCodeForToken(ctx context.Context, code string) (string, error) {
	token, err := a.config.Exchange(ctx, code)
	if err != nil {
		return "", err
	}
	return token.AccessToken, nil
}

// azureUserInfoResponse represents the user information retrieved from Microsoft Graph API.
type azureUserInfoResponse struct {
	ID                string `json:"id"`
	UserPrincipalName string `json:"userPrincipalName"` // This is often the email
	DisplayName       string `json:"displayName"`
}

// GetUserInfo uses the access token to fetch user information from Microsoft Graph API.
func (a *AzureProvider) GetUserInfo(ctx context.Context, token string) (*domain.OAuthUserInfo, error) {
	const microsoftGraphAPI = "https://graph.microsoft.com/v1.0/me"

	req, err := http.NewRequestWithContext(ctx, "GET", microsoftGraphAPI, nil)
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
		bodyBytes, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("failed to get user info from Azure AD: status %d, body: %s", resp.StatusCode, string(bodyBytes))
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var userInfo azureUserInfoResponse
	if err := json.Unmarshal(body, &userInfo); err != nil {
		return nil, err
	}

	return &domain.OAuthUserInfo{
		ID:    userInfo.ID,
		Email: userInfo.UserPrincipalName,
		Name:  userInfo.DisplayName,
	}, nil
}
