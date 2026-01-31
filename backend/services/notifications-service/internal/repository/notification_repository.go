package repository

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"time"

	"github.com/gofund/notifications-service/internal/models"
	"github.com/jmoiron/sqlx"
)

// NotificationRepository handles database operations for notifications
type NotificationRepository interface {
	Create(notification *models.Notification) error
	GetByID(id string) (*models.Notification, error)
	GetByUserID(userID string, page, pageSize int) ([]models.Notification, int64, error)
	List(query models.ListNotificationsQuery) ([]models.Notification, int64, error)
	MarkAsRead(id string) error
	MarkAsEmailSent(id string) error
	MarkAsEmailFailed(id string, reason string) error
	IncrementRetryCount(id string) error
	Delete(id string) error
	GetUnreadCount(userID string) (int64, error)
}

type notificationRepository struct {
	db *sqlx.DB
}

// NewNotificationRepository creates a new notification repository
func NewNotificationRepository(db *sqlx.DB) NotificationRepository {
	return &notificationRepository{db: db}
}

// Create creates a new notification
func (r *notificationRepository) Create(notification *models.Notification) error {
	dataJSON, err := json.Marshal(notification.Data)
	if err != nil {
		return fmt.Errorf("failed to marshal data: %w", err)
	}

	query := `
		INSERT INTO notifications (user_id, type, title, message, data, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		RETURNING id
	`

	now := time.Now()
	err = r.db.QueryRow(
		query,
		notification.UserID,
		notification.Type,
		notification.Title,
		notification.Message,
		dataJSON,
		now,
		now,
	).Scan(&notification.ID)

	if err != nil {
		return fmt.Errorf("failed to create notification: %w", err)
	}

	notification.CreatedAt = now
	notification.UpdatedAt = now
	return nil
}

// GetByID retrieves a notification by ID
func (r *notificationRepository) GetByID(id string) (*models.Notification, error) {
	query := `
		SELECT id, user_id, type, title, message, data, email_sent, email_sent_at, 
		       email_failed_reason, retry_count, is_read, read_at, created_at, updated_at
		FROM notifications
		WHERE id = $1
	`

	var notification models.Notification
	var dataJSON []byte

	err := r.db.QueryRow(query, id).Scan(
		&notification.ID,
		&notification.UserID,
		&notification.Type,
		&notification.Title,
		&notification.Message,
		&dataJSON,
		&notification.EmailSent,
		&notification.EmailSentAt,
		&notification.EmailFailedReason,
		&notification.RetryCount,
		&notification.IsRead,
		&notification.ReadAt,
		&notification.CreatedAt,
		&notification.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("notification not found")
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get notification: %w", err)
	}

	if err := json.Unmarshal(dataJSON, &notification.Data); err != nil {
		return nil, fmt.Errorf("failed to unmarshal data: %w", err)
	}

	return &notification, nil
}

// GetByUserID retrieves notifications for a user with pagination
func (r *notificationRepository) GetByUserID(userID string, page, pageSize int) ([]models.Notification, int64, error) {
	offset := (page - 1) * pageSize

	// Get total count
	var total int64
	countQuery := `SELECT COUNT(*) FROM notifications WHERE user_id = $1`
	if err := r.db.Get(&total, countQuery, userID); err != nil {
		return nil, 0, fmt.Errorf("failed to count notifications: %w", err)
	}

	// Get notifications
	query := `
		SELECT id, user_id, type, title, message, data, email_sent, email_sent_at, 
		       email_failed_reason, retry_count, is_read, read_at, created_at, updated_at
		FROM notifications
		WHERE user_id = $1
		ORDER BY created_at DESC
		LIMIT $2 OFFSET $3
	`

	rows, err := r.db.Query(query, userID, pageSize, offset)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to get notifications: %w", err)
	}
	defer rows.Close()

	var notifications []models.Notification
	for rows.Next() {
		var notification models.Notification
		var dataJSON []byte

		err := rows.Scan(
			&notification.ID,
			&notification.UserID,
			&notification.Type,
			&notification.Title,
			&notification.Message,
			&dataJSON,
			&notification.EmailSent,
			&notification.EmailSentAt,
			&notification.EmailFailedReason,
			&notification.RetryCount,
			&notification.IsRead,
			&notification.ReadAt,
			&notification.CreatedAt,
			&notification.UpdatedAt,
		)
		if err != nil {
			return nil, 0, fmt.Errorf("failed to scan notification: %w", err)
		}

		if err := json.Unmarshal(dataJSON, &notification.Data); err != nil {
			return nil, 0, fmt.Errorf("failed to unmarshal data: %w", err)
		}

		notifications = append(notifications, notification)
	}

	return notifications, total, nil
}

