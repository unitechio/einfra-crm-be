package domain

import "time"

// FeatureFlag represents a feature that can be enabled or disabled in the system.
// This allows for granular control over the available functionality.
// swagger:model
type FeatureFlag struct {
	ID          string `json:"id" gorm:"primary_key"`
	Name        string `json:"name" gorm:"uniqueIndex"` // e.g., "Customer Management", "PDF Reports"
	Key         string `json:"key" gorm:"uniqueIndex"`  // e.g., "customer_management", "pdf_reports"
	Category    string `json:"category"`                // e.g., "Core Feature", "Reporting"
	Enabled     bool   `json:"enabled"`                 // Whether the feature is currently active globally
	Description string `json:"description"`

	// License Integration
	RequiredTier LicenseTier `json:"required_tier" gorm:"type:varchar(50);default:'free'"` // Minimum tier required
	IsPremium    bool        `json:"is_premium" gorm:"default:false"`                      // Is this a premium feature?

	// Usage Limits (0 = unlimited)
	MaxUsagePerMonth int `json:"max_usage_per_month" gorm:"default:0"` // For API-based features

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// FeatureFlagRepository defines the interface for feature flag persistence.
type FeatureFlagRepository interface {
	Create(flag *FeatureFlag) (*FeatureFlag, error)
	GetByName(name string) (*FeatureFlag, error)
	GetByKey(key string) (*FeatureFlag, error)
	GetAll() ([]*FeatureFlag, error)
	GetByCategory(category string) ([]*FeatureFlag, error)
	GetByTier(tier LicenseTier) ([]*FeatureFlag, error) // Get features available for a tier
	Update(flag *FeatureFlag) (*FeatureFlag, error)
	Delete(id string) error
}

// IsAvailableForTier checks if the feature is available for a given license tier
func (f *FeatureFlag) IsAvailableForTier(tier LicenseTier) bool {
	if !f.Enabled {
		return false
	}

	// Tier hierarchy: Free < Pro < Enterprise < Custom
	tierOrder := map[LicenseTier]int{
		TierFree:       1,
		TierPro:        2,
		TierEnterprise: 3,
		TierCustom:     4,
	}

	return tierOrder[tier] >= tierOrder[f.RequiredTier]
}

// Common feature keys
const (
	FeatureUserManagement     = "user_management"
	FeatureRoleManagement     = "role_management"
	FeatureDocumentManagement = "document_management"
	FeatureEmailNotifications = "email_notifications"
	FeatureAuditLogs          = "audit_logs"
	FeatureAPIAccess          = "api_access"
	FeatureAdvancedReporting  = "advanced_reporting"
	FeatureCustomIntegrations = "custom_integrations"
	FeatureWhiteLabeling      = "white_labeling"
	FeaturePrioritySupport    = "priority_support"
	FeatureCustomFeatures     = "custom_features"
)
