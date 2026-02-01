package service

import (
	"context"
	"fmt"
	"time"

	"github.com/gofund/payments-service/internal/dto"
	"github.com/gofund/payments-service/internal/repository"
	"github.com/gofund/shared/events"
	"github.com/gofund/shared/logger"
	"github.com/gofund/shared/messaging"
	"github.com/gofund/shared/metrics"
	"github.com/gofund/shared/models"
	"github.com/google/uuid"
)

// WebhookService handles webhook processing
type WebhookService struct {
	webhookRepo     *repository.WebhookRepository
	paymentRepo     *repository.PaymentRepository
	eventPublisher  *messaging.EventPublisher
}

// NewWebhookService creates a new webhook service
func NewWebhookService(
	webhookRepo *repository.WebhookRepository,
	paymentRepo *repository.PaymentRepository,
	eventPublisher *messaging.EventPublisher,
) *WebhookService {
	return &WebhookService{
		webhookRepo:    webhookRepo,
		paymentRepo:    paymentRepo,
		eventPublisher: eventPublisher,
	}
}

// ProcessWebhook processes a Paystack webhook event
func (ws *WebhookService) ProcessWebhook(ctx context.Context, payload *dto.WebhookPayload, signature string) error {
	// Generate event ID from Paystack data
	eventID := ws.generateEventID(payload)

	logger.Info("Processing webhook event", map[string]interface{}{
		"event_id":   eventID,
		"event_type": payload.Event,
	})

	// Step 1: Check idempotency (has this event been processed?)
	if ws.webhookRepo.IsEventProcessed(ctx, eventID) {
		logger.Info("Webhook event already processed, skipping", map[string]interface{}{
			"event_id":   eventID,
			"event_type": payload.Event,
		})
		metrics.IncrementCounter("webhook.duplicate.count")
		return nil // Not an error - just already processed
	}

	// Step 2: Save webhook event
	webhookEvent := &models.WebhookEvent{
		EventID:   eventID,
		Event:     payload.Event,
		Data:      payload.Data,
		Signature: signature,
		Processed: false,
	}

	if err := ws.webhookRepo.SaveWebhookEvent(ctx, webhookEvent); err != nil {
		logger.Error("Failed to save webhook event", map[string]interface{}{
			"error":      err.Error(),
			"event_id":   eventID,
			"event_type": payload.Event,
		})
		return fmt.Errorf("failed to save webhook: %w", err)
	}

	// Step 3: Process based on event type
	var processErr error
	switch payload.Event {
	case "charge.success":
		processErr = ws.processChargeSuccess(ctx, payload.Data)
	case "charge.failed":
		processErr = ws.processChargeFailed(ctx, payload.Data)
	case "transfer.success":
		processErr = ws.processTransferSuccess(ctx, payload.Data)
	case "transfer.failed":
		processErr = ws.processTransferFailed(ctx, payload.Data)
	default:
		logger.Info("Unhandled webhook event type", map[string]interface{}{
			"event_type": payload.Event,
		})
		// Mark as processed even if we don't handle it
		ws.webhookRepo.MarkWebhookProcessed(ctx, eventID)
		return nil
	}

	if processErr != nil {
		logger.Error("Failed to process webhook event", map[string]interface{}{
			"error":      processErr.Error(),
			"event_id":   eventID,
			"event_type": payload.Event,
		})
		return processErr
	}

	// Step 4: Mark webhook as processed
	if err := ws.webhookRepo.MarkWebhookProcessed(ctx, eventID); err != nil {
		logger.Error("Failed to mark webhook as processed", map[string]interface{}{
			"error":    err.Error(),
			"event_id": eventID,
		})
		// Don't return error - event was processed successfully
	}

	metrics.IncrementCounter("webhook.processed.count")
	logger.Info("Webhook event processed successfully", map[string]interface{}{
		"event_id":   eventID,
		"event_type": payload.Event,
	})

	return nil
}

