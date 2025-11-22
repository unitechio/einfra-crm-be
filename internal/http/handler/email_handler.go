
package http

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/unitechio/einfra-be/internal/domain"
)

// EmailHandler handles HTTP requests related to emails.
type EmailHandler struct {
	emailUsecase domain.EmailUsecase
}

// NewEmailHandler creates a new EmailHandler.
func NewEmailHandler(e *echo.Echo, emailUsecase domain.EmailUsecase) {
	h := &EmailHandler{
		emailUsecase: emailUsecase,
	}

	e.POST("/emails/send", h.SendCustomEmail)
}

// sendEmailRequest represents the request body for sending a custom email.
type sendEmailRequest struct {
	To       []string               `json:"to"`
	Subject  string                 `json:"subject"`
	Template string                 `json:"template"`
	Data     map[string]interface{} `json:"data"`
}

// SendCustomEmail handles the request to send a custom email.
func (h *EmailHandler) SendCustomEmail(c echo.Context) error {
	var req sendEmailRequest
	if err := c.Bind(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	// Here you might want to add validation for the request data.

	ctx := c.Request().Context()
	emailData := domain.EmailData{
		To:       req.To,
		Subject:  req.Subject,
		Template: req.Template,
		Data:     req.Data,
	}

	if err := h.emailUsecase.SendCustomEmail(ctx, emailData); err != nil {
		// In a real application, you might want to log this error and return a more generic message.
		return echo.NewHTTPError(http.StatusInternalServerError, "failed to send email")
	}

	return c.JSON(http.StatusOK, map[string]string{"message": "Email sent successfully"})
}
