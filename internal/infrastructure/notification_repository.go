package infrastructure

import (
	"mymodule/internal/domain"

	"gorm.io/gorm"
)

// NewNotificationRepository creates a new notification repository.
func NewNotificationRepository(db *gorm.DB) domain.NotificationRepository {
	return &notificationRepository{db: db}
}

type notificationRepository struct {
	db *gorm.DB
}

// GetAll retrieves all notifications.
func (r *notificationRepository) GetAll() ([]*domain.Notification, error) {
	var notifications []*domain.Notification
	if err := r.db.Find(&notifications).Error; err != nil {
		return nil, err
	}
	return notifications, nil
}

// GetByID retrieves a notification by its ID.
func (r *notificationRepository) GetByID(id string) (*domain.Notification, error) {
	var notification domain.Notification
	if err := r.db.First(&notification, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &notification, nil
}

// Create creates a new notification.
func (r *notificationRepository) Create(notification *domain.Notification) error {
	return r.db.Create(notification).Error
}

// Update updates an existing notification.
func (r *notificationRepository) Update(notification *domain.Notification) error {
	return r.db.Save(notification).Error
}

// Delete deletes a notification by its ID.
func (r *notificationRepository) Delete(id string) error {
	return r.db.Delete(&domain.Notification{}, "id = ?", id).Error
}

// GetUnread retrieves all unread notifications.
func (r *notificationRepository) GetUnread() ([]*domain.Notification, error) {
	var notifications []*domain.Notification
	if err := r.db.Where("read = ?", false).Find(&notifications).Error; err != nil {
		return nil, err
	}
	return notifications, nil
}
