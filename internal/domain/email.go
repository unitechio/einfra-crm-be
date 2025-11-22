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
