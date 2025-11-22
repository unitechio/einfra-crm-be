package socket

import (
	"time"

	"github.com/unitechio/einfra-be/internal/domain"
)

// NotificationMessage represents a notification message for WebSocket
type NotificationMessage struct {
	ID          string                      `json:"id"`
	Type        domain.NotificationType     `json:"type"`
	Priority    domain.NotificationPriority `json:"priority"`
	Title       string                      `json:"title"`
	Message     string                      `json:"message"`
	Data        interface{}                 `json:"data,omitempty"`
	ActionURL   string                      `json:"action_url,omitempty"`
	ActionLabel string                      `json:"action_label,omitempty"`
	Icon        string                      `json:"icon,omitempty"`
	CreatedAt   time.Time                   `json:"created_at"`
}

// NewNotificationMessage creates a Message from a domain.Notification
func NewNotificationMessage(notification *domain.Notification) Message {
	return Message{
		Type: MessageTypeNotification,
		Data: NotificationMessage{
			ID:          notification.ID,
			Type:        notification.Type,
			Priority:    notification.Priority,
			Title:       notification.Title,
			Message:     notification.Message,
			Data:        notification.Data,
			ActionURL:   notification.ActionURL,
			ActionLabel: notification.ActionLabel,
			Icon:        notification.Icon,
			CreatedAt:   notification.CreatedAt,
		},
		Timestamp: time.Now(),
	}
}

// SystemMessage represents a system message
type SystemMessage struct {
	Code    string      `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

// NewSystemMessage creates a system Message
func NewSystemMessage(code, message string, data interface{}) Message {
	return Message{
		Type: MessageTypeSystem,
		Data: SystemMessage{
			Code:    code,
			Message: message,
			Data:    data,
		},
		Timestamp: time.Now(),
	}
}

// ErrorMessage represents an error message
type ErrorMessage struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

// NewErrorMessage creates an error Message
func NewErrorMessage(code, message string) Message {
	return Message{
		Type: MessageTypeError,
		Data: ErrorMessage{
			Code:    code,
			Message: message,
		},
		Timestamp: time.Now(),
	}
}
