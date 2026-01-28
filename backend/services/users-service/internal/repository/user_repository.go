package repository

import (
	"errors"
	"strings"
	"time"

	"github.com/gofund/shared/models"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

// UserRepository handles user database operations
type UserRepository struct {
	db *gorm.DB
}

// NewUserRepository creates a new user repository instance
func NewUserRepository(db *gorm.DB) *UserRepository {
	return &UserRepository{db: db}
}

// CreateUser creates a new user
func (r *UserRepository) CreateUser(user *models.User) error {
	if err := r.db.Create(user).Error; err != nil {
		// Check for unique constraint violations
		if strings.Contains(err.Error(), "duplicate key") {
			if strings.Contains(err.Error(), "email") {
				return errors.New("email already exists")
			}
			if strings.Contains(err.Error(), "username") {
				return errors.New("username already exists")
			}
		}
		return err
	}
	return nil
}

// GetUserByID retrieves a user by ID
func (r *UserRepository) GetUserByID(id uuid.UUID) (*models.User, error) {
	var user models.User
	if err := r.db.First(&user, "id = ?", id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("user not found")
		}
		return nil, err
	}
	return &user, nil
}

// GetUserByEmail retrieves a user by email
func (r *UserRepository) GetUserByEmail(email string) (*models.User, error) {
	var user models.User
	if err := r.db.First(&user, "email = ?", email).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("user not found")
		}
		return nil, err
	}
	return &user, nil
}

// GetUserByUsername retrieves a user by username
func (r *UserRepository) GetUserByUsername(username string) (*models.User, error) {
	var user models.User
	if err := r.db.First(&user, "username = ?", username).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("user not found")
		}
		return nil, err
	}
	return &user, nil
}

// UpdateUser updates an existing user
func (r *UserRepository) UpdateUser(user *models.User) error {
	return r.db.Save(user).Error
}

// DeleteUser deletes a user by ID
func (r *UserRepository) DeleteUser(id uuid.UUID) error {
	result := r.db.Delete(&models.User{}, "id = ?", id)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return errors.New("user not found")
	}
	return nil
}

// EmailExists checks if an email already exists
func (r *UserRepository) EmailExists(email string) (bool, error) {
	var count int64
	if err := r.db.Model(&models.User{}).Where("email = ?", email).Count(&count).Error; err != nil {
		return false, err
	}
	return count > 0, nil
}

// UsernameExists checks if a username already exists
func (r *UserRepository) UsernameExists(username string) (bool, error) {
	var count int64
	if err := r.db.Model(&models.User{}).Where("username = ?", username).Count(&count).Error; err != nil {
		return false, err
	}
	return count > 0, nil
}

// CreatePasswordResetToken creates a new password reset token
func (r *UserRepository) CreatePasswordResetToken(token *models.PasswordResetToken) error {
	return r.db.Create(token).Error
}

// GetPasswordResetToken retrieves a password reset token by token hash
func (r *UserRepository) GetPasswordResetToken(tokenHash string) (*models.PasswordResetToken, error) {
	var token models.PasswordResetToken
	if err := r.db.First(&token, "token_hash = ? AND used = false", tokenHash).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("token not found")
		}
		return nil, err
	}
	return &token, nil
}

// MarkPasswordResetTokenUsed marks a password reset token as used
func (r *UserRepository) MarkPasswordResetTokenUsed(tokenHash string) error {
	return r.db.Model(&models.PasswordResetToken{}).
		Where("token_hash = ?", tokenHash).
		Update("used", true).Error
}

// DeleteExpiredPasswordResetTokens deletes expired password reset tokens
func (r *UserRepository) DeleteExpiredPasswordResetTokens() error {
	return r.db.Where("expires_at < ?", time.Now()).Delete(&models.PasswordResetToken{}).Error
}