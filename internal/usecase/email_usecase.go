
package usecase

import (
	"context"
	"mymodule/internal/domain"
)

// emailUsecase implements the domain.EmailUsecase interface.
type emailUsecase struct {
	emailService domain.EmailService
}

// NewEmailUsecase creates a new email usecase.
func NewEmailUsecase(emailService domain.EmailService) domain.EmailUsecase {
	return &emailUsecase{
		emailService: emailService,
	}
}

// SendWelcomeEmail sends a welcome email to a new user.
func (uc *emailUsecase) SendWelcomeEmail(ctx context.Context, user *domain.User) error {
	data := domain.EmailData{
		To:       []string{user.Email},
		Subject:  "Welcome to Our Platform!",
		Template: "welcome.html",
		Data: map[string]interface{}{
			"Name": user.Username,
		},
	}
	return uc.emailService.SendEmail(ctx, data)
}

// SendCustomEmail sends a generic email based on the provided data.
func (uc *emailUsecase) SendCustomEmail(ctx context.Context, data domain.EmailData) error {
	return uc.emailService.SendEmail(ctx, data)
}
