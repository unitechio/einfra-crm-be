package repository

import (
	"context"

	"github.com/unitechio/einfra-be/internal/domain"
	"gorm.io/gorm"
)

type UserSettingsRepository interface {
	GetByUserID(ctx context.Context, userID string) (*domain.UserSettings, error)
	Update(ctx context.Context, userID string, settings *domain.UserSettings) error
}

type userSettingsRepository struct {
	db *gorm.DB
}

func NewUserSettingsRepository(db *gorm.DB) UserSettingsRepository {
	return &userSettingsRepository{db: db}
}

func (r *userSettingsRepository) GetByUserID(ctx context.Context, userID string) (*domain.UserSettings, error) {
	var settings domain.UserSettings
	if err := r.db.WithContext(ctx).Where("user_id = ?", userID).First(&settings).Error; err != nil {
		return nil, err
	}
	return &settings, nil
}

func (r *userSettingsRepository) Update(ctx context.Context, userID string, settings *domain.UserSettings) error {
	return r.db.WithContext(ctx).Model(&domain.UserSettings{}).Where("user_id = ?", userID).Updates(settings).Error
}
