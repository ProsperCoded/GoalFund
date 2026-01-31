package service

import (
	"bytes"
	"fmt"
	"html/template"
	"log"
	"net/smtp"
	"time"

	"github.com/gofund/notifications-service/internal/config"
	"github.com/gofund/shared/metrics"
)

// EmailService handles email sending
type EmailService interface {
	SendEmail(to, subject, htmlBody, textBody string) error
	SendTemplatedEmail(to, subject, templateName string, data interface{}) error
}

type emailService struct {
	config    *config.Config
	templates *template.Template
}

// NewEmailService creates a new email service
func NewEmailService(cfg *config.Config) (EmailService, error) {
	// Load email templates
	templates, err := template.ParseGlob("internal/templates/emails/*.html")
	if err != nil {
		log.Printf("Warning: Failed to load email templates: %v", err)
		// Continue without templates - we can still send plain emails
	}

	return &emailService{
		config:    cfg,
		templates: templates,
	}, nil
}

// SendEmail sends a plain email
func (s *emailService) SendEmail(to, subject, htmlBody, textBody string) error {
	start := time.Now()

	// Build email message
	from := fmt.Sprintf("%s <%s>", s.config.SMTPFromName, s.config.SMTPFrom)
	
	headers := make(map[string]string)
	headers["From"] = from
	headers["To"] = to
	headers["Subject"] = subject
	headers["MIME-Version"] = "1.0"
	headers["Content-Type"] = "text/html; charset=UTF-8"

	message := ""
	for k, v := range headers {
		message += fmt.Sprintf("%s: %s\r\n", k, v)
	}
	message += "\r\n" + htmlBody

	// SMTP authentication
	auth := smtp.PlainAuth(
		"",
		s.config.SMTPUsername,
		s.config.SMTPPassword,
		s.config.SMTPHost,
	)

	// Send email
	addr := fmt.Sprintf("%s:%s", s.config.SMTPHost, s.config.SMTPPort)
	err := smtp.SendMail(
		addr,
		auth,
		s.config.SMTPFrom,
		[]string{to},
		[]byte(message),
	)

	duration := time.Since(start)

	if err != nil {
		metrics.TrackEmailSent(false, duration)
		return fmt.Errorf("failed to send email: %w", err)
	}

	metrics.TrackEmailSent(true, duration)
	log.Printf("Email sent to %s (subject: %s, duration: %v)", to, subject, duration)
	return nil
}

// SendTemplatedEmail sends an email using a template
func (s *emailService) SendTemplatedEmail(to, subject, templateName string, data interface{}) error {
	if s.templates == nil {
		return fmt.Errorf("email templates not loaded")
	}

	// Execute template
	var htmlBody bytes.Buffer
	err := s.templates.ExecuteTemplate(&htmlBody, templateName, data)
	if err != nil {
		return fmt.Errorf("failed to execute template: %w", err)
	}

	// Send email
	return s.SendEmail(to, subject, htmlBody.String(), "")
}

// EmailTemplateData represents common data for email templates
type EmailTemplateData struct {
	UserName    string
	GoalTitle   string
	Amount      string
	Message     string
	ActionURL   string
	CurrentYear int
}
