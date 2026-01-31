package service

import (
	"fmt"

	"github.com/gofund/payments-service/internal/dto"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

// RefundDisbursementService handles actual fund disbursement through payment providers (e.g., Paystack)
type RefundDisbursementService struct {
	db *gorm.DB
	// paystackClient *paystack.Client // Will be added when implementing Paystack
}

// NewRefundDisbursementService creates a new refund disbursement service instance
func NewRefundDisbursementService(db *gorm.DB) *RefundDisbursementService {
	return &RefundDisbursementService{
		db: db,
	}
}

// InitiateDisbursement initiates a refund disbursement to a user's settlement account
// This uses Paystack's Transfer API to send money back to contributors
func (rds *RefundDisbursementService) InitiateDisbursement(req *dto.DisbursementRequest) (*dto.DisbursementResponse, error) {
	// Generate unique reference for this disbursement
	reference := fmt.Sprintf("REFUND-%s", req.DisbursementID.String())

	// Create transfer recipient
	recipientCode, err := rds.createTransferRecipient(
		req.AccountName,
		req.AccountNumber,
		req.BankCode,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create transfer recipient: %w", err)
	}

	// Initiate transfer
	transferCode, err := rds.initiateTransfer(
		recipientCode,
		req.Amount,
		reference,
		req.Reason,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to initiate transfer: %w", err)
	}

	return &dto.DisbursementResponse{
		TransferCode: transferCode,
		Reference:    reference,
		Status:       "PENDING",
	}, nil
}

// createTransferRecipient creates a transfer recipient on Paystack
func (rds *RefundDisbursementService) createTransferRecipient(name, accountNumber, bankCode string) (string, error) {
	// TODO: Implement Paystack API call
	// POST https://api.paystack.co/transferrecipient
	// Headers: Authorization: Bearer SECRET_KEY
	
	// For now, return a mock recipient code
	return fmt.Sprintf("RCP_%s", uuid.New().String()[:8]), nil
}

// initiateTransfer initiates a transfer on Paystack
func (rds *RefundDisbursementService) initiateTransfer(recipientCode string, amount int64, reference, reason string) (string, error) {
	// TODO: Implement Paystack API call
	// POST https://api.paystack.co/transfer
	// Headers: Authorization: Bearer SECRET_KEY
	
	// For now, return a mock transfer code
	return fmt.Sprintf("TRF_%s", uuid.New().String()[:8]), nil
}

// VerifyDisbursement verifies the status of a disbursement
func (rds *RefundDisbursementService) VerifyDisbursement(transferCode string) (string, error) {
	// TODO: Implement Paystack API call
	// GET https://api.paystack.co/transfer/{transfer_code}
	// Headers: Authorization: Bearer SECRET_KEY
	
	// Return status: "pending", "success", "failed", "reversed"
	return "success", nil
}

// GetBankList retrieves list of supported banks from Paystack
func (rds *RefundDisbursementService) GetBankList() ([]dto.Bank, error) {
	// TODO: Implement Paystack API call
	// GET https://api.paystack.co/bank
	// Headers: Authorization: Bearer SECRET_KEY
	
	return []dto.Bank{
		{Code: "057", Name: "Zenith Bank"},
		{Code: "058", Name: "GTBank"},
		{Code: "033", Name: "United Bank for Africa"},
		// ... more banks
	}, nil
}

// ResolveAccountNumber resolves an account number to get the account name
func (rds *RefundDisbursementService) ResolveAccountNumber(accountNumber, bankCode string) (string, error) {
	// TODO: Implement Paystack API call
	// GET https://api.paystack.co/bank/resolve?account_number=xxx&bank_code=057
	// Headers: Authorization: Bearer SECRET_KEY
	
	// For now, return a mock account name
	return "Account Holder Name", nil
}
