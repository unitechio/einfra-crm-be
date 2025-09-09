package service

import (
	"encoding/json"
	"log"
	"time"

	"github.com/gin-gonic/gin"
	"mymodule/internal/domain"
	"mymodule/internal/adapter/http/middleware"
)

// jsonAuditService is an implementation of the AuditService that logs to stdout as JSON.
type jsonAuditService struct{}

// NewJSONAuditService creates a new instance of the audit service.
func NewJSONAuditService() domain.AuditService {
	return &jsonAuditService{}
}

// Log records an audit event.
func (s *jsonAuditService) Log(c *gin.Context, entry domain.AuditEntry) (*domain.AuditEntry, error) {
	// Populate common fields from the context.
	entry.Timestamp = time.Now().UTC().Format(time.RFC3339)
	entry.ClientIP = c.ClientIP()

	// It's good practice to get user info from the context, which is set by the auth middleware.
	if userID, exists := c.Get("userID"); exists {
		if id, ok := userID.(string); ok {
			entry.UserID = id
		}
	}
	if userEmail, exists := c.Get("userEmail"); exists {
		if email, ok := userEmail.(string); ok {
			entry.UserEmail = email
		}
	}
    if role, exists := c.Get(middleware.UserRoleKey); exists {
        if r, ok := role.(string); ok {
            entry.UserRole = r
        }
    }

	// Use a JSON encoder to ensure the output is a single, valid JSON line.
	encoder := json.NewEncoder(log.Writer())
	return &entry, encoder.Encode(entry)
}
func (s *jsonAuditService) GetAll(c *gin.Context) ([]*domain.AuditEntry, error) {
	return nil, nil
}
func (s *jsonAuditService) GetByID(c *gin.Context, id string) (*domain.AuditEntry, error) {
	return nil, nil
}
func (s *jsonAuditService) Update(c *gin.Context, id string, entry domain.AuditEntry) (*domain.AuditEntry, error) {
	return nil, nil
}
func (s *jsonAuditService) Delete(c *gin.Context, id string) error {
	return nil
}
