
package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"mymodule/internal/domain"
	"mymodule/internal/errorx"
)

type AuditHandler struct {
	auditService domain.AuditService
}

func NewAuditHandler(auditService domain.AuditService) *AuditHandler {
	return &AuditHandler{
		auditService: auditService,
	}
}

func (h *AuditHandler) Log(c *gin.Context) {
	var entry domain.AuditEntry
	if err := c.ShouldBindJSON(&entry); err != nil {
		c.Error(errorx.New(http.StatusBadRequest, "Invalid request body"))
		return
	}

	newEntry, err := h.auditService.Log(c, entry)
	if err != nil {
		c.Error(errorx.New(http.StatusInternalServerError, "Failed to log audit entry").WithStack())
		return
	}

	c.JSON(http.StatusCreated, newEntry)
}

func (h *AuditHandler) GetAll(c *gin.Context) {
	audits, err := h.auditService.GetAll(c)
	if err != nil {
		c.Error(errorx.New(http.StatusInternalServerError, "Failed to get audit entries").WithStack())
		return
	}
	c.JSON(http.StatusOK, audits)
}

func (h *AuditHandler) GetByID(c *gin.Context) {
	id := c.Param("id")
	audit, err := h.auditService.GetByID(c, id)
	if err != nil {
		c.Error(errorx.New(http.StatusNotFound, "Audit entry not found"))
		return
	}
	c.JSON(http.StatusOK, audit)
}

func (h *AuditHandler) Update(c *gin.Context) {
	id := c.Param("id")
	var entry domain.AuditEntry
	if err := c.ShouldBindJSON(&entry); err != nil {
		c.Error(errorx.New(http.StatusBadRequest, "Invalid request body"))
		return
	}

	updatedEntry, err := h.auditService.Update(c, id, entry)
	if err != nil {
		c.Error(errorx.New(http.StatusInternalServerError, "Failed to update audit entry").WithStack())
		return
	}

	c.JSON(http.StatusOK, updatedEntry)
}

func (h *AuditHandler) Delete(c *gin.Context) {
	id := c.Param("id")
	if err := h.auditService.Delete(c, id); err != nil {
		c.Error(errorx.New(http.StatusInternalServerError, "Failed to delete audit entry").WithStack())
		return
	}

	c.JSON(http.StatusNoContent, nil)
}
