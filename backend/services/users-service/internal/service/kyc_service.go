package service

import (
	"errors"
	"regexp"
	"time"

	"github.com/gofund/shared/metrics"
	"github.com/gofund/users-service/internal/dto"
	"github.com/gofund/users-service/internal/repository"
	"github.com/google/uuid"
)

// KYCService handles KYC verification business logic
type KYCService struct {
	userRepo     *repository.UserRepository
	eventService *EventService
}

// NewKYCService creates a new KYC service instance
func NewKYCService(userRepo *repository.UserRepository, eventService *EventService) *KYCService {
	return &KYCService{
		userRepo:     userRepo,
		eventService: eventService,
	}
}

// SubmitNIN submits a NIN for KYC verification
// In this dummy implementation, we automatically verify the user
func (s *KYCService) SubmitNIN(userID string, req *dto.SubmitNINRequest) (*dto.KYCStatusResponse, error) {
	// Parse user ID
	id, err := uuid.Parse(userID)
	if err != nil {
		return nil, errors.New("invalid user ID")
	}

	// Get user
	user, err := s.userRepo.GetUserByID(id)
	if err != nil {
		return nil, errors.New("user not found")
	}

	// Check if already verified
	if user.KYCVerified {
		return nil, errors.New("user is already KYC verified")
	}

	// Validate NIN format (11 digits)
	if !s.isValidNIN(req.NIN) {
		return nil, errors.New("invalid NIN format - must be 11 digits")
	}

	// Check if NIN is already used by another user
	existingUser, err := s.userRepo.GetUserByNIN(req.NIN)
	if err == nil && existingUser != nil && existingUser.ID != user.ID {
		return nil, errors.New("NIN already registered to another account")
	}

	// Update user with NIN and mark as verified (dummy verification)
	now := time.Now()
	user.NIN = req.NIN
	user.KYCVerified = true
	user.KYCVerifiedAt = &now

	if err := s.userRepo.UpdateUser(user); err != nil {
		return nil, errors.New("failed to update KYC status")
	}

	// Publish KYC verified event
	if s.eventService != nil {
		if err := s.eventService.PublishKYCVerified(user); err != nil {
			// Log error but don't fail the verification
		}
	}

	// Track KYC verification metric
	metrics.TrackKYCVerification(user.ID.String())

	return &dto.KYCStatusResponse{
		KYCVerified:   user.KYCVerified,
		KYCVerifiedAt: user.KYCVerifiedAt,
		NIN:           s.maskNIN(user.NIN),
	}, nil
}

// GetKYCStatus retrieves the KYC status for a user
func (s *KYCService) GetKYCStatus(userID string) (*dto.KYCStatusResponse, error) {
	// Parse user ID
	id, err := uuid.Parse(userID)
	if err != nil {
		return nil, errors.New("invalid user ID")
	}

	// Get user
	user, err := s.userRepo.GetUserByID(id)
	if err != nil {
		return nil, errors.New("user not found")
	}

	return &dto.KYCStatusResponse{
		KYCVerified:   user.KYCVerified,
		KYCVerifiedAt: user.KYCVerifiedAt,
		NIN:           s.maskNIN(user.NIN),
	}, nil
}

// isValidNIN validates NIN format (11 digits)
func (s *KYCService) isValidNIN(nin string) bool {
	// NIN should be exactly 11 digits
	matched, _ := regexp.MatchString(`^\d{11}$`, nin)
	return matched
}

// maskNIN masks NIN for privacy (shows only last 4 digits)
func (s *KYCService) maskNIN(nin string) string {
	if nin == "" {
		return ""
	}
	if len(nin) != 11 {
		return "***"
	}
	return "*******" + nin[7:]
}
