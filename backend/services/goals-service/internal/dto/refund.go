package dto

// InitiateRefundRequest represents a refund initiation request
type InitiateRefundRequest struct {
	GoalID           string  `json:"goal_id" binding:"required,uuid"`
	RefundPercentage float64 `json:"refund_percentage" binding:"required,min=1,max=100"`
	Reason           string  `json:"reason"`
}
