
package repository

import (
	"github.com/google/uuid"
	"gorm.io/gorm"

	"mymodule/internal/domain"
)

// PostgresFeatureFlagRepository implements the domain.FeatureFlagRepository interface using GORM.
type PostgresFeatureFlagRepository struct {
	DB *gorm.DB
}

// NewPostgresFeatureFlagRepository creates a new instance of PostgresFeatureFlagRepository.
func NewPostgresFeatureFlagRepository(db *gorm.DB) *PostgresFeatureFlagRepository {
	return &PostgresFeatureFlagRepository{DB: db}
}

// Create adds a new feature flag to the database.
func (r *PostgresFeatureFlagRepository) Create(flag *domain.FeatureFlag) (*domain.FeatureFlag, error) {
	flag.ID = uuid.New().String()
	if err := r.DB.Create(flag).Error; err != nil {
		return nil, err
	}
	return flag, nil
}

// GetByName retrieves a feature flag by its name.
func (r *PostgresFeatureFlagRepository) GetByName(name string) (*domain.FeatureFlag, error) {
	var flag domain.FeatureFlag
	if err := r.DB.Where("name = ?", name).First(&flag).Error; err != nil {
		return nil, err
	}
	return &flag, nil
}

// GetAll retrieves all feature flags from the database.
func (r *PostgresFeatureFlagRepository) GetAll() ([]*domain.FeatureFlag, error) {
	var flags []*domain.FeatureFlag
	if err := r.DB.Find(&flags).Error; err != nil {
		return nil, err
	}
	return flags, nil
}

// GetByCategory retrieves all feature flags of a specific category.
func (r *PostgresFeatureFlagRepository) GetByCategory(category string) ([]*domain.FeatureFlag, error) {
	var flags []*domain.FeatureFlag
	if err := r.DB.Where("category = ?", category).Find(&flags).Error; err != nil {
		return nil, err
	}
	return flags, nil
}

// Update modifies an existing feature flag in the database.
func (r *PostgresFeatureFlagRepository) Update(flag *domain.FeatureFlag) (*domain.FeatureFlag, error) {
	if err := r.DB.Save(flag).Error; err != nil {
		return nil, err
	}
	return flag, nil
}

// Delete removes a feature flag from the database by its ID.
func (r *PostgresFeatureFlagRepository) Delete(id string) error {
	return r.DB.Delete(&domain.FeatureFlag{}, "id = ?", id).Error
}
