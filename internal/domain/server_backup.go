package domain

import (
	"context"
	"time"
)

// BackupType represents the type of backup
type BackupType string

const (
	// BackupTypeFull represents a full backup
	BackupTypeFull BackupType = "full"
	// BackupTypeIncremental represents an incremental backup
	BackupTypeIncremental BackupType = "incremental"
	// BackupTypeDifferential represents a differential backup
	BackupTypeDifferential BackupType = "differential"
)

// BackupStatus represents the status of a backup
type BackupStatus string

const (
	// BackupStatusPending indicates backup is queued
	BackupStatusPending BackupStatus = "pending"
	// BackupStatusInProgress indicates backup is running
	BackupStatusInProgress BackupStatus = "in_progress"
	// BackupStatusCompleted indicates backup completed successfully
	BackupStatusCompleted BackupStatus = "completed"
	// BackupStatusFailed indicates backup failed
	BackupStatusFailed BackupStatus = "failed"
)

// ServerBackup represents a backup of a server
// @Description Server backup entity with metadata and status
type ServerBackup struct {
	ID          string       `json:"id" gorm:"primaryKey;type:uuid;default:gen_random_uuid()" example:"550e8400-e29b-41d4-a716-446655440000"`
	ServerID    string       `json:"server_id" gorm:"type:uuid;not null;index" validate:"required" example:"550e8400-e29b-41d4-a716-446655440000"`
	Name        string       `json:"name" gorm:"type:varchar(255);not null" validate:"required" example:"daily-backup-2024-01-01"`
	Description string       `json:"description" gorm:"type:text" example:"Daily automated backup"`
	Type        BackupType   `json:"type" gorm:"type:varchar(50);not null" validate:"required,oneof=full incremental differential" example:"full"`
	Status      BackupStatus `json:"status" gorm:"type:varchar(50);not null;index" validate:"required" example:"completed"`

	// Backup details
	BackupPath string  `json:"backup_path" gorm:"type:varchar(500)" example:"/backups/server-01/2024-01-01-full.tar.gz"`
	SizeBytes  int64   `json:"size_bytes" gorm:"type:bigint" example:"1073741824"`
	SizeGB     float64 `json:"size_gb" gorm:"-" example:"1.00"` // Calculated field
	Compressed bool    `json:"compressed" gorm:"type:boolean;default:true" example:"true"`
	Encrypted  bool    `json:"encrypted" gorm:"type:boolean;default:false" example:"false"`

	// Timing
	StartedAt   *time.Time `json:"started_at,omitempty" gorm:"type:timestamp" example:"2024-01-01T00:00:00Z"`
	CompletedAt *time.Time `json:"completed_at,omitempty" gorm:"type:timestamp" example:"2024-01-01T01:00:00Z"`

	// Error tracking
	ErrorMessage string `json:"error_message,omitempty" gorm:"type:text"`

	// Retention
	ExpiresAt *time.Time `json:"expires_at,omitempty" gorm:"type:timestamp;index" example:"2024-02-01T00:00:00Z"`

	CreatedAt time.Time  `json:"created_at" gorm:"autoCreateTime" example:"2024-01-01T00:00:00Z"`
	UpdatedAt time.Time  `json:"updated_at" gorm:"autoUpdateTime" example:"2024-01-01T00:00:00Z"`
	DeletedAt *time.Time `json:"deleted_at,omitempty" gorm:"index" swaggertype:"string"`
}

// TableName specifies the table name for ServerBackup model
func (ServerBackup) TableName() string {
	return "server_backups"
}

// BackupFilter represents filtering options for backup queries
type BackupFilter struct {
	ServerID string       `json:"server_id,omitempty"`
	Type     BackupType   `json:"type,omitempty"`
	Status   BackupStatus `json:"status,omitempty"`
	Page     int          `json:"page" validate:"min=1"`
	PageSize int          `json:"page_size" validate:"min=1,max=100"`
}

// ServerBackupRepository defines the interface for backup data persistence
type ServerBackupRepository interface {
	// Create creates a new backup record
	Create(ctx context.Context, backup *ServerBackup) error

	// GetByID retrieves a backup by its ID
	GetByID(ctx context.Context, id string) (*ServerBackup, error)

	// List retrieves all backups with pagination and filtering
	List(ctx context.Context, filter BackupFilter) ([]*ServerBackup, int64, error)

	// Update updates an existing backup
	Update(ctx context.Context, backup *ServerBackup) error

	// Delete soft deletes a backup
	Delete(ctx context.Context, id string) error

	// DeleteExpired deletes expired backups
	DeleteExpired(ctx context.Context) (int64, error)

	// GetByServerID retrieves all backups for a server
	GetByServerID(ctx context.Context, serverID string) ([]*ServerBackup, error)
}

// ServerBackupUsecase defines the business logic for backup management
type ServerBackupUsecase interface {
	// CreateBackup creates a new backup
	CreateBackup(ctx context.Context, backup *ServerBackup) error

	// GetBackup retrieves a backup by ID
	GetBackup(ctx context.Context, id string) (*ServerBackup, error)

	// ListBackups retrieves backups with filtering and pagination
	ListBackups(ctx context.Context, filter BackupFilter) ([]*ServerBackup, int64, error)

	// RestoreBackup restores a server from a backup
	RestoreBackup(ctx context.Context, backupID string) error

	// DeleteBackup deletes a backup
	DeleteBackup(ctx context.Context, id string) error

	// CleanupExpiredBackups removes expired backups
	CleanupExpiredBackups(ctx context.Context) (int64, error)

	// GetBackupStatus gets the current status of a backup
	GetBackupStatus(ctx context.Context, backupID string) (*ServerBackup, error)
}
