package domain

import (
	"time"
)

// Environment represents a deployment environment in the infrastructure
// @Description Deployment environment (dev, staging, production, etc.)
type Environment struct {
	ID          string     `json:"id" gorm:"primaryKey;type:uuid;default:gen_random_uuid()" example:"550e8400-e29b-41d4-a716-446655440000"`
	Name        string     `json:"name" gorm:"type:varchar(50);not null;uniqueIndex" validate:"required,min=2,max=50,oneof=dev staging production qa uat" example:"production"`
	DisplayName string     `json:"display_name" gorm:"type:varchar(100);not null" validate:"required" example:"Production Environment"`
	Description string     `json:"description" gorm:"type:text" example:"Production environment for live services"`
	Color       string     `json:"color" gorm:"type:varchar(20)" example:"#FF0000"` // For UI display
	IsActive    bool       `json:"is_active" gorm:"type:boolean;default:true;index" example:"true"`
	SortOrder   int        `json:"sort_order" gorm:"type:int;default:0" example:"1"` // For ordering in UI
	CreatedAt   time.Time  `json:"created_at" gorm:"autoCreateTime" example:"2024-01-01T00:00:00Z"`
	UpdatedAt   time.Time  `json:"updated_at" gorm:"autoUpdateTime" example:"2024-01-01T00:00:00Z"`
	DeletedAt   *time.Time `json:"deleted_at,omitempty" gorm:"index" swaggertype:"string" example:"2024-01-01T00:00:00Z"`
}

// TableName specifies the table name for Environment model
func (Environment) TableName() string {
	return "environments"
}

// EnvironmentType represents standard environment types
type EnvironmentType string

const (
	// EnvironmentDev represents development environment
	EnvironmentDev EnvironmentType = "dev"
	// EnvironmentStaging represents staging environment
	EnvironmentStaging EnvironmentType = "staging"
	// EnvironmentProduction represents production environment
	EnvironmentProduction EnvironmentType = "production"
	// EnvironmentQA represents QA/testing environment
	EnvironmentQA EnvironmentType = "qa"
	// EnvironmentUAT represents user acceptance testing environment
	EnvironmentUAT EnvironmentType = "uat"
)

// EnvironmentFilter represents filtering options for environment queries
type EnvironmentFilter struct {
	IsActive *bool  `json:"is_active,omitempty"`
	Name     string `json:"name,omitempty"`
	Page     int    `json:"page" validate:"min=1"`
	PageSize int    `json:"page_size" validate:"min=1,max=100"`
}
