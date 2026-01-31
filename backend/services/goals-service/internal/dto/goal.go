package dto

import (
	"time"

	"github.com/gofund/shared/models"
)

// CreateGoalRequest represents a request to create a goal
type CreateGoalRequest struct {
	Title         string
	Description   string
	TargetAmount  int64
	Currency      string
	Deadline      *time.Time
	BankName      string
	AccountNumber string
	AccountName   string
	Milestones    []CreateMilestoneRequest
}

// CreateMilestoneRequest represents a request to create a milestone
type CreateMilestoneRequest struct {
	Title              string
	Description        string
	TargetAmount       int64
	OrderIndex         int
	IsRecurring        bool
	RecurrenceType     *models.RecurrenceType
	RecurrenceInterval int
	NextDueDate        *time.Time
}

// UpdateGoalRequest represents a request to update a goal
type UpdateGoalRequest struct {
	Title         *string
	Description   *string
	BankName      *string
	AccountNumber *string
	AccountName   *string
}

// GoalProgress represents goal progress information
type GoalProgress struct {
	Goal               models.Goal
	TotalContributions int64
	TotalWithdrawals   int64
	AvailableBalance   int64
	ProgressPercent    float64
	ContributorCount   int64
	Milestones         []MilestoneProgress
}

// MilestoneProgress represents milestone progress information
type MilestoneProgress struct {
	Milestone       models.Milestone
	CurrentAmount   int64
	ProgressPercent float64
}
