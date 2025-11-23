package domain

import (
	"time"
)

// EmailData contains all information needed to send an email
type EmailData struct {
	To          []string               // Recipient email addresses
	CC          []string               // Carbon copy recipients
	BCC         []string               // Blind carbon copy recipients
	From        string                 // Sender email (optional, uses default if empty)
	ReplyTo     string                 // Reply-to email address
	Subject     string                 // Email subject
	Body        string                 // Plain text body
	HTMLBody    string                 // HTML body
	Template    string                 // Template name (if using templates)
	Data        map[string]interface{} // Template data
	Attachments []EmailAttachment      // File attachments
	Headers     map[string]string      // Custom headers
	Priority    EmailPriority          // Email priority
}

// EmailAttachment represents a file attachment
type EmailAttachment struct {
	Filename    string // Name of the file
	Content     []byte // File content
	ContentType string // MIME type (e.g., "application/pdf")
	Inline      bool   // Whether to embed inline (for images)
	ContentID   string // Content ID for inline attachments
}

// EmailPriority represents email priority levels
type EmailPriority int

const (
	PriorityNormal EmailPriority = iota
	PriorityLow
	PriorityHigh
)

// EmailLog represents a sent email record
type EmailLog struct {
	ID        string                 `json:"id"`
	To        []string               `json:"to"`
	CC        []string               `json:"cc"`
	BCC       []string               `json:"bcc"`
	From      string                 `json:"from"`
	Subject   string                 `json:"subject"`
	Template  string                 `json:"template,omitempty"`
	Status    EmailStatus            `json:"status"`
	Error     string                 `json:"error,omitempty"`
	Metadata  map[string]interface{} `json:"metadata,omitempty"`
	SentAt    time.Time              `json:"sent_at"`
	CreatedAt time.Time              `json:"created_at"`
	UpdatedAt time.Time              `json:"updated_at"`
}

// EmailStatus represents the status of an email
type EmailStatus string

const (
	EmailStatusPending   EmailStatus = "pending"
	EmailStatusSent      EmailStatus = "sent"
	EmailStatusFailed    EmailStatus = "failed"
	EmailStatusBounced   EmailStatus = "bounced"
	EmailStatusDelivered EmailStatus = "delivered"
)

// EmailLogFilter represents filters for querying email logs
type EmailLogFilter struct {
	Status   EmailStatus
	From     string
	To       string
	DateFrom *time.Time
	DateTo   *time.Time
	Template string
	Limit    int
	Offset   int
}
