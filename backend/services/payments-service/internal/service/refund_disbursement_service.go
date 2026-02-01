package service

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/gofund/payments-service/internal/dto"
	"github.com/gofund/shared/logger"
	"github.com/gofund/shared/metrics"
)

// RefundDisbursementService handles actual fund disbursement through payment providers (e.g., Paystack)
type RefundDisbursementService struct {
	paystackClient *PaystackClient
}

// NewRefundDisbursementService creates a new refund disbursement service instance
func NewRefundDisbursementService(paystackClient *PaystackClient) *RefundDisbursementService {
	return &RefundDisbursementService{
		paystackClient: paystackClient,
	}
}

// PaystackTransferRecipientRequest represents the request to create a transfer recipient
type PaystackTransferRecipientRequest struct {
	Type          string `json:"type"`
	Name          string `json:"name"`
	AccountNumber string `json:"account_number"`
	BankCode      string `json:"bank_code"`
	Currency      string `json:"currency"`
}

// PaystackTransferRecipientResponse represents the response from creating a transfer recipient
type PaystackTransferRecipientResponse struct {
	Status  bool   `json:"status"`
	Message string `json:"message"`
	Data    struct {
		RecipientCode string `json:"recipient_code"`
		Type          string `json:"type"`
		Name          string `json:"name"`
		AccountNumber string `json:"account_number"`
		BankCode      string `json:"bank_code"`
	} `json:"data"`
}

// PaystackTransferRequest represents the request to initiate a transfer
type PaystackTransferRequest struct {
	Source    string `json:"source"`
	Amount    int64  `json:"amount"`
	Recipient string `json:"recipient"`
	Reason    string `json:"reason"`
	Reference string `json:"reference"`
	Currency  string `json:"currency"`
}

// PaystackTransferResponse represents the response from initiating a transfer
type PaystackTransferResponse struct {
	Status  bool   `json:"status"`
	Message string `json:"message"`
	Data    struct {
		TransferCode string `json:"transfer_code"`
		Reference    string `json:"reference"`
		Status       string `json:"status"`
		Amount       int64  `json:"amount"`
		CreatedAt    string `json:"created_at"`
	} `json:"data"`
}

// PaystackVerifyTransferResponse represents the response from verifying a transfer
type PaystackVerifyTransferResponse struct {
	Status  bool   `json:"status"`
	Message string `json:"message"`
	Data    struct {
		TransferCode string `json:"transfer_code"`
		Reference    string `json:"reference"`
		Status       string `json:"status"`
		Amount       int64  `json:"amount"`
		Reason       string `json:"reason"`
		CreatedAt    string `json:"created_at"`
		UpdatedAt    string `json:"updated_at"`
	} `json:"data"`
}

// InitiateDisbursement initiates a refund disbursement to a user's settlement account
// This uses Paystack's Transfer API to send money back to contributors
func (rds *RefundDisbursementService) InitiateDisbursement(req *dto.DisbursementRequest) (*dto.DisbursementResponse, error) {
	// Generate unique reference for this disbursement
	reference := fmt.Sprintf("REFUND-%s", req.DisbursementID.String())

	logger.Info("Initiating refund disbursement", map[string]interface{}{
		"disbursement_id": req.DisbursementID.String(),
		"user_id":         req.UserID.String(),
		"amount":          req.Amount,
		"account_number":  req.AccountNumber,
		"bank_code":       req.BankCode,
	})

	// Create transfer recipient
	recipientCode, err := rds.createTransferRecipient(
		req.AccountName,
		req.AccountNumber,
		req.BankCode,
		req.Currency,
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
		req.Currency,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to initiate transfer: %w", err)
	}

	metrics.IncrementCounter("refund.disbursement.initiated")

	return &dto.DisbursementResponse{
		TransferCode: transferCode,
		Reference:    reference,
		Status:       "PENDING",
	}, nil
}

// createTransferRecipient creates a transfer recipient on Paystack
func (rds *RefundDisbursementService) createTransferRecipient(name, accountNumber, bankCode, currency string) (string, error) {
	url := fmt.Sprintf("%s/transferrecipient", rds.paystackClient.baseURL)

	reqBody := PaystackTransferRecipientRequest{
		Type:          "nuban",
		Name:          name,
		AccountNumber: accountNumber,
		BankCode:      bankCode,
		Currency:      currency,
	}

	body, err := json.Marshal(reqBody)
	if err != nil {
		return "", fmt.Errorf("failed to marshal request: %w", err)
	}

	httpReq, err := http.NewRequest("POST", url, bytes.NewBuffer(body))
	if err != nil {
		return "", fmt.Errorf("failed to create request: %w", err)
	}

	httpReq.Header.Set("Authorization", fmt.Sprintf("Bearer %s", rds.paystackClient.secretKey))
	httpReq.Header.Set("Content-Type", "application/json")

	startTime := time.Now()
	resp, err := rds.paystackClient.client.Do(httpReq)
	duration := time.Since(startTime).Milliseconds()

	metrics.RecordHistogram("paystack.api.create_recipient.duration", float64(duration))

	if err != nil {
		metrics.IncrementCounter("paystack.api.create_recipient.error")
		return "", fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read response: %w", err)
	}

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		metrics.IncrementCounter("paystack.api.create_recipient.failed")
		logger.Error("Failed to create transfer recipient", map[string]interface{}{
			"status_code": resp.StatusCode,
			"response":    string(respBody),
		})
		return "", fmt.Errorf("paystack API error: status %d, body: %s", resp.StatusCode, string(respBody))
	}

	var paystackResp PaystackTransferRecipientResponse
	if err := json.Unmarshal(respBody, &paystackResp); err != nil {
		return "", fmt.Errorf("failed to unmarshal response: %w", err)
	}

	if !paystackResp.Status {
		metrics.IncrementCounter("paystack.api.create_recipient.failed")
		return "", fmt.Errorf("failed to create recipient: %s", paystackResp.Message)
	}

	metrics.IncrementCounter("paystack.api.create_recipient.success")
	logger.Info("Transfer recipient created successfully", map[string]interface{}{
		"recipient_code": paystackResp.Data.RecipientCode,
		"account_number": accountNumber,
	})

	return paystackResp.Data.RecipientCode, nil
}

