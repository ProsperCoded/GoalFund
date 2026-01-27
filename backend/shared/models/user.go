package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// User represents a user in the system
type User struct {
	ID              uuid.UUID `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	Email           string    `gorm:"uniqueIndex;not null;size:255" json:"email"`
	Username        string    `gorm:"uniqueIndex;not null;size:100" json:"username"`
	PasswordHash    string    `gorm:"not null;size:255" json:"-"` // Never serialize password
	FirstName       string    `gorm:"size:100" json:"first_name"`
	LastName        string    `gorm:"size:100" json:"last_name"`
	Phone           string    `gorm:"size:20" json:"phone"`
	EmailVerified   bool      `gorm:"default:false" json:"email_verified"`
	PhoneVerified   bool      `gorm:"default:false" json:"phone_verified"`
	CreatedAt       time.Time `gorm:"not null" json:"created_at"`
	UpdatedAt       time.Time `gorm:"not null" json:"updated_at"`
	
	// Relationships
	Roles    []Role    `gorm:"many2many:user_roles;" json:"roles,omitempty"`
	Sessions []Session `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE" json:"-"`
}

// BeforeCreate sets UUID before creating user
func (u *User) BeforeCreate(tx *gorm.DB) error {
	if u.ID == uuid.Nil {
		u.ID = uuid.New()
	}
	return nil
}

// Role represents a user role
type Role struct {
	ID          uuid.UUID              `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	Name        string                 `gorm:"uniqueIndex;not null;size:50" json:"name"`
	Description string                 `gorm:"size:255" json:"description"`
	Permissions map[string]interface{} `gorm:"type:jsonb" json:"permissions"`
	CreatedAt   time.Time              `gorm:"not null" json:"created_at"`
	
	// Relationships
	Users []User `gorm:"many2many:user_roles;" json:"-"`
}

// BeforeCreate sets UUID before creating role
func (r *Role) BeforeCreate(tx *gorm.DB) error {
	if r.ID == uuid.Nil {
		r.ID = uuid.New()
	}
	return nil
}

// UserRole represents the many-to-many relationship between users and roles
type UserRole struct {
	UserID     uuid.UUID `gorm:"type:uuid;primaryKey" json:"user_id"`
	RoleID     uuid.UUID `gorm:"type:uuid;primaryKey" json:"role_id"`
	AssignedAt time.Time `gorm:"not null;default:CURRENT_TIMESTAMP" json:"assigned_at"`
	
	// Relationships
	User User `gorm:"constraint:OnDelete:CASCADE"`
	Role Role `gorm:"constraint:OnDelete:CASCADE"`
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