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
)

// FontSize represents the font size.
type FontSize string

const (
	SmallFontSize FontSize = "small"
	LargeFontSize FontSize = "large"
)

// NotificationPreference defines how a user wants to be notified.
type NotificationPreference string

const (
	AllMessages       NotificationPreference = "all"
	DirectAndMentions NotificationPreference = "direct"
	NoMessages        NotificationPreference = "none"
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

// UserSettings represents all user-specific settings.
type UserSettings struct {
	// Display Settings
	Theme    Theme           `json:"theme"`
	FontSize FontSize        `json:"font_size"`
	Sidebar  SidebarSettings `json:"sidebar"`

	// Notification Settings
	NotificationPreference NotificationPreference    `json:"notification_preference"`
	EmailNotifications     EmailNotificationSettings `json:"email_notifications"`
	UseMobileSettings      bool                      `json:"use_mobile_settings"`

	// Account Settings
	Name        string    `json:"name"`
	DateOfBirth time.Time `json:"date_of_birth"`
	Language    string    `json:"language"`

	// Security Settings
	TwoFactorEnabled bool `json:"two_factor_enabled"`
}
