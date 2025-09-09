
package domain

import (
	"context"
)

// EmailData defines the structure for the data to be injected into an email template.
type EmailData struct {
	To       []string
	Subject  string
	Template string
	Data     map[string]interface{}
}

// EmailService defines the interface for sending emails.
// This is an infrastructure-level interface.
type EmailService interface {
	SendEmail(ctx context.Context, data EmailData) error
}

// EmailUsecase defines the interface for email-related use cases.
// This is a business-logic-level interface.
type EmailUsecase interface {
	// SendWelcomeEmail sends a welcome email to a new user.
	SendWelcomeEmail(ctx context.Context, user *User) error
	// SendCustomEmail sends a generic email based on the provided data.
	SendCustomEmail(ctx context.Context, data EmailData) error
}
