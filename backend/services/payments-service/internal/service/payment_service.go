package service

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/gofund/payments-service/internal/dto"
	"github.com/gofund/payments-service/internal/repository"
	"github.com/gofund/shared/events"
	"github.com/gofund/shared/messaging"
	"github.com/gofund/shared/metrics"
	"github.com/gofund/shared/models"
	"github.com/google/uuid"
)

// PaymentService handles payment business logic
type PaymentService struct {
	paymentRepo     *repository.PaymentRepository
	idempotencyRepo *repository.IdempotencyRepository
	paystackClient  *PaystackClient
	eventPublisher  messaging.Publisher
}

// NewPaymentService creates a new payment service
func NewPaymentService(
	paymentRepo *repository.PaymentRepository,
	idempotencyRepo *repository.IdempotencyRepository,
	paystackClient *PaystackClient,
	eventPublisher messaging.Publisher,
) *PaymentService {
	return &PaymentService{
		paymentRepo:     paymentRepo,
		idempotencyRepo: idempotencyRepo,
		paystackClient:  paystackClient,
		eventPublisher:  eventPublisher,
	}
}

// InitializePayment initializes a new payment with Paystack
func (ps *PaymentService) InitializePayment(ctx context.Context, req *dto.InitializePaymentRequest) (*dto.InitializePaymentResponse, error) {
	// Generate unique payment ID and reference
	paymentID := uuid.New().String()
	reference := fmt.Sprintf("PAY-%s", uuid.New().String()[:13])

	// Create payment record with INITIATED status
	payment := &models.Payment{
		PaymentID:         paymentID,
		PaystackReference: reference,
		UserID:            req.UserID.String(),
		GoalID:            req.GoalID.String(),
		Amount:            req.Amount,
		Currency:          req.Currency,
		Status:            models.PaymentStatusInitiated,
	}

	if err := ps.paymentRepo.CreatePayment(ctx, payment); err != nil {
		log.Printf("[ERROR] Failed to create payment record: %v (user_id: %s, goal_id: %s)",
			err, req.UserID.String(), req.GoalID.String())
		metrics.IncrementCounter("payment.creation.failed")
		return nil, fmt.Errorf("failed to create payment: %w", err)
	}

	// Prepare Paystack initialization request
	paystackReq := &dto.PaystackInitializeRequest{
		Email:       req.Email,
		Amount:      req.Amount,
		Currency:    req.Currency,
		Reference:   reference,
		CallbackURL: req.CallbackURL,
		Metadata: map[string]interface{}{
			"payment_id": paymentID,
			"user_id":    req.UserID.String(),
			"goal_id":    req.GoalID.String(),
		},
		Channels: []string{"card", "bank", "ussd", "qr", "mobile_money", "bank_transfer"},
	}

	// Add custom metadata if provided
	if req.Metadata != nil {
		for k, v := range req.Metadata {
			paystackReq.Metadata[k] = v
		}
	}

	// Initialize transaction with Paystack
	paystackResp, err := ps.paystackClient.InitializeTransaction(paystackReq)
	if err != nil {
		log.Printf("[ERROR] Paystack initialization failed: %v (payment_id: %s)", err, paymentID)
		metrics.IncrementCounter("payment.initialization.failed")

		// Update payment status to FAILED
		payment.Status = models.PaymentStatusFailed
		ps.paymentRepo.UpdatePayment(ctx, payment)

		return nil, fmt.Errorf("failed to initialize payment with Paystack: %w", err)
	}

	// Update payment with Paystack data
	payment.Status = models.PaymentStatusPending
	payment.PaystackData = map[string]interface{}{
		"authorization_url": paystackResp.Data.AuthorizationURL,
		"access_code":       paystackResp.Data.AccessCode,
	}

	if err := ps.paymentRepo.UpdatePayment(ctx, payment); err != nil {
		log.Printf("[ERROR] Failed to update payment with Paystack data: %v (payment_id: %s)",
			err, paymentID)
		return nil, fmt.Errorf("failed to update payment: %w", err)
	}

	// Track metrics
	metrics.IncrementCounter("payment.initialized.count")

	log.Printf("[INFO] Payment initialized successfully (payment_id: %s, reference: %s, url: %s)",
		paymentID, reference, paystackResp.Data.AuthorizationURL)

	return &dto.InitializePaymentResponse{
		PaymentID:        paymentID,
		AuthorizationURL: paystackResp.Data.AuthorizationURL,
		AccessCode:       paystackResp.Data.AccessCode,
		Reference:        reference,
	}, nil
}

