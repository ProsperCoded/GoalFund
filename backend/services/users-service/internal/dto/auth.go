package dto

import (
	"time"

	"github.com/gofund/shared/models"
)

// LoginRequest represents a login request
type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=8"`
}

// RegisterRequest represents a registration request
type RegisterRequest struct {
	Email                   string `json:"email" binding:"required,email"`
	Username                string `json:"username" binding:"omitempty,min=3,max=50"`
	Password                string `json:"password" binding:"omitempty,min=8"`
	FirstName               string `json:"first_name" binding:"omitempty,min=2,max=50"`
	LastName                string `json:"last_name" binding:"omitempty,min=2,max=50"`
	Phone                   string `json:"phone"`
	SettlementBankName      string `json:"settlement_bank_name"`
	SettlementAccountNumber string `json:"settlement_account_number"`
	SettlementAccountName   string `json:"settlement_account_name"`
}

// SetPasswordRequest represents a first-time password setup request
type SetPasswordRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=8"`
}

// AuthResponse represents authentication response
type AuthResponse struct {
	User         *UserResponse `json:"user"`
	AccessToken  string        `json:"access_token"`
	RefreshToken string        `json:"refresh_token"`
	TokenType    string        `json:"token_type"`
	ExpiresIn    int64         `json:"expires_in"`
}

// UserResponse represents user data in response
type UserResponse struct {
	ID            string          `json:"id"`
	Email         string          `json:"email"`
	Username      string          `json:"username"`
	FirstName     string          `json:"first_name"`
	LastName      string          `json:"last_name"`
	Phone         string          `json:"phone"`
	EmailVerified bool            `json:"email_verified"`
	PhoneVerified bool            `json:"phone_verified"`
	KYCVerified   bool            `json:"kyc_verified"`
	KYCVerifiedAt *time.Time      `json:"kyc_verified_at,omitempty"`
	Role          models.UserRole `json:"role"`
	CreatedAt     time.Time       `json:"created_at"`
}

// UpdateProfileRequest represents a profile update request
type UpdateProfileRequest struct {
	FirstName string `json:"first_name" binding:"omitempty,min=2,max=50"`
	LastName  string `json:"last_name" binding:"omitempty,min=2,max=50"`
	Phone     string `json:"phone" binding:"omitempty"`
}

// ForgotPasswordRequest represents a forgot password request
type ForgotPasswordRequest struct {
	Email string `json:"email" binding:"required,email"`
}

// ResetPasswordRequest represents a reset password request
type ResetPasswordRequest struct {
	Token       string `json:"token" binding:"required"`
	NewPassword string `json:"new_password" binding:"required,min=8"`
}

// RefreshRequest represents a token refresh request
type RefreshRequest struct {
	RefreshToken string `json:"refresh_token" binding:"required"`
}

// Claims represents user claims for validation
type Claims struct {
	UserID string   `json:"user_id"`
	Email  string   `json:"email"`
	Roles  []string `json:"roles"`
}
