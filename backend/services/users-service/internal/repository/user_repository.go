package repository

import (
	"errors"
	"strings"

	"github.com/google/uuid"
	"github.com/gofund/shared/models"
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
	if err := r.db.Preload("Roles").First(&user, "id = ?", id).Error; err != nil {
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
	if err := r.db.Preload("Roles").First(&user, "email = ?", email).Error; err != nil {
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
	if err := r.db.Preload("Roles").First(&user, "username = ?", username).Error; err != nil {
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

// GetUserRoles retrieves user roles as string slice
func (r *UserRepository) GetUserRoles(userID uuid.UUID) ([]string, error) {
	var roles []string
	if err := r.db.Table("roles").
		Select("roles.name").
		Joins("JOIN user_roles ON user_roles.role_id = roles.id").
		Where("user_roles.user_id = ?", userID).
		Pluck("name", &roles).Error; err != nil {
		return nil, err
	}
	return roles, nil
}

// AssignRole assigns a role to a user
func (r *UserRepository) AssignRole(userID, roleID uuid.UUID) error {
	userRole := models.UserRole{
		UserID: userID,
		RoleID: roleID,
	}
	return r.db.Create(&userRole).Error
}

// RemoveRole removes a role from a user
func (r *UserRepository) RemoveRole(userID, roleID uuid.UUID) error {
	return r.db.Where("user_id = ? AND role_id = ?", userID, roleID).Delete(&models.UserRole{}).Error
}

// GetRoleByName retrieves a role by name
func (r *UserRepository) GetRoleByName(name string) (*models.Role, error) {
	var role models.Role
	if err := r.db.First(&role, "name = ?", name).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("role not found")
		}
		return nil, err
	}
	return &role, nil
}