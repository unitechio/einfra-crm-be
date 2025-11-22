
package handler

import (
	"github.com/unitechio/einfra-be/internal/domain"
	"github.com/unitechio/einfra-be/internal/errorx"
	"github.com/unitechio/einfra-be/internal/logger"
	"github.com/unitechio/einfra-be/internal/usecase"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// NotificationHandler handles notification API requests.
type NotificationHandler struct {
	uc  usecase.NotificationUsecase
	log logger.Logger
}

// NewNotificationHandler creates a new notification handler.
func NewNotificationHandler(uc usecase.NotificationUsecase, log logger.Logger) *NotificationHandler {
	return &NotificationHandler{uc: uc, log: log}
}

// createNotificationRequest represents the request body for creating a notification.
type createNotificationRequest struct {
	Message string `json:"message" binding:"required"`
}

// GetAll godoc
// @Summary Get all notifications
// @Description Get all notifications
// @Tags notifications
// @Produce json
// @Success 200 {array} domain.Notification
// @Router /notifications [get]
func (h *NotificationHandler) GetAll(c *gin.Context) {
	h.log.Info(c.Request.Context(), "GetAll request received")
	notifications, err := h.uc.GetAll()
	if err != nil {
		h.log.Error(c.Request.Context(), "Failed to get notifications", logger.LogField{Key: "error", Value: err})
		c.Error(errorx.New(http.StatusInternalServerError, "Failed to get notifications").WithStack())
		return
	}
	h.log.Info(c.Request.Context(), "GetAll request successful")
	c.JSON(http.StatusOK, notifications)
}

// GetByID godoc
// @Summary Get a notification by ID
// @Description Get a notification by its ID
// @Tags notifications
// @Produce json
// @Param id path string true "Notification ID"
// @Success 200 {object} domain.Notification
// @Failure 404 {object} gin.H
// @Router /notifications/{id} [get]
func (h *NotificationHandler) GetByID(c *gin.Context) {
	id := c.Param("id")
	h.log.Info(c.Request.Context(), "GetByID request received", logger.LogField{Key: "id", Value: id})
	notification, err := h.uc.GetByID(id)
	if err != nil {
		h.log.Error(c.Request.Context(), "Notification not found", logger.LogField{Key: "id", Value: id}, logger.LogField{Key: "error", Value: err})
		c.Error(errorx.New(http.StatusNotFound, "Notification not found"))
		return
	}
	h.log.Info(c.Request.Context(), "GetByID request successful", logger.LogField{Key: "id", Value: id})
	c.JSON(http.StatusOK, notification)
}

// Create godoc
// @Summary Create a new notification
// @Description Create a new notification
// @Tags notifications
// @Accept json
// @Produce json
// @Param notification body createNotificationRequest true "Notification message"
// @Success 201 {object} domain.Notification
// @Failure 400 {object} gin.H
// @Failure 500 {object} gin.H
// @Router /notifications [post]
func (h *NotificationHandler) Create(c *gin.Context) {
	h.log.Info(c.Request.Context(), "Create request received")
	var req createNotificationRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		h.log.Error(c.Request.Context(), "Invalid request body", logger.LogField{Key: "error", Value: err})
		c.Error(errorx.New(http.StatusBadRequest, "Invalid request body"))
		return
	}

	notification := &domain.Notification{
		ID:        uuid.New().String(),
		Message:   req.Message,
		CreatedAt: time.Now(),
		Read:      false,
	}

	if err := h.uc.Create(c.Request.Context(), notification); err != nil {
		h.log.Error(c.Request.Context(), "Failed to create notification", logger.LogField{Key: "error", Value: err})
		c.Error(errorx.New(http.StatusInternalServerError, "Failed to create notification").WithStack())
		return
	}

	h.log.Info(c.Request.Context(), "Create request successful", logger.LogField{Key: "notification_id", Value: notification.ID})
	c.JSON(http.StatusCreated, notification)
}

// MarkAsRead godoc
// @Summary Mark a notification as read
// @Description Mark a notification as read by its ID
// @Tags notifications
// @Produce json
// @Param id path string true "Notification ID"
// @Success 200 {object} gin.H
// @Failure 500 {object} gin.H
// @Router /notifications/{id}/read [put]
func (h *NotificationHandler) MarkAsRead(c *gin.Context) {
	id := c.Param("id")
	h.log.Info(c.Request.Context(), "MarkAsRead request received", logger.LogField{Key: "id", Value: id})

	if err := h.uc.MarkAsRead(id); err != nil {
		h.log.Error(c.Request.Context(), "Failed to mark notification as read", logger.LogField{Key: "id", Value: id}, logger.LogField{Key: "error", Value: err})
		c.Error(errorx.New(http.StatusInternalServerError, "Failed to mark notification as read").WithStack())
		return
	}

	h.log.Info(c.Request.Context(), "MarkAsRead request successful", logger.LogField{Key: "id", Value: id})
	c.JSON(http.StatusOK, gin.H{"message": "Notification marked as read"})
}

// GetUnread godoc
// @Summary Get unread notifications
// @Description Get all unread notifications
// @Tags notifications
// @Produce json
// @Success 200 {array} domain.Notification
// @Failure 500 {object} gin.H
// @Router /notifications/unread [get]
func (h *NotificationHandler) GetUnread(c *gin.Context) {
	h.log.Info(c.Request.Context(), "GetUnread request received")
	notifications, err := h.uc.GetUnread()
	if err != nil {
		h.log.Error(c.Request.Context(), "Failed to get unread notifications", logger.LogField{Key: "error", Value: err})
		c.Error(errorx.New(http.StatusInternalServerError, "Failed to get unread notifications").WithStack())
		return
	}
	h.log.Info(c.Request.Context(), "GetUnread request successful")
	c.JSON(http.StatusOK, notifications)
}

// Delete godoc
// @Summary Delete a notification
// @Description Delete a notification by its ID
// @Tags notifications
// @Produce json
// @Param id path string true "Notification ID"
// @Success 200 {object} gin.H
// @Failure 500 {object} gin.H
// @Router /notifications/{id} [delete]
func (h *NotificationHandler) Delete(c *gin.Context) {
	id := c.Param("id")
	h.log.Info(c.Request.Context(), "Delete request received", logger.LogField{Key: "id", Value: id})

	if err := h.uc.Delete(id); err != nil {
		h.log.Error(c.Request.Context(), "Failed to delete notification", logger.LogField{Key: "id", Value: id}, logger.LogField{Key: "error", Value: err})
		c.Error(errorx.New(http.StatusInternalServerError, "Failed to delete notification").WithStack())
		return
	}

	h.log.Info(c.Request.Context(), "Delete request successful", logger.LogField{Key: "id", Value: id})
	c.JSON(http.StatusOK, gin.H{"message": "Notification deleted"})
}
