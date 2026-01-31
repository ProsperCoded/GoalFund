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
type UpdatePreferencesRequest = models.UpdatePreferencesRequest

// ListNotificationsQuery represents query parameters for listing notifications
type ListNotificationsQuery = models.ListNotificationsQuery

// PaginatedNotifications represents paginated notification results
type PaginatedNotifications struct {
	Notifications []models.Notification `json:"notifications"`
	Total         int64                 `json:"total"`
	Page          int                   `json:"page"`
	PageSize      int                   `json:"page_size"`
	TotalPages    int                   `json:"total_pages"`
}
