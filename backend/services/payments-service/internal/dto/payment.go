package dto

import "github.com/google/uuid"

// InitializePaymentRequest represents a request to initialize a payment
type InitializePaymentRequest struct {
	UserID      uuid.UUID              `json:"user_id" binding:"required"`
	GoalID      uuid.UUID              `json:"goal_id" binding:"required"`
	Amount      int64                  `json:"amount" binding:"required,min=100"` // Minimum 1 NGN (100 kobo)
	Currency    string                 `json:"currency" binding:"required"`
	Email       string                 `json:"email" binding:"required,email"`
	CallbackURL string                 `json:"callback_url"`
	Metadata    map[string]interface{} `json:"metadata"`
}

// InitializePaymentResponse represents the response from payment initialization
type InitializePaymentResponse struct {
	PaymentID        string `json:"payment_id"`
	AuthorizationURL string `json:"authorization_url"`
	AccessCode       string `json:"access_code"`
	Reference        string `json:"reference"`
}

// PaymentStatusResponse represents the payment status
type PaymentStatusResponse struct {
	PaymentID string                 `json:"payment_id"`
	Reference string                 `json:"reference"`
	Status    string                 `json:"status"`
	Amount    int64                  `json:"amount"`
	Currency  string                 `json:"currency"`
	PaidAt    *string                `json:"paid_at,omitempty"`
	CreatedAt string                 `json:"created_at"`
	UpdatedAt string                 `json:"updated_at"`
	Metadata  map[string]interface{} `json:"metadata,omitempty"`
}

// VerifyPaymentResponse represents the response from payment verification
type VerifyPaymentResponse struct {
	PaymentID string                 `json:"payment_id"`
	Reference string                 `json:"reference"`
	Status    string                 `json:"status"`
	Amount    int64                  `json:"amount"`
	Currency  string                 `json:"currency"`
	PaidAt    *string                `json:"paid_at,omitempty"`
	Channel   string                 `json:"channel,omitempty"`
	Metadata  map[string]interface{} `json:"metadata,omitempty"`
}

// PaystackInitializeRequest represents Paystack transaction initialization request
type PaystackInitializeRequest struct {
	Email       string                 `json:"email"`
	Amount      int64                  `json:"amount"` // Amount in kobo
	Currency    string                 `json:"currency"`
	Reference   string                 `json:"reference"`
	CallbackURL string                 `json:"callback_url,omitempty"`
	Metadata    map[string]interface{} `json:"metadata,omitempty"`
	Channels    []string               `json:"channels,omitempty"`
}

// PaystackInitializeResponse represents Paystack transaction initialization response
type PaystackInitializeResponse struct {
	Status  bool   `json:"status"`
	Message string `json:"message"`
	Data    struct {
		AuthorizationURL string `json:"authorization_url"`
		AccessCode       string `json:"access_code"`
		Reference        string `json:"reference"`
	} `json:"data"`
}

// PaystackVerifyResponse represents Paystack transaction verification response
type PaystackVerifyResponse struct {
	Status  bool   `json:"status"`
	Message string `json:"message"`
	Data    struct {
		ID              int64                  `json:"id"`
		Domain          string                 `json:"domain"`
		Status          string                 `json:"status"`
		Reference       string                 `json:"reference"`
		Amount          int64                  `json:"amount"`
		Message         string                 `json:"message"`
		GatewayResponse string                 `json:"gateway_response"`
		PaidAt          string                 `json:"paid_at"`
		CreatedAt       string                 `json:"created_at"`
		Channel         string                 `json:"channel"`
		Currency        string                 `json:"currency"`
		IPAddress       string                 `json:"ip_address"`
		Metadata        map[string]interface{} `json:"metadata"`
		Customer        struct {
			ID           int64  `json:"id"`
			Email        string `json:"email"`
			CustomerCode string `json:"customer_code"`
		} `json:"customer"`
	} `json:"data"`
}

// PaystackBankListResponse represents Paystack bank list response
type PaystackBankListResponse struct {
	Status  bool   `json:"status"`
	Message string `json:"message"`
	Data    []struct {
		ID       int64  `json:"id"`
		Name     string `json:"name"`
		Slug     string `json:"slug"`
		Code     string `json:"code"`
		Longcode string `json:"longcode"`
		Gateway  string `json:"gateway"`
		Active   bool   `json:"active"`
		IsDeleted bool  `json:"is_deleted"`
		Country  string `json:"country"`
		Currency string `json:"currency"`
		Type     string `json:"type"`
	} `json:"data"`
}

// PaystackResolveAccountResponse represents Paystack account resolution response
type PaystackResolveAccountResponse struct {
	Status  bool   `json:"status"`
	Message string `json:"message"`
	Data    struct {
		AccountNumber string `json:"account_number"`
		AccountName   string `json:"account_name"`
		BankID        int64  `json:"bank_id"`
	} `json:"data"`
}

// ResolveAccountRequest represents account resolution request
type ResolveAccountRequest struct {
	AccountNumber string `json:"account_number" binding:"required"`
	BankCode      string `json:"bank_code" binding:"required"`
}

// ResolveAccountResponse represents account resolution response
type ResolveAccountResponse struct {
	AccountNumber string `json:"account_number"`
	AccountName   string `json:"account_name"`
	BankCode      string `json:"bank_code"`
	BankName      string `json:"bank_name"`
}
