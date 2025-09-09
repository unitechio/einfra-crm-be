package domain

import "time"

// Notification represents a notification message.
// swagger:model
type Notification struct {
	ID        string    `json:"id" gorm:"primary_key"`
	Message   string    `json:"message"`
	CreatedAt time.Time `json:"created_at"`
	Read      bool      `json:"read"`
}

// NotificationRepository defines the interface for notification persistence.
// swagger:model
type NotificationRepository interface {
	GetAll() ([]*Notification, error)
	GetByID(id string) (*Notification, error)
	Create(notification *Notification) error
	Update(notification *Notification) error
	Delete(id string) error
	GetUnread() ([]*Notification, error)
}
