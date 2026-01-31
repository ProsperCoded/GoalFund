package repository

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/gofund/notifications-service/internal/models"
	"github.com/jmoiron/sqlx"
)

// PreferenceRepository handles database operations for notification preferences
type PreferenceRepository interface {
	Create(preferences *models.NotificationPreferences) error
	GetByUserID(userID string) (*models.NotificationPreferences, error)
	Update(userID string, updates models.UpdatePreferencesRequest) error
	CreateDefault(userID string) error
}

type preferenceRepository struct {
	db *sqlx.DB
}

// NewPreferenceRepository creates a new preference repository
func NewPreferenceRepository(db *sqlx.DB) PreferenceRepository {
	return &preferenceRepository{db: db}
}

// Create creates new notification preferences
func (r *preferenceRepository) Create(preferences *models.NotificationPreferences) error {
	query := `
		INSERT INTO notification_preferences (
			user_id, email_enabled, payment_notifications, contribution_notifications,
			withdrawal_notifications, proof_notifications, goal_notifications,
			marketing_emails, created_at, updated_at
		)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
		RETURNING id
	`

	now := time.Now()
	err := r.db.QueryRow(
		query,
		preferences.UserID,
		preferences.EmailEnabled,
		preferences.PaymentNotifications,
		preferences.ContributionNotifications,
		preferences.WithdrawalNotifications,
		preferences.ProofNotifications,
		preferences.GoalNotifications,
		preferences.MarketingEmails,
		now,
		now,
	).Scan(&preferences.ID)

	if err != nil {
		return fmt.Errorf("failed to create preferences: %w", err)
	}

	preferences.CreatedAt = now
	preferences.UpdatedAt = now
	return nil
}

// GetByUserID retrieves preferences by user ID
func (r *preferenceRepository) GetByUserID(userID string) (*models.NotificationPreferences, error) {
	query := `
		SELECT id, user_id, email_enabled, payment_notifications, contribution_notifications,
		       withdrawal_notifications, proof_notifications, goal_notifications,
		       marketing_emails, created_at, updated_at
		FROM notification_preferences
		WHERE user_id = $1
	`

	var preferences models.NotificationPreferences
	err := r.db.QueryRow(query, userID).Scan(
		&preferences.ID,
		&preferences.UserID,
		&preferences.EmailEnabled,
		&preferences.PaymentNotifications,
		&preferences.ContributionNotifications,
		&preferences.WithdrawalNotifications,
		&preferences.ProofNotifications,
		&preferences.GoalNotifications,
		&preferences.MarketingEmails,
		&preferences.CreatedAt,
		&preferences.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("preferences not found")
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get preferences: %w", err)
	}

	return &preferences, nil
}

// Update updates notification preferences
func (r *preferenceRepository) Update(userID string, updates models.UpdatePreferencesRequest) error {
	query := `
		UPDATE notification_preferences
		SET 
			email_enabled = COALESCE($1, email_enabled),
			payment_notifications = COALESCE($2, payment_notifications),
			contribution_notifications = COALESCE($3, contribution_notifications),
			withdrawal_notifications = COALESCE($4, withdrawal_notifications),
			proof_notifications = COALESCE($5, proof_notifications),
			goal_notifications = COALESCE($6, goal_notifications),
			marketing_emails = COALESCE($7, marketing_emails),
			updated_at = $8
		WHERE user_id = $9
	`

	now := time.Now()
	_, err := r.db.Exec(
		query,
		updates.EmailEnabled,
		updates.PaymentNotifications,
		updates.ContributionNotifications,
		updates.WithdrawalNotifications,
		updates.ProofNotifications,
		updates.GoalNotifications,
		updates.MarketingEmails,
		now,
		userID,
	)

	if err != nil {
		return fmt.Errorf("failed to update preferences: %w", err)
	}

	return nil
}

// CreateDefault creates default preferences for a new user
func (r *preferenceRepository) CreateDefault(userID string) error {
	preferences := &models.NotificationPreferences{
		UserID:                    userID,
		EmailEnabled:              true,
		PaymentNotifications:      true,
		ContributionNotifications: true,
		WithdrawalNotifications:   true,
		ProofNotifications:        true,
		GoalNotifications:         true,
		MarketingEmails:           false,
	}

	return r.Create(preferences)
}
