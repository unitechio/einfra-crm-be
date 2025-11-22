
package email

import (
	"bytes"
	"context"
	"fmt"
	"html/template"
	"os"
	"path/filepath"

	"github.com/jaytaylor/html2text"
	"gopkg.in/gomail.v2"
	"github.com/unitechio/einfra-be/internal/config"
	"github.com/unitechio/einfra-be/internal/domain"
)

// smtpService implements the domain.EmailService interface using an SMTP server.
type smtpService struct {
	dialer   *gomail.Dialer
	from     string
	templates *template.Template
}

// NewSmtpService creates a new SMTP email service.
func NewSmtpService(cfg *config.SMTPConfig) (domain.EmailService, error) {
	d := gomail.NewDialer(cfg.Host, cfg.Port, cfg.Username, cfg.Password)

	templates, err := parseTemplates("internal/templates/emails")
	if err != nil {
		return nil, fmt.Errorf("failed to parse email templates: %w", err)
	}

	return &smtpService{
		dialer:   d,
		from:     cfg.From,
		templates: templates,
	}, nil
}

// SendEmail sends an email using the configured SMTP server.
func (s *smtpService) SendEmail(ctx context.Context, emailData domain.EmailData) error {
	// Render the HTML body from the template.
	htmlBody, err := s.renderTemplate(emailData.Template, emailData.Data)
	if err != nil {
		return fmt.Errorf("failed to render email template: %w", err)
	}

	// Generate a plain text version of the HTML.
	plainTextBody, err := html2text.FromString(htmlBody)
	if err != nil {
		return fmt.Errorf("failed to convert HTML to plain text: %w", err)
	}

	// Compose the email.
	m := gomail.NewMessage()
	m.SetHeader("From", s.from)
	m.SetHeader("To", emailData.To...)
	m.SetHeader("Subject", emailData.Subject)
	m.SetBody("text/plain", plainTextBody)
	m.AddAlternative("text/html", htmlBody)

	// Send the email.
	return s.dialer.DialAndSend(m)
}

// renderTemplate executes a template with the given data.
func (s *smtpService) renderTemplate(templateName string, data interface{}) (string, error) {
	var tpl bytes.Buffer
	if err := s.templates.ExecuteTemplate(&tpl, templateName, data); err != nil {
		return "", err
	}
	return tpl.String(), nil
}

// parseTemplates loads and parses all HTML templates from a directory.
func parseTemplates(dir string) (*template.Template, error) {
	t := template.New("")
	err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() && filepath.Ext(path) == ".html" {
			_, err := t.ParseFiles(path)
			return err
		}
		return nil
	})

	if err != nil {
		return nil, err
	}
	return t, nil
}
