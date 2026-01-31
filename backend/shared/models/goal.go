package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// GoalStatus represents the status of a goal
type GoalStatus string

const (
	GoalStatusOpen      GoalStatus = "OPEN"
	GoalStatusClosed    GoalStatus = "CLOSED"
	GoalStatusCancelled GoalStatus = "CANCELLED"
)

// Goal represents a funding goal with milestone support
type Goal struct {
	ID           uuid.UUID  `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	OwnerID      uuid.UUID  `gorm:"type:uuid;not null;index" json:"owner_id"`
	Title        string     `gorm:"not null;size:255" json:"title"`
	Description  string     `gorm:"type:text" json:"description"`
	TargetAmount int64      `gorm:"not null" json:"target_amount"` // Amount in smallest currency unit
	Currency     string     `gorm:"not null;size:3;default:'NGN'" json:"currency"`
	Deadline     *time.Time `json:"deadline,omitempty"`
	Status       GoalStatus `gorm:"not null;default:'OPEN';size:20" json:"status"`

	// Deposit account details (where goal owner receives withdrawals)
	DepositBankName      string `gorm:"size:100" json:"deposit_bank_name,omitempty"`
	DepositAccountNumber string `gorm:"size:20" json:"deposit_account_number,omitempty"`
	DepositAccountName   string `gorm:"size:255" json:"deposit_account_name,omitempty"`

	CreatedAt time.Time `gorm:"not null" json:"created_at"`
	UpdatedAt time.Time `gorm:"not null" json:"updated_at"`

	// Relationships
	Milestones    []Milestone    `gorm:"foreignKey:GoalID;constraint:OnDelete:CASCADE" json:"milestones,omitempty"`
	Contributions []Contribution `gorm:"foreignKey:GoalID;constraint:OnDelete:CASCADE" json:"contributions,omitempty"`
	Withdrawals   []Withdrawal   `gorm:"foreignKey:GoalID;constraint:OnDelete:CASCADE" json:"withdrawals,omitempty"`
	Proofs        []Proof        `gorm:"foreignKey:GoalID;constraint:OnDelete:CASCADE" json:"proofs,omitempty"`
	Refunds       []Refund       `gorm:"foreignKey:GoalID;constraint:OnDelete:CASCADE" json:"refunds,omitempty"`
}

// BeforeCreate sets UUID before creating goal
func (g *Goal) BeforeCreate(tx *gorm.DB) error {
	if g.ID == uuid.Nil {
		g.ID = uuid.New()
	}
	return nil
}

// TableName specifies the table name for Goal
func (Goal) TableName() string {
	return "goals"
}
// RefundStatus represents the status of a refund
type RefundStatus string

const (
	RefundStatusPending    RefundStatus = "PENDING"
	RefundStatusProcessing RefundStatus = "PROCESSING"
	RefundStatusCompleted  RefundStatus = "COMPLETED"
	RefundStatusFailed     RefundStatus = "FAILED"
)

// Refund represents a refund request for a goal
type Refund struct {
	ID                uuid.UUID    `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	GoalID            uuid.UUID    `gorm:"type:uuid;not null;index" json:"goal_id"`
	InitiatedBy       uuid.UUID    `gorm:"type:uuid;not null;index" json:"initiated_by"`
	RefundPercentage  float64      `gorm:"type:decimal(5,2);not null" json:"refund_percentage"` // 1-100
	TotalRefundAmount int64        `gorm:"not null" json:"total_refund_amount"`
	Currency          string       `gorm:"not null;size:3;default:'NGN'" json:"currency"`
	Reason            string       `gorm:"type:text" json:"reason,omitempty"`
	Status            RefundStatus `gorm:"not null;default:'PENDING';size:20" json:"status"`
	CreatedAt         time.Time    `gorm:"not null" json:"created_at"`
	CompletedAt       *time.Time   `json:"completed_at,omitempty"`

	// Relationships
	Goal         Goal                `gorm:"constraint:OnDelete:CASCADE"`
	Disbursements []RefundDisbursement `gorm:"foreignKey:RefundID;constraint:OnDelete:CASCADE" json:"disbursements,omitempty"`
}

