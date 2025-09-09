package infrastructure

import (
	"context"
	"fmt"

	"github.com/Nerzal/gocloak/v13"
	"mymodule/internal/config"
	"mymodule/internal/domain"
)

// keycloakAuthRepository is the implementation of AuthRepository using Keycloak.
type keycloakAuthRepository struct {
	client *gocloak.GoCloak
	cfg    config.KeycloakConfig
}

// NewKeycloakAuthRepository creates a new KeycloakAuthRepository.
func NewKeycloakAuthRepository(cfg config.KeycloakConfig) domain.AuthRepository {
	client := gocloak.NewClient(cfg.URL)
	return &keycloakAuthRepository{
		client: client,
		cfg:    cfg,
	}
}

// GenerateToken for Keycloak is not implemented as tokens are issued by Keycloak itself.
func (r *keycloakAuthRepository) GenerateToken(userID string, roles []string) (string, error) {
	return "", fmt.Errorf("token generation is handled by Keycloak")
}

// ValidateToken validates a token using Keycloak's introspection endpoint.
func (r *keycloakAuthRepository) ValidateToken(tokenString string) (*domain.Claims, error) {
	ctx := context.Background()

	rptResult, err := r.client.RetrospectToken(ctx, tokenString, r.cfg.ClientID, r.cfg.ClientSecret, r.cfg.Realm)
	if err != nil {
		return nil, fmt.Errorf("token introspection failed: %w", err)
	}

	if !*rptResult.Active {
		return nil, fmt.Errorf("token is not active")
	}

	// Decode the token to get the claims without verifying the signature (already done by Keycloak).
	_, claims, err := r.client.DecodeAccessToken(ctx, tokenString, r.cfg.Realm)
	if err != nil {
		return nil, fmt.Errorf("failed to decode token: %w", err)
	}

	// Extract roles from the claims.
	// This depends on how roles are configured in Keycloak.
	// We'll look for them in 'realm_access.roles'.
	var roles []string
	if realmAccess, ok := (*claims)["realm_access"].(map[string]interface{}); ok {
		if realmRoles, ok := realmAccess["roles"].([]interface{}); ok {
			for _, role := range realmRoles {
				if r, ok := role.(string); ok {
					roles = append(roles, r)
				}
			}
		}
	}

	// Create our custom claims struct.
	customClaims := &domain.Claims{
		UserID: (*claims)["sub"].(string),
		Roles:  roles,
	}

	return customClaims, nil
}
