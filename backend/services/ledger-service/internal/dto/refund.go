package dto

import "github.com/google/uuid"

// ProcessRefundRequest represents a request to process a refund in the ledger
type ProcessRefundRequest struct {
	RefundID       uuid.UUID
	GoalID         uuid.UUID
	UserID         uuid.UUID
	ContributionID uuid.UUID
	Amount         int64
	Currency       string
}
