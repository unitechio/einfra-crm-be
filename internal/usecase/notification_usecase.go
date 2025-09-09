
package usecase

import (
	"context"
	"encoding/json"
	"mymodule/internal/domain"
	"mymodule/internal/logger"
	"mymodule/internal/messaging"
	"mymodule/internal/realtime"
)

// NotificationUsecase handles notification business logic.
type NotificationUsecase interface {
	GetAll() ([]*domain.Notification, error)
	GetByID(id string) (*domain.Notification, error)
	Create(ctx context.Context, notification *domain.Notification) error
	MarkAsRead(id string) error
	GetUnread() ([]*domain.Notification, error)
	Delete(id string) error
}

// NewNotificationUsecase creates a new notification usecase.
func NewNotificationUsecase(repo domain.NotificationRepository, hub *realtime.Hub, publisher messaging.Publisher, log logger.Logger) NotificationUsecase {
	return &notificationUsecase{repo: repo, hub: hub, publisher: publisher, log: log}
}

type notificationUsecase struct {
	repo      domain.NotificationRepository
	hub       *realtime.Hub
	publisher messaging.Publisher
	log       logger.Logger
}

// GetAll retrieves all notifications.
func (uc *notificationUsecase) GetAll() ([]*domain.Notification, error) {
	uc.log.Info(context.Background(), "GetAll usecase called")
	return uc.repo.GetAll()
}

// GetByID retrieves a notification by its ID.
func (uc *notificationUsecase) GetByID(id string) (*domain.Notification, error) {
	uc.log.Info(context.Background(), "GetByID usecase called", logger.LogField{Key: "id", Value: id})
	return uc.repo.GetByID(id)
}

// Create creates a new notification, broadcasts it, and publishes it to the message broker.
func (uc *notificationUsecase) Create(ctx context.Context, notification *domain.Notification) error {
	uc.log.Info(ctx, "Create usecase called", logger.LogField{Key: "notification_id", Value: notification.ID})
	if err := uc.repo.Create(notification); err != nil {
		uc.log.Error(ctx, "Failed to create notification in repo", logger.LogField{Key: "error", Value: err})
		return err
	}

	// Send to WebSocket clients
	msgBytes, err := json.Marshal(notification)
	if err != nil {
		uc.log.Error(ctx, "Error marshalling notification for websocket", logger.LogField{Key: "error", Value: err})
	} else {
		uc.hub.Broadcast(msgBytes)
	}

	// Publish event to the message broker
	msg := messaging.Message{
		ID:      notification.ID,
		Topic:   "notifications", // This could be configured
		Payload: msgBytes,
		Headers: make(map[string]string),
	}
	if err := uc.publisher.Publish(ctx, msg); err != nil {
		// Log the error but don't fail the entire request.
		// The notification is saved in the DB and sent to websockets.
		// The system should be resilient to the message broker being temporarily down.
		uc.log.Error(ctx, "Error publishing notification event", logger.LogField{Key: "error", Value: err})
	}

	return nil
}

// MarkAsRead marks a notification as read.
func (uc *notificationUsecase) MarkAsRead(id string) error {
	uc.log.Info(context.Background(), "MarkAsRead usecase called", logger.LogField{Key: "id", Value: id})
	notification, err := uc.repo.GetByID(id)
	if err != nil {
		uc.log.Error(context.Background(), "Failed to get notification by ID in MarkAsRead", logger.LogField{Key: "id", Value: id}, logger.LogField{Key: "error", Value: err})
		return err
	}

	notification.Read = true
	return uc.repo.Update(notification)
}

// GetUnread retrieves all unread notifications.
func (uc *notificationUsecase) GetUnread() ([]*domain.Notification, error) {
	uc.log.Info(context.Background(), "GetUnread usecase called")
	return uc.repo.GetUnread()
}

// Delete deletes a notification.
func (uc *notificationUsecase) Delete(id string) error {
	uc.log.Info(context.Background(), "Delete usecase called", logger.LogField{Key: "id", Value: id})
	return uc.repo.Delete(id)
}
