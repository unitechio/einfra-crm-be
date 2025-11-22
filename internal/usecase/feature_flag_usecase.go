package usecase

import (
	"github.com/unitechio/einfra-be/internal/domain"
)

// FeatureFlagUseCase provides the business logic for feature flags.
type FeatureFlagUseCase struct {
	repo domain.FeatureFlagRepository
}

// NewFeatureFlagUseCase creates a new instance of FeatureFlagUseCase.
func NewFeatureFlagUseCase(repo domain.FeatureFlagRepository) *FeatureFlagUseCase {
	return &FeatureFlagUseCase{repo: repo}
}

// CreateFeatureFlag creates a new feature flag.
func (uc *FeatureFlagUseCase) CreateFeatureFlag(flag *domain.FeatureFlag) (*domain.FeatureFlag, error) {
	return uc.repo.Create(flag)
}

// GetFeatureFlagByName retrieves a feature flag by its name.
func (uc *FeatureFlagUseCase) GetFeatureFlagByName(name string) (*domain.FeatureFlag, error) {
	return uc.repo.GetByName(name)
}

// GetAllFeatureFlags retrieves all feature flags.
func (uc *FeatureFlagUseCase) GetAllFeatureFlags() ([]*domain.FeatureFlag, error) {
	return uc.repo.GetAll()
}

// GetFeatureFlagsByCategory retrieves all feature flags of a specific category.
func (uc *FeatureFlagUseCase) GetFeatureFlagsByCategory(category string) ([]*domain.FeatureFlag, error) {
	return uc.repo.GetByCategory(category)
}

// UpdateFeatureFlag updates an existing feature flag.
func (uc *FeatureFlagUseCase) UpdateFeatureFlag(flag *domain.FeatureFlag) (*domain.FeatureFlag, error) {
	return uc.repo.Update(flag)
}

// DeleteFeatureFlag deletes a feature flag by its ID.
func (uc *FeatureFlagUseCase) DeleteFeatureFlag(id string) error {
	return uc.repo.Delete(id)
}

// EnableFeature enables a feature flag.
func (uc *FeatureFlagUseCase) EnableFeature(name string) (*domain.FeatureFlag, error) {
    flag, err := uc.repo.GetByName(name)
    if err != nil {
        return nil, err
    }
    flag.Enabled = true
    return uc.repo.Update(flag)
}

// DisableFeature disables a feature flag.
func (uc *FeatureFlagUseCase) DisableFeature(name string) (*domain.FeatureFlag, error) {
    flag, err := uc.repo.GetByName(name)
    if err != nil {
        return nil, err
    }
    flag.Enabled = false
    return uc.repo.Update(flag)
}

// IsFeatureEnabled checks if a feature is enabled.
func (uc *FeatureFlagUseCase) IsFeatureEnabled(name string) (bool, error) {
    flag, err := uc.repo.GetByName(name)
    if err != nil {
        return false, err
    }
    return flag.Enabled, nil
}
