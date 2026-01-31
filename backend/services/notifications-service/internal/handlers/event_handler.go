package handlers

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/gofund/notifications-service/internal/dto"
	"github.com/gofund/notifications-service/internal/models"
	"github.com/gofund/notifications-service/internal/service"
	"github.com/gofund/shared/events"
)

// EventHandler handles events from RabbitMQ
type EventHandler struct {
	notificationService service.NotificationService
}

// NewEventHandler creates a new event handler
func NewEventHandler(notificationService service.NotificationService) *EventHandler {
	return &EventHandler{
		notificationService: notificationService,
	}
}

// HandlePaymentVerified handles PaymentVerified events
func (h *EventHandler) HandlePaymentVerified(data []byte) error {
	var event events.PaymentVerified
	if err := json.Unmarshal(data, &event); err != nil {
		return fmt.Errorf("failed to unmarshal event: %w", err)
	}

	log.Printf("Processing PaymentVerified event: %s for user %s", event.ID, event.UserID)

	// Create notification
	req := dto.CreateNotificationRequest{
		UserID:  event.UserID,
		Type:    models.NotificationTypePaymentVerified,
		Title:   "Payment Successful",
		Message: fmt.Sprintf("Your payment of ₦%.2f has been verified successfully.", float64(event.Amount)/100),
		Data: map[string]interface{}{
			"payment_id": event.PaymentID,
			"goal_id":    event.GoalID,
			"amount":     event.Amount,
			"email":      "", // This should be fetched from user service
		},
	}

	_, err := h.notificationService.CreateNotification(req)
	if err != nil {
		return fmt.Errorf("failed to create notification: %w", err)
	}

	log.Printf("PaymentVerified notification created for user %s", event.UserID)
	return nil
}

// HandleContributionConfirmed handles ContributionConfirmed events
func (h *EventHandler) HandleContributionConfirmed(data []byte) error {
	// Define a simple event structure for contribution confirmed
	var event struct {
		ID            string `json:"id"`
		UserID        string `json:"user_id"`
		GoalID        string `json:"goal_id"`
		GoalOwnerID   string `json:"goal_owner_id"`
		Amount        int64  `json:"amount"`
		ContributorName string `json:"contributor_name"`
		GoalTitle     string `json:"goal_title"`
		CreatedAt     int64  `json:"created_at"`
	}

	if err := json.Unmarshal(data, &event); err != nil {
		return fmt.Errorf("failed to unmarshal event: %w", err)
	}

	log.Printf("Processing ContributionConfirmed event: %s", event.ID)

	// Notify goal owner
	ownerReq := dto.CreateNotificationRequest{
		UserID:  event.GoalOwnerID,
		Type:    models.NotificationTypeContributionConfirmed,
		Title:   "New Contribution Received",
		Message: fmt.Sprintf("You received a contribution of ₦%.2f for your goal '%s'.", float64(event.Amount)/100, event.GoalTitle),
		Data: map[string]interface{}{
			"goal_id":          event.GoalID,
			"contributor_id":   event.UserID,
			"contributor_name": event.ContributorName,
			"amount":           event.Amount,
			"email":            "", // Should be fetched from user service
		},
	}

	_, err := h.notificationService.CreateNotification(ownerReq)
	if err != nil {
		return fmt.Errorf("failed to create notification for goal owner: %w", err)
	}

	log.Printf("ContributionConfirmed notification created for goal owner %s", event.GoalOwnerID)
	return nil
}

// HandleWithdrawalRequested handles WithdrawalRequested events
func (h *EventHandler) HandleWithdrawalRequested(data []byte) error {
	var event struct {
		ID        string `json:"id"`
		GoalID    string `json:"goal_id"`
		OwnerID   string `json:"owner_id"`
		Amount    int64  `json:"amount"`
		GoalTitle string `json:"goal_title"`
		CreatedAt int64  `json:"created_at"`
	}

	if err := json.Unmarshal(data, &event); err != nil {
		return fmt.Errorf("failed to unmarshal event: %w", err)
	}

	log.Printf("Processing WithdrawalRequested event: %s", event.ID)

	// Notify goal owner
	req := dto.CreateNotificationRequest{
		UserID:  event.OwnerID,
		Type:    models.NotificationTypeWithdrawalRequested,
		Title:   "Withdrawal Requested",
		Message: fmt.Sprintf("Your withdrawal request of ₦%.2f for '%s' is being processed.", float64(event.Amount)/100, event.GoalTitle),
		Data: map[string]interface{}{
			"goal_id": event.GoalID,
			"amount":  event.Amount,
			"email":   "", // Should be fetched from user service
		},
	}

	_, err := h.notificationService.CreateNotification(req)
	if err != nil {
		return fmt.Errorf("failed to create notification: %w", err)
	}

	log.Printf("WithdrawalRequested notification created for user %s", event.OwnerID)
	return nil
}