// initiateTransfer initiates a transfer on Paystack
func (rds *RefundDisbursementService) initiateTransfer(recipientCode string, amount int64, reference, reason, currency string) (string, error) {
	url := fmt.Sprintf("%s/transfer", rds.paystackClient.baseURL)

	reqBody := PaystackTransferRequest{
		Source:    "balance",
		Amount:    amount,
		Recipient: recipientCode,
		Reason:    reason,
		Reference: reference,
		Currency:  currency,
	}

	body, err := json.Marshal(reqBody)
	if err != nil {
		return "", fmt.Errorf("failed to marshal request: %w", err)
	}

	httpReq, err := http.NewRequest("POST", url, bytes.NewBuffer(body))
	if err != nil {
		return "", fmt.Errorf("failed to create request: %w", err)
	}

	httpReq.Header.Set("Authorization", fmt.Sprintf("Bearer %s", rds.paystackClient.secretKey))
	httpReq.Header.Set("Content-Type", "application/json")

	startTime := time.Now()
	resp, err := rds.paystackClient.client.Do(httpReq)
	duration := time.Since(startTime).Milliseconds()

	metrics.RecordHistogram("paystack.api.initiate_transfer.duration", float64(duration))

	if err != nil {
		metrics.IncrementCounter("paystack.api.initiate_transfer.error")
		return "", fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read response: %w", err)
	}

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		metrics.IncrementCounter("paystack.api.initiate_transfer.failed")
		logger.Error("Failed to initiate transfer", map[string]interface{}{
			"status_code": resp.StatusCode,
			"response":    string(respBody),
		})
		return "", fmt.Errorf("paystack API error: status %d, body: %s", resp.StatusCode, string(respBody))
	}

	var paystackResp PaystackTransferResponse
	if err := json.Unmarshal(respBody, &paystackResp); err != nil {
		return "", fmt.Errorf("failed to unmarshal response: %w", err)
	}

	if !paystackResp.Status {
		metrics.IncrementCounter("paystack.api.initiate_transfer.failed")
		return "", fmt.Errorf("failed to initiate transfer: %s", paystackResp.Message)
	}

	metrics.IncrementCounter("paystack.api.initiate_transfer.success")
	logger.Info("Transfer initiated successfully", map[string]interface{}{
		"transfer_code": paystackResp.Data.TransferCode,
		"reference":     reference,
		"amount":        amount,
	})

	return paystackResp.Data.TransferCode, nil
}

// VerifyDisbursement verifies the status of a disbursement
func (rds *RefundDisbursementService) VerifyDisbursement(transferCode string) (string, error) {
	url := fmt.Sprintf("%s/transfer/%s", rds.paystackClient.baseURL, transferCode)

	httpReq, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return "", fmt.Errorf("failed to create request: %w", err)
	}

	httpReq.Header.Set("Authorization", fmt.Sprintf("Bearer %s", rds.paystackClient.secretKey))

	resp, err := rds.paystackClient.client.Do(httpReq)
	if err != nil {
		return "", fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read response: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("paystack API error: status %d, body: %s", resp.StatusCode, string(respBody))
	}

	var paystackResp PaystackVerifyTransferResponse
	if err := json.Unmarshal(respBody, &paystackResp); err != nil {
		return "", fmt.Errorf("failed to unmarshal response: %w", err)
	}

	if !paystackResp.Status {
		return "", fmt.Errorf("failed to verify transfer: %s", paystackResp.Message)
	}

	logger.Info("Transfer verified", map[string]interface{}{
		"transfer_code": transferCode,
		"status":        paystackResp.Data.Status,
	})

	// Return status: "pending", "success", "failed", "reversed"
	return paystackResp.Data.Status, nil
}

// GetBankList retrieves list of supported banks from Paystack
func (rds *RefundDisbursementService) GetBankList() ([]dto.Bank, error) {
	// Delegate to the Paystack client's ListBanks method
	paystackResp, err := rds.paystackClient.ListBanks("nigeria")
	if err != nil {
		return nil, err
	}

	banks := make([]dto.Bank, 0, len(paystackResp.Data))
	for _, bank := range paystackResp.Data {
		if bank.Active && !bank.IsDeleted {
			banks = append(banks, dto.Bank{
				ID:   int(bank.ID),
				Code: bank.Code,
				Name: bank.Name,
			})
		}
	}

	return banks, nil
}

// ResolveAccountNumber resolves an account number to get the account name
func (rds *RefundDisbursementService) ResolveAccountNumber(accountNumber, bankCode string) (string, error) {
	// Delegate to the Paystack client's ResolveAccountNumber method
	paystackResp, err := rds.paystackClient.ResolveAccountNumber(accountNumber, bankCode)
	if err != nil {
		return "", err
	}

	return paystackResp.Data.AccountName, nil
}