// BeforeCreate sets UUID before creating refund
func (r *Refund) BeforeCreate(tx *gorm.DB) error {
	if r.ID == uuid.Nil {
		r.ID = uuid.New()
	}
	return nil
}

// TableName specifies the table name for Refund
func (Refund) TableName() string {
	return "refunds"
}

// RefundDisbursement represents an individual refund to a contributor
type RefundDisbursement struct {
	ID              uuid.UUID    `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	RefundID        uuid.UUID    `gorm:"type:uuid;not null;index" json:"refund_id"`
	ContributionID  uuid.UUID    `gorm:"type:uuid;not null;index" json:"contribution_id"`
	UserID          uuid.UUID    `gorm:"type:uuid;not null;index" json:"user_id"`
	Amount          int64        `gorm:"not null" json:"amount"`
	Currency        string       `gorm:"not null;size:3;default:'NGN'" json:"currency"`

	// Settlement account snapshot (at time of refund)
	SettlementBankName      string `gorm:"size:100" json:"settlement_bank_name,omitempty"`
	SettlementAccountNumber string `gorm:"size:20" json:"settlement_account_number,omitempty"`
	SettlementAccountName   string `gorm:"size:255" json:"settlement_account_name,omitempty"`

	Status                RefundStatus `gorm:"not null;default:'PENDING';size:20" json:"status"`
	LedgerTransactionID   *uuid.UUID   `gorm:"type:uuid" json:"ledger_transaction_id,omitempty"`
	CreatedAt             time.Time    `gorm:"not null" json:"created_at"`
	CompletedAt           *time.Time   `json:"completed_at,omitempty"`

	// Relationships
	Refund       Refund       `gorm:"constraint:OnDelete:CASCADE"`
	Contribution Contribution `gorm:"constraint:OnDelete:CASCADE"`
}

// BeforeCreate sets UUID before creating refund disbursement
func (rd *RefundDisbursement) BeforeCreate(tx *gorm.DB) error {
	if rd.ID == uuid.Nil {
		rd.ID = uuid.New()
	}
	return nil
}

// TableName specifies the table name for RefundDisbursement
func (RefundDisbursement) TableName() string {
	return "refund_disbursements"
}
// RecurrenceType represents the type of milestone recurrence
type RecurrenceType string

const (
	RecurrenceWeekly   RecurrenceType = "WEEKLY"
	RecurrenceMonthly  RecurrenceType = "MONTHLY"
	RecurrenceSemester RecurrenceType = "SEMESTER"
	RecurrenceYearly   RecurrenceType = "YEARLY"
)

// MilestoneStatus represents the status of a milestone
type MilestoneStatus string

const (
	MilestoneStatusPending   MilestoneStatus = "PENDING"
	MilestoneStatusActive    MilestoneStatus = "ACTIVE"
	MilestoneStatusCompleted MilestoneStatus = "COMPLETED"
)

// Milestone represents a goal milestone (supports recurring milestones)
type Milestone struct {
	ID           uuid.UUID       `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	GoalID       uuid.UUID       `gorm:"type:uuid;not null;index" json:"goal_id"`
	Title        string          `gorm:"not null;size:255" json:"title"`
	Description  string          `gorm:"type:text" json:"description"`
	TargetAmount int64           `gorm:"not null" json:"target_amount"`
	OrderIndex   int             `gorm:"not null" json:"order_index"`

	// Recurring support
	IsRecurring        bool            `gorm:"default:false" json:"is_recurring"`
	RecurrenceType     *RecurrenceType `gorm:"size:20" json:"recurrence_type,omitempty"`
	RecurrenceInterval int             `json:"recurrence_interval,omitempty"`
	NextDueDate        *time.Time      `json:"next_due_date,omitempty"`

	Status      MilestoneStatus `gorm:"not null;default:'PENDING';size:20" json:"status"`
	CompletedAt *time.Time      `json:"completed_at,omitempty"`

	CreatedAt time.Time `gorm:"not null" json:"created_at"`
	UpdatedAt time.Time `gorm:"not null" json:"updated_at"`

	// Relationships
	Goal          Goal           `gorm:"constraint:OnDelete:CASCADE"`
	Contributions []Contribution `gorm:"foreignKey:MilestoneID;constraint:OnDelete:SET NULL" json:"contributions,omitempty"`
	Withdrawals   []Withdrawal   `gorm:"foreignKey:MilestoneID;constraint:OnDelete:SET NULL" json:"withdrawals,omitempty"`
	Proofs        []Proof        `gorm:"foreignKey:MilestoneID;constraint:OnDelete:SET NULL" json:"proofs,omitempty"`
}

