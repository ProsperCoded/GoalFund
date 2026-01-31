package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/gofund/notifications-service/internal/models"
	"github.com/gofund/notifications-service/internal/service"
)

// NotificationHandler handles HTTP requests for notifications
type NotificationHandler struct {
	notificationService service.NotificationService
}

// NewNotificationHandler creates a new notification handler
func NewNotificationHandler(notificationService service.NotificationService) *NotificationHandler {
	return &NotificationHandler{
		notificationService: notificationService,
	}
}

// GetNotifications godoc
// @Summary Get user notifications
// @Description Get paginated notifications for the authenticated user
// @Tags notifications
// @Accept json
// @Produce json
// @Param page query int false "Page number" default(1)
// @Param page_size query int false "Page size" default(20)
// @Param is_read query bool false "Filter by read status"
// @Param type query string false "Filter by notification type"
// @Success 200 {object} models.PaginatedNotifications
// @Failure 401 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/v1/notifications [get]
func (h *NotificationHandler) GetNotifications(c *gin.Context) {
	// Get user ID from context (set by auth middleware)
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	// Parse query parameters
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))

	var query models.ListNotificationsQuery
	query.UserID = userID.(string)
	query.Page = page
	query.PageSize = pageSize

	if isReadStr := c.Query("is_read"); isReadStr != "" {
		isRead, _ := strconv.ParseBool(isReadStr)
		query.IsRead = &isRead
	}

	if notifType := c.Query("type"); notifType != "" {
		query.Type = notifType
	}

	// Get notifications
	result, err := h.notificationService.ListNotifications(query)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, result)
}

// GetNotification godoc
// @Summary Get a specific notification
// @Description Get a notification by ID
// @Tags notifications
// @Accept json
// @Produce json
// @Param id path string true "Notification ID"
// @Success 200 {object} models.Notification
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/v1/notifications/{id} [get]
func (h *NotificationHandler) GetNotification(c *gin.Context) {
	id := c.Param("id")

	notification, err := h.notificationService.GetNotification(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "notification not found"})
		return
	}

	// Verify user owns this notification
	userID, _ := c.Get("user_id")
	if notification.UserID != userID.(string) {
		c.JSON(http.StatusForbidden, gin.H{"error": "forbidden"})
		return
	}

	c.JSON(http.StatusOK, notification)
}

// MarkAsRead godoc
// @Summary Mark notification as read
// @Description Mark a notification as read
// @Tags notifications
// @Accept json
// @Produce json
// @Param id path string true "Notification ID"
// @Success 200 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/v1/notifications/{id}/read [put]
func (h *NotificationHandler) MarkAsRead(c *gin.Context) {
	id := c.Param("id")

	// Verify user owns this notification
	notification, err := h.notificationService.GetNotification(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "notification not found"})
		return
	}

	userID, _ := c.Get("user_id")
	if notification.UserID != userID.(string) {
		c.JSON(http.StatusForbidden, gin.H{"error": "forbidden"})
		return
	}

	// Mark as read
	if err := h.notificationService.MarkAsRead(id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "notification marked as read"})
}

// DeleteNotification godoc
// @Summary Delete a notification
// @Description Delete a notification by ID
// @Tags notifications
// @Accept json
// @Produce json
// @Param id path string true "Notification ID"
// @Success 200 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/v1/notifications/{id} [delete]
func (h *NotificationHandler) DeleteNotification(c *gin.Context) {
	id := c.Param("id")

	// Verify user owns this notification
	notification, err := h.notificationService.GetNotification(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "notification not found"})
		return
	}

	userID, _ := c.Get("user_id")
	if notification.UserID != userID.(string) {
		c.JSON(http.StatusForbidden, gin.H{"error": "forbidden"})
		return
	}

	// Delete notification
	if err := h.notificationService.DeleteNotification(id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "notification deleted"})
}

// GetUnreadCount godoc
// @Summary Get unread notification count
// @Description Get the count of unread notifications for the authenticated user
// @Tags notifications
// @Accept json
// @Produce json
// @Success 200 {object} map[string]int64
// @Failure 401 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/v1/notifications/unread/count [get]
func (h *NotificationHandler) GetUnreadCount(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	count, err := h.notificationService.GetUnreadCount(userID.(string))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"count": count})
}

// GetPreferences godoc
// @Summary Get notification preferences
// @Description Get notification preferences for the authenticated user
// @Tags preferences
// @Accept json
// @Produce json
// @Success 200 {object} models.NotificationPreferences
// @Failure 401 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/v1/notifications/preferences [get]
func (h *NotificationHandler) GetPreferences(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	preferences, err := h.notificationService.GetUserPreferences(userID.(string))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, preferences)
}

// UpdatePreferences godoc
// @Summary Update notification preferences
// @Description Update notification preferences for the authenticated user
// @Tags preferences
// @Accept json
// @Produce json
// @Param preferences body models.UpdatePreferencesRequest true "Preferences to update"
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/v1/notifications/preferences [put]
func (h *NotificationHandler) UpdatePreferences(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	var req models.UpdatePreferencesRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.notificationService.UpdateUserPreferences(userID.(string), req); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "preferences updated successfully"})
}

// HealthCheck godoc
// @Summary Health check
// @Description Check if the service is healthy
// @Tags health
// @Accept json
// @Produce json
// @Success 200 {object} map[string]string
// @Router /api/v1/notifications/health [get]
func (h *NotificationHandler) HealthCheck(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status":  "healthy",
		"service": "notifications-service",
	})
}
