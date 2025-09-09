package domain

import (
	"context"
	"time"

	"github.com/gin-gonic/gin"
)

// AuditAction represents the type of action being performed.
type AuditAction string

const (
	ActionCreate AuditAction = "CREATE"
	ActionUpdate AuditAction = "UPDATE"
	ActionDelete AuditAction = "DELETE"
	ActionLogin  AuditAction = "LOGIN"
	ActionRead   AuditAction = "READ"
)

// AuditEntry represents a single audit log record.
type AuditEntry struct {
	ID          string      `json:"id"`
	UserID      string      `json:"user_id"`
	UserEmail   string      `json:"user_email"`
	UserRole    string      `json:"user_role"`
	Action      AuditAction `json:"action"`
	TargetType  string      `json:"target_type"`  // e.g., "Product", "User"
	TargetID    string      `json:"target_id"`    // e.g., product ID, user ID
	Details     interface{} `json:"details,omitempty"` // Additional details (e.g., fields changed)
	Timestamp   time.Time   `json:"timestamp"`
	ClientIP    string      `json:"client_ip"`
}

// AuditService defines the interface for the audit logging service.
type AuditService interface {
	Log(c *gin.Context, entry AuditEntry) (AuditEntry, error)
	GetAll(c *gin.Context) ([]AuditEntry, error)
	GetByID(c *gin.Context, id string) (AuditEntry, error)
	Update(c *gin.Context, id string, entry AuditEntry) (AuditEntry, error)
	Delete(c *gin.Context, id string) error
}

// AuditRepository defines the interface for the audit repository.
type AuditRepository interface {
	Add(ctx context.Context, entry AuditEntry) (AuditEntry, error)
	GetAll(ctx context.Context) ([]AuditEntry, error)
	GetByID(ctx context.Context, id string) (AuditEntry, error)
	Update(ctx context.Context, id string, entry AuditEntry) (AuditEntry, error)
	Delete(ctx context.Context, id string) error
}
