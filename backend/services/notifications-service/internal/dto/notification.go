package dto

import (
	"github.com/gofund/notifications-service/internal/models"
)

// CreateNotificationRequest represents a request to create a notification
type CreateNotificationRequest struct {
	UserID  string                  `json:"user_id" binding:"required"`
	Type    models.NotificationType `json:"type" binding:"required"`
	Title   string                  `json:"title" binding:"required"`
	Message string                  `json:"message" binding:"required"`
	Data    map[string]interface{}  `json:"data"`
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
	Notifications []models.Notification `json:"notifications"`
	Total         int64                 `json:"total"`
	Page          int                   `json:"page"`
	PageSize      int                   `json:"page_size"`
	TotalPages    int                   `json:"total_pages"`
}
