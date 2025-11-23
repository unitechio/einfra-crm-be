package domain

import (
	"time"
)

// Theme represents the UI theme.
type Theme string

const (
	LightTheme  Theme = "light"
	DarkTheme   Theme = "dark"
	SystemTheme Theme = "system" // Automatically switch theme based on system preferences.
	AutoTheme   Theme = "auto"   // Auto switch based on time of day
)

// FontSize represents the font size.
type FontSize string

const (
	SmallFontSize      FontSize = "small"
	MediumFontSize     FontSize = "medium"
	LargeFontSize      FontSize = "large"
	ExtraLargeFontSize FontSize = "extra-large"
)

// TableDensity represents table row spacing
type TableDensity string

const (
	CompactDensity     TableDensity = "compact"
	ComfortableDensity TableDensity = "comfortable"
	SpaciousDensity    TableDensity = "spacious"
)

// DigestFrequency represents notification digest frequency
type DigestFrequency string

const (
	RealtimeDigest DigestFrequency = "realtime"
	HourlyDigest   DigestFrequency = "hourly"
	DailyDigest    DigestFrequency = "daily"
	WeeklyDigest   DigestFrequency = "weekly"
	NeverDigest    DigestFrequency = "never"
)

// NotificationLevel represents notification delivery preference
type NotificationLevel string

const (
	AllNotificationsLevel  NotificationLevel = "all"
	ImportantOnlyLevel     NotificationLevel = "important"
	MentionsOnlyLevel      NotificationLevel = "mentions"
	NoneNotificationsLevel NotificationLevel = "none"
)

// EmailNotificationSettings holds the user's email notification preferences.
type EmailNotificationSettings struct {
	Communication bool `json:"communication"` // Emails about account activity.
	Marketing     bool `json:"marketing"`     // Emails about new products and features.
	Social        bool `json:"social"`        // Emails for friend requests, follows, etc.
	Security      bool `json:"security"`      // Emails about account security.
}

// SidebarSettings defines which items to display in the sidebar.
type SidebarSettings struct {
	Recents      bool `json:"recents"`
	Home         bool `json:"home"`
	Applications bool `json:"applications"`
	Desktop      bool `json:"desktop"`
	Downloads    bool `json:"downloads"`
	Documents    bool `json:"documents"`
}

