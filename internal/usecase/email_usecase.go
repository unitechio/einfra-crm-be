package usecase

import (
	"context"

	"github.com/unitechio/einfra-be/internal/domain"
)

type EmailUsecase interface {
	SendEmail(ctx context.Context, to []string, subject, body string) error
	SendEmailWithTemplate(ctx context.Context, to []string, templateName string, data interface{}) error
}

type emailUsecase struct {
	emailService domain.EmailService
}

func NewEmailUsecase(emailService domain.EmailService) EmailUsecase {
	return &emailUsecase{
		emailService: emailService,
	}
}

func (u *emailUsecase) SendEmail(ctx context.Context, to []string, subject, body string) error {
	data := domain.EmailData{
		To:      to,
		Subject: subject,
		Data:    map[string]interface{}{"body": body}, // Simplified
	}
	return u.emailService.SendEmail(ctx, data)
}

func (u *emailUsecase) SendEmailWithTemplate(ctx context.Context, to []string, templateName string, data interface{}) error {
	emailData := domain.EmailData{
		To:       to,
		Template: templateName,
		Data:     data.(map[string]interface{}), // Type assertion, might need better handling
	}
	return u.emailService.SendEmail(ctx, emailData)
}
