package domain

import (
	"time"
)

// AuditAction represents the type of action performed
type AuditAction string

const (
	AuditActionCreate           AuditAction = "create"
	AuditActionUpdate           AuditAction = "update"
	AuditActionDelete           AuditAction = "delete"
	AuditActionGet              AuditAction = "get"
	AuditActionRead             AuditAction = "read"
	AuditActionLogin            AuditAction = "login"
	AuditActionLogout           AuditAction = "logout"
	AuditActionPasswordChange   AuditAction = "password_change"
	AuditActionPermissionChange AuditAction = "permission_change"
)

// AuditLog represents an audit log entry
// @Description Audit log for tracking user actions and system events
type AuditLog struct {
	ID            string      `json:"id" gorm:"primaryKey;type:uuid;default:gen_random_uuid()" example:"550e8400-e29b-41d4-a716-446655440000"`
	UserID        *string     `json:"user_id,omitempty" gorm:"type:uuid;index" example:"550e8400-e29b-41d4-a716-446655440000"`
	Username      string      `json:"username" gorm:"type:varchar(100);index" example:"john.doe"`
	Action        AuditAction `json:"action" gorm:"type:varchar(50);not null;index" validate:"required" example:"create"`
	Resource      string      `json:"resource" gorm:"type:varchar(100);not null;index" validate:"required" example:"server"`
	ResourceID    *string     `json:"resource_id,omitempty" gorm:"type:varchar(255);index" example:"550e8400-e29b-41d4-a716-446655440000"`
	Description   string      `json:"description" gorm:"type:text" example:"Created new server: web-server-01"`
	IPAddress     string      `json:"ip_address" gorm:"type:varchar(45);index" example:"192.168.1.100"`
	UserAgent     string      `json:"user_agent" gorm:"type:text" example:"Mozilla/5.0..."`
	RequestMethod string      `json:"request_method" gorm:"type:varchar(10)" example:"POST"`
	RequestPath   string      `json:"request_path" gorm:"type:varchar(500)" example:"/api/servers"`
	StatusCode    int         `json:"status_code" gorm:"type:int;index" example:"201"`
	Duration      int64       `json:"duration" gorm:"type:bigint" example:"150"` // milliseconds
	Changes       interface{} `json:"changes,omitempty" gorm:"type:jsonb"`       // Before/After data
	Metadata      interface{} `json:"metadata,omitempty" gorm:"type:jsonb"`      // Additional context
	Success       bool        `json:"success" gorm:"type:boolean;index" example:"true"`
	ErrorMessage  string      `json:"error_message,omitempty" gorm:"type:text"`
	CreatedAt     time.Time   `json:"created_at" gorm:"autoCreateTime;index" example:"2024-01-01T00:00:00Z"`
}

// AuditChange represents before/after data for audit logs
// @Description Before and after data for tracking changes
type AuditChange struct {
	Before interface{} `json:"before,omitempty"`
	After  interface{} `json:"after,omitempty"`
}

type AuditFilter struct {
	UserID     *string      `json:"user_id,omitempty"`
	Username   string       `json:"username,omitempty"`
	Action     *AuditAction `json:"action,omitempty"`
	Resource   string       `json:"resource,omitempty"`
	ResourceID string       `json:"resource_id,omitempty"`
	IPAddress  string       `json:"ip_address,omitempty"`
	Success    *bool        `json:"success,omitempty"`
	StartDate  *time.Time   `json:"start_date,omitempty"`
	EndDate    *time.Time   `json:"end_date,omitempty"`
	Page       int          `json:"page" validate:"min=1"`
	PageSize   int          `json:"page_size" validate:"min=1,max=100"`
	SortBy     string       `json:"sort_by,omitempty"`    // created_at, action, resource
	SortOrder  string       `json:"sort_order,omitempty"` // asc, desc
}

// AuditStatistics represents audit log statistics
// @Description Statistics about audit logs
type AuditStatistics struct {
	TotalLogs         int64                 `json:"total_logs" example:"1000"`
	SuccessfulActions int64                 `json:"successful_actions" example:"950"`
	FailedActions     int64                 `json:"failed_actions" example:"50"`
	UniqueUsers       int64                 `json:"unique_users" example:"25"`
	ActionBreakdown   map[AuditAction]int64 `json:"action_breakdown"`
	ResourceBreakdown map[string]int64      `json:"resource_breakdown"`
	HourlyActivity    []HourlyActivityStats `json:"hourly_activity"`
}

// HourlyActivityStats represents hourly activity statistics
type HourlyActivityStats struct {
	Hour  int   `json:"hour" example:"14"`
	Count int64 `json:"count" example:"50"`
}
