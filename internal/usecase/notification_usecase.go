package usecase

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/unitechio/einfra-be/internal/domain"
	"github.com/unitechio/einfra-be/internal/logger"
	"github.com/unitechio/einfra-be/internal/realtime"
	"github.com/unitechio/einfra-be/internal/repository"
)

type NotificationUsecase interface {
	SendNotification(ctx context.Context, notification *domain.Notification) error
	SendNotificationFromTemplate(ctx context.Context, userID, templateName string, variables map[string]string) error
	SendBulkNotification(ctx context.Context, userIDs []string, notification *domain.Notification) error
	GetNotification(ctx context.Context, id string) (*domain.Notification, error)
	GetUserNotifications(ctx context.Context, userID string, filter domain.NotificationFilter) ([]*domain.Notification, int64, error)
	GetUnreadCount(ctx context.Context, userID string) (int64, error)
	MarkAsRead(ctx context.Context, id string) error
	MarkAllAsRead(ctx context.Context, userID string) error
	DeleteNotification(ctx context.Context, id string) error
	CleanupOldNotifications(ctx context.Context, retentionDays int) error
	GetUserPreferences(ctx context.Context, userID string) (*domain.NotificationPreference, error)
	UpdateUserPreferences(ctx context.Context, userID string, preferences *domain.NotificationPreference) error
}

type notificationUsecase struct {
	repo         repository.NotificationRepository
	templateRepo domain.NotificationTemplateRepository
	prefRepo     domain.NotificationPreferenceRepository
	userRepo     repository.UserRepository
	hub          *realtime.Hub
	log          logger.Logger
}

func NewNotificationUsecase(
	repo repository.NotificationRepository,
	templateRepo domain.NotificationTemplateRepository,
	prefRepo domain.NotificationPreferenceRepository,
	userRepo repository.UserRepository,
	hub *realtime.Hub,
	log logger.Logger,
) NotificationUsecase {
	return &notificationUsecase{
		repo:         repo,
		templateRepo: templateRepo,
		prefRepo:     prefRepo,
		userRepo:     userRepo,
		hub:          hub,
		log:          log,
	}
}

func (u *notificationUsecase) SendNotification(ctx context.Context, notification *domain.Notification) error {
	// Check user preferences
	prefs, err := u.prefRepo.GetByUserID(ctx, notification.UserID)
	if err == nil && prefs != nil {
		// Check if notification type is enabled
		// This logic can be more complex based on NotificationTypes array
		// For now, basic checks
		if notification.Channel == domain.NotificationChannelInApp && !prefs.EnableInApp {
			return nil // Skip
		}
		if notification.Channel == domain.NotificationChannelEmail && !prefs.EnableEmail {
			return nil // Skip
		}
		// Check quiet hours
		if prefs.IsInQuietHours() && notification.Priority != domain.NotificationPriorityUrgent {
			// Maybe delay or skip? For now, let's just log and proceed or skip.
			// u.log.Info(ctx, "Notification skipped due to quiet hours", logger.LogField{Key: "user_id", Value: notification.UserID})
			// return nil
		}
	}

	// Save to DB
	if err := u.repo.Create(ctx, notification); err != nil {
		return err
	}

	// Send to Realtime Hub
	if u.hub != nil {
		u.hub.SendToUser(notification.UserID, notification)
	}

	return nil
}

func (u *notificationUsecase) SendNotificationFromTemplate(ctx context.Context, userID, templateName string, variables map[string]string) error {
	template, err := u.templateRepo.GetByName(ctx, templateName)
	if err != nil {
		return fmt.Errorf("template not found: %w", err)
	}

	if !template.IsActive {
		return fmt.Errorf("template is inactive")
	}

	// Replace variables in Subject and Body
	subject := template.Subject
	body := template.BodyText // Or BodyHTML
	for k, v := range variables {
		subject = strings.ReplaceAll(subject, "{{"+k+"}}", v)
		body = strings.ReplaceAll(body, "{{"+k+"}}", v)
	}

	notification := &domain.Notification{
		UserID:   userID,
		Type:     template.Type,
		Channel:  template.Channel,
		Priority: template.Priority,
		Title:    subject,
		Message:  body,
		IsSent:   false,
	}

	return u.SendNotification(ctx, notification)
}

func (u *notificationUsecase) SendBulkNotification(ctx context.Context, userIDs []string, notification *domain.Notification) error {
	for _, userID := range userIDs {
		n := *notification // Copy
		n.UserID = userID
		n.ID = "" // Reset ID to let DB generate new one
		if err := u.SendNotification(ctx, &n); err != nil {
			u.log.Error(ctx, "Failed to send bulk notification", logger.LogField{Key: "user_id", Value: userID}, logger.LogField{Key: "error", Value: err})
			// Continue with others
		}
	}
	return nil
}

func (u *notificationUsecase) GetNotification(ctx context.Context, id string) (*domain.Notification, error) {
	return u.repo.GetByID(ctx, id)
}

func (u *notificationUsecase) GetUserNotifications(ctx context.Context, userID string, filter domain.NotificationFilter) ([]*domain.Notification, int64, error) {
	return u.repo.GetByUserID(ctx, userID, filter)
}

func (u *notificationUsecase) GetUnreadCount(ctx context.Context, userID string) (int64, error) {
	return u.repo.GetUnreadCount(ctx, userID)
}

func (u *notificationUsecase) MarkAsRead(ctx context.Context, id string) error {
	return u.repo.MarkAsRead(ctx, id)
}

func (u *notificationUsecase) MarkAllAsRead(ctx context.Context, userID string) error {
	return u.repo.MarkAllAsRead(ctx, userID)
}

func (u *notificationUsecase) DeleteNotification(ctx context.Context, id string) error {
	return u.repo.Delete(ctx, id)
}

func (u *notificationUsecase) CleanupOldNotifications(ctx context.Context, retentionDays int) error {
	duration := time.Duration(retentionDays) * 24 * time.Hour
	return u.repo.DeleteOlderThan(ctx, duration)
}

func (u *notificationUsecase) GetUserPreferences(ctx context.Context, userID string) (*domain.NotificationPreference, error) {
	pref, err := u.prefRepo.GetByUserID(ctx, userID)
	if err != nil {
		// If not found, return default preferences
		return &domain.NotificationPreference{
			UserID:      userID,
			EnableInApp: true,
			EnableEmail: true,
			EnablePush:  true,
		}, nil
	}
	return pref, nil
}

func (u *notificationUsecase) UpdateUserPreferences(ctx context.Context, userID string, preferences *domain.NotificationPreference) error {
	existing, err := u.prefRepo.GetByUserID(ctx, userID)
	if err != nil {
		// Create if not exists
		preferences.UserID = userID
		return u.prefRepo.Create(ctx, preferences)
	}
	preferences.ID = existing.ID
	preferences.UserID = userID
	return u.prefRepo.Update(ctx, preferences)
}
