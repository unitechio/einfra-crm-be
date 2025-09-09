package handler

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"mymodule/internal/domain"
)

// PingHandler is a simple handler for health checks that also demonstrates auditing.
type PingHandler struct {
	auditService domain.AuditService
}

// NewPingHandler creates a new handler with its dependencies.
func NewPingHandler(as domain.AuditService) *PingHandler {
	return &PingHandler{auditService: as}
}

// Ping is a simple endpoint that responds with pong and logs an audit event.
func (h *PingHandler) Ping(c *gin.Context) {
    // For public endpoints, user info might not be available.
    // The audit service is designed to handle this gracefully.
	auditEntry := domain.AuditEntry{
		Action:     domain.ActionRead,
		TargetType: "SystemHealth",
		Details:    "Ping-pong health check was performed",
	}

	 if _, err := h.auditService.Log(c, auditEntry); err != nil {
		log.Printf("Failed to write audit log for ping: %v", err)
	}

	c.JSON(http.StatusOK, gin.H{"message": "pong"})
}
