package service

import (
	"fmt"
	"log"
	"net/smtp"
	"time"

	"github.com/gofund/notifications-service/internal/config"
	"github.com/gofund/shared/metrics"
	"github.com/gofund/shared/models"
)

// EmailService handles email sending
type EmailService interface {
	Send(payload models.EmailPayload) error
}

type emailService struct {
	config        *config.Config
	renderService RenderService
}

// NewEmailService creates a new email service
func NewEmailService(cfg *config.Config, renderService RenderService) EmailService {
	return &emailService{
		config:        cfg,
		renderService: renderService,
	}
}

// Send renders and sends an email
func (s *emailService) Send(payload models.EmailPayload) error {
	start := time.Now()

	// 1. Render HTML Body
	htmlBody, err := s.renderService.Render(payload.Type, payload.Data)
	if err != nil {
		return fmt.Errorf("failed to render email: %w", err)
	}

	// 2. Build email message
	from := fmt.Sprintf("%s <%s>", s.config.SMTPFromName, s.config.SMTPFrom)
	
	headers := make(map[string]string)
	headers["From"] = from
	headers["To"] = payload.Recipient
	headers["Subject"] = payload.Subject
	headers["MIME-Version"] = "1.0"
	headers["Content-Type"] = "text/html; charset=UTF-8"

	message := ""
	for k, v := range headers {
		message += fmt.Sprintf("%s: %s\r\n", k, v)
	}
	message += "\r\n" + htmlBody

	// 3. SMTP authentication (Gmail App Password)
	auth := smtp.PlainAuth(
		"",
		s.config.SMTPUsername,
		s.config.SMTPPassword,
		s.config.SMTPHost,
	)

	// 4. Send email
	addr := fmt.Sprintf("%s:%s", s.config.SMTPHost, s.config.SMTPPort)
	err = smtp.SendMail(
		addr,
		auth,
		s.config.SMTPFrom,
		[]string{payload.Recipient},
		[]byte(message),
	)

	duration := time.Since(start)

	if err != nil {
		metrics.TrackEmailSent(false, duration)
		return fmt.Errorf("failed to send email via SMTP: %w", err)
	}

	metrics.TrackEmailSent(true, duration)
	log.Printf("Email [%s] sent to %s (duration: %v)", payload.Type, payload.Recipient, duration)
	return nil
}
