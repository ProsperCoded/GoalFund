package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// UserRole enum type
type UserRole string

const (
	UserRoleUser  UserRole = "user"
	UserRoleAdmin UserRole = "admin"
)

// User represents a user in the system
type User struct {
	ID              uuid.UUID  `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	Email           string     `gorm:"uniqueIndex;not null;size:255" json:"email"`
	Username        string     `gorm:"uniqueIndex;not null;size:100" json:"username"`
	PasswordHash    string     `gorm:"not null;size:255" json:"-"` // Never serialize password
	FirstName       string     `gorm:"size:100" json:"first_name"`
	LastName        string     `gorm:"size:100" json:"last_name"`
	Phone           string     `gorm:"size:20" json:"phone"`
	EmailVerified   bool       `gorm:"default:false" json:"email_verified"`
	PhoneVerified   bool       `gorm:"default:false" json:"phone_verified"`
	
	// Password Setup Tracking
	HasSetPassword  bool       `gorm:"default:true" json:"has_set_password"` // False for email-only contributions
	
	// KYC Verification Fields
	NIN             string     `gorm:"size:11;index" json:"nin,omitempty"` // National Identification Number (11 digits)
	KYCVerified     bool       `gorm:"default:false" json:"kyc_verified"`
	KYCVerifiedAt   *time.Time `gorm:"index" json:"kyc_verified_at,omitempty"`
	
	// Settlement Account Details (for refunds/withdrawals)
	SettlementBankName      string `gorm:"size:100" json:"settlement_bank_name,omitempty"`
	SettlementAccountNumber string `gorm:"size:20;index" json:"settlement_account_number,omitempty"`
	SettlementAccountName   string `gorm:"size:255" json:"settlement_account_name,omitempty"`
	
	Role            UserRole   `gorm:"type:user_role;default:user" json:"role"`
	CreatedAt       time.Time  `gorm:"not null" json:"created_at"`
	UpdatedAt       time.Time  `gorm:"not null" json:"updated_at"`
	
	// Relationships
	Sessions []Session `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE" json:"-"`
}

// BeforeCreate sets UUID before creating user
func (u *User) BeforeCreate(tx *gorm.DB) error {
	if u.ID == uuid.Nil {
		u.ID = uuid.New()
	}
	return nil
}

// PasswordResetToken represents a password reset token
type PasswordResetToken struct {
	ID        uuid.UUID `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	UserID    uuid.UUID `gorm:"type:uuid;not null;index" json:"user_id"`
	TokenHash string    `gorm:"uniqueIndex;not null;size:255" json:"-"`
	ExpiresAt time.Time `gorm:"not null" json:"expires_at"`
	Used      bool      `gorm:"default:false" json:"used"`
	CreatedAt time.Time `gorm:"not null" json:"created_at"`
	
	// Relationships
	User User `gorm:"constraint:OnDelete:CASCADE"`
}

// BeforeCreate sets UUID before creating password reset token
func (p *PasswordResetToken) BeforeCreate(tx *gorm.DB) error {
	if p.ID == uuid.Nil {
		p.ID = uuid.New()
	}
	return nil
}

// IsExpired checks if the token has expired
func (p *PasswordResetToken) IsExpired() bool {
	return time.Now().After(p.ExpiresAt)
}

// Session represents a user session
type Session struct {
	ID        uuid.UUID              `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	UserID    uuid.UUID              `gorm:"type:uuid;not null;index" json:"user_id"`
	TokenHash string                 `gorm:"uniqueIndex;not null;size:255" json:"-"`
	ExpiresAt time.Time              `gorm:"not null" json:"expires_at"`
	Metadata  map[string]interface{} `gorm:"type:jsonb" json:"metadata"`
	CreatedAt time.Time              `gorm:"not null" json:"created_at"`
	
	// Relationships
	User User `gorm:"constraint:OnDelete:CASCADE"`
}

// BeforeCreate sets UUID before creating session
func (s *Session) BeforeCreate(tx *gorm.DB) error {
	if s.ID == uuid.Nil {
		s.ID = uuid.New()
	}
	return nil
}

// IsExpired checks if the session has expired
func (s *Session) IsExpired() bool {
	return time.Now().After(s.ExpiresAt)
}