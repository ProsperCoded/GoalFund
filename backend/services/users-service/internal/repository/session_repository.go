package repository

import (
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/gofund/shared/models"
	"gorm.io/gorm"
)

// SessionRepository handles session database operations
type SessionRepository struct {
	db *gorm.DB
}

// NewSessionRepository creates a new session repository instance
func NewSessionRepository(db *gorm.DB) *SessionRepository {
	return &SessionRepository{db: db}
}

// CreateSession creates a new session
func (r *SessionRepository) CreateSession(session *models.Session) error {
	return r.db.Create(session).Error
}

// GetSessionByTokenHash retrieves a session by token hash
func (r *SessionRepository) GetSessionByTokenHash(tokenHash string) (*models.Session, error) {
	var session models.Session
	if err := r.db.Preload("User").First(&session, "token_hash = ?", tokenHash).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("session not found")
		}
		return nil, err
	}
	return &session, nil
}

// GetUserSessions retrieves all sessions for a user
func (r *SessionRepository) GetUserSessions(userID uuid.UUID) ([]models.Session, error) {
	var sessions []models.Session
	if err := r.db.Where("user_id = ?", userID).Find(&sessions).Error; err != nil {
		return nil, err
	}
	return sessions, nil
}

// DeleteSession deletes a session by token hash
func (r *SessionRepository) DeleteSession(tokenHash string) error {
	result := r.db.Where("token_hash = ?", tokenHash).Delete(&models.Session{})
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return errors.New("session not found")
	}
	return nil
}

// DeleteUserSessions deletes all sessions for a user
func (r *SessionRepository) DeleteUserSessions(userID uuid.UUID) error {
	return r.db.Where("user_id = ?", userID).Delete(&models.Session{}).Error
}

// DeleteExpiredSessions deletes all expired sessions
func (r *SessionRepository) DeleteExpiredSessions() error {
	return r.db.Where("expires_at < ?", time.Now()).Delete(&models.Session{}).Error
}

// UpdateSession updates an existing session
func (r *SessionRepository) UpdateSession(session *models.Session) error {
	return r.db.Save(session).Error
}

// SessionExists checks if a session exists by token hash
func (r *SessionRepository) SessionExists(tokenHash string) (bool, error) {
	var count int64
	if err := r.db.Model(&models.Session{}).Where("token_hash = ?", tokenHash).Count(&count).Error; err != nil {
		return false, err
	}
	return count > 0, nil
}