// BeforeCreate sets UUID before creating milestone
func (m *Milestone) BeforeCreate(tx *gorm.DB) error {
	if m.ID == uuid.Nil {
		m.ID = uuid.New()
	}
	return nil
}

// TableName specifies the table name for Milestone
func (Milestone) TableName() string {
	return "milestones"
}

// ContributionStatus represents the status of a contribution
type ContributionStatus string

const (
	ContributionStatusPending   ContributionStatus = "PENDING"
	ContributionStatusConfirmed ContributionStatus = "CONFIRMED"
	ContributionStatusFailed    ContributionStatus = "FAILED"
	ContributionStatusRefunded  ContributionStatus = "REFUNDED"
)

// Contribution represents a user's contribution to a goal
type Contribution struct {
	ID          uuid.UUID          `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	GoalID      uuid.UUID          `gorm:"type:uuid;not null;index" json:"goal_id"`
	MilestoneID *uuid.UUID         `gorm:"type:uuid;index" json:"milestone_id,omitempty"`
	UserID      uuid.UUID          `gorm:"type:uuid;not null;index" json:"user_id"`
	PaymentID   *uuid.UUID         `gorm:"type:uuid;index" json:"payment_id,omitempty"` // Reference to payment service
	Amount      int64              `gorm:"not null" json:"amount"`
	Currency    string             `gorm:"not null;size:3;default:'NGN'" json:"currency"`
	Status      ContributionStatus `gorm:"not null;default:'PENDING';size:20" json:"status"`
	CreatedAt   time.Time          `gorm:"not null" json:"created_at"`
	UpdatedAt   time.Time          `gorm:"not null" json:"updated_at"`

	// Relationships
	Goal      Goal       `gorm:"constraint:OnDelete:CASCADE"`
	Milestone *Milestone `gorm:"constraint:OnDelete:SET NULL"`
}

// BeforeCreate sets UUID before creating contribution
func (c *Contribution) BeforeCreate(tx *gorm.DB) error {
	if c.ID == uuid.Nil {
		c.ID = uuid.New()
	}
	return nil
}

// TableName specifies the table name for Contribution
func (Contribution) TableName() string {
	return "contributions"
}

// WithdrawalStatus represents the status of a withdrawal
type WithdrawalStatus string

const (
	WithdrawalStatusPending    WithdrawalStatus = "PENDING"
	WithdrawalStatusProcessing WithdrawalStatus = "PROCESSING"
	WithdrawalStatusCompleted  WithdrawalStatus = "COMPLETED"
	WithdrawalStatusFailed     WithdrawalStatus = "FAILED"
)

// Withdrawal represents a withdrawal request by goal owner
type Withdrawal struct {
	ID          uuid.UUID        `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	GoalID      uuid.UUID        `gorm:"type:uuid;not null;index" json:"goal_id"`
	MilestoneID *uuid.UUID       `gorm:"type:uuid;index" json:"milestone_id,omitempty"`
	OwnerID     uuid.UUID        `gorm:"type:uuid;not null;index" json:"owner_id"`
	Amount      int64            `gorm:"not null" json:"amount"`
	Currency    string           `gorm:"not null;size:3;default:'NGN'" json:"currency"`

	// Bank details snapshot (at time of withdrawal)
	BankName      string `gorm:"not null;size:100" json:"bank_name"`
	AccountNumber string `gorm:"not null;size:20" json:"account_number"`
	AccountName   string `gorm:"not null;size:255" json:"account_name"`

	Status              WithdrawalStatus `gorm:"not null;default:'PENDING';size:20" json:"status"`
	LedgerTransactionID *uuid.UUID       `gorm:"type:uuid" json:"ledger_transaction_id,omitempty"`

	RequestedAt time.Time  `gorm:"not null" json:"requested_at"`
	CompletedAt *time.Time `json:"completed_at,omitempty"`

	// Relationships
	Goal      Goal       `gorm:"constraint:OnDelete:CASCADE"`
	Milestone *Milestone `gorm:"constraint:OnDelete:SET NULL"`
}

