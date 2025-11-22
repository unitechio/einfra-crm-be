package domain

import (
	"context"
	"time"
)

// NotificationType represents the type of notification
type NotificationType string

const (
	// NotificationTypeInfo represents an informational notification
	NotificationTypeInfo NotificationType = "info"
	// NotificationTypeSuccess represents a success notification
	NotificationTypeSuccess NotificationType = "success"
	// NotificationTypeWarning represents a warning notification
	NotificationTypeWarning NotificationType = "warning"
	// NotificationTypeError represents an error notification
	NotificationTypeError NotificationType = "error"
	// NotificationTypeSystem represents a system notification
	NotificationTypeSystem NotificationType = "system"
)

// NotificationChannel represents the delivery channel for notifications
type NotificationChannel string

const (
	// NotificationChannelInApp represents in-app notifications
	NotificationChannelInApp NotificationChannel = "in_app"
	// NotificationChannelEmail represents email notifications
	NotificationChannelEmail NotificationChannel = "email"
	// NotificationChannelSMS represents SMS notifications
	NotificationChannelSMS NotificationChannel = "sms"
	// NotificationChannelWebhook represents webhook notifications
	NotificationChannelWebhook NotificationChannel = "webhook"
	// NotificationChannelPush represents push notifications
	NotificationChannelPush NotificationChannel = "push"
)

// NotificationPriority represents the priority of a notification
type NotificationPriority string

const (
	// NotificationPriorityLow represents low priority
	NotificationPriorityLow NotificationPriority = "low"
	// NotificationPriorityNormal represents normal priority
	NotificationPriorityNormal NotificationPriority = "normal"
	// NotificationPriorityHigh represents high priority
	NotificationPriorityHigh NotificationPriority = "high"
	// NotificationPriorityUrgent represents urgent priority
	NotificationPriorityUrgent NotificationPriority = "urgent"
)

// Notification represents a notification in the system
// @Description Notification entity for user alerts and messages
type Notification struct {
	ID           string               `json:"id" gorm:"primaryKey;type:uuid;default:gen_random_uuid()" example:"550e8400-e29b-41d4-a716-446655440000"`
	UserID       string               `json:"user_id" gorm:"type:uuid;not null;index" validate:"required" example:"550e8400-e29b-41d4-a716-446655440000"`
	Type         NotificationType     `json:"type" gorm:"type:varchar(50);not null;index" validate:"required" example:"info"`
	Channel      NotificationChannel  `json:"channel" gorm:"type:varchar(50);not null" validate:"required" example:"in_app"`
	Priority     NotificationPriority `json:"priority" gorm:"type:varchar(50);not null;index" validate:"required" example:"normal"`
	Title        string               `json:"title" gorm:"type:varchar(255);not null" validate:"required" example:"Server Alert"`
	Message      string               `json:"message" gorm:"type:text;not null" validate:"required" example:"Server web-01 is down"`
	Data         interface{}          `json:"data,omitempty" gorm:"type:jsonb"` // Additional structured data
	ActionURL    string               `json:"action_url,omitempty" gorm:"type:varchar(500)" example:"/servers/550e8400"`
	ActionLabel  string               `json:"action_label,omitempty" gorm:"type:varchar(100)" example:"View Server"`
	Icon         string               `json:"icon,omitempty" gorm:"type:varchar(100)" example:"server-alert"`
	IsRead       bool                 `json:"is_read" gorm:"type:boolean;default:false;index" example:"false"`
	ReadAt       *time.Time           `json:"read_at,omitempty" swaggertype:"string" example:"2024-01-01T00:00:00Z"`
	IsSent       bool                 `json:"is_sent" gorm:"type:boolean;default:false;index" example:"true"`
	SentAt       *time.Time           `json:"sent_at,omitempty" swaggertype:"string" example:"2024-01-01T00:00:00Z"`
	FailedAt     *time.Time           `json:"failed_at,omitempty" swaggertype:"string" example:"2024-01-01T00:00:00Z"`
	ErrorMessage string               `json:"error_message,omitempty" gorm:"type:text"`
	ExpiresAt    *time.Time           `json:"expires_at,omitempty" swaggertype:"string" example:"2024-01-08T00:00:00Z"`
	CreatedAt    time.Time            `json:"created_at" gorm:"autoCreateTime;index" example:"2024-01-01T00:00:00Z"`
	UpdatedAt    time.Time            `json:"updated_at" gorm:"autoUpdateTime" example:"2024-01-01T00:00:00Z"`
	DeletedAt    *time.Time           `json:"deleted_at,omitempty" gorm:"index" swaggertype:"string" example:"2024-01-01T00:00:00Z"`
}

