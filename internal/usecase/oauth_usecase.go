
package usecase

import (
	"context"
	"errors"
	"fmt"

	"mymodule/internal/domain"
)

// oauthUsecase implements the domain.OAuthUsecase interface.
type oauthUsecase struct {
	userRepo    domain.UserRepository
	tokenRepo   domain.TokenRepository
	oauthProviders map[string]domain.OAuthProvider
}

// NewOAuthUsecase creates a new OAuth use case.
func NewOAuthUsecase(userRepo domain.UserRepository, tokenRepo domain.TokenRepository, providers map[string]domain.OAuthProvider) domain.OAuthUsecase {
	return &oauthUsecase{
		userRepo:    userRepo,
		tokenRepo:   tokenRepo,
		oauthProviders: providers,
	}
}

// HandleCallback handles the OAuth2 callback, creates or finds a user, and returns a JWT.
func (uc *oauthUsecase) HandleCallback(ctx context.Context, providerName, state, code string) (*domain.User, string, error) {
	provider, ok := uc.oauthProviders[providerName]
	if !ok {
		return nil, "", fmt.Errorf("provider %s not supported", providerName)
	}

	// TODO: Validate state to prevent CSRF attacks.

	// Exchange the authorization code for a token.
	accessToken, err := provider.ExchangeCodeForToken(ctx, code)
	if err != nil {
		return nil, "", err
	}

	// Get user info from the provider.
	userInfo, err := provider.GetUserInfo(ctx, accessToken)
	if err != nil {
		return nil, "", err
	}

	// Check if the user already exists.
	user, err := uc.userRepo.FindByProviderID(domain.AuthProvider(providerName), userInfo.ID)
	if err != nil {
		// If the user is not found, create a new one.
		if errors.Is(err, errors.New("user not found")) {
			newUser := &domain.User{
				Username:       userInfo.Email, // Use email as the username
				Role:           "user",         // Default role
				AuthProvider:   domain.AuthProvider(providerName),
				AuthProviderID: userInfo.ID,
			}
			user, err = uc.userRepo.Create(newUser)
			if err != nil {
				return nil, "", err
			}
		} else {
			// Handle other potential errors from the repository.
			return nil, "", err
		}
	}

	// Generate a JWT token for the user.
	jwtToken, err := uc.tokenRepo.GenerateToken(user.ID, user.Role)
	if err != nil {
		return nil, "", err
	}

	return user, jwtToken, nil
}
