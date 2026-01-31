package service

import (
	"time"

	"github.com/gofund/shared/events"
	"github.com/gofund/shared/messaging"
	"github.com/gofund/shared/models"
	"github.com/google/uuid"
)

// EventService handles publishing events for user-related actions
type EventService struct {
	publisher messaging.Publisher
}

// NewEventService creates a new event service instance
func NewEventService(publisher messaging.Publisher) *EventService {
	return &EventService{
		publisher: publisher,
	}
}

// PublishUserSignedUp publishes a UserSignedUp event
func (s *EventService) PublishUserSignedUp(user *models.User) error {
	event := events.UserSignedUp{
		ID:        uuid.New().String(),
		UserID:    user.ID.String(),
		Email:     user.Email,
		Username:  user.Username,
		CreatedAt: time.Now().Unix(),
	}

	return s.publisher.Publish("UserSignedUp", event)
}

// PublishPasswordResetRequested publishes a PasswordResetRequested event
func (s *EventService) PublishPasswordResetRequested(user *models.User, token string) error {
	event := events.PasswordResetRequested{
		ID:        uuid.New().String(),
		UserID:    user.ID.String(),
		Email:     user.Email,
		Token:     token,
		CreatedAt: time.Now().Unix(),
	}

	return s.publisher.Publish("PasswordResetRequested", event)
}

// PublishEmailVerificationRequested publishes an EmailVerificationRequested event
func (s *EventService) PublishEmailVerificationRequested(user *models.User, token string) error {
	event := events.EmailVerificationRequested{
		ID:        uuid.New().String(),
		UserID:    user.ID.String(),
		Email:     user.Email,
		Token:     token,
		CreatedAt: time.Now().Unix(),
	}

	return s.publisher.Publish("EmailVerificationRequested", event)
}

// PublishKYCVerified publishes a KYCVerified event
func (s *EventService) PublishKYCVerified(user *models.User) error {
	event := events.KYCVerified{
		ID:        uuid.New().String(),
		UserID:    user.ID.String(),
		Email:     user.Email,
		Username:  user.Username,
		CreatedAt: time.Now().Unix(),
	}

	return s.publisher.Publish("KYCVerified", event)
}