// TableName specifies the table name for Notification model
func (Notification) TableName() string {
	return "notifications"
}

// IsExpired checks if the notification is expired
func (n *Notification) IsExpired() bool {
	if n.ExpiresAt == nil {
		return false
	}
	return time.Now().After(*n.ExpiresAt)
}

// MarkAsRead marks the notification as read
func (n *Notification) MarkAsRead() {
	n.IsRead = true
	now := time.Now()
	n.ReadAt = &now
}

// NotificationTemplate represents a notification template
// @Description Template for creating notifications
type NotificationTemplate struct {
	ID          string               `json:"id" gorm:"primaryKey;type:uuid;default:gen_random_uuid()" example:"550e8400-e29b-41d4-a716-446655440000"`
	Name        string               `json:"name" gorm:"type:varchar(100);not null;uniqueIndex" validate:"required" example:"server_down_alert"`
	DisplayName string               `json:"display_name" gorm:"type:varchar(255);not null" example:"Server Down Alert"`
	Description string               `json:"description" gorm:"type:text" example:"Template for server down notifications"`
	Type        NotificationType     `json:"type" gorm:"type:varchar(50);not null" validate:"required" example:"error"`
	Channel     NotificationChannel  `json:"channel" gorm:"type:varchar(50);not null" validate:"required" example:"email"`
	Priority    NotificationPriority `json:"priority" gorm:"type:varchar(50);not null" validate:"required" example:"high"`
	Subject     string               `json:"subject" gorm:"type:varchar(255)" example:"Server {{server_name}} is Down"`
	BodyHTML    string               `json:"body_html" gorm:"type:text" example:"<p>Server {{server_name}} is down...</p>"`
	BodyText    string               `json:"body_text" gorm:"type:text" example:"Server {{server_name}} is down..."`
	Variables   []string             `json:"variables" gorm:"type:jsonb" example:"server_name,server_ip"`
	IsActive    bool                 `json:"is_active" gorm:"type:boolean;default:true;index" example:"true"`
	CreatedAt   time.Time            `json:"created_at" gorm:"autoCreateTime" example:"2024-01-01T00:00:00Z"`
	UpdatedAt   time.Time            `json:"updated_at" gorm:"autoUpdateTime" example:"2024-01-01T00:00:00Z"`
	DeletedAt   *time.Time           `json:"deleted_at,omitempty" gorm:"index" swaggertype:"string" example:"2024-01-01T00:00:00Z"`
}

// TableName specifies the table name for NotificationTemplate model
func (NotificationTemplate) TableName() string {
	return "notification_templates"
}

