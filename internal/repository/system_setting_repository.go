package repository

import (
	"github.com/google/uuid"
	"gorm.io/gorm"

	"github.com/unitechio/einfra-be/internal/domain"
)

// SystemSettingRepository implements the domain.SystemSettingRepository interface using GORM.
type SystemSettingRepository struct {
	DB *gorm.DB
}

// NewSystemSettingRepository creates a new instance of SystemSettingRepository.
func NewSystemSettingRepository(db *gorm.DB) *SystemSettingRepository {
	return &SystemSettingRepository{DB: db}
}

// Create adds a new system setting to the database.
func (r *SystemSettingRepository) Create(setting *domain.SystemSetting) (*domain.SystemSetting, error) {
	setting.ID = uuid.New().String()
	if err := r.DB.Create(setting).Error; err != nil {
		return nil, err
	}
	return setting, nil
}

// GetByKey retrieves a system setting by its key.
func (r *SystemSettingRepository) GetByKey(key string) (*domain.SystemSetting, error) {
	var setting domain.SystemSetting
	if err := r.DB.Where("`key` = ?", key).First(&setting).Error; err != nil {
		return nil, err
	}
	return &setting, nil
}

// GetAll retrieves all system settings from the database.
func (r *SystemSettingRepository) GetAll() ([]*domain.SystemSetting, error) {
	var settings []*domain.SystemSetting
	if err := r.DB.Find(&settings).Error; err != nil {
		return nil, err
	}
	return settings, nil
}

// GetByCategory retrieves all system settings of a specific category.
func (r *SystemSettingRepository) GetByCategory(category string) ([]*domain.SystemSetting, error) {
	var settings []*domain.SystemSetting
	if err := r.DB.Where("category = ?", category).Find(&settings).Error; err != nil {
		return nil, err
	}
	return settings, nil
}

// Update modifies an existing system setting in the database.
func (r *SystemSettingRepository) Update(setting *domain.SystemSetting) (*domain.SystemSetting, error) {
	if err := r.DB.Save(setting).Error; err != nil {
		return nil, err
	}
	return setting, nil
}

// Delete removes a system setting from the database by its ID.
func (r *SystemSettingRepository) Delete(id string) error {
	return r.DB.Delete(&domain.SystemSetting{}, "id = ?", id).Error
}
