
package domain

import "context"

// AuthProvider is a type for different authentication providers.
type AuthProvider string

const (
	// AuthProviderGoogle represents the Google authentication provider.
	AuthProviderGoogle AuthProvider = "google"
	// AuthProviderAzure represents the Azure AD authentication provider.
	AuthProviderAzure AuthProvider = "azure"
)

// AuthUsecase defines the interface for authentication-related use cases.
type AuthUsecase interface {
	// Register creates a new user with the given username and password.
	Register(ctx context.Context, username, password string) (*User, error)

	// Login authenticates a user and returns a JWT token.
	Login(ctx context.Context, username, password string) (string, error)
}
