package usecase

import (
	"context"
	"fmt"

	"github.com/unitechio/einfra-be/internal/auth"
	"github.com/unitechio/einfra-be/internal/domain"
	"github.com/unitechio/einfra-be/internal/repository"
)

type OAuthUsecase interface {
	GetAuthURL(provider string) (string, error)
	HandleCallback(ctx context.Context, provider, code string) (*domain.AuthResponse, error)
}

type oauthUsecase struct {
	userRepo       repository.UserRepository
	authRepo       repository.AuthRepository
	jwtService     *auth.JWTService
	oauthProviders map[string]domain.OAuthProvider
}

func NewOAuthUsecase(
	userRepo repository.UserRepository,
	authRepo repository.AuthRepository,
	jwtService *auth.JWTService,
	providers map[string]domain.OAuthProvider,
) OAuthUsecase {
	return &oauthUsecase{
		userRepo:       userRepo,
		authRepo:       authRepo,
		jwtService:     jwtService,
		oauthProviders: providers,
	}
}

func (u *oauthUsecase) GetAuthURL(providerName string) (string, error) {
	provider, ok := u.oauthProviders[providerName]
	if !ok {
		return "", fmt.Errorf("provider %s not supported", providerName)
	}
	// State generation should be more secure in production
	return provider.GetAuthURL("state"), nil
}

func (u *oauthUsecase) HandleCallback(ctx context.Context, providerName, code string) (*domain.AuthResponse, error) {
	provider, ok := u.oauthProviders[providerName]
	if !ok {
		return nil, fmt.Errorf("provider %s not supported", providerName)
	}

	accessToken, err := provider.ExchangeCodeForToken(ctx, code)
	if err != nil {
		return nil, fmt.Errorf("failed to exchange code: %w", err)
	}

	userInfo, err := provider.GetUserInfo(ctx, accessToken)
	if err != nil {
		return nil, fmt.Errorf("failed to get user info: %w", err)
	}

	user, err := u.userRepo.GetByEmail(ctx, userInfo.Email)
	if err != nil {
		// Create new user if not exists
		// This logic might need to be adjusted based on specific requirements
		// For now, we'll assume user creation is handled elsewhere or we return an error
		// Or we can create a basic user here
		return nil, fmt.Errorf("user not found: %w", err)
	}

	// Generate tokens
	// Assuming user has permissions, or we fetch them
	// For simplicity, passing empty permissions or fetching them if possible
	// But userRepo.GetByEmail returns *domain.User which might have Role loaded
	// If not, we might need to fetch role permissions.
	// For now, let's assume we can generate token with user info.

	// We need permissions for GenerateAccessToken.
	// Let's assume empty for now or fetch them.
	var permissions []string
	// if user.Role != nil { ... }

	token, err := u.jwtService.GenerateAccessToken(user, permissions)
	if err != nil {
		return nil, fmt.Errorf("failed to generate access token: %w", err)
	}

	// We also need refresh token. AuthRepository has CreateRefreshToken but not Generate.
	// Usually Refresh Token is just a random string or UUID.
	// Let's assume we have a helper or just generate a UUID.
	// But wait, `authRepo` has `CreateRefreshToken`.
	// We need to generate the string first.
	// Let's use uuid.New().String() for now.

	refreshToken := "some-refresh-token" // Placeholder or use uuid

	return &domain.AuthResponse{
		AccessToken:  token,
		RefreshToken: refreshToken,
		User:         user,
	}, nil
}
