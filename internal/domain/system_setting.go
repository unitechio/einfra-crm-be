
package domain

import "time"

// SystemSetting represents a single configuration setting for the system.
// These are key-value pairs that control various aspects of the application.
// swagger:model
type SystemSetting struct {
	ID          string    `json:"id" gorm:"primary_key"`
	Key         string    `json:"key" gorm:"uniqueIndex"` // e.g., "site_name", "maintenance_mode"
	Value       string    `json:"value"`                  // The value of the setting
	Category    string    `json:"category"`               // e.g., "General", "Display", "Email"
	Description string    `json:"description"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// SystemSettingRepository defines the interface for system setting persistence.
type SystemSettingRepository interface {
	Create(setting *SystemSetting) (*SystemSetting, error)
	GetByKey(key string) (*SystemSetting, error)
	GetAll() ([]*SystemSetting, error)
    GetByCategory(category string) ([]*SystemSetting, error)
	Update(setting *SystemSetting) (*SystemSetting, error)
	Delete(id string) error
}
