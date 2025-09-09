
package repository

import (
	"github.com/google/uuid"
	"gorm.io/gorm"

	"mymodule/internal/domain"
)

// PostgresSystemSettingRepository implements the domain.SystemSettingRepository interface using GORM.
type PostgresSystemSettingRepository struct {
	DB *gorm.DB
}

// NewPostgresSystemSettingRepository creates a new instance of PostgresSystemSettingRepository.
func NewPostgresSystemSettingRepository(db *gorm.DB) *PostgresSystemSettingRepository {
	return &PostgresSystemSettingRepository{DB: db}
}

// Create adds a new system setting to the database.
func (r *PostgresSystemSettingRepository) Create(setting *domain.SystemSetting) (*domain.SystemSetting, error) {
	setting.ID = uuid.New().String()
	if err := r.DB.Create(setting).Error; err != nil {
		return nil, err
	}
	return setting, nil
}

// GetByKey retrieves a system setting by its key.
func (r *PostgresSystemSettingRepository) GetByKey(key string) (*domain.SystemSetting, error) {
	var setting domain.SystemSetting
	if err := r.DB.Where("`key` = ?", key).First(&setting).Error; err != nil {
		return nil, err
	}
	return &setting, nil
}

// GetAll retrieves all system settings from the database.
func (r *PostgresSystemSettingRepository) GetAll() ([]*domain.SystemSetting, error) {
	var settings []*domain.SystemSetting
	if err := r.DB.Find(&settings).Error; err != nil {
		return nil, err
	}
	return settings, nil
}

// GetByCategory retrieves all system settings of a specific category.
func (r *PostgresSystemSettingRepository) GetByCategory(category string) ([]*domain.SystemSetting, error) {
	var settings []*domain.SystemSetting
	if err := r.DB.Where("category = ?", category).Find(&settings).Error; err != nil {
		return nil, err
	}
	return settings, nil
}


// Update modifies an existing system setting in the database.
func (r *PostgresSystemSettingRepository) Update(setting *domain.SystemSetting) (*domain.SystemSetting, error) {
	if err := r.DB.Save(setting).Error; err != nil {
		return nil, err
	}
	return setting, nil
}

// Delete removes a system setting from the database by its ID.
func (r *PostgresSystemSettingRepository) Delete(id string) error {
	return r.DB.Delete(&domain.SystemSetting{}, "id = ?", id).Error
}
