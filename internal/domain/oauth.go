
package domain

import (
	"context"
)

// OAuthProvider defines the interface for an OAuth2 provider.
type OAuthProvider interface {
	GetAuthURL(state string) string
	ExchangeCodeForToken(ctx context.Context, code string) (string, error)
	GetUserInfo(ctx context.Context, token string) (*OAuthUserInfo, error)
}

// OAuthUsecase defines the interface for handling the OAuth2 callback.
type OAuthUsecase interface {
	HandleCallback(ctx context.Context, provider, state, code string) (*User, string, error)
}

// OAuthUserInfo represents the user information retrieved from an OAuth provider.
type OAuthUserInfo struct {
	ID    string `json:"id"`
	Email string `json:"email"`
	Name  string `json:"name"`
}
