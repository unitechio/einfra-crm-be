package usecase

import (
	"context"

	"github.com/unitechio/einfra-be/internal/domain"
	"github.com/unitechio/einfra-be/internal/repository"
)

type UserSettingsUseCase interface {
	GetUserSettings(ctx context.Context, userID string) (*domain.UserSettings, error)
	UpdateUserSettings(ctx context.Context, userID string, settings *domain.UserSettings) error
}

type userSettingsUsecase struct {
	userSettingsRepo repository.UserSettingsRepository
}

func NewUserSettingsUseCase(userSettingsRepo repository.UserSettingsRepository) UserSettingsUseCase {
	return &userSettingsUsecase{
		userSettingsRepo: userSettingsRepo,
	}
}

func (u *userSettingsUsecase) GetUserSettings(ctx context.Context, userID string) (*domain.UserSettings, error) {
	return u.userSettingsRepo.GetByUserID(ctx, userID)
}

func (u *userSettingsUsecase) UpdateUserSettings(ctx context.Context, userID string, settings *domain.UserSettings) error {
	return u.userSettingsRepo.Update(ctx, userID, settings)
}
