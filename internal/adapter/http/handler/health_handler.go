
package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// HealthHandler handles health check API requests.
type HealthHandler struct{}

// NewHealthHandler creates a new health handler.
func NewHealthHandler() *HealthHandler {
	return &HealthHandler{}
}

// HealthCheck godoc
// @Summary Check the health of the service
// @Description Returns the status of the service
// @Tags health
// @Produce json
// @Success 200 {object} gin.H
// @Router /health [get]
func (h *HealthHandler) HealthCheck(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status": "ok",
	})
}
