package domain

import (
	"context"
	"time"
)

// ServiceStatus represents the status of a service
type ServiceStatus string

const (
	// ServiceStatusRunning indicates service is running
	ServiceStatusRunning ServiceStatus = "running"
	// ServiceStatusStopped indicates service is stopped
	ServiceStatusStopped ServiceStatus = "stopped"
	// ServiceStatusFailed indicates service has failed
	ServiceStatusFailed ServiceStatus = "failed"
	// ServiceStatusUnknown indicates service status is unknown
	ServiceStatusUnknown ServiceStatus = "unknown"
)

// ServiceAction represents an action to perform on a service
type ServiceAction string

const (
	// ServiceActionStart starts a service
	ServiceActionStart ServiceAction = "start"
	// ServiceActionStop stops a service
	ServiceActionStop ServiceAction = "stop"
	// ServiceActionRestart restarts a service
	ServiceActionRestart ServiceAction = "restart"
	// ServiceActionReload reloads a service configuration
	ServiceActionReload ServiceAction = "reload"
	// ServiceActionEnable enables a service at boot
	ServiceActionEnable ServiceAction = "enable"
	// ServiceActionDisable disables a service at boot
	ServiceActionDisable ServiceAction = "disable"
)

// ServerService represents a service running on a server
// @Description Server service entity with status and configuration
type ServerService struct {
	ID          string        `json:"id" gorm:"primaryKey;type:uuid;default:gen_random_uuid()" example:"550e8400-e29b-41d4-a716-446655440000"`
	ServerID    string        `json:"server_id" gorm:"type:uuid;not null;index" validate:"required" example:"550e8400-e29b-41d4-a716-446655440000"`
	Name        string        `json:"name" gorm:"type:varchar(255);not null;index" validate:"required" example:"nginx"`
	DisplayName string        `json:"display_name" gorm:"type:varchar(255)" example:"Nginx Web Server"`
	Description string        `json:"description" gorm:"type:text" example:"High-performance HTTP server"`
	Status      ServiceStatus `json:"status" gorm:"type:varchar(50);not null" validate:"required" example:"running"`

	// Service details
	Enabled    bool   `json:"enabled" gorm:"type:boolean;default:false" example:"true"` // Auto-start on boot
	PID        int    `json:"pid,omitempty" gorm:"type:int" example:"1234"`
	Port       int    `json:"port,omitempty" gorm:"type:int" example:"80"`
	ConfigPath string `json:"config_path,omitempty" gorm:"type:varchar(500)" example:"/etc/nginx/nginx.conf"`
	LogPath    string `json:"log_path,omitempty" gorm:"type:varchar(500)" example:"/var/log/nginx/error.log"`

	// Monitoring
	MemoryUsageMB int       `json:"memory_usage_mb,omitempty" gorm:"type:int" example:"128"`
	CPUUsage      float64   `json:"cpu_usage,omitempty" gorm:"type:decimal(5,2)" example:"5.50"`
	Uptime        int64     `json:"uptime,omitempty" gorm:"type:bigint" example:"86400"` // Seconds
	LastCheckedAt time.Time `json:"last_checked_at" gorm:"type:timestamp" example:"2024-01-01T00:00:00Z"`

	CreatedAt time.Time  `json:"created_at" gorm:"autoCreateTime" example:"2024-01-01T00:00:00Z"`
	UpdatedAt time.Time  `json:"updated_at" gorm:"autoUpdateTime" example:"2024-01-01T00:00:00Z"`
	DeletedAt *time.Time `json:"deleted_at,omitempty" gorm:"index" swaggertype:"string"`
}

// TableName specifies the table name for ServerService model
func (ServerService) TableName() string {
	return "server_services"
}

// ServiceFilter represents filtering options for service queries
type ServiceFilter struct {
	ServerID string        `json:"server_id,omitempty"`
	Status   ServiceStatus `json:"status,omitempty"`
	Enabled  *bool         `json:"enabled,omitempty"`
	Page     int           `json:"page" validate:"min=1"`
	PageSize int           `json:"page_size" validate:"min=1,max=100"`
}

// ServiceActionRequest represents a request to perform an action on a service
type ServiceActionRequest struct {
	Action ServiceAction `json:"action" validate:"required,oneof=start stop restart reload enable disable" example:"restart"`
}

// ServerServiceRepository defines the interface for service data persistence
type ServerServiceRepository interface {
	// Create creates a new service record
	Create(ctx context.Context, service *ServerService) error

	// GetByID retrieves a service by its ID
	GetByID(ctx context.Context, id string) (*ServerService, error)

	// GetByServerAndName retrieves a service by server ID and service name
	GetByServerAndName(ctx context.Context, serverID, name string) (*ServerService, error)

	// List retrieves all services with pagination and filtering
	List(ctx context.Context, filter ServiceFilter) ([]*ServerService, int64, error)

	// Update updates an existing service
	Update(ctx context.Context, service *ServerService) error

	// Delete soft deletes a service
	Delete(ctx context.Context, id string) error

	// UpdateStatus updates only the status of a service
	UpdateStatus(ctx context.Context, id string, status ServiceStatus) error
}

// ServerServiceUsecase defines the business logic for service management
type ServerServiceUsecase interface {
	// ListServices retrieves all services running on a server
	ListServices(ctx context.Context, serverID string) ([]*ServerService, error)

	// GetService retrieves a service by ID
	GetService(ctx context.Context, id string) (*ServerService, error)

	// GetServiceStatus gets the current status of a service
	GetServiceStatus(ctx context.Context, serverID, serviceName string) (*ServerService, error)

	// PerformAction performs an action on a service (start, stop, restart, etc.)
	PerformAction(ctx context.Context, serverID, serviceName string, action ServiceAction) error

	// GetServiceLogs retrieves recent logs for a service
	GetServiceLogs(ctx context.Context, serverID, serviceName string, lines int) ([]string, error)

	// RefreshServices refreshes the list of services from the server
	RefreshServices(ctx context.Context, serverID string) error
}
