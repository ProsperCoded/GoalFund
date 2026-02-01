package dto

// WebhookPayload represents the incoming webhook payload from Paystack
type WebhookPayload struct {
	Event string                 `json:"event"`
	Data  map[string]interface{} `json:"data"`
}

// WebhookResponse represents the response to a webhook
type WebhookResponse struct {
	Status  string `json:"status"`
	Message string `json:"message,omitempty"`
}

// ChargeSuccessData represents the data in a charge.success webhook
type ChargeSuccessData struct {
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
}

// TransferSuccessData represents the data in a transfer.success webhook
type TransferSuccessData struct {
	ID            int64  `json:"id"`
	Domain        string `json:"domain"`
	Status        string `json:"status"`
	Reference     string `json:"reference"`
	Amount        int64  `json:"amount"`
	Currency      string `json:"currency"`
	Reason        string `json:"reason"`
	TransferCode  string `json:"transfer_code"`
	TransferredAt string `json:"transferred_at"`
	CreatedAt     string `json:"created_at"`
	UpdatedAt     string `json:"updated_at"`
	Recipient     struct {
		RecipientCode string `json:"recipient_code"`
		Name          string `json:"name"`
		Email         string `json:"email"`
	} `json:"recipient"`
}
