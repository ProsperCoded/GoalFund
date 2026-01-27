package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// PaymentStatus represents the status of a payment
type PaymentStatus string

const (
	PaymentStatusInitiated PaymentStatus = "INITIATED"
	PaymentStatusPending   PaymentStatus = "PENDING"
	PaymentStatusVerified  PaymentStatus = "VERIFIED"
	PaymentStatusFailed    PaymentStatus = "FAILED"
)

// Payment represents a payment record in MongoDB
type Payment struct {
	ID                 primitive.ObjectID     `bson:"_id,omitempty" json:"id"`
	PaymentID          string                 `bson:"paymentId" json:"payment_id"` // Unique payment identifier
	PaystackReference  string                 `bson:"paystackReference,omitempty" json:"paystack_reference"`
	UserID             string                 `bson:"userId" json:"user_id"`
	GoalID             string                 `bson:"goalId,omitempty" json:"goal_id"`
	Amount             int64                  `bson:"amount" json:"amount"`
	Currency           string                 `bson:"currency" json:"currency"`
	Status             PaymentStatus          `bson:"status" json:"status"`
	PaystackData       map[string]interface{} `bson:"paystackData,omitempty" json:"paystack_data"`
	CreatedAt          time.Time              `bson:"createdAt" json:"created_at"`
	UpdatedAt          time.Time              `bson:"updatedAt" json:"updated_at"`
}

// WebhookEvent represents a webhook event from Paystack
type WebhookEvent struct {
	ID          primitive.ObjectID     `bson:"_id,omitempty" json:"id"`
	EventID     string                 `bson:"eventId" json:"event_id"` // Unique event identifier
	Event       string                 `bson:"event" json:"event"`      // Event type (charge.success, etc.)
	Data        map[string]interface{} `bson:"data" json:"data"`        // Raw webhook payload
	Signature   string                 `bson:"signature,omitempty" json:"signature"`
	Processed   bool                   `bson:"processed" json:"processed"`
	ReceivedAt  time.Time              `bson:"receivedAt" json:"received_at"`
	ProcessedAt *time.Time             `bson:"processedAt,omitempty" json:"processed_at"`
}

// IdempotencyKey represents an idempotency key for preventing duplicate processing
type IdempotencyKey struct {
	ID        primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Key       string             `bson:"key" json:"key"` // Unique idempotency key
	PaymentID string             `bson:"paymentId" json:"payment_id"`
	CreatedAt time.Time          `bson:"createdAt" json:"created_at"`
	ExpiresAt time.Time          `bson:"expiresAt" json:"expires_at"` // TTL index
}

// PaymentAttempt represents a payment attempt (for retry logic)
type PaymentAttempt struct {
	ID        primitive.ObjectID     `bson:"_id,omitempty" json:"id"`
	PaymentID string                 `bson:"paymentId" json:"payment_id"`
	Attempt   int                    `bson:"attempt" json:"attempt"`
	Status    PaymentStatus          `bson:"status" json:"status"`
	Error     string                 `bson:"error,omitempty" json:"error"`
	Metadata  map[string]interface{} `bson:"metadata,omitempty" json:"metadata"`
	CreatedAt time.Time              `bson:"createdAt" json:"created_at"`
}

// PaymentMethod represents supported payment methods
type PaymentMethod string

const (
	PaymentMethodCard         PaymentMethod = "CARD"
	PaymentMethodBankTransfer PaymentMethod = "BANK_TRANSFER"
	PaymentMethodUSSD         PaymentMethod = "USSD"
	PaymentMethodQR           PaymentMethod = "QR"
)

// PaymentRequest represents a payment initialization request
type PaymentRequest struct {
	ID            primitive.ObjectID     `bson:"_id,omitempty" json:"id"`
	PaymentID     string                 `bson:"paymentId" json:"payment_id"`
	UserID        string                 `bson:"userId,omitempty" json:"user_id"`
	GoalID        string                 `bson:"goalId,omitempty" json:"goal_id"`
	Amount        int64                  `bson:"amount" json:"amount"`
	Currency      string                 `bson:"currency" json:"currency"`
	Email         string                 `bson:"email" json:"email"`
	CallbackURL   string                 `bson:"callbackUrl,omitempty" json:"callback_url"`
	Method        PaymentMethod          `bson:"method,omitempty" json:"method"`
	Metadata      map[string]interface{} `bson:"metadata,omitempty" json:"metadata"`
	CreatedAt     time.Time              `bson:"createdAt" json:"created_at"`
	ExpiresAt     time.Time              `bson:"expiresAt" json:"expires_at"`
}