// NotificationPreference represents user notification preferences
// @Description User preferences for receiving notifications
type NotificationPreference struct {
	ID                string    `json:"id" gorm:"primaryKey;type:uuid;default:gen_random_uuid()" example:"550e8400-e29b-41d4-a716-446655440000"`
	UserID            string    `json:"user_id" gorm:"type:uuid;not null;uniqueIndex" validate:"required" example:"550e8400-e29b-41d4-a716-446655440000"`
	EnableInApp       bool      `json:"enable_in_app" gorm:"type:boolean;default:true" example:"true"`
	EnableEmail       bool      `json:"enable_email" gorm:"type:boolean;default:true" example:"true"`
	EnableSMS         bool      `json:"enable_sms" gorm:"type:boolean;default:false" example:"false"`
	EnablePush        bool      `json:"enable_push" gorm:"type:boolean;default:true" example:"true"`
	EmailDigest       bool      `json:"email_digest" gorm:"type:boolean;default:false" example:"false"`
	DigestFrequency   string    `json:"digest_frequency" gorm:"type:varchar(50)" example:"daily"` // daily, weekly
	QuietHoursStart   *int      `json:"quiet_hours_start,omitempty" example:"22"`                 // Hour 0-23
	QuietHoursEnd     *int      `json:"quiet_hours_end,omitempty" example:"8"`                    // Hour 0-23
	NotificationTypes []string  `json:"notification_types" gorm:"type:jsonb"`                     // Types to receive
	CreatedAt         time.Time `json:"created_at" gorm:"autoCreateTime" example:"2024-01-01T00:00:00Z"`
	UpdatedAt         time.Time `json:"updated_at" gorm:"autoUpdateTime" example:"2024-01-01T00:00:00Z"`
}

// TableName specifies the table name for NotificationPreference model
func (NotificationPreference) TableName() string {
	return "notification_preferences"
}

// IsInQuietHours checks if current time is within quiet hours
func (np *NotificationPreference) IsInQuietHours() bool {
	if np.QuietHoursStart == nil || np.QuietHoursEnd == nil {
		return false
	}

	currentHour := time.Now().Hour()
	start := *np.QuietHoursStart
	end := *np.QuietHoursEnd

	if start < end {
		return currentHour >= start && currentHour < end
	}
	// Handle overnight quiet hours (e.g., 22:00 to 08:00)
	return currentHour >= start || currentHour < end
}

// NotificationFilter represents filtering options for notification queries
type NotificationFilter struct {
	UserID   *string               `json:"user_id,omitempty"`
	Type     *NotificationType     `json:"type,omitempty"`
	Channel  *NotificationChannel  `json:"channel,omitempty"`
	Priority *NotificationPriority `json:"priority,omitempty"`
	IsRead   *bool                 `json:"is_read,omitempty"`
	IsSent   *bool                 `json:"is_sent,omitempty"`
	Page     int                   `json:"page" validate:"min=1"`
	PageSize int                   `json:"page_size" validate:"min=1,max=100"`
}

// NotificationTemplateRepository defines the interface for notification template storage
type NotificationTemplateRepository interface {
	// Create creates a new notification template
	Create(ctx context.Context, template *NotificationTemplate) error

	// GetByID retrieves a template by ID
	GetByID(ctx context.Context, id string) (*NotificationTemplate, error)

	// GetByName retrieves a template by name
	GetByName(ctx context.Context, name string) (*NotificationTemplate, error)

	// List retrieves all templates
	List(ctx context.Context, filter NotificationTemplateFilter) ([]*NotificationTemplate, int64, error)

	// Update updates a template
	Update(ctx context.Context, template *NotificationTemplate) error

	// Delete soft deletes a template
	Delete(ctx context.Context, id string) error
}

// NotificationTemplateFilter represents filtering options for template queries
type NotificationTemplateFilter struct {
	Type     *NotificationType    `json:"type,omitempty"`
	Channel  *NotificationChannel `json:"channel,omitempty"`
	IsActive *bool                `json:"is_active,omitempty"`
	Page     int                  `json:"page" validate:"min=1"`
	PageSize int                  `json:"page_size" validate:"min=1,max=100"`
}

// NotificationPreferenceRepository defines the interface for notification preference storage
type NotificationPreferenceRepository interface {
	// Create creates notification preferences for a user
	Create(ctx context.Context, preference *NotificationPreference) error

	// GetByUserID retrieves preferences for a user
	GetByUserID(ctx context.Context, userID string) (*NotificationPreference, error)

	// Update updates notification preferences
	Update(ctx context.Context, preference *NotificationPreference) error

	// Delete deletes notification preferences
	Delete(ctx context.Context, userID string) error
}
