
package domain

import "time"

// FeatureFlag represents a feature that can be enabled or disabled in the system.
// This allows for granular control over the available functionality.
// swagger:model
type FeatureFlag struct {
	ID          string    `json:"id" gorm:"primary_key"`
	Name        string    `json:"name" gorm:"uniqueIndex"` // e.g., "Customer Management", "PDF Reports"
	Category    string    `json:"category"`               // e.g., "Core Feature", "Reporting"
	Enabled     bool      `json:"enabled"`                // Whether the feature is currently active
	Description string    `json:"description"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// FeatureFlagRepository defines the interface for feature flag persistence.
type FeatureFlagRepository interface {
	Create(flag *FeatureFlag) (*FeatureFlag, error)
	GetByName(name string) (*FeatureFlag, error)
	GetAll() ([]*FeatureFlag, error)
    GetByCategory(category string) ([]*FeatureFlag, error)
	Update(flag *FeatureFlag) (*FeatureFlag, error)
	Delete(id string) error
}
