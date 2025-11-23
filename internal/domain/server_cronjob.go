package domain

import (
	"context"
	"time"
)

// CronjobStatus represents the status of a cronjob
type CronjobStatus string

const (
	// CronjobStatusActive indicates cronjob is active
	CronjobStatusActive CronjobStatus = "active"
	// CronjobStatusInactive indicates cronjob is inactive
	CronjobStatusInactive CronjobStatus = "inactive"
	// CronjobStatusFailed indicates cronjob has failed
	CronjobStatusFailed CronjobStatus = "failed"
)

// ServerCronjob represents a scheduled task on a server
// @Description Server cronjob entity with cron expression and execution history
type ServerCronjob struct {
	ID          string        `json:"id" gorm:"primaryKey;type:uuid;default:gen_random_uuid()" example:"550e8400-e29b-41d4-a716-446655440000"`
	ServerID    string        `json:"server_id" gorm:"type:uuid;not null;index" validate:"required" example:"550e8400-e29b-41d4-a716-446655440000"`
	Name        string        `json:"name" gorm:"type:varchar(255);not null" validate:"required" example:"daily-backup"`
	Description string        `json:"description" gorm:"type:text" example:"Daily backup task"`
	Status      CronjobStatus `json:"status" gorm:"type:varchar(50);not null" validate:"required" example:"active"`

	// Cron configuration
	CronExpression string `json:"cron_expression" gorm:"type:varchar(100);not null" validate:"required" example:"0 2 * * *"`
	Command        string `json:"command" gorm:"type:text;not null" validate:"required" example:"/usr/local/bin/backup.sh"`
	WorkingDir     string `json:"working_dir,omitempty" gorm:"type:varchar(500)" example:"/var/backups"`
	User           string `json:"user,omitempty" gorm:"type:varchar(100)" example:"root"`

	// Execution tracking
	LastRunAt      *time.Time `json:"last_run_at,omitempty" gorm:"type:timestamp" example:"2024-01-01T02:00:00Z"`
	NextRunAt      *time.Time `json:"next_run_at,omitempty" gorm:"type:timestamp" example:"2024-01-02T02:00:00Z"`
	LastExitCode   int        `json:"last_exit_code,omitempty" gorm:"type:int" example:"0"`
	LastOutput     string     `json:"last_output,omitempty" gorm:"type:text"`
	LastError      string     `json:"last_error,omitempty" gorm:"type:text"`
	ExecutionCount int        `json:"execution_count" gorm:"type:int;default:0" example:"100"`
	FailureCount   int        `json:"failure_count" gorm:"type:int;default:0" example:"2"`

	// Notifications
	NotifyOnFailure bool   `json:"notify_on_failure" gorm:"type:boolean;default:true" example:"true"`
	NotifyEmail     string `json:"notify_email,omitempty" gorm:"type:varchar(255)" example:"admin@example.com"`

	CreatedAt time.Time  `json:"created_at" gorm:"autoCreateTime" example:"2024-01-01T00:00:00Z"`
	UpdatedAt time.Time  `json:"updated_at" gorm:"autoUpdateTime" example:"2024-01-01T00:00:00Z"`
	DeletedAt *time.Time `json:"deleted_at,omitempty" gorm:"index" swaggertype:"string"`
}

// TableName specifies the table name for ServerCronjob model
func (ServerCronjob) TableName() string {
	return "server_cronjobs"
}

// CronjobExecution represents a single execution of a cronjob
// @Description Cronjob execution history record
type CronjobExecution struct {
	ID         string    `json:"id" gorm:"primaryKey;type:uuid;default:gen_random_uuid()" example:"550e8400-e29b-41d4-a716-446655440000"`
	CronjobID  string    `json:"cronjob_id" gorm:"type:uuid;not null;index" validate:"required" example:"550e8400-e29b-41d4-a716-446655440000"`
	StartedAt  time.Time `json:"started_at" gorm:"type:timestamp;not null" example:"2024-01-01T02:00:00Z"`
	FinishedAt time.Time `json:"finished_at" gorm:"type:timestamp" example:"2024-01-01T02:05:00Z"`
	ExitCode   int       `json:"exit_code" gorm:"type:int" example:"0"`
	Output     string    `json:"output,omitempty" gorm:"type:text"`
	Error      string    `json:"error,omitempty" gorm:"type:text"`
	Duration   int       `json:"duration" gorm:"type:int" example:"300"` // Seconds
	Success    bool      `json:"success" gorm:"type:boolean" example:"true"`
	CreatedAt  time.Time `json:"created_at" gorm:"autoCreateTime" example:"2024-01-01T02:00:00Z"`
}

// TableName specifies the table name for CronjobExecution model
func (CronjobExecution) TableName() string {
	return "cronjob_executions"
}

// CronjobFilter represents filtering options for cronjob queries
type CronjobFilter struct {
	ServerID string        `json:"server_id,omitempty"`
	Status   CronjobStatus `json:"status,omitempty"`
	Page     int           `json:"page" validate:"min=1"`
	PageSize int           `json:"page_size" validate:"min=1,max=100"`
}

// ServerCronjobRepository defines the interface for cronjob data persistence
type ServerCronjobRepository interface {
	// Create creates a new cronjob record
	Create(ctx context.Context, cronjob *ServerCronjob) error

	// GetByID retrieves a cronjob by its ID
	GetByID(ctx context.Context, id string) (*ServerCronjob, error)

	// List retrieves all cronjobs with pagination and filtering
	List(ctx context.Context, filter CronjobFilter) ([]*ServerCronjob, int64, error)

	// Update updates an existing cronjob
	Update(ctx context.Context, cronjob *ServerCronjob) error

	// Delete soft deletes a cronjob
	Delete(ctx context.Context, id string) error

	// GetByServerID retrieves all cronjobs for a server
	GetByServerID(ctx context.Context, serverID string) ([]*ServerCronjob, error)

	// CreateExecution creates a new execution record
	CreateExecution(ctx context.Context, execution *CronjobExecution) error

	// GetExecutions retrieves execution history for a cronjob
	GetExecutions(ctx context.Context, cronjobID string, limit int) ([]*CronjobExecution, error)
}

// ServerCronjobUsecase defines the business logic for cronjob management
type ServerCronjobUsecase interface {
	// CreateCronjob creates a new cronjob
	CreateCronjob(ctx context.Context, cronjob *ServerCronjob) error

	// GetCronjob retrieves a cronjob by ID
	GetCronjob(ctx context.Context, id string) (*ServerCronjob, error)

	// ListCronjobs retrieves cronjobs with filtering and pagination
	ListCronjobs(ctx context.Context, filter CronjobFilter) ([]*ServerCronjob, int64, error)

	// UpdateCronjob updates a cronjob
	UpdateCronjob(ctx context.Context, cronjob *ServerCronjob) error

	// DeleteCronjob deletes a cronjob
	DeleteCronjob(ctx context.Context, id string) error

	// ExecuteCronjob manually executes a cronjob
	ExecuteCronjob(ctx context.Context, cronjobID string) error

	// ValidateCronExpression validates a cron expression
	ValidateCronExpression(expression string) error

	// GetExecutionHistory retrieves execution history for a cronjob
	GetExecutionHistory(ctx context.Context, cronjobID string, limit int) ([]*CronjobExecution, error)
}
