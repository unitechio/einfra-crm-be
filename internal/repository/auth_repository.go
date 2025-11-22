package repository

import (
	"context"
	"fmt"
	"time"

	"gorm.io/gorm"

	"github.com/unitechio/einfra-be/internal/auth"
	"github.com/unitechio/einfra-be/internal/domain"
)

type AuthRepository interface {
	HashPassword(password string) (string, error)
	ComparePassword(hashedPassword, password string) error
	CreateRefreshToken(ctx context.Context, token *domain.RefreshToken) error
	GetRefreshTokenByToken(ctx context.Context, token string) (*domain.RefreshToken, error)
	GetRefreshTokensByUserID(ctx context.Context, userID string) ([]*domain.RefreshToken, error)
	RevokeRefreshToken(ctx context.Context, tokenID string) error
	RevokeAllRefreshTokensForUser(ctx context.Context, userID string) error
	DeleteExpiredRefreshTokens(ctx context.Context) error
	GeneratePasswordResetToken(ctx context.Context, email string) (string, error)
	ValidatePasswordResetToken(ctx context.Context, token string) (string, error)
	GenerateEmailVerificationToken(ctx context.Context, userID string) (string, error)
	ValidateEmailVerificationToken(ctx context.Context, token string) (string, error)
}

type authRepository struct {
	db *gorm.DB
}

func NewAuthRepository(db *gorm.DB) AuthRepository {
	return &authRepository{db: db}
}

func (r *authRepository) HashPassword(password string) (string, error) {
	return auth.HashPassword(password)
}

func (r *authRepository) ComparePassword(hashedPassword, password string) error {
	return auth.CheckPassword(hashedPassword, password)
}

func (r *authRepository) CreateRefreshToken(ctx context.Context, token *domain.RefreshToken) error {
	return r.db.WithContext(ctx).Create(token).Error
}

func (r *authRepository) GetRefreshTokenByToken(ctx context.Context, token string) (*domain.RefreshToken, error) {
	var refreshToken domain.RefreshToken
	err := r.db.WithContext(ctx).
		First(&refreshToken, "token = ?", token).Error

	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("refresh token not found")
		}
		return nil, err
	}

	return &refreshToken, nil
}

func (r *authRepository) GetRefreshTokensByUserID(ctx context.Context, userID string) ([]*domain.RefreshToken, error) {
	var tokens []*domain.RefreshToken
	err := r.db.WithContext(ctx).
		Where("user_id = ? AND is_revoked = false", userID).
		Find(&tokens).Error

	if err != nil {
		return nil, err
	}

	return tokens, nil
}

func (r *authRepository) RevokeRefreshToken(ctx context.Context, tokenID string) error {
	now := time.Now()
	return r.db.WithContext(ctx).
		Model(&domain.RefreshToken{}).
		Where("id = ?", tokenID).
		Updates(map[string]interface{}{
			"is_revoked": true,
			"revoked_at": now,
		}).Error
}

func (r *authRepository) RevokeAllRefreshTokensForUser(ctx context.Context, userID string) error {
	now := time.Now()
	return r.db.WithContext(ctx).
		Model(&domain.RefreshToken{}).
		Where("user_id = ? AND is_revoked = false", userID).
		Updates(map[string]interface{}{
			"is_revoked": true,
			"revoked_at": now,
		}).Error
}

func (r *authRepository) DeleteExpiredRefreshTokens(ctx context.Context) error {
	return r.db.WithContext(ctx).
		Where("expires_at < ?", time.Now()).
		Delete(&domain.RefreshToken{}).Error
}

func (r *authRepository) GeneratePasswordResetToken(ctx context.Context, email string) (string, error) {
	return "", nil
}

func (r *authRepository) ValidatePasswordResetToken(ctx context.Context, token string) (string, error) {
	return "", nil
}

func (r *authRepository) GenerateEmailVerificationToken(ctx context.Context, userID string) (string, error) {
	return "", nil
}

func (r *authRepository) ValidateEmailVerificationToken(ctx context.Context, token string) (string, error) {
	return "", nil
}