// HandleWithdrawalCompleted handles WithdrawalCompleted events
func (h *EventHandler) HandleWithdrawalCompleted(data []byte) error {
	var event struct {
		ID        string `json:"id"`
		GoalID    string `json:"goal_id"`
		OwnerID   string `json:"owner_id"`
		Amount    int64  `json:"amount"`
		GoalTitle string `json:"goal_title"`
		CreatedAt int64  `json:"created_at"`
	}

	if err := json.Unmarshal(data, &event); err != nil {
		return fmt.Errorf("failed to unmarshal event: %w", err)
	}

	log.Printf("Processing WithdrawalCompleted event: %s", event.ID)

	// Notify goal owner
	req := dto.CreateNotificationRequest{
		UserID:  event.OwnerID,
		Type:    models.NotificationTypeWithdrawalCompleted,
		Title:   "Withdrawal Completed",
		Message: fmt.Sprintf("Your withdrawal of ₦%.2f for '%s' has been completed successfully.", float64(event.Amount)/100, event.GoalTitle),
		Data: map[string]interface{}{
			"goal_id": event.GoalID,
			"amount":  event.Amount,
			"email":   "", // Should be fetched from user service
		},
	}

	_, err := h.notificationService.CreateNotification(req)
	if err != nil {
		return fmt.Errorf("failed to create notification: %w", err)
	}

	log.Printf("WithdrawalCompleted notification created for user %s", event.OwnerID)
	return nil
}

// HandleProofSubmitted handles ProofSubmitted events
func (h *EventHandler) HandleProofSubmitted(data []byte) error {
	var event events.ProofSubmitted
	if err := json.Unmarshal(data, &event); err != nil {
		return fmt.Errorf("failed to unmarshal event: %w", err)
	}

	log.Printf("Processing ProofSubmitted event: %s", event.ID)

	// Note: In a real implementation, you'd fetch contributors from the goals service
	// For now, we'll just log this
	log.Printf("ProofSubmitted for goal %s - contributors should be notified", event.GoalID)
	
	return nil
}

// HandleProofVoted handles ProofVoted events
func (h *EventHandler) HandleProofVoted(data []byte) error {
	var event struct {
		ID        string `json:"id"`
		GoalID    string `json:"goal_id"`
		ProofID   string `json:"proof_id"`
		VoterID   string `json:"voter_id"`
		OwnerID   string `json:"owner_id"`
		Vote      bool   `json:"vote"`
		CreatedAt int64  `json:"created_at"`
	}

	if err := json.Unmarshal(data, &event); err != nil {
		return fmt.Errorf("failed to unmarshal event: %w", err)
	}

	log.Printf("Processing ProofVoted event: %s", event.ID)

	voteText := "satisfied"
	if !event.Vote {
		voteText = "not satisfied"
	}

	// Notify goal owner
	req := dto.CreateNotificationRequest{
		UserID:  event.OwnerID,
		Type:    models.NotificationTypeProofVoted,
		Title:   "New Vote on Your Proof",
		Message: fmt.Sprintf("A contributor voted %s with your proof.", voteText),
		Data: map[string]interface{}{
			"goal_id":  event.GoalID,
			"proof_id": event.ProofID,
			"voter_id": event.VoterID,
			"vote":     event.Vote,
			"email":    "", // Should be fetched from user service
		},
	}

	_, err := h.notificationService.CreateNotification(req)
	if err != nil {
		return fmt.Errorf("failed to create notification: %w", err)
	}

	log.Printf("ProofVoted notification created for user %s", event.OwnerID)
	return nil
}

// HandleGoalFunded handles GoalFunded events
func (h *EventHandler) HandleGoalFunded(data []byte) error {
	var event events.GoalFunded
	if err := json.Unmarshal(data, &event); err != nil {
		return fmt.Errorf("failed to unmarshal event: %w", err)
	}

	log.Printf("Processing GoalFunded event: %s", event.ID)

	// Note: In a real implementation, you'd fetch the goal owner and contributors
	// For now, we'll just log this
	log.Printf("GoalFunded for goal %s - owner and contributors should be notified", event.GoalID)
	
	return nil
}

