package models

// EmailType represents the unique identifier for an email template
type EmailType string

const (
	EmailTypeUserSignedUp          EmailType = "user_signed_up"
	EmailTypeEmailVerification     EmailType = "email_verification"
	EmailTypePasswordReset         EmailType = "password_reset"
	EmailTypePaymentVerified       EmailType = "payment_verified"
	EmailTypeContributionConfirmed EmailType = "contribution_confirmed"
	EmailTypeWithdrawalRequested   EmailType = "withdrawal_requested"
	EmailTypeWithdrawalCompleted   EmailType = "withdrawal_completed"
	EmailTypeProofSubmitted        EmailType = "proof_submitted"
	EmailTypeProofVoted            EmailType = "proof_voted"
	EmailTypeGoalFunded            EmailType = "goal_funded"
	EmailTypeKYCVerified           EmailType = "kyc_verified"
)

// EmailPayload represents the data sent to the notification service
type EmailPayload struct {
	Type      EmailType              `json:"type"`
	Recipient string                 `json:"recipient"`
	Subject   string                 `json:"subject"`
	Data      map[string]interface{} `json:"data"`
}
