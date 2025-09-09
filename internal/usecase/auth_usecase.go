
package usecase

import (
	"context"
	"errors"

	"mymodule/internal/domain"
)

// authUsecase implements both AuthUsecase and OAuthUsecase.
type authUsecase struct {
	userRepo        domain.UserRepository
	tokenRepo       domain.TokenRepository
	oauthProviders  map[string]domain.OAuthProvider
}

// NewAuthUsecase creates a new combined AuthUsecase.
func NewAuthUsecase(userRepo domain.UserRepository, tokenRepo domain.TokenRepository, providers map[string]domain.OAuthProvider) domain.AuthUsecase {
	return &authUsecase{
		userRepo:        userRepo,
		tokenRepo:       tokenRepo,
		oauthProviders:  providers,
	}
}

// Login validates credentials for local users.
func (uc *authUsecase) Login(username, password string) (*domain.User, string, error) {
	user, err := uc.userRepo.FindByUsername(username)
	if err != nil {
		return nil, "", err
	}

	if user.AuthProvider != domain.AuthProviderLocal {
		return nil, "", errors.New("user is not a local user")
	}

	if user.PasswordHash != password { // Example check, use bcrypt in production
		return nil, "", errors.New("invalid credentials")
	}

	token, err := uc.tokenRepo.GenerateToken(user.ID, user.Role)
	if err != nil {
		return nil, "", err
	}

	return user, token, nil
}

// HandleLogin gets the authentication URL from the specified provider.
func (uc *authUsecase) HandleLogin(provider string) (string, error) {
	p, ok := uc.oauthProviders[provider]
	if !ok {
		return "", errors.New("invalid oauth provider")
	}
	state := "random-state-string" // Should be random and validated
	return p.GetAuthURL(state), nil
}

// HandleCallback handles the callback from the OAuth provider.
func (uc *authUsecase) HandleCallback(ctx context.Context, providerName, code string) (*domain.User, string, error) {
	p, ok := uc.oauthProviders[providerName]
	if !ok {
		return nil, "", errors.New("invalid oauth provider")
	}

	providerToken, err := p.ExchangeCodeForToken(ctx, code)
	if err != nil {
		return nil, "", err
	}

	userInfo, err := p.GetUserInfo(ctx, providerToken)
	if err != nil {
		return nil, "", err
	}

	provider := domain.AuthProvider(providerName)
	user, err := uc.userRepo.FindByAuthProvider(provider, userInfo.ID)
	if err != nil {
		newUser := &domain.User{
			Username:       userInfo.Email,
			Role:           "user",
			AuthProvider:   provider,
			AuthProviderID: userInfo.ID,
		}
		user, err = uc.userRepo.Create(newUser)
		if err != nil {
			return nil, "", err
		}
	}

	jwtToken, err := uc.tokenRepo.GenerateToken(user.ID, user.Role)
	if err != nil {
		return nil, "", err
	}

	return user, jwtToken, nil
}

func (uc *authUsecase) ProtectedData(userID string) (string, error) {
	return "This is protected data for user " + userID, nil
}
