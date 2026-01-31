package events

// Event represents the base event interface
type Event interface {
	EventType() string
	EventID() string
	Timestamp() int64
}

// PaymentVerified event is emitted when a payment is verified
type PaymentVerified struct {
	ID        string
	PaymentID string
	UserID    string
	GoalID    string
	Amount    int64 // Amount in smallest currency unit (e.g., kobo for NGN)
	CreatedAt int64
}

func (e PaymentVerified) EventType() string { return "PaymentVerified" }
func (e PaymentVerified) EventID() string  { return e.ID }
func (e PaymentVerified) Timestamp() int64 { return e.CreatedAt }

// LedgerEntryCreated event is emitted when a ledger entry is created
type LedgerEntryCreated struct {
	ID           string
	LedgerEntryID string
	AccountID    string
	Amount       int64
	EntryType    string
	CreatedAt    int64
}

func (e LedgerEntryCreated) EventType() string { return "LedgerEntryCreated" }
func (e LedgerEntryCreated) EventID() string    { return e.ID }
func (e LedgerEntryCreated) Timestamp() int64    { return e.CreatedAt }

// GoalFunded event is emitted when a goal reaches its target
type GoalFunded struct {
	ID        string
	GoalID    string
	Amount    int64
	CreatedAt int64
}

func (e GoalFunded) EventType() string { return "GoalFunded" }
func (e GoalFunded) EventID() string   { return e.ID }
func (e GoalFunded) Timestamp() int64  { return e.CreatedAt }

// ProofSubmitted event is emitted when proof is submitted
type ProofSubmitted struct {
	ID        string
	GoalID    string
	ProofID   string
	CreatedAt int64
}

func (e ProofSubmitted) EventType() string { return "ProofSubmitted" }
func (e ProofSubmitted) EventID() string   { return e.ID }
func (e ProofSubmitted) Timestamp() int64  { return e.CreatedAt }

// ProofVerified event is emitted when proof is verified
type ProofVerified struct {
	ID        string
	GoalID    string
	ProofID   string
	CreatedAt int64
}

func (e ProofVerified) EventType() string { return "ProofVerified" }
func (e ProofVerified) EventID() string   { return e.ID }
func (e ProofVerified) Timestamp() int64 { return e.CreatedAt }

// UserSignedUp event is emitted when a user signs up
type UserSignedUp struct {
	ID        string
	UserID    string
	Email     string
	Username  string
	CreatedAt int64
}

func (e UserSignedUp) EventType() string { return "UserSignedUp" }
func (e UserSignedUp) EventID() string   { return e.ID }
func (e UserSignedUp) Timestamp() int64  { return e.CreatedAt }

// PasswordResetRequested event is emitted when a password reset is requested
type PasswordResetRequested struct {
	ID        string
	UserID    string
	Email     string
	Token     string // This will be the actual token, not the hash
	CreatedAt int64
}

func (e PasswordResetRequested) EventType() string { return "PasswordResetRequested" }
func (e PasswordResetRequested) EventID() string   { return e.ID }
func (e PasswordResetRequested) Timestamp() int64  { return e.CreatedAt }

// EmailVerificationRequested event is emitted when email verification is requested
type EmailVerificationRequested struct {
	ID        string
	UserID    string
	Email     string
	Token     string
	CreatedAt int64
}

func (e EmailVerificationRequested) EventType() string { return "EmailVerificationRequested" }
func (e EmailVerificationRequested) EventID() string   { return e.ID }
func (e EmailVerificationRequested) Timestamp() int64  { return e.CreatedAt }

// KYCVerified event is emitted when a user completes KYC verification
type KYCVerified struct {
	ID        string
	UserID    string
	Email     string
	Username  string
	CreatedAt int64
}

func (e KYCVerified) EventType() string { return "KYCVerified" }
func (e KYCVerified) EventID() string   { return e.ID }
func (e KYCVerified) Timestamp() int64  { return e.CreatedAt }

// RefundInitiated event is emitted when a refund is initiated
type RefundInitiated struct {
	ID                string
	RefundID          string
	GoalID            string
	InitiatedBy       string
	RefundPercentage  float64
	TotalRefundAmount int64
	CreatedAt         int64
}

func (e RefundInitiated) EventType() string { return "RefundInitiated" }
func (e RefundInitiated) EventID() string   { return e.ID }
func (e RefundInitiated) Timestamp() int64  { return e.CreatedAt }

// RefundCompleted event is emitted when a refund is completed
type RefundCompleted struct {
	ID                string
	RefundID          string
	GoalID            string
	TotalRefundAmount int64
	CompletedAt       int64
}

func (e RefundCompleted) EventType() string { return "RefundCompleted" }
func (e RefundCompleted) EventID() string   { return e.ID }
func (e RefundCompleted) Timestamp() int64  { return e.CompletedAt }

// ContributionRefunded event is emitted when a contribution is refunded
type ContributionRefunded struct {
	ID             string
	ContributionID string
	UserID         string
	GoalID         string
	RefundAmount   int64
	CreatedAt      int64
}

func (e ContributionRefunded) EventType() string { return "ContributionRefunded" }
func (e ContributionRefunded) EventID() string   { return e.ID }
func (e ContributionRefunded) Timestamp() int64  { return e.CreatedAt }