// VerifyPayment verifies a payment with Paystack (instant verification)
func (ps *PaymentService) VerifyPayment(ctx context.Context, reference string) (*dto.VerifyPaymentResponse, error) {
	// Step 1: Get payment by reference
	payment, err := ps.paymentRepo.GetPaymentByReference(ctx, reference)
	if err != nil {
		log.Printf("[ERROR] Payment not found: %v (reference: %s)", err, reference)
		return nil, fmt.Errorf("payment not found: %w", err)
	}

	// Step 2: Check if already verified (idempotency)
	if payment.Status == models.PaymentStatusVerified {
		log.Printf("[INFO] Payment already verified (payment_id: %s, reference: %s)",
			payment.PaymentID, reference)

		return ps.mapPaymentToVerifyResponse(payment), nil
	}

	// Step 3: Verify with Paystack
	paystackResp, err := ps.paystackClient.VerifyTransaction(reference)
	if err != nil {
		log.Printf("[ERROR] Paystack verification failed: %v (payment_id: %s, reference: %s)",
			err, payment.PaymentID, reference)
		return nil, fmt.Errorf("failed to verify payment with Paystack: %w", err)
	}

	// Step 4: Update payment status based on Paystack response
	if paystackResp.Data.Status == "success" {
		payment.Status = models.PaymentStatusVerified
		payment.PaystackData = map[string]interface{}{
			"id":               paystackResp.Data.ID,
			"status":           paystackResp.Data.Status,
			"reference":        paystackResp.Data.Reference,
			"amount":           paystackResp.Data.Amount,
			"paid_at":          paystackResp.Data.PaidAt,
			"channel":          paystackResp.Data.Channel,
			"currency":         paystackResp.Data.Currency,
			"gateway_response": paystackResp.Data.GatewayResponse,
			"customer":         paystackResp.Data.Customer,
		}

		if err := ps.paymentRepo.UpdatePayment(ctx, payment); err != nil {
			log.Printf("[ERROR] Failed to update payment status: %v (payment_id: %s)",
				err, payment.PaymentID)
			return nil, fmt.Errorf("failed to update payment: %w", err)
		}

		// Step 5: Emit PaymentVerified event
		if err := ps.emitPaymentVerifiedEvent(payment); err != nil {
			// Log error but don't fail the request
			log.Printf("[ERROR] Failed to emit PaymentVerified event: %v (payment_id: %s)",
				err, payment.PaymentID)
		}

		// Track metrics
		metrics.IncrementCounter("payment.verified.count")

		log.Printf("[INFO] Payment verified successfully (payment_id: %s, reference: %s, amount: %d, channel: %s)",
			payment.PaymentID, reference, payment.Amount, paystackResp.Data.Channel)

	} else {
		// Payment failed
		payment.Status = models.PaymentStatusFailed
		payment.PaystackData = map[string]interface{}{
			"status":           paystackResp.Data.Status,
			"gateway_response": paystackResp.Data.GatewayResponse,
		}

		ps.paymentRepo.UpdatePayment(ctx, payment)

		metrics.IncrementCounter("payment.failed.count")

		log.Printf("[INFO] Payment verification failed (payment_id: %s, reference: %s, status: %s)",
			payment.PaymentID, reference, paystackResp.Data.Status)
	}

	return ps.mapPaymentToVerifyResponse(payment), nil
}

// GetPaymentStatus retrieves the current status of a payment
func (ps *PaymentService) GetPaymentStatus(ctx context.Context, paymentID string) (*dto.PaymentStatusResponse, error) {
	payment, err := ps.paymentRepo.GetPaymentByID(ctx, paymentID)
	if err != nil {
		return nil, fmt.Errorf("payment not found: %w", err)
	}

	return ps.mapPaymentToStatusResponse(payment), nil
}

