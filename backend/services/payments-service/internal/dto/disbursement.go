package dto

import "github.com/google/uuid"

// DisbursementRequest represents a request to disburse funds
type DisbursementRequest struct {
	DisbursementID  uuid.UUID
	UserID          uuid.UUID
	Amount          int64
	Currency        string
	BankCode        string
	AccountNumber   string
	AccountName     string
	Reason          string
}

// DisbursementResponse represents the response from a disbursement provider
type DisbursementResponse struct {
	TransferCode string
	Reference    string
	Status       string
}
