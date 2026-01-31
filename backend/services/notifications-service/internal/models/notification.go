package models

import (
	"time"
)

// NotificationType represents the type of notification
type NotificationType string

const (
	NotificationTypePaymentVerified       NotificationType = "payment_verified"
	NotificationTypeContributionConfirmed NotificationType = "contribution_confirmed"
	NotificationTypeWithdrawalRequested   NotificationType = "withdrawal_requested"
	NotificationTypeWithdrawalCompleted   NotificationType = "withdrawal_completed"
	NotificationTypeProofSubmitted        NotificationType = "proof_submitted"
	NotificationTypeProofVoted            NotificationType = "proof_voted"
	NotificationTypeGoalFunded            NotificationType = "goal_funded"
	NotificationTypeUserSignedUp          NotificationType = "user_signed_up"
	NotificationTypePasswordReset         NotificationType = "password_reset"
	NotificationTypeEmailVerification     NotificationType = "email_verification"
	NotificationTypeKYCVerified           NotificationType = "kyc_verified"
)

// Notification represents a notification record
type Notification struct {
	ID                 string                 `json:"id" db:"id"`
	UserID             string                 `json:"user_id" db:"user_id"`
	Type               NotificationType       `json:"type" db:"type"`
	Title              string                 `json:"title" db:"title"`
	Message            string                 `json:"message" db:"message"`
	Data               map[string]interface{} `json:"data" db:"data"`
	EmailSent          bool                   `json:"email_sent" db:"email_sent"`
	EmailSentAt        *time.Time             `json:"email_sent_at,omitempty" db:"email_sent_at"`
	EmailFailedReason  *string                `json:"email_failed_reason,omitempty" db:"email_failed_reason"`
	RetryCount         int                    `json:"retry_count" db:"retry_count"`
	IsRead             bool                   `json:"is_read" db:"is_read"`
	ReadAt             *time.Time             `json:"read_at,omitempty" db:"read_at"`
	CreatedAt          time.Time              `json:"created_at" db:"created_at"`
	UpdatedAt          time.Time              `json:"updated_at" db:"updated_at"`
}

// NotificationPreferences represents user notification preferences
type NotificationPreferences struct {
	ID                        string    `json:"id" db:"id"`
	UserID                    string    `json:"user_id" db:"user_id"`
	EmailEnabled              bool      `json:"email_enabled" db:"email_enabled"`
	PaymentNotifications      bool      `json:"payment_notifications" db:"payment_notifications"`
	ContributionNotifications bool      `json:"contribution_notifications" db:"contribution_notifications"`
	WithdrawalNotifications   bool      `json:"withdrawal_notifications" db:"withdrawal_notifications"`
	ProofNotifications        bool      `json:"proof_notifications" db:"proof_notifications"`
	GoalNotifications         bool      `json:"goal_notifications" db:"goal_notifications"`
	MarketingEmails           bool      `json:"marketing_emails" db:"marketing_emails"`
	CreatedAt                 time.Time `json:"created_at" db:"created_at"`
	UpdatedAt                 time.Time `json:"updated_at" db:"updated_at"`
}

// CreateNotificationRequest represents a request to create a notification
type CreateNotificationRequest struct {
	UserID  string                 `json:"user_id" binding:"required"`
	Type    NotificationType       `json:"type" binding:"required"`
	Title   string                 `json:"title" binding:"required"`
	Message string                 `json:"message" binding:"required"`
	Data    map[string]interface{} `json:"data"`
}

// UpdatePreferencesRequest represents a request to update notification preferences
type UpdatePreferencesRequest struct {
	EmailEnabled              *bool `json:"email_enabled"`
	PaymentNotifications      *bool `json:"payment_notifications"`
	ContributionNotifications *bool `json:"contribution_notifications"`
	WithdrawalNotifications   *bool `json:"withdrawal_notifications"`
	ProofNotifications        *bool `json:"proof_notifications"`
	GoalNotifications         *bool `json:"goal_notifications"`
	MarketingEmails           *bool `json:"marketing_emails"`
}

// ListNotificationsQuery represents query parameters for listing notifications
type ListNotificationsQuery struct {
	UserID   string `form:"user_id"`
	IsRead   *bool  `form:"is_read"`
	Type     string `form:"type"`
	Page     int    `form:"page"`
	PageSize int    `form:"page_size"`
}

// PaginatedNotifications represents paginated notification results
type PaginatedNotifications struct {
	Notifications []Notification `json:"notifications"`
	Total         int64          `json:"total"`
	Page          int            `json:"page"`
	PageSize      int            `json:"page_size"`
	TotalPages    int            `json:"total_pages"`
}
