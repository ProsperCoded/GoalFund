package events

import (
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/gofund/goals-service/internal/service"
	"github.com/gofund/shared/events"
	"github.com/gofund/shared/messaging"
	"github.com/google/uuid"
)

// EventHandler handles incoming events from RabbitMQ
type EventHandler struct {
	contributionService *service.ContributionService
	goalService         *service.GoalService
	publisher           messaging.Publisher
}

// NewEventHandler creates a new event handler instance
func NewEventHandler(
	contributionService *service.ContributionService,
	goalService *service.GoalService,
	publisher messaging.Publisher,
) *EventHandler {
	return &EventHandler{
		contributionService: contributionService,
		goalService:         goalService,
		publisher:           publisher,
	}
}

// HandlePaymentVerified handles the PaymentVerified event
func (h *EventHandler) HandlePaymentVerified(data []byte) error {
	var event events.PaymentVerified
	if err := json.Unmarshal(data, &event); err != nil {
		return fmt.Errorf("failed to unmarshal PaymentVerified event: %w", err)
	}

	log.Printf("Received PaymentVerified event: GoalID=%s, UserID=%s, Amount=%d", event.GoalID, event.UserID, event.Amount)

	// In a real system, we'd have a contribution ID in the event metadata.
	// Since it's missing from the contract, we'll try to find a matching pending contribution.
	// Note: This is a simplified approach.
	
	goalID, err := uuid.Parse(event.GoalID)
	if err != nil {
		return fmt.Errorf("invalid goal ID in event: %w", err)
	}

	userID, err := uuid.Parse(event.UserID)
	if err != nil {
		return fmt.Errorf("invalid user ID in event: %w", err)
	}

	paymentID, err := uuid.Parse(event.PaymentID)
	if err != nil {
		// Log error but maybe we can proceed if paymentID is not strictly needed for business logic
		log.Printf("Warning: invalid payment ID in event: %v", err)
	}

	// Find contributions for this goal and user
	contributions, err := h.contributionService.GetContributionsByGoal(goalID)
	if err != nil {
		return fmt.Errorf("failed to fetch contributions for goal: %w", err)
	}

	var targetContributionID uuid.UUID
	for _, c := range contributions {
		if c.UserID == userID && c.Amount == event.Amount && c.Status == "PENDING" {
			targetContributionID = c.ID
			break
		}
	}

	if targetContributionID == uuid.Nil {
		log.Printf("No matching pending contribution found for GoalID=%s, UserID=%s, Amount=%d", event.GoalID, event.UserID, event.Amount)
		return nil
	}

	// Confirm contribution
	if err := h.contributionService.ConfirmContribution(targetContributionID, paymentID); err != nil {
		return fmt.Errorf("failed to confirm contribution: %w", err)
	}

	log.Printf("Confirmed contribution %s for goal %s", targetContributionID, goalID)

	// Check if goal reached its target and emit event if needed
	progress, err := h.goalService.GetGoalProgress(goalID)
	if err == nil && progress.ProgressPercent >= 100 && progress.Goal.Status == "OPEN" {
		if h.publisher != nil {
			goalFundedEvent := events.GoalFunded{
				ID:        uuid.New().String(),
				GoalID:    goalID.String(),
				Amount:    progress.TotalContributions,
				CreatedAt: time.Now().Unix(),
			}
			if err := h.publisher.Publish("GoalFunded", goalFundedEvent); err != nil {
				log.Printf("Failed to publish GoalFunded event: %v", err)
			} else {
				log.Printf("Goal %s is now fully funded! GoalFunded event published.", goalID)
			}
		}
	}

	return nil
}
