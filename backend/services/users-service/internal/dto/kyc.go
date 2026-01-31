package dto

import "time"

// SubmitNINRequest represents a NIN submission request
type SubmitNINRequest struct {
	NIN string `json:"nin" binding:"required,len=11"`
}

// KYCStatusResponse represents KYC status information
type KYCStatusResponse struct {
	KYCVerified   bool       `json:"kyc_verified"`
	KYCVerifiedAt *time.Time `json:"kyc_verified_at,omitempty"`
	NIN           string     `json:"nin,omitempty"` // Masked for privacy
}
