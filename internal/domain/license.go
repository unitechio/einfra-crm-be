package domain

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"time"
)

// LicenseTier represents the tier/plan level of a license
type LicenseTier string

const (
	TierFree       LicenseTier = "free"
	TierPro        LicenseTier = "professional"
	TierEnterprise LicenseTier = "enterprise"
	TierCustom     LicenseTier = "custom"
)

// LicenseStatus represents the current status of a license
type LicenseStatus string

const (
	LicenseStatusActive    LicenseStatus = "active"
	LicenseStatusExpired   LicenseStatus = "expired"
	LicenseStatusSuspended LicenseStatus = "suspended"
	LicenseStatusRevoked   LicenseStatus = "revoked"
)

// License represents a software license for the system
type License struct {
	ID         string        `json:"id" gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	LicenseKey string        `json:"license_key" gorm:"type:varchar(512);uniqueIndex;not null"`
	Tier       LicenseTier   `json:"tier" gorm:"type:varchar(50);not null;index"`
	Status     LicenseStatus `json:"status" gorm:"type:varchar(50);not null;default:'active'"`

	// Organization/Customer Info
	OrganizationID   string `json:"organization_id" gorm:"type:uuid;index"`
	OrganizationName string `json:"organization_name" gorm:"type:varchar(255)"`
	ContactEmail     string `json:"contact_email" gorm:"type:varchar(255)"`

	// License Limits
	MaxUsers   int `json:"max_users" gorm:"default:0"`     // 0 = unlimited
	MaxAPICall int `json:"max_api_calls" gorm:"default:0"` // per month, 0 = unlimited
	MaxStorage int `json:"max_storage" gorm:"default:0"`   // in GB, 0 = unlimited

	// Usage Tracking
	CurrentUsers    int `json:"current_users" gorm:"default:0"`
	CurrentAPICalls int `json:"current_api_calls" gorm:"default:0"`
	CurrentStorage  int `json:"current_storage" gorm:"default:0"`

	// Time Management
	IssuedAt    time.Time  `json:"issued_at" gorm:"not null"`
	ExpiresAt   *time.Time `json:"expires_at,omitempty"` // nil = perpetual
	ActivatedAt *time.Time `json:"activated_at,omitempty"`
	SuspendedAt *time.Time `json:"suspended_at,omitempty"`

	// Metadata
	Metadata string `json:"metadata" gorm:"type:jsonb"` // Additional custom data
	Notes    string `json:"notes" gorm:"type:text"`

	// Audit
	CreatedAt time.Time  `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt time.Time  `json:"updated_at" gorm:"autoUpdateTime"`
	DeletedAt *time.Time `json:"deleted_at,omitempty" gorm:"index"`
}

// LicenseUsageLog tracks API usage and other metrics
type LicenseUsageLog struct {
	ID         string    `json:"id" gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	LicenseID  string    `json:"license_id" gorm:"type:uuid;not null;index"`
	UsageType  string    `json:"usage_type" gorm:"type:varchar(50);not null"` // api_call, storage, user_login
	Count      int       `json:"count" gorm:"default:1"`
	Metadata   string    `json:"metadata" gorm:"type:jsonb"`
	RecordedAt time.Time `json:"recorded_at" gorm:"autoCreateTime;index"`
}

// LicenseRepository defines the interface for license persistence
type LicenseRepository interface {
	Create(license *License) (*License, error)
	GetByID(id string) (*License, error)
	GetByKey(key string) (*License, error)
	GetByOrganization(orgID string) (*License, error)
	GetAll() ([]*License, error)
	GetByTier(tier LicenseTier) ([]*License, error)
	GetByStatus(status LicenseStatus) ([]*License, error)
	Update(license *License) (*License, error)
	Delete(id string) error

	// Usage tracking
	LogUsage(log *LicenseUsageLog) error
	GetUsageStats(licenseID string, from, to time.Time) ([]*LicenseUsageLog, error)
	ResetMonthlyUsage(licenseID string) error
}

// IsValid checks if the license is currently valid
func (l *License) IsValid() bool {
	if l.Status != LicenseStatusActive {
		return false
	}

	// Check expiry
	if l.ExpiresAt != nil && l.ExpiresAt.Before(time.Now()) {
		return false
	}

	return true
}

// IsExpired checks if the license has expired
func (l *License) IsExpired() bool {
	if l.ExpiresAt == nil {
		return false // Perpetual license
	}
	return l.ExpiresAt.Before(time.Now())
}

// CanAddUser checks if the license allows adding more users
func (l *License) CanAddUser() bool {
	if l.MaxUsers == 0 {
		return true // Unlimited
	}
	return l.CurrentUsers < l.MaxUsers
}

// CanMakeAPICall checks if the license allows more API calls
func (l *License) CanMakeAPICall() bool {
	if l.MaxAPICall == 0 {
		return true // Unlimited
	}
	return l.CurrentAPICalls < l.MaxAPICall
}

// GetTierLimits returns the default limits for a tier
func GetTierLimits(tier LicenseTier) (maxUsers, maxAPICalls, maxStorage int) {
	switch tier {
	case TierFree:
		return 5, 1000, 1 // 5 users, 1k API calls/month, 1GB storage
	case TierPro:
		return 50, 50000, 50 // 50 users, 50k API calls/month, 50GB storage
	case TierEnterprise:
		return 0, 0, 0 // Unlimited
	case TierCustom:
		return 0, 0, 0 // Will be set manually
	default:
		return 5, 1000, 1
	}
}

// GenerateLicenseKey generates a random license key
func GenerateLicenseKey() (string, error) {
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}

	key := base64.URLEncoding.EncodeToString(b)

	// Format: EINFRA-XXXX-XXXX-XXXX-XXXX
	formatted := fmt.Sprintf("EINFRA-%s-%s-%s-%s",
		key[0:4], key[4:8], key[8:12], key[12:16])

	return formatted, nil
}

// LicenseActivationRequest represents a request to activate a license
type LicenseActivationRequest struct {
	LicenseKey       string `json:"license_key" binding:"required"`
	OrganizationID   string `json:"organization_id" binding:"required"`
	OrganizationName string `json:"organization_name" binding:"required"`
	ContactEmail     string `json:"contact_email" binding:"required,email"`
}

// LicenseValidationResponse represents the response of license validation
type LicenseValidationResponse struct {
	Valid     bool          `json:"valid"`
	License   *License      `json:"license,omitempty"`
	Tier      LicenseTier   `json:"tier"`
	Status    LicenseStatus `json:"status"`
	ExpiresAt *time.Time    `json:"expires_at,omitempty"`
	DaysLeft  int           `json:"days_left,omitempty"`
	Message   string        `json:"message,omitempty"`
	Limits    LicenseLimits `json:"limits"`
}

// LicenseLimits represents the current usage vs limits
type LicenseLimits struct {
	MaxUsers        int `json:"max_users"`
	CurrentUsers    int `json:"current_users"`
	MaxAPICalls     int `json:"max_api_calls"`
	CurrentAPICalls int `json:"current_api_calls"`
	MaxStorage      int `json:"max_storage"`
	CurrentStorage  int `json:"current_storage"`
}
