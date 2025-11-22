package repository

import (
	"github.com/google/uuid"
	"gorm.io/gorm"

	"github.com/unitechio/einfra-be/internal/domain"
)

type FeatureFlagRepository struct {
	DB *gorm.DB
}

func NewFeatureFlagRepository(db *gorm.DB) *FeatureFlagRepository {
	return &FeatureFlagRepository{DB: db}
}

func (r *FeatureFlagRepository) Create(flag *domain.FeatureFlag) (*domain.FeatureFlag, error) {
	flag.ID = uuid.New().String()
	if err := r.DB.Create(flag).Error; err != nil {
		return nil, err
	}
	return flag, nil
}

func (r *FeatureFlagRepository) GetByName(name string) (*domain.FeatureFlag, error) {
	var flag domain.FeatureFlag
	if err := r.DB.Where("name = ?", name).First(&flag).Error; err != nil {
		return nil, err
	}
	return &flag, nil
}

func (r *FeatureFlagRepository) GetAll() ([]*domain.FeatureFlag, error) {
	var flags []*domain.FeatureFlag
	if err := r.DB.Find(&flags).Error; err != nil {
		return nil, err
	}
	return flags, nil
}

func (r *FeatureFlagRepository) GetByCategory(category string) ([]*domain.FeatureFlag, error) {
	var flags []*domain.FeatureFlag
	if err := r.DB.Where("category = ?", category).Find(&flags).Error; err != nil {
		return nil, err
	}
	return flags, nil
}

func (r *FeatureFlagRepository) Update(flag *domain.FeatureFlag) (*domain.FeatureFlag, error) {
	if err := r.DB.Save(flag).Error; err != nil {
		return nil, err
	}
	return flag, nil
}

func (r *FeatureFlagRepository) Delete(id string) error {
	return r.DB.Delete(&domain.FeatureFlag{}, "id = ?", id).Error
}
