package service

import (
	"errors"

	"github.com/gofund/shared/models"
	"github.com/gofund/users-service/internal/dto"
	"github.com/gofund/users-service/internal/repository"
	"github.com/google/uuid"
)

// UserService handles user-related business logic
type UserService struct {
	userRepo *repository.UserRepository
}

// NewUserService creates a new user service instance
func NewUserService(userRepo *repository.UserRepository) *UserService {
	return &UserService{
		userRepo: userRepo,
	}
}

// GetProfile retrieves user profile by ID
func (s *UserService) GetProfile(userID string) (*dto.UserResponse, error) {
	id, err := uuid.Parse(userID)
	if err != nil {
		return nil, errors.New("invalid user ID")
	}

	user, err := s.userRepo.GetUserByID(id)
	if err != nil {
		return nil, errors.New("user not found")
	}

	return mapUserToResponse(user), nil
}

// UpdateProfile updates user profile
func (s *UserService) UpdateProfile(userID string, req *dto.UpdateProfileRequest) (*dto.UserResponse, error) {
	id, err := uuid.Parse(userID)
	if err != nil {
		return nil, errors.New("invalid user ID")
	}

	user, err := s.userRepo.GetUserByID(id)
	if err != nil {
		return nil, errors.New("user not found")
	}

	// Update fields if provided
	if req.FirstName != "" {
		user.FirstName = req.FirstName
	}
	if req.LastName != "" {
		user.LastName = req.LastName
	}
	if req.Phone != "" {
		user.Phone = req.Phone
	}

	if err := s.userRepo.UpdateUser(user); err != nil {
		return nil, errors.New("failed to update profile")
	}

	return mapUserToResponse(user), nil
}

// UpdateSettlementAccount updates user's settlement account details
func (s *UserService) UpdateSettlementAccount(userID string, bankName, accountNumber, accountName string) error {
	id, err := uuid.Parse(userID)
	if err != nil {
		return errors.New("invalid user ID")
	}

	user, err := s.userRepo.GetUserByID(id)
	if err != nil {
		return errors.New("user not found")
	}

	user.SettlementBankName = bankName
	user.SettlementAccountNumber = accountNumber
	user.SettlementAccountName = accountName

	return s.userRepo.UpdateUser(user)
}

// mapUserToResponse converts user model to response (Internal helper)
func mapUserToResponse(user *models.User) *dto.UserResponse {
	if user == nil {
		return nil
	}
	return &dto.UserResponse{
		ID:            user.ID.String(),
		Email:         user.Email,
		Username:      user.Username,
		FirstName:     user.FirstName,
		LastName:      user.LastName,
		Phone:         user.Phone,
		EmailVerified: user.EmailVerified,
		PhoneVerified: user.PhoneVerified,
		KYCVerified:   user.KYCVerified,
		KYCVerifiedAt: user.KYCVerifiedAt,
		Role:          user.Role,
		CreatedAt:     user.CreatedAt,
	}
}