// HandleUserSignedUp handles UserSignedUp events
func (h *EventHandler) HandleUserSignedUp(data []byte) error {
	var event events.UserSignedUp
	if err := json.Unmarshal(data, &event); err != nil {
		return fmt.Errorf("failed to unmarshal event: %w", err)
	}

	log.Printf("Processing UserSignedUp event: %s for user %s", event.ID, event.UserID)

	// Create default notification preferences
	if err := h.notificationService.CreateDefaultPreferences(event.UserID); err != nil {
		log.Printf("Failed to create default preferences for user %s: %v", event.UserID, err)
	}

	// Create welcome notification
	req := dto.CreateNotificationRequest{
		UserID:  event.UserID,
		Type:    models.NotificationTypeUserSignedUp,
		Title:   "Welcome to GoFund!",
		Message: fmt.Sprintf("Welcome %s! Thank you for joining GoFund. Start creating goals or contributing to existing ones.", event.Username),
		Data: map[string]interface{}{
			"email":    event.Email,
			"username": event.Username,
		},
	}

	_, err := h.notificationService.CreateNotification(req)
	if err != nil {
		return fmt.Errorf("failed to create notification: %w", err)
	}

	log.Printf("UserSignedUp notification created for user %s", event.UserID)
	return nil
}

// HandlePasswordResetRequested handles PasswordResetRequested events
func (h *EventHandler) HandlePasswordResetRequested(data []byte) error {
	var event events.PasswordResetRequested
	if err := json.Unmarshal(data, &event); err != nil {
		return fmt.Errorf("failed to unmarshal event: %w", err)
	}

	log.Printf("Processing PasswordResetRequested event: %s for user %s", event.ID, event.UserID)

	// Create notification
	req := dto.CreateNotificationRequest{
		UserID:  event.UserID,
		Type:    models.NotificationTypePasswordReset,
		Title:   "Password Reset Requested",
		Message: "A password reset was requested for your account. If you didn't request this, please ignore this message.",
		Data: map[string]interface{}{
			"email": event.Email,
			"token": event.Token,
		},
	}

	_, err := h.notificationService.CreateNotification(req)
	if err != nil {
		return fmt.Errorf("failed to create notification: %w", err)
	}

	log.Printf("PasswordResetRequested notification created for user %s", event.UserID)
	return nil
}

// HandleEmailVerificationRequested handles EmailVerificationRequested events
func (h *EventHandler) HandleEmailVerificationRequested(data []byte) error {
	var event events.EmailVerificationRequested
	if err := json.Unmarshal(data, &event); err != nil {
		return fmt.Errorf("failed to unmarshal event: %w", err)
	}

	log.Printf("Processing EmailVerificationRequested event: %s for user %s", event.ID, event.UserID)

	// Create notification
	req := dto.CreateNotificationRequest{
		UserID:  event.UserID,
		Type:    models.NotificationTypeEmailVerification,
		Title:   "Verify Your Email",
		Message: "Please verify your email address to complete your registration.",
		Data: map[string]interface{}{
			"email": event.Email,
			"token": event.Token,
		},
	}

	_, err := h.notificationService.CreateNotification(req)
	if err != nil {
		return fmt.Errorf("failed to create notification: %w", err)
	}

	log.Printf("EmailVerificationRequested notification created for user %s", event.UserID)
	return nil
}

// HandleKYCVerified handles KYCVerified events
func (h *EventHandler) HandleKYCVerified(data []byte) error {
	var event events.KYCVerified
	if err := json.Unmarshal(data, &event); err != nil {
		return fmt.Errorf("failed to unmarshal event: %w", err)
	}

	log.Printf("Processing KYCVerified event: %s for user %s", event.ID, event.UserID)

	// Create notification
	req := dto.CreateNotificationRequest{
		UserID:  event.UserID,
		Type:    models.NotificationTypeKYCVerified,
		Title:   "KYC Verification Complete",
		Message: "Your KYC verification has been completed successfully. You can now access all features.",
		Data: map[string]interface{}{
			"email":    event.Email,
			"username": event.Username,
		},
	}

	_, err := h.notificationService.CreateNotification(req)
	if err != nil {
		return fmt.Errorf("failed to create notification: %w", err)
	}

	log.Printf("KYCVerified notification created for user %s", event.UserID)
	return nil
}