// UserSettings represents all user-specific settings with database persistence
type UserSettings struct {
	ID     string `json:"id" gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	UserID string `json:"user_id" gorm:"type:uuid;not null;uniqueIndex"`

	// Display Settings
	Theme            Theme           `json:"theme" gorm:"type:varchar(20);default:'light'"`
	FontSize         FontSize        `json:"font_size" gorm:"type:varchar(20);default:'medium'"`
	FontFamily       string          `json:"font_family" gorm:"type:varchar(50);default:'Inter'"` // Inter, Roboto, Arial, etc.
	CompactMode      bool            `json:"compact_mode" gorm:"default:false"`
	SidebarCollapsed bool            `json:"sidebar_collapsed" gorm:"default:false"`
	Sidebar          SidebarSettings `json:"sidebar" gorm:"type:jsonb"`

	// Localization Settings
	Language   string `json:"language" gorm:"type:varchar(10);default:'en'"`            // en, vi, ja, etc.
	Timezone   string `json:"timezone" gorm:"type:varchar(50);default:'UTC'"`           // UTC, Asia/Ho_Chi_Minh, etc.
	DateFormat string `json:"date_format" gorm:"type:varchar(20);default:'YYYY-MM-DD'"` // YYYY-MM-DD, DD/MM/YYYY, MM/DD/YYYY
	TimeFormat string `json:"time_format" gorm:"type:varchar(10);default:'24h'"`        // 12h, 24h
	Currency   string `json:"currency" gorm:"type:varchar(10);default:'USD'"`           // USD, VND, JPY, etc.

	// Notification Settings
	NotificationLevel    NotificationLevel         `json:"notification_level" gorm:"type:varchar(20);default:'all'"`
	EmailNotifications   EmailNotificationSettings `json:"email_notifications" gorm:"type:jsonb"`
	PushNotifications    bool                      `json:"push_notifications" gorm:"default:true"`
	DesktopNotifications bool                      `json:"desktop_notifications" gorm:"default:false"`
	NotificationSound    bool                      `json:"notification_sound" gorm:"default:true"`
	NotifyOnMention      bool                      `json:"notify_on_mention" gorm:"default:true"`
	NotifyOnAssignment   bool                      `json:"notify_on_assignment" gorm:"default:true"`
	NotifyOnStatusChange bool                      `json:"notify_on_status_change" gorm:"default:true"`
	NotifyOnComment      bool                      `json:"notify_on_comment" gorm:"default:true"`
	DigestFrequency      DigestFrequency           `json:"digest_frequency" gorm:"type:varchar(20);default:'daily'"`
	UseMobileSettings    bool                      `json:"use_mobile_settings" gorm:"default:false"`

	// Dashboard Settings
	DefaultDashboard string `json:"default_dashboard" gorm:"type:varchar(50)"` // Dashboard ID or name
	WidgetLayout     string `json:"widget_layout" gorm:"type:text"`            // JSON string of widget positions
	FavoritePages    string `json:"favorite_pages" gorm:"type:text"`           // JSON array of page URLs
	RecentlyViewed   string `json:"recently_viewed" gorm:"type:text"`          // JSON array of recently viewed items

	// Table/List Settings
	DefaultPageSize int          `json:"default_page_size" gorm:"default:20"` // 10, 20, 50, 100
	TableDensity    TableDensity `json:"table_density" gorm:"type:varchar(20);default:'comfortable'"`
	ShowRowNumbers  bool         `json:"show_row_numbers" gorm:"default:false"`

	// Accessibility Settings
	HighContrast          bool `json:"high_contrast" gorm:"default:false"`
	ReduceMotion          bool `json:"reduce_motion" gorm:"default:false"`
	ScreenReaderOptimized bool `json:"screen_reader_optimized" gorm:"default:false"`
	KeyboardShortcuts     bool `json:"keyboard_shortcuts" gorm:"default:true"`

	// Privacy Settings
	ShowOnlineStatus    bool `json:"show_online_status" gorm:"default:true"`
	ShowLastSeen        bool `json:"show_last_seen" gorm:"default:true"`
	AllowDataCollection bool `json:"allow_data_collection" gorm:"default:true"`

	// Security Settings
	TwoFactorEnabled bool `json:"two_factor_enabled" gorm:"default:false"`

	// Advanced Settings
	DeveloperMode bool   `json:"developer_mode" gorm:"default:false"`
	BetaFeatures  bool   `json:"beta_features" gorm:"default:false"`
	CustomCSS     string `json:"custom_css" gorm:"type:text"` // Custom CSS for power users

	// Account Settings
	Name        string     `json:"name" gorm:"type:varchar(100)"`
	DateOfBirth *time.Time `json:"date_of_birth,omitempty"`

	// Metadata
	CreatedAt time.Time  `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt time.Time  `json:"updated_at" gorm:"autoUpdateTime"`
	DeletedAt *time.Time `json:"deleted_at,omitempty" gorm:"index"`
}

