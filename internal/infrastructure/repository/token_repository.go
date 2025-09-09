
package repository

import (
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"mymodule/internal/domain"
)

// jwtTokenRepository is the implementation of TokenRepository.
type jwtTokenRepository struct {
	secretKey []byte
}

// NewTokenRepository creates a new TokenRepository.
func NewTokenRepository() domain.TokenRepository {
	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		panic("JWT_SECRET environment variable not set")
	}
	return &jwtTokenRepository{secretKey: []byte(secret)}
}

// GenerateToken generates a new JWT for the given user ID and role.
func (r *jwtTokenRepository) GenerateToken(userID, userRole string) (string, error) {
	expirationTime := time.Now().Add(24 * time.Hour) // Token valid for 24 hours

	claims := &domain.Claims{
		UserID: userID,
		Role:   userRole,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Subject:   userID,
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(r.secretKey)
}

// ValidateToken validates the given JWT string.
func (r *jwtTokenRepository) ValidateToken(tokenString string) (*domain.Claims, error) {
	claims := &domain.Claims{}

	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		// Don't forget to validate the alg is what you expect:
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, jwt.ErrSignatureInvalid
		}
		return r.secretKey, nil
	})

	if err != nil {
		return nil, err
	}

	if !token.Valid {
		return nil, errors.New("invalid token")
	}

	return claims, nil
}
