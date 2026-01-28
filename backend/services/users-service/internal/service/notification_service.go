package service

import (
	"encoding/json"
	"log"

	"github.com/gofund/shared/events"
	"github.com/gofund/shared/messaging"
)

// NotificationService handles consuming notification events
type NotificationService struct {
	consumer messaging.Consumer
}

// NewNotificationService creates a new notification service instance
func NewNotificationService(consumer messaging.Consumer) *NotificationService {
	return &NotificationService{
		consumer: consumer,
	}
}

// StartConsumers starts all event consumers for notifications
func (s *NotificationService) StartConsumers() error {
	// Start consuming UserSignedUp events
	if err := s.consumer.Consume("UserSignedUp", s.handleUserSignedUp); err != nil {
		return err
	}

	// Start consuming PasswordResetRequested events
	if err := s.consumer.Consume("PasswordResetRequested", s.handlePasswordResetRequested); err != nil {
		return err
	}

	// Start consuming EmailVerificationRequested events
	if err := s.consumer.Consume("EmailVerificationRequested", s.handleEmailVerificationRequested); err != nil {
		return err
	}

	return nil
}

// handleUserSignedUp handles UserSignedUp events
func (s *NotificationService) handleUserSignedUp(body []byte) error {
	var event events.UserSignedUp
	if err := json.Unmarshal(body, &event); err != nil {
		log.Printf("Error unmarshaling UserSignedUp event: %v", err)
		return err
	}

	// Log the email sending (placeholder for actual email service)
	log.Printf("📧 NOTIFICATION: Sending welcome email to %s (UserID: %s, Username: %s)", 
		event.Email, event.UserID, event.Username)
	log.Printf("   Subject: Welcome to GoFund!")
	log.Printf("   Content: Welcome %s! Your account has been created successfully.", event.Username)
	
	return nil
}

// handlePasswordResetRequested handles PasswordResetRequested events
func (s *NotificationService) handlePasswordResetRequested(body []byte) error {
	var event events.PasswordResetRequested
	if err := json.Unmarshal(body, &event); err != nil {
		log.Printf("Error unmarshaling PasswordResetRequested event: %v", err)
		return err
	}

	// Log the email sending (placeholder for actual email service)
	log.Printf("📧 NOTIFICATION: Sending password reset email to %s (UserID: %s)", 
		event.Email, event.UserID)
	log.Printf("   Subject: Reset Your GoFund Password")
	log.Printf("   Content: Click this link to reset your password: /reset-password?token=%s", event.Token)
	
	return nil
}

// handleEmailVerificationRequested handles EmailVerificationRequested events
func (s *NotificationService) handleEmailVerificationRequested(body []byte) error {
	var event events.EmailVerificationRequested
	if err := json.Unmarshal(body, &event); err != nil {
		log.Printf("Error unmarshaling EmailVerificationRequested event: %v", err)
		return err
	}

	// Log the email sending (placeholder for actual email service)
	log.Printf("📧 NOTIFICATION: Sending email verification to %s (UserID: %s)", 
		event.Email, event.UserID)
	log.Printf("   Subject: Verify Your GoFund Email")
	log.Printf("   Content: Click this link to verify your email: /verify-email?token=%s", event.Token)
	
	return nil
}