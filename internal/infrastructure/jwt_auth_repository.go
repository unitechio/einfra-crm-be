
package infrastructure

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"mymodule/internal/config"
	"mymodule/internal/domain"
)

// jwtAuthRepository is the implementation of AuthRepository using JWT.
type jwtAuthRepository struct {
	cfg config.JWTConfig
}

// NewJWTAuthRepository creates a new JWTAuthRepository.
func NewJWTAuthRepository(cfg config.JWTConfig) domain.AuthRepository {
	return &jwtAuthRepository{cfg: cfg}
}

// GenerateToken generates a new JWT token.
func (r *jwtAuthRepository) GenerateToken(userID string, roles []string) (string, error) {
	expirationTime := time.Now().Add(time.Duration(r.cfg.ExpirationHours) * time.Hour)
	claims := &domain.Claims{
		UserID: userID,
		Roles:  roles,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(r.cfg.Secret))
}

// ValidateToken validates a JWT token.
func (r *jwtAuthRepository) ValidateToken(tokenString string) (*domain.Claims, error) {
	claims := &domain.Claims{}

	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(r.cfg.Secret), nil
	})

	if err != nil {
		return nil, err
	}

	if !token.Valid {
		return nil, fmt.Errorf("invalid token")
	}

	return claims, nil
}
