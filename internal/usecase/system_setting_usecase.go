
package usecase

import (
	"mymodule/internal/domain"
)

// SystemSettingUseCase provides the business logic for system settings.
type SystemSettingUseCase struct {
	repo domain.SystemSettingRepository
}

// NewSystemSettingUseCase creates a new instance of SystemSettingUseCase.
func NewSystemSettingUseCase(repo domain.SystemSettingRepository) *SystemSettingUseCase {
	return &SystemSettingUseCase{repo: repo}
}

// CreateSystemSetting creates a new system setting.
func (uc *SystemSettingUseCase) CreateSystemSetting(setting *domain.SystemSetting) (*domain.SystemSetting, error) {
	return uc.repo.Create(setting)
}

// GetSystemSettingByKey retrieves a system setting by its key.
func (uc *SystemSettingUseCase) GetSystemSettingByKey(key string) (*domain.SystemSetting, error) {
	return uc.repo.GetByKey(key)
}

// GetAllSystemSettings retrieves all system settings.
func (uc *SystemSettingUseCase) GetAllSystemSettings() ([]*domain.SystemSetting, error) {
	return uc.repo.GetAll()
}

// GetSystemSettingsByCategory retrieves all system settings of a specific category.
func (uc *SystemSettingUseCase) GetSystemSettingsByCategory(category string) ([]*domain.SystemSetting, error) {
	return uc.repo.GetByCategory(category)
}

// UpdateSystemSetting updates an existing system setting.
func (uc *SystemSettingUseCase) UpdateSystemSetting(setting *domain.SystemSetting) (*domain.SystemSetting, error) {
	return uc.repo.Update(setting)
}

// DeleteSystemSetting deletes a system setting by its ID.
func (uc *SystemSettingUseCase) DeleteSystemSetting(id string) error {
	return uc.repo.Delete(id)
}
