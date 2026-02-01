package service

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"

	"github.com/gofund/payments-service/internal/dto"
	"github.com/gofund/shared/metrics"
)

// PaystackClient handles communication with Paystack API
type PaystackClient struct {
	secretKey string
	baseURL   string
	client    *http.Client
}

// NewPaystackClient creates a new Paystack API client
func NewPaystackClient(secretKey, baseURL string) *PaystackClient {
	return &PaystackClient{
		secretKey: secretKey,
		baseURL:   baseURL,
		client: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

// InitializeTransaction initializes a transaction with Paystack
func (pc *PaystackClient) InitializeTransaction(req *dto.PaystackInitializeRequest) (*dto.PaystackInitializeResponse, error) {
	url := fmt.Sprintf("%s/transaction/initialize", pc.baseURL)

	// Marshal request body
	body, err := json.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	// Create HTTP request
	httpReq, err := http.NewRequest("POST", url, bytes.NewBuffer(body))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	// Set headers
	httpReq.Header.Set("Authorization", fmt.Sprintf("Bearer %s", pc.secretKey))
	httpReq.Header.Set("Content-Type", "application/json")

	// Log request
	log.Printf("[INFO] Initializing Paystack transaction (reference: %s, amount: %d, email: %s)",
		req.Reference, req.Amount, req.Email)

	// Send request
	startTime := time.Now()
	resp, err := pc.client.Do(httpReq)
	duration := time.Since(startTime).Milliseconds()

	// Track metrics
	metrics.RecordHistogram("paystack.api.initialize.duration", float64(duration))

	if err != nil {
		metrics.IncrementCounter("paystack.api.initialize.error")
		return nil, fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	// Read response body
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	// Check status code
	if resp.StatusCode != http.StatusOK {
		metrics.IncrementCounter("paystack.api.initialize.failed")
		log.Printf("[ERROR] Paystack initialization failed (status: %d, response: %s)",
			resp.StatusCode, string(respBody))
		return nil, fmt.Errorf("paystack API error: status %d, body: %s", resp.StatusCode, string(respBody))
	}

	// Parse response
	var paystackResp dto.PaystackInitializeResponse
	if err := json.Unmarshal(respBody, &paystackResp); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	if !paystackResp.Status {
		metrics.IncrementCounter("paystack.api.initialize.failed")
		return nil, fmt.Errorf("paystack initialization failed: %s", paystackResp.Message)
	}

	metrics.IncrementCounter("paystack.api.initialize.success")
	log.Printf("[INFO] Paystack transaction initialized successfully (reference: %s, access_code: %s)",
		req.Reference, paystackResp.Data.AccessCode)

	return &paystackResp, nil
}

// VerifyTransaction verifies a transaction with Paystack
func (pc *PaystackClient) VerifyTransaction(reference string) (*dto.PaystackVerifyResponse, error) {
	url := fmt.Sprintf("%s/transaction/verify/%s", pc.baseURL, reference)

	// Create HTTP request
	httpReq, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	// Set headers
	httpReq.Header.Set("Authorization", fmt.Sprintf("Bearer %s", pc.secretKey))

	// Log request
	log.Printf("[INFO] Verifying Paystack transaction (reference: %s)", reference)

	// Send request
	startTime := time.Now()
	resp, err := pc.client.Do(httpReq)
	duration := time.Since(startTime).Milliseconds()

	// Track metrics
	metrics.RecordHistogram("paystack.api.verify.duration", float64(duration))

	if err != nil {
		metrics.IncrementCounter("paystack.api.verify.error")
		return nil, fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	// Read response body
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	// Check status code
	if resp.StatusCode != http.StatusOK {
		metrics.IncrementCounter("paystack.api.verify.failed")
		log.Printf("[ERROR] Paystack verification failed (status: %d, response: %s)",
			resp.StatusCode, string(respBody))
		return nil, fmt.Errorf("paystack API error: status %d, body: %s", resp.StatusCode, string(respBody))
	}

	// Parse response
	var paystackResp dto.PaystackVerifyResponse
	if err := json.Unmarshal(respBody, &paystackResp); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	if !paystackResp.Status {
		metrics.IncrementCounter("paystack.api.verify.failed")
		return nil, fmt.Errorf("paystack verification failed: %s", paystackResp.Message)
	}

	metrics.IncrementCounter("paystack.api.verify.success")
	log.Printf("[INFO] Paystack transaction verified successfully (reference: %s, status: %s, amount: %d)",
		reference, paystackResp.Data.Status, paystackResp.Data.Amount)

	return &paystackResp, nil
}

// ListBanks retrieves the list of supported banks from Paystack
func (pc *PaystackClient) ListBanks(country string) (*dto.PaystackBankListResponse, error) {
	url := fmt.Sprintf("%s/bank?country=%s&perPage=100", pc.baseURL, country)

	// Create HTTP request
	httpReq, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	// Set headers
	httpReq.Header.Set("Authorization", fmt.Sprintf("Bearer %s", pc.secretKey))

	// Send request
	resp, err := pc.client.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	// Read response body
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	// Check status code
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("paystack API error: status %d, body: %s", resp.StatusCode, string(respBody))
	}

	// Parse response
	var paystackResp dto.PaystackBankListResponse
	if err := json.Unmarshal(respBody, &paystackResp); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	if !paystackResp.Status {
		return nil, fmt.Errorf("paystack list banks failed: %s", paystackResp.Message)
	}

	log.Printf("[INFO] Retrieved bank list from Paystack (country: %s, count: %d)",
		country, len(paystackResp.Data))

	return &paystackResp, nil
}

// ResolveAccountNumber resolves an account number to get account name
func (pc *PaystackClient) ResolveAccountNumber(accountNumber, bankCode string) (*dto.PaystackResolveAccountResponse, error) {
	url := fmt.Sprintf("%s/bank/resolve?account_number=%s&bank_code=%s", pc.baseURL, accountNumber, bankCode)

	// Create HTTP request
	httpReq, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	// Set headers
	httpReq.Header.Set("Authorization", fmt.Sprintf("Bearer %s", pc.secretKey))

	// Log request
	log.Printf("[INFO] Resolving account number (account: %s, bank: %s)",
		accountNumber, bankCode)

	// Send request
	startTime := time.Now()
	resp, err := pc.client.Do(httpReq)
	duration := time.Since(startTime).Milliseconds()

	// Track metrics
	metrics.RecordHistogram("paystack.api.resolve_account.duration", float64(duration))

	if err != nil {
		metrics.IncrementCounter("paystack.api.resolve_account.error")
		return nil, fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	// Read response body
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	// Check status code
	if resp.StatusCode != http.StatusOK {
		metrics.IncrementCounter("paystack.api.resolve_account.failed")
		log.Printf("[ERROR] Account resolution failed (status: %d, account: %s, bank: %s)",
			resp.StatusCode, accountNumber, bankCode)
		return nil, fmt.Errorf("paystack API error: status %d, body: %s", resp.StatusCode, string(respBody))
	}

	// Parse response
	var paystackResp dto.PaystackResolveAccountResponse
	if err := json.Unmarshal(respBody, &paystackResp); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	if !paystackResp.Status {
		metrics.IncrementCounter("paystack.api.resolve_account.failed")
		return nil, fmt.Errorf("account resolution failed: %s", paystackResp.Message)
	}

	metrics.IncrementCounter("paystack.api.resolve_account.success")
	log.Printf("[INFO] Account resolved successfully (account: %s, name: %s)",
		accountNumber, paystackResp.Data.AccountName)

	return &paystackResp, nil
}
