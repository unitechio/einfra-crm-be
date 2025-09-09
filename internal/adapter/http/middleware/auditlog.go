package middleware

import (
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

func AuditLog(auditService domain.AuditService) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()

		action := getAuditAction(c.Request.Method)
		targetType, targetID := getTarget(c.Request.URL.Path)

		entry := domain.AuditEntry{
			UserID:     c.GetString("userID"), // Assuming userID is set by auth middleware
			UserEmail:  c.GetString("userEmail"), // Assuming userEmail is set by auth middleware
			UserRole:   c.GetString("userRole"),  // Assuming userRole is set by auth middleware
			Action:     action,
			TargetType: targetType,
			TargetID:   targetID,
			Timestamp:  time.Now(),
			ClientIP:   c.ClientIP(),
		}

		go auditService.Log(c.Copy(), entry) // Run in a goroutine to not block the request
	}
}

func getAuditAction(method string) domain.AuditAction {
	switch method {
	case "GET":
		return domain.ActionRead
	case "POST":
		return domain.ActionCreate
	case "PUT", "PATCH":
		return domain.ActionUpdate
	case "DELETE":
		return domain.ActionDelete
	default:
		return ""
	}
}

func getTarget(path string) (string, string) {
	parts := strings.Split(strings.Trim(path, "/"), "/")
	if len(parts) >= 3 {
		// Assuming path is like /api/v1/resource/id
		return parts[2], parts[3]
	}
	if len(parts) >= 2 {
		// Assuming path is like /api/v1/resource
		return parts[2], ""
	}
	return "", ""
}
