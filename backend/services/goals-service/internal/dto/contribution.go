package dto

import "github.com/google/uuid"

// CreateContributionRequest represents a request to create a contribution
type CreateContributionRequest struct {
	GoalID      uuid.UUID
	MilestoneID *uuid.UUID
	Amount      int64
}

// CreateWithdrawalRequest represents a request to create a withdrawal
type CreateWithdrawalRequest struct {
	GoalID        uuid.UUID
	MilestoneID   *uuid.UUID
	Amount        int64
	BankName      string
	AccountNumber string
	AccountName   string
}

// CreateProofRequest represents a request to create a proof
type CreateProofRequest struct {
	GoalID      uuid.UUID
	MilestoneID *uuid.UUID
	Title       string
	Description string
	MediaURLs   []string
}

// CreateVoteRequest represents a request to create a vote
type CreateVoteRequest struct {
	ProofID     uuid.UUID
	IsSatisfied bool
	Comment     string
}

// VoteStats represents vote statistics
type VoteStats struct {
	TotalVotes       int64
	SatisfiedVotes   int64
	UnsatisfiedVotes int64
	SatisfactionRate float64
}
