package auth

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"

	"github.com/unitechio/einfra-be/internal/config"
	"github.com/unitechio/einfra-be/internal/domain"
)

type JWTService struct {
	cfg *config.AuthConfig
}

func NewJWTService(cfg *config.AuthConfig) *JWTService {
	return &JWTService{cfg: cfg}
}

func (s *JWTService) GenerateAccessToken(user *domain.User, permissions []string) (string, error) {
	expiresAt := time.Now().Add(time.Duration(s.cfg.JWTExpiration) * time.Second)

	claims := jwt.MapClaims{
		"user_id":     user.ID,
		"username":    user.Username,
		"email":       user.Email,
		"role_id":     user.RoleID,
		"role_name":   "",
		"permissions": permissions,
		"token_type":  string(domain.TokenTypeAccess),
		"iat":         time.Now().Unix(),
		"exp":         expiresAt.Unix(),
	}

	if user.Role != nil {
		claims["role_name"] = user.Role.Name
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(s.cfg.JWTSecret))
}

func (s *JWTService) ValidateAccessToken(tokenString string) (*domain.TokenClaims, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(s.cfg.JWTSecret), nil
	})

	if err != nil {
		return nil, err
	}

	if !token.Valid {
		return nil, fmt.Errorf("invalid token")
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, fmt.Errorf("invalid token claims")
	}

	tokenClaims := &domain.TokenClaims{
		UserID:   claims["user_id"].(string),
		Username: claims["username"].(string),
		Email:    claims["email"].(string),
		RoleID:   claims["role_id"].(string),
		RoleName: claims["role_name"].(string),
	}

	if perms, ok := claims["permissions"].([]interface{}); ok {
		for _, perm := range perms {
			if permStr, ok := perm.(string); ok {
				tokenClaims.Permissions = append(tokenClaims.Permissions, permStr)
			}
		}
	}

	if iat, ok := claims["iat"].(float64); ok {
		tokenClaims.IssuedAt = time.Unix(int64(iat), 0)
	}
	if exp, ok := claims["exp"].(float64); ok {
		tokenClaims.ExpiresAt = time.Unix(int64(exp), 0)
	}

	if tokenTypeStr, ok := claims["token_type"].(string); ok {
		tokenClaims.TokenType = domain.TokenType(tokenTypeStr)
	}

	return tokenClaims, nil
}

func (s *JWTService) GenerateRefreshToken() (string, error) {
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return hex.EncodeToString(b), nil
}

func (s *JWTService) HashPassword(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hash), nil
}

func (s *JWTService) CheckPassword(hashedPassword, password string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
}