// BeforeCreate sets UUID before creating withdrawal
func (w *Withdrawal) BeforeCreate(tx *gorm.DB) error {
	if w.ID == uuid.Nil {
		w.ID = uuid.New()
	}
	if w.RequestedAt.IsZero() {
		w.RequestedAt = time.Now()
	}
	return nil
}

// TableName specifies the table name for Withdrawal
func (Withdrawal) TableName() string {
	return "withdrawals"
}

// Proof represents proof of goal accomplishment
type Proof struct {
	ID          uuid.UUID  `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	GoalID      uuid.UUID  `gorm:"type:uuid;not null;index" json:"goal_id"`
	MilestoneID *uuid.UUID `gorm:"type:uuid;index" json:"milestone_id,omitempty"`
	SubmittedBy uuid.UUID  `gorm:"type:uuid;not null;index" json:"submitted_by"`
	Title       string     `gorm:"not null;size:255" json:"title"`
	Description string     `gorm:"type:text" json:"description"`
	MediaURLs   []string   `gorm:"type:jsonb;serializer:json" json:"media_urls,omitempty"`
	SubmittedAt time.Time  `gorm:"not null" json:"submitted_at"`

	// Relationships
	Goal      Goal       `gorm:"constraint:OnDelete:CASCADE"`
	Milestone *Milestone `gorm:"constraint:OnDelete:SET NULL"`
	Votes     []Vote     `gorm:"foreignKey:ProofID;constraint:OnDelete:CASCADE" json:"votes,omitempty"`
}

// BeforeCreate sets UUID before creating proof
func (p *Proof) BeforeCreate(tx *gorm.DB) error {
	if p.ID == uuid.Nil {
		p.ID = uuid.New()
	}
	if p.SubmittedAt.IsZero() {
		p.SubmittedAt = time.Now()
	}
	return nil
}

// TableName specifies the table name for Proof
func (Proof) TableName() string {
	return "proofs"
}

// Vote represents a community vote on proof verification
type Vote struct {
	ID          uuid.UUID `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	ProofID     uuid.UUID `gorm:"type:uuid;not null;index" json:"proof_id"`
	VoterID     uuid.UUID `gorm:"type:uuid;not null;index" json:"voter_id"`
	IsSatisfied bool      `gorm:"not null" json:"is_satisfied"` // TRUE = satisfied, FALSE = not satisfied
	Comment     string    `gorm:"type:text" json:"comment,omitempty"`
	VotedAt     time.Time `gorm:"not null" json:"voted_at"`

	// Relationships
	Proof Proof `gorm:"constraint:OnDelete:CASCADE"`
}

// BeforeCreate sets UUID before creating vote
func (v *Vote) BeforeCreate(tx *gorm.DB) error {
	if v.ID == uuid.Nil {
		v.ID = uuid.New()
	}
	if v.VotedAt.IsZero() {
		v.VotedAt = time.Now()
	}
	return nil
}

// TableName specifies the table name for Vote
func (Vote) TableName() string {
	return "votes"
}