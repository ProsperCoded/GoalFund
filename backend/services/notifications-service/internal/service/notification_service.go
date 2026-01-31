package service

import (
	"fmt"
	"log"
	"math"

	"github.com/gofund/notifications-service/internal/dto"
	"github.com/gofund/notifications-service/internal/models"
	"github.com/gofund/notifications-service/internal/repository"
)

// NotificationService handles notification business logic
type NotificationService interface {
	CreateNotification(req dto.CreateNotificationRequest) (*models.Notification, error)
	GetNotification(id string) (*models.Notification, error)
	GetUserNotifications(userID string, page, pageSize int) (*dto.PaginatedNotifications, error)
	ListNotifications(query dto.ListNotificationsQuery) (*dto.PaginatedNotifications, error)
	MarkAsRead(id string) error
	DeleteNotification(id string) error
	GetUnreadCount(userID string) (int64, error)
	
	// Preference methods
	GetUserPreferences(userID string) (*models.NotificationPreferences, error)
	UpdateUserPreferences(userID string, updates dto.UpdatePreferencesRequest) error
	CreateDefaultPreferences(userID string) error
}

type notificationService struct {
	notificationRepo repository.NotificationRepository
	preferenceRepo   repository.PreferenceRepository
	emailService     EmailService
}

// NewNotificationService creates a new notification service
func NewNotificationService(
	notificationRepo repository.NotificationRepository,
	preferenceRepo repository.PreferenceRepository,
	emailService EmailService,
) NotificationService {
	return &notificationService{
		notificationRepo: notificationRepo,
		preferenceRepo:   preferenceRepo,
		emailService:     emailService,
	}
}

// CreateNotification creates a new notification
func (s *notificationService) CreateNotification(req dto.CreateNotificationRequest) (*models.Notification, error) {
	notification := &models.Notification{
		UserID:  req.UserID,
		Type:    req.Type,
		Title:   req.Title,
		Message: req.Message,
		Data:    req.Data,
	}

	if err := s.notificationRepo.Create(notification); err != nil {
		return nil, fmt.Errorf("failed to create notification: %w", err)
	}

	// Try to send email asynchronously
	go s.sendEmailNotification(notification)

	return notification, nil
}

// sendEmailNotification sends an email for a notification
func (s *notificationService) sendEmailNotification(notification *models.Notification) {
	// Check user preferences
	preferences, err := s.preferenceRepo.GetByUserID(notification.UserID)
	if err != nil {
		log.Printf("Failed to get preferences for user %s: %v", notification.UserID, err)
		// Create default preferences if not found
		if err := s.preferenceRepo.CreateDefault(notification.UserID); err != nil {
			log.Printf("Failed to create default preferences: %v", err)
		}
		// Assume email is enabled by default
	} else if !preferences.EmailEnabled {
		log.Printf("Email notifications disabled for user %s", notification.UserID)
		return
	}

	// Get user email from notification data
	email, ok := notification.Data["email"].(string)
	if !ok || email == "" {
		log.Printf("No email found in notification data for user %s", notification.UserID)
		s.notificationRepo.MarkAsEmailFailed(notification.ID, "no email address")
		return
	}

	// Send email
	subject := notification.Title
	htmlBody := fmt.Sprintf(`
		<!DOCTYPE html>
		<html>
		<head>
			<meta charset="UTF-8">
			<style>
				body { font-family: Arial, sans-serif; line-height: 1.6; color: #333; }
				.container { max-width: 600px; margin: 0 auto; padding: 20px; }
				.header { background-color: #4CAF50; color: white; padding: 20px; text-align: center; }
				.content { padding: 20px; background-color: #f9f9f9; }
				.footer { text-align: center; padding: 20px; font-size: 12px; color: #666; }
			</style>
		</head>
		<body>
			<div class="container">
				<div class="header">
					<h1>GoFund</h1>
				</div>
				<div class="content">
					<h2>%s</h2>
					<p>%s</p>
				</div>
				<div class="footer">
					<p>&copy; 2026 GoFund. All rights reserved.</p>
					<p>This is an automated message, please do not reply.</p>
				</div>
			</div>
		</body>
		</html>
	`, notification.Title, notification.Message)

	err = s.emailService.SendEmail(email, subject, htmlBody, notification.Message)
	if err != nil {
		log.Printf("Failed to send email for notification %s: %v", notification.ID, err)
		s.notificationRepo.MarkAsEmailFailed(notification.ID, err.Error())
		s.notificationRepo.IncrementRetryCount(notification.ID)
		return
	}

	// Mark as sent
	if err := s.notificationRepo.MarkAsEmailSent(notification.ID); err != nil {
		log.Printf("Failed to mark notification as sent: %v", err)
	}
}

// GetNotification retrieves a notification by ID
func (s *notificationService) GetNotification(id string) (*models.Notification, error) {
	return s.notificationRepo.GetByID(id)
}

// GetUserNotifications retrieves notifications for a user with pagination
func (s *notificationService) GetUserNotifications(userID string, page, pageSize int) (*dto.PaginatedNotifications, error) {
	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = 20
	}

	notifications, total, err := s.notificationRepo.GetByUserID(userID, page, pageSize)
	if err != nil {
		return nil, fmt.Errorf("failed to get user notifications: %w", err)
	}

	totalPages := int(math.Ceil(float64(total) / float64(pageSize)))

	return &dto.PaginatedNotifications{
		Notifications: notifications,
		Total:         total,
		Page:          page,
		PageSize:      pageSize,
		TotalPages:    totalPages,
	}, nil
}

// ListNotifications retrieves notifications with filters
func (s *notificationService) ListNotifications(query dto.ListNotificationsQuery) (*dto.PaginatedNotifications, error) {
	notifications, total, err := s.notificationRepo.List(query)
	if err != nil {
		return nil, fmt.Errorf("failed to list notifications: %w", err)
	}

	if query.PageSize <= 0 {
		query.PageSize = 20
	}
	totalPages := int(math.Ceil(float64(total) / float64(query.PageSize)))

	return &dto.PaginatedNotifications{
		Notifications: notifications,
		Total:         total,
		Page:          query.Page,
		PageSize:      query.PageSize,
		TotalPages:    totalPages,
	}, nil
}

// MarkAsRead marks a notification as read
func (s *notificationService) MarkAsRead(id string) error {
	return s.notificationRepo.MarkAsRead(id)
}

// DeleteNotification deletes a notification
func (s *notificationService) DeleteNotification(id string) error {
	return s.notificationRepo.Delete(id)
}

// GetUnreadCount gets the count of unread notifications
func (s *notificationService) GetUnreadCount(userID string) (int64, error) {
	return s.notificationRepo.GetUnreadCount(userID)
}

// GetUserPreferences retrieves user notification preferences
func (s *notificationService) GetUserPreferences(userID string) (*models.NotificationPreferences, error) {
	preferences, err := s.preferenceRepo.GetByUserID(userID)
	if err != nil {
		// Create default preferences if not found
		if err := s.preferenceRepo.CreateDefault(userID); err != nil {
			return nil, fmt.Errorf("failed to create default preferences: %w", err)
		}
		return s.preferenceRepo.GetByUserID(userID)
	}
	return preferences, nil
}

// UpdateUserPreferences updates user notification preferences
func (s *notificationService) UpdateUserPreferences(userID string, updates dto.UpdatePreferencesRequest) error {
	return s.preferenceRepo.Update(userID, updates)
}

// CreateDefaultPreferences creates default preferences for a user
func (s *notificationService) CreateDefaultPreferences(userID string) error {
	return s.preferenceRepo.CreateDefault(userID)
}
