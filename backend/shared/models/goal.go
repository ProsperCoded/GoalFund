package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// GoalStatus represents the status of a goal
type GoalStatus string

const (
	GoalStatusOpen           GoalStatus = "OPEN"
	GoalStatusFunded         GoalStatus = "FUNDED"
	GoalStatusWithdrawn      GoalStatus = "WITHDRAWN"
	GoalStatusProofSubmitted GoalStatus = "PROOF_SUBMITTED"
	GoalStatusVerified       GoalStatus = "VERIFIED"
	GoalStatusCancelled      GoalStatus = "CANCELLED"
)

// Goal represents a funding goal
type Goal struct {
	ID           uuid.UUID  `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	OwnerID      uuid.UUID  `gorm:"type:uuid;not null;index" json:"owner_id"`
	Title        string     `gorm:"not null;size:255" json:"title"`
	Description  string     `gorm:"type:text" json:"description"`
	TargetAmount int64      `gorm:"not null" json:"target_amount"` // Amount in smallest currency unit
	Currency     string     `gorm:"not null;size:3;default:'NGN'" json:"currency"`
	Deadline     *time.Time `json:"deadline"`
	Status       GoalStatus `gorm:"not null;default:'OPEN';size:20" json:"status"`
	CreatedAt    time.Time  `gorm:"not null" json:"created_at"`
	UpdatedAt    time.Time  `gorm:"not null" json:"updated_at"`

	// Relationships
	Contributions []Contribution `gorm:"foreignKey:GoalID;constraint:OnDelete:CASCADE" json:"contributions,omitempty"`
	Proofs        []Proof        `gorm:"foreignKey:GoalID;constraint:OnDelete:CASCADE" json:"proofs,omitempty"`
}

// BeforeCreate sets UUID before creating goal
func (g *Goal) BeforeCreate(tx *gorm.DB) error {
	if g.ID == uuid.Nil {
		g.ID = uuid.New()
	}
	return nil
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
	ID        uuid.UUID          `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	GoalID    uuid.UUID          `gorm:"type:uuid;not null;index" json:"goal_id"`
	UserID    uuid.UUID          `gorm:"type:uuid;not null;index" json:"user_id"`
	PaymentID uuid.UUID          `gorm:"type:uuid;index" json:"payment_id"` // Reference to payment service
	Amount    int64              `gorm:"not null" json:"amount"`
	Currency  string             `gorm:"not null;size:3;default:'NGN'" json:"currency"`
	Status    ContributionStatus `gorm:"not null;default:'PENDING';size:20" json:"status"`
	CreatedAt time.Time          `gorm:"not null" json:"created_at"`
	UpdatedAt time.Time          `gorm:"not null" json:"updated_at"`

	// Relationships
	Goal Goal `gorm:"constraint:OnDelete:CASCADE"`
}

// BeforeCreate sets UUID before creating contribution
func (c *Contribution) BeforeCreate(tx *gorm.DB) error {
	if c.ID == uuid.Nil {
		c.ID = uuid.New()
	}
	return nil
}

// ProofStatus represents the status of a proof submission
type ProofStatus string

const (
	ProofStatusPending  ProofStatus = "PENDING"
	ProofStatusApproved ProofStatus = "APPROVED"
	ProofStatusRejected ProofStatus = "REJECTED"
)

// Proof represents proof of goal accomplishment
type Proof struct {
	ID          uuid.UUID   `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	GoalID      uuid.UUID   `gorm:"type:uuid;not null;index" json:"goal_id"`
	SubmittedBy uuid.UUID   `gorm:"type:uuid;not null;index" json:"submitted_by"`
	Title       string      `gorm:"not null;size:255" json:"title"`
	Description string      `gorm:"type:text" json:"description"`
	MediaURLs   []string    `gorm:"type:jsonb" json:"media_urls"`
	Status      ProofStatus `gorm:"not null;default:'PENDING';size:20" json:"status"`
	SubmittedAt time.Time   `gorm:"not null" json:"submitted_at"`
	VerifiedAt  *time.Time  `json:"verified_at"`

	// Relationships
	Goal  Goal   `gorm:"constraint:OnDelete:CASCADE"`
	Votes []Vote `gorm:"foreignKey:ProofID;constraint:OnDelete:CASCADE" json:"votes,omitempty"`
}

// BeforeCreate sets UUID before creating proof
func (p *Proof) BeforeCreate(tx *gorm.DB) error {
	if p.ID == uuid.Nil {
		p.ID = uuid.New()
	}
	return nil
}

// Vote represents a community vote on proof verification
type Vote struct {
	ID         uuid.UUID `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	ProofID    uuid.UUID `gorm:"type:uuid;not null;index" json:"proof_id"`
	VoterID    uuid.UUID `gorm:"type:uuid;not null;index" json:"voter_id"`
	IsApproved bool      `gorm:"not null" json:"is_approved"`
	Comment    string    `gorm:"type:text" json:"comment"`
	VotedAt    time.Time `gorm:"not null" json:"voted_at"`

	// Relationships
	Proof Proof `gorm:"constraint:OnDelete:CASCADE"`
}

// BeforeCreate sets UUID before creating vote
func (v *Vote) BeforeCreate(tx *gorm.DB) error {
	if v.ID == uuid.Nil {
		v.ID = uuid.New()
	}
	return nil
}

// Unique constraint to prevent duplicate votes
func (Vote) TableName() string {
	return "votes"
}

// Add unique constraint on proof_id and voter_id
type VoteConstraint struct {
	ProofID uuid.UUID `gorm:"uniqueIndex:idx_proof_voter"`
	VoterID uuid.UUID `gorm:"uniqueIndex:idx_proof_voter"`
}