// ListBanks retrieves the list of supported banks
func (ps *PaymentService) ListBanks(ctx context.Context, country string) ([]dto.Bank, error) {
	if country == "" {
		country = "nigeria" // Default to Nigeria
	}

	paystackResp, err := ps.paystackClient.ListBanks(country)
	if err != nil {
		log.Printf("[ERROR] Failed to list banks: %v (country: %s)", err, country)
		return nil, fmt.Errorf("failed to list banks: %w", err)
	}

	// Map Paystack banks to our DTO
	banks := make([]dto.Bank, 0, len(paystackResp.Data))
	for _, bank := range paystackResp.Data {
		if bank.Active && !bank.IsDeleted {
			banks = append(banks, dto.Bank{
				ID:   int(bank.ID),
				Name: bank.Name,
				Code: bank.Code,
			})
		}
	}

	return banks, nil
}

// ResolveAccount resolves an account number to get account name
func (ps *PaymentService) ResolveAccount(ctx context.Context, req *dto.ResolveAccountRequest) (*dto.ResolveAccountResponse, error) {
	paystackResp, err := ps.paystackClient.ResolveAccountNumber(req.AccountNumber, req.BankCode)
	if err != nil {
		log.Printf("[ERROR] Failed to resolve account: %v (account: %s, bank: %s)",
			err, req.AccountNumber, req.BankCode)
		return nil, fmt.Errorf("failed to resolve account: %w", err)
	}

	// Get bank name from bank list
	banks, _ := ps.ListBanks(ctx, "nigeria")
	bankName := ""
	for _, bank := range banks {
		if bank.Code == req.BankCode {
			bankName = bank.Name
			break
		}
	}

	return &dto.ResolveAccountResponse{
		AccountNumber: paystackResp.Data.AccountNumber,
		AccountName:   paystackResp.Data.AccountName,
		BankCode:      req.BankCode,
		BankName:      bankName,
	}, nil
}

// emitPaymentVerifiedEvent emits a PaymentVerified event
func (ps *PaymentService) emitPaymentVerifiedEvent(payment *models.Payment) error {
	event := events.PaymentVerified{
		ID:        uuid.New().String(),
		PaymentID: payment.PaymentID,
		UserID:    payment.UserID,
		GoalID:    payment.GoalID,
		Amount:    payment.Amount,
		CreatedAt: time.Now().Unix(),
	}

	if err := ps.eventPublisher.Publish("PaymentVerified", event); err != nil {
		return fmt.Errorf("failed to publish event: %w", err)
	}

	log.Printf("[INFO] PaymentVerified event emitted (event_id: %s, payment_id: %s, user_id: %s, goal_id: %s, amount: %d)",
		event.ID, payment.PaymentID, payment.UserID, payment.GoalID, payment.Amount)

	return nil
}

// mapPaymentToVerifyResponse maps a payment to verify response
func (ps *PaymentService) mapPaymentToVerifyResponse(payment *models.Payment) *dto.VerifyPaymentResponse {
	resp := &dto.VerifyPaymentResponse{
		PaymentID: payment.PaymentID,
		Reference: payment.PaystackReference,
		Status:    string(payment.Status),
		Amount:    payment.Amount,
		Currency:  payment.Currency,
		Metadata:  payment.PaystackData,
	}

	// Extract paid_at if available
	if payment.PaystackData != nil {
		if paidAt, ok := payment.PaystackData["paid_at"].(string); ok && paidAt != "" {
			resp.PaidAt = &paidAt
		}
		if channel, ok := payment.PaystackData["channel"].(string); ok {
			resp.Channel = channel
		}
	}

	return resp
}

// mapPaymentToStatusResponse maps a payment to status response
func (ps *PaymentService) mapPaymentToStatusResponse(payment *models.Payment) *dto.PaymentStatusResponse {
	resp := &dto.PaymentStatusResponse{
		PaymentID: payment.PaymentID,
		Reference: payment.PaystackReference,
		Status:    string(payment.Status),
		Amount:    payment.Amount,
		Currency:  payment.Currency,
		CreatedAt: payment.CreatedAt.Format(time.RFC3339),
		UpdatedAt: payment.UpdatedAt.Format(time.RFC3339),
		Metadata:  payment.PaystackData,
	}

	// Extract paid_at if available
	if payment.PaystackData != nil {
		if paidAt, ok := payment.PaystackData["paid_at"].(string); ok && paidAt != "" {
			resp.PaidAt = &paidAt
		}
	}

	return resp
}
