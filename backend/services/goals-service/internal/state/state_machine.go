package state

import (
	"errors"

	"github.com/gofund/shared/models"
)

var (
	ErrInvalidStateTransition = errors.New("invalid state transition")
)

// GoalStateMachine manages the state transitions for goals
type GoalStateMachine struct{}

// NewGoalStateMachine creates a new state machine instance
func NewGoalStateMachine() *GoalStateMachine {
	return &GoalStateMachine{}
}

// CanTransition checks if a transition from current to next is valid
func (sm *GoalStateMachine) CanTransition(current, next models.GoalStatus) bool {
	switch current {
	case models.GoalStatusOpen:
		return next == models.GoalStatusFunded || next == models.GoalStatusCancelled || next == models.GoalStatusClosed
	case models.GoalStatusFunded:
		return next == models.GoalStatusWithdrawn || next == models.GoalStatusCancelled
	case models.GoalStatusWithdrawn:
		return next == models.GoalStatusProofSubmitted
	case models.GoalStatusProofSubmitted:
		return next == models.GoalStatusVerified || next == models.GoalStatusOpen // If proof rejected, maybe back to open or funded
	case models.GoalStatusVerified:
		return false // Terminal state
	case models.GoalStatusCancelled:
		return false // Terminal state
	case models.GoalStatusClosed:
		return next == models.GoalStatusOpen || next == models.GoalStatusCancelled
	default:
		return false
	}
}

// ValidateTransition returns an error if the transition is invalid
func (sm *GoalStateMachine) ValidateTransition(current, next models.GoalStatus) error {
	if !sm.CanTransition(current, next) {
		return ErrInvalidStateTransition
	}
	return nil
}
