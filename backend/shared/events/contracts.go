package events

// Event represents the base event interface
type Event interface {
	EventType() string
	EventID() string
	Timestamp() int64
}

// PaymentVerified event is emitted when a payment is verified
type PaymentVerified struct {
	EventID    string
	PaymentID  string
	UserID     string
	GoalID     string
	Amount     int64 // Amount in smallest currency unit (e.g., kobo for NGN)
	Timestamp  int64
}

func (e PaymentVerified) EventType() string { return "PaymentVerified" }
func (e PaymentVerified) EventID() string   { return e.EventID }
func (e PaymentVerified) Timestamp() int64  { return e.Timestamp }

// LedgerEntryCreated event is emitted when a ledger entry is created
type LedgerEntryCreated struct {
	EventID        string
	LedgerEntryID  string
	AccountID      string
	Amount         int64
	EntryType      string
	Timestamp      int64
}

func (e LedgerEntryCreated) EventType() string { return "LedgerEntryCreated" }
func (e LedgerEntryCreated) EventID() string   { return e.EventID }
func (e LedgerEntryCreated) Timestamp() int64  { return e.Timestamp }

// GoalFunded event is emitted when a goal reaches its target
type GoalFunded struct {
	EventID   string
	GoalID    string
	Amount    int64
	Timestamp int64
}

func (e GoalFunded) EventType() string { return "GoalFunded" }
func (e GoalFunded) EventID() string   { return e.EventID }
func (e GoalFunded) Timestamp() int64  { return e.Timestamp }

// ProofSubmitted event is emitted when proof is submitted
type ProofSubmitted struct {
	EventID   string
	GoalID    string
	ProofID   string
	Timestamp int64
}

func (e ProofSubmitted) EventType() string { return "ProofSubmitted" }
func (e ProofSubmitted) EventID() string   { return e.EventID }
func (e ProofSubmitted) Timestamp() int64  { return e.Timestamp }

// ProofVerified event is emitted when proof is verified
type ProofVerified struct {
	EventID   string
	GoalID    string
	ProofID   string
	Timestamp int64
}

func (e ProofVerified) EventType() string { return "ProofVerified" }
func (e ProofVerified) EventID() string   { return e.EventID }
func (e ProofVerified) Timestamp() int64  { return e.Timestamp }