// UserSettingsUpdate represents partial update request
type UserSettingsUpdate struct {
	// Display Settings
	Theme            *Theme           `json:"theme,omitempty"`
	FontSize         *FontSize        `json:"font_size,omitempty"`
	FontFamily       *string          `json:"font_family,omitempty"`
	CompactMode      *bool            `json:"compact_mode,omitempty"`
	SidebarCollapsed *bool            `json:"sidebar_collapsed,omitempty"`
	Sidebar          *SidebarSettings `json:"sidebar,omitempty"`

	// Localization Settings
	Language   *string `json:"language,omitempty"`
	Timezone   *string `json:"timezone,omitempty"`
	DateFormat *string `json:"date_format,omitempty"`
	TimeFormat *string `json:"time_format,omitempty"`
	Currency   *string `json:"currency,omitempty"`

	// Notification Settings
	NotificationLevel    *NotificationLevel         `json:"notification_level,omitempty"`
	EmailNotifications   *EmailNotificationSettings `json:"email_notifications,omitempty"`
	PushNotifications    *bool                      `json:"push_notifications,omitempty"`
	DesktopNotifications *bool                      `json:"desktop_notifications,omitempty"`
	NotificationSound    *bool                      `json:"notification_sound,omitempty"`
	NotifyOnMention      *bool                      `json:"notify_on_mention,omitempty"`
	NotifyOnAssignment   *bool                      `json:"notify_on_assignment,omitempty"`
	NotifyOnStatusChange *bool                      `json:"notify_on_status_change,omitempty"`
	NotifyOnComment      *bool                      `json:"notify_on_comment,omitempty"`
	DigestFrequency      *DigestFrequency           `json:"digest_frequency,omitempty"`
	UseMobileSettings    *bool                      `json:"use_mobile_settings,omitempty"`

	// Dashboard Settings
	DefaultDashboard *string `json:"default_dashboard,omitempty"`
	WidgetLayout     *string `json:"widget_layout,omitempty"`
	FavoritePages    *string `json:"favorite_pages,omitempty"`
	RecentlyViewed   *string `json:"recently_viewed,omitempty"`

	// Table/List Settings
	DefaultPageSize *int          `json:"default_page_size,omitempty"`
	TableDensity    *TableDensity `json:"table_density,omitempty"`
	ShowRowNumbers  *bool         `json:"show_row_numbers,omitempty"`

	// Accessibility Settings
	HighContrast          *bool `json:"high_contrast,omitempty"`
	ReduceMotion          *bool `json:"reduce_motion,omitempty"`
	ScreenReaderOptimized *bool `json:"screen_reader_optimized,omitempty"`
	KeyboardShortcuts     *bool `json:"keyboard_shortcuts,omitempty"`

	// Privacy Settings
	ShowOnlineStatus    *bool `json:"show_online_status,omitempty"`
	ShowLastSeen        *bool `json:"show_last_seen,omitempty"`
	AllowDataCollection *bool `json:"allow_data_collection,omitempty"`

	// Security Settings
	TwoFactorEnabled *bool `json:"two_factor_enabled,omitempty"`

	// Advanced Settings
	DeveloperMode *bool   `json:"developer_mode,omitempty"`
	BetaFeatures  *bool   `json:"beta_features,omitempty"`
	CustomCSS     *string `json:"custom_css,omitempty"`

	// Account Settings
	Name        *string    `json:"name,omitempty"`
	DateOfBirth *time.Time `json:"date_of_birth,omitempty"`
}

// GetDefaultSettings returns default settings for a new user
func GetDefaultSettings(userID string) *UserSettings {
	return &UserSettings{
		UserID:           userID,
		Theme:            LightTheme,
		FontSize:         MediumFontSize,
		FontFamily:       "Inter",
		CompactMode:      false,
		SidebarCollapsed: false,
		Sidebar: SidebarSettings{
			Recents:      true,
			Home:         true,
			Applications: true,
			Desktop:      true,
			Downloads:    true,
			Documents:    true,
		},
		Language:          "en",
		Timezone:          "UTC",
		DateFormat:        "YYYY-MM-DD",
		TimeFormat:        "24h",
		Currency:          "USD",
		NotificationLevel: AllNotificationsLevel,
		EmailNotifications: EmailNotificationSettings{
			Communication: true,
			Marketing:     false,
			Social:        true,
			Security:      true,
		},
		PushNotifications:     true,
		DesktopNotifications:  false,
		NotificationSound:     true,
		NotifyOnMention:       true,
		NotifyOnAssignment:    true,
		NotifyOnStatusChange:  true,
		NotifyOnComment:       true,
		DigestFrequency:       DailyDigest,
		UseMobileSettings:     false,
		DefaultPageSize:       20,
		TableDensity:          ComfortableDensity,
		ShowRowNumbers:        false,
		HighContrast:          false,
		ReduceMotion:          false,
		ScreenReaderOptimized: false,
		KeyboardShortcuts:     true,
		ShowOnlineStatus:      true,
		ShowLastSeen:          true,
		AllowDataCollection:   true,
		TwoFactorEnabled:      false,
		DeveloperMode:         false,
		BetaFeatures:          false,
	}
}
