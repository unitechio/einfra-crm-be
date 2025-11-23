package domain

import "time"

// SystemSetting represents a system-wide configuration setting
type SystemSetting struct {
	ID          string     `json:"id" gorm:"type:varchar(36);primary_key"`
	Key         string     `json:"key" gorm:"type:varchar(255);uniqueIndex;not null"`
	Value       string     `json:"value" gorm:"type:text"`
	Category    string     `json:"category" gorm:"type:varchar(255);index"`
	Description string     `json:"description" gorm:"type:text"`
	CreatedAt   time.Time  `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt   time.Time  `json:"updated_at" gorm:"autoUpdateTime"`
	DeletedAt   *time.Time `json:"deleted_at,omitempty" gorm:"index"`
}

// SystemSettingRepository defines the interface for system setting data access
type SystemSettingRepository interface {
	Create(setting *SystemSetting) (*SystemSetting, error)
	GetByKey(key string) (*SystemSetting, error)
	GetAll() ([]*SystemSetting, error)
	GetByCategory(category string) ([]*SystemSetting, error)
	Update(setting *SystemSetting) (*SystemSetting, error)
	Delete(id string) error
}

// Common system setting categories
const (
	CategoryGeneral      = "general"
	CategoryEmail        = "email"
	CategorySecurity     = "security"
	CategoryNotification = "notification"
	CategoryIntegration  = "integration"
	CategoryMaintenance  = "maintenance"
)

// Common system setting keys
const (
	KeyMaintenanceMode    = "maintenance_mode"
	KeySiteName           = "site_name"
	KeySiteURL            = "site_url"
	KeyEmailFrom          = "email_from"
	KeyEmailSMTPHost      = "email_smtp_host"
	KeyEmailSMTPPort      = "email_smtp_port"
	KeySessionTimeout     = "session_timeout"
	KeyMaxLoginAttempts   = "max_login_attempts"
	KeyPasswordMinLength  = "password_min_length"
	KeyEnableRegistration = "enable_registration"
	KeyEnableTwoFactor    = "enable_two_factor"
)