// processChargeSuccess handles charge.success webhook
func (ws *WebhookService) processChargeSuccess(ctx context.Context, data map[string]interface{}) error {
	// Extract reference
	reference, ok := data["reference"].(string)
	if !ok {
		return fmt.Errorf("missing or invalid reference in webhook data")
	}

	logger.Info("Processing charge.success webhook", map[string]interface{}{
		"reference": reference,
	})

	// Get payment by reference
	payment, err := ws.paymentRepo.GetPaymentByReference(ctx, reference)
	if err != nil {
		logger.Error("Payment not found for webhook", map[string]interface{}{
			"error":     err.Error(),
			"reference": reference,
		})
		return fmt.Errorf("payment not found: %w", err)
	}

	// Check if already verified (idempotency - might have been verified via API)
	if payment.Status == models.PaymentStatusVerified {
		logger.Info("Payment already verified, webhook is backup confirmation", map[string]interface{}{
			"payment_id": payment.PaymentID,
			"reference":  reference,
		})
		return nil // Success - no action needed
	}

	// Update payment status to VERIFIED
	payment.Status = models.PaymentStatusVerified
	payment.PaystackData = data

	if err := ws.paymentRepo.UpdatePayment(ctx, payment); err != nil {
		logger.Error("Failed to update payment status", map[string]interface{}{
			"error":      err.Error(),
			"payment_id": payment.PaymentID,
		})
		return fmt.Errorf("failed to update payment: %w", err)
	}

	// Emit PaymentVerified event
	if err := ws.emitPaymentVerifiedEvent(payment); err != nil {
		logger.Error("Failed to emit PaymentVerified event", map[string]interface{}{
			"error":      err.Error(),
			"payment_id": payment.PaymentID,
		})
		// Don't fail - payment is already verified
	}

	metrics.IncrementCounter("webhook.payment.verified.count")
	logger.Info("Payment verified via webhook", map[string]interface{}{
		"payment_id": payment.PaymentID,
		"reference":  reference,
		"amount":     payment.Amount,
	})

	return nil
}

// processChargeFailed handles charge.failed webhook
func (ws *WebhookService) processChargeFailed(ctx context.Context, data map[string]interface{}) error {
	reference, ok := data["reference"].(string)
	if !ok {
		return fmt.Errorf("missing or invalid reference in webhook data")
	}

	logger.Info("Processing charge.failed webhook", map[string]interface{}{
		"reference": reference,
	})

	// Get payment by reference
	payment, err := ws.paymentRepo.GetPaymentByReference(ctx, reference)
	if err != nil {
		return fmt.Errorf("payment not found: %w", err)
	}

	// Update payment status to FAILED
	payment.Status = models.PaymentStatusFailed
	payment.PaystackData = data

	if err := ws.paymentRepo.UpdatePayment(ctx, payment); err != nil {
		return fmt.Errorf("failed to update payment: %w", err)
	}

	metrics.IncrementCounter("webhook.payment.failed.count")
	logger.Info("Payment marked as failed via webhook", map[string]interface{}{
		"payment_id": payment.PaymentID,
		"reference":  reference,
	})

	return nil
}

// processTransferSuccess handles transfer.success webhook (for refunds)
func (ws *WebhookService) processTransferSuccess(ctx context.Context, data map[string]interface{}) error {
	reference, ok := data["reference"].(string)
	if !ok {
		return fmt.Errorf("missing or invalid reference in webhook data")
	}

	logger.Info("Processing transfer.success webhook", map[string]interface{}{
		"reference": reference,
	})

	// TODO: Update refund disbursement status
	// This will be handled by the refund disbursement service

	metrics.IncrementCounter("webhook.transfer.success.count")
	return nil
}

// processTransferFailed handles transfer.failed webhook (for refunds)
func (ws *WebhookService) processTransferFailed(ctx context.Context, data map[string]interface{}) error {
	reference, ok := data["reference"].(string)
	if !ok {
		return fmt.Errorf("missing or invalid reference in webhook data")
	}

	logger.Info("Processing transfer.failed webhook", map[string]interface{}{
		"reference": reference,
	})

	// TODO: Update refund disbursement status
	// This will be handled by the refund disbursement service

	metrics.IncrementCounter("webhook.transfer.failed.count")
	return nil
}

// emitPaymentVerifiedEvent emits a PaymentVerified event
func (ws *WebhookService) emitPaymentVerifiedEvent(payment *models.Payment) error {
	event := events.PaymentVerified{
		ID:        uuid.New().String(),
		PaymentID: payment.PaymentID,
		UserID:    payment.UserID,
		GoalID:    payment.GoalID,
		Amount:    payment.Amount,
		CreatedAt: time.Now().Unix(),
	}

	if err := ws.eventPublisher.Publish("PaymentVerified", event); err != nil {
		return fmt.Errorf("failed to publish event: %w", err)
	}

	logger.Info("PaymentVerified event emitted from webhook", map[string]interface{}{
		"event_id":   event.ID,
		"payment_id": payment.PaymentID,
		"user_id":    payment.UserID,
		"goal_id":    payment.GoalID,
		"amount":     payment.Amount,
	})

	return nil
}

// generateEventID generates a unique event ID from webhook data
func (ws *WebhookService) generateEventID(payload *dto.WebhookPayload) string {
	// Try to extract Paystack's event ID if available
	if id, ok := payload.Data["id"]; ok {
		return fmt.Sprintf("%s-%v", payload.Event, id)
	}

	// Fallback to reference-based ID
	if reference, ok := payload.Data["reference"].(string); ok {
		return fmt.Sprintf("%s-%s", payload.Event, reference)
	}

	// Last resort: generate UUID
	return fmt.Sprintf("%s-%s", payload.Event, uuid.New().String())
}