// List retrieves notifications with filters
func (r *notificationRepository) List(query models.ListNotificationsQuery) ([]models.Notification, int64, error) {
	// Set defaults
	if query.Page <= 0 {
		query.Page = 1
	}
	if query.PageSize <= 0 {
		query.PageSize = 20
	}

	offset := (query.Page - 1) * query.PageSize

	// Build WHERE clause
	whereClause := "WHERE 1=1"
	args := []interface{}{}
	argCount := 1

	if query.UserID != "" {
		whereClause += fmt.Sprintf(" AND user_id = $%d", argCount)
		args = append(args, query.UserID)
		argCount++
	}

	if query.IsRead != nil {
		whereClause += fmt.Sprintf(" AND is_read = $%d", argCount)
		args = append(args, *query.IsRead)
		argCount++
	}

	if query.Type != "" {
		whereClause += fmt.Sprintf(" AND type = $%d", argCount)
		args = append(args, query.Type)
		argCount++
	}

	// Get total count
	var total int64
	countQuery := fmt.Sprintf("SELECT COUNT(*) FROM notifications %s", whereClause)
	if err := r.db.Get(&total, countQuery, args...); err != nil {
		return nil, 0, fmt.Errorf("failed to count notifications: %w", err)
	}

	// Get notifications
	sqlQuery := fmt.Sprintf(`
		SELECT id, user_id, type, title, message, data, email_sent, email_sent_at, 
		       email_failed_reason, retry_count, is_read, read_at, created_at, updated_at
		FROM notifications
		%s
		ORDER BY created_at DESC
		LIMIT $%d OFFSET $%d
	`, whereClause, argCount, argCount+1)

	args = append(args, query.PageSize, offset)

	rows, err := r.db.Query(sqlQuery, args...)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to get notifications: %w", err)
	}
	defer rows.Close()

	var notifications []models.Notification
	for rows.Next() {
		var notification models.Notification
		var dataJSON []byte

		err := rows.Scan(
			&notification.ID,
			&notification.UserID,
			&notification.Type,
			&notification.Title,
			&notification.Message,
			&dataJSON,
			&notification.EmailSent,
			&notification.EmailSentAt,
			&notification.EmailFailedReason,
			&notification.RetryCount,
			&notification.IsRead,
			&notification.ReadAt,
			&notification.CreatedAt,
			&notification.UpdatedAt,
		)
		if err != nil {
			return nil, 0, fmt.Errorf("failed to scan notification: %w", err)
		}

		if err := json.Unmarshal(dataJSON, &notification.Data); err != nil {
			return nil, 0, fmt.Errorf("failed to unmarshal data: %w", err)
		}

		notifications = append(notifications, notification)
	}

	return notifications, total, nil
}

// MarkAsRead marks a notification as read
func (r *notificationRepository) MarkAsRead(id string) error {
	query := `
		UPDATE notifications
		SET is_read = true, read_at = $1, updated_at = $2
		WHERE id = $3
	`

	now := time.Now()
	_, err := r.db.Exec(query, now, now, id)
	if err != nil {
		return fmt.Errorf("failed to mark notification as read: %w", err)
	}

	return nil
}

// MarkAsEmailSent marks a notification as email sent
func (r *notificationRepository) MarkAsEmailSent(id string) error {
	query := `
		UPDATE notifications
		SET email_sent = true, email_sent_at = $1, updated_at = $2
		WHERE id = $3
	`

	now := time.Now()
	_, err := r.db.Exec(query, now, now, id)
	if err != nil {
		return fmt.Errorf("failed to mark notification as email sent: %w", err)
	}

	return nil
}

// MarkAsEmailFailed marks a notification as email failed
func (r *notificationRepository) MarkAsEmailFailed(id string, reason string) error {
	query := `
		UPDATE notifications
		SET email_failed_reason = $1, updated_at = $2
		WHERE id = $3
	`

	now := time.Now()
	_, err := r.db.Exec(query, reason, now, id)
	if err != nil {
		return fmt.Errorf("failed to mark notification as email failed: %w", err)
	}

	return nil
}

// IncrementRetryCount increments the retry count
func (r *notificationRepository) IncrementRetryCount(id string) error {
	query := `
		UPDATE notifications
		SET retry_count = retry_count + 1, updated_at = $1
		WHERE id = $2
	`

	now := time.Now()
	_, err := r.db.Exec(query, now, id)
	if err != nil {
		return fmt.Errorf("failed to increment retry count: %w", err)
	}

	return nil
}

// Delete deletes a notification
func (r *notificationRepository) Delete(id string) error {
	query := `DELETE FROM notifications WHERE id = $1`

	_, err := r.db.Exec(query, id)
	if err != nil {
		return fmt.Errorf("failed to delete notification: %w", err)
	}

	return nil
}

// GetUnreadCount gets the count of unread notifications for a user
func (r *notificationRepository) GetUnreadCount(userID string) (int64, error) {
	var count int64
	query := `SELECT COUNT(*) FROM notifications WHERE user_id = $1 AND is_read = false`

	if err := r.db.Get(&count, query, userID); err != nil {
		return 0, fmt.Errorf("failed to get unread count: %w", err)
	}

	return count, nil
}
