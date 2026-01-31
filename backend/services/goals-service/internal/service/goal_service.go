package service

import (
	"errors"
	"fmt"
	"time"

	"github.com/gofund/services/goals-service/internal/repository"
	"github.com/gofund/shared/models"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

var (
	ErrGoalNotFound          = errors.New("goal not found")
	ErrMilestoneNotFound     = errors.New("milestone not found")
	ErrUnauthorized          = errors.New("unauthorized")
	ErrInsufficientBalance   = errors.New("insufficient balance")
	ErrBankDetailsRequired   = errors.New("bank details required for withdrawal")
	ErrInvalidGoalStatus     = errors.New("invalid goal status for this operation")
	ErrContributionNotFound  = errors.New("contribution not found")
	ErrProofNotFound         = errors.New("proof not found")
	ErrNotContributor        = errors.New("only contributors can vote")
	ErrAlreadyVoted          = errors.New("you have already voted on this proof")
)

// GoalService handles business logic for goals
type GoalService struct {
	repo *repository.Repository
}

// NewGoalService creates a new goal service
func NewGoalService(repo *repository.Repository) *GoalService {
	return &GoalService{repo: repo}
}

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

// CreateGoal creates a new goal with optional milestones
func (s *GoalService) CreateGoal(ownerID uuid.UUID, req CreateGoalRequest) (*models.Goal, error) {
	// Validate
	if req.TargetAmount <= 0 {
		return nil, errors.New("target amount must be greater than 0")
	}

	goal := &models.Goal{
		OwnerID:       ownerID,
		Title:         req.Title,
		Description:   req.Description,
		TargetAmount:  req.TargetAmount,
		Currency:      req.Currency,
		Deadline:      req.Deadline,
		Status:        models.GoalStatusOpen,
		BankName:      req.BankName,
		AccountNumber: req.AccountNumber,
		AccountName:   req.AccountName,
	}

	if err := s.repo.Goal.CreateGoal(goal); err != nil {
		return nil, err
	}

	// Create milestones if provided
	if len(req.Milestones) > 0 {
		for _, mReq := range req.Milestones {
			milestone := &models.Milestone{
				GoalID:             goal.ID,
				Title:              mReq.Title,
				Description:        mReq.Description,
				TargetAmount:       mReq.TargetAmount,
				OrderIndex:         mReq.OrderIndex,
				IsRecurring:        mReq.IsRecurring,
				RecurrenceType:     mReq.RecurrenceType,
				RecurrenceInterval: mReq.RecurrenceInterval,
				NextDueDate:        mReq.NextDueDate,
				Status:             models.MilestoneStatusPending,
			}

			if err := s.repo.Milestone.CreateMilestone(milestone); err != nil {
				return nil, err
			}
		}
	}

	// Reload with relationships
	return s.repo.Goal.GetGoalByID(goal.ID)
}

// GetGoal retrieves a goal by ID
func (s *GoalService) GetGoal(id uuid.UUID) (*models.Goal, error) {
	return s.repo.Goal.GetGoalByID(id)
}

// GetGoalsByOwner retrieves all goals for an owner
func (s *GoalService) GetGoalsByOwner(ownerID uuid.UUID) ([]models.Goal, error) {
	return s.repo.Goal.GetGoalsByOwnerID(ownerID)
}

// UpdateGoalRequest represents a request to update a goal
type UpdateGoalRequest struct {
	Title         *string
	Description   *string
	BankName      *string
	AccountNumber *string
	AccountName   *string
}

// UpdateGoal updates a goal
func (s *GoalService) UpdateGoal(goalID, userID uuid.UUID, req UpdateGoalRequest) (*models.Goal, error) {
	goal, err := s.repo.Goal.GetGoalByIDSimple(goalID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrGoalNotFound
		}
		return nil, err
	}

	// Check ownership
	if goal.OwnerID != userID {
		return nil, ErrUnauthorized
	}

	// Update fields
	if req.Title != nil {
		goal.Title = *req.Title
	}
	if req.Description != nil {
		goal.Description = *req.Description
	}
	if req.BankName != nil {
		goal.BankName = *req.BankName
	}
	if req.AccountNumber != nil {
		goal.AccountNumber = *req.AccountNumber
	}
	if req.AccountName != nil {
		goal.AccountName = *req.AccountName
	}

	if err := s.repo.Goal.UpdateGoal(goal); err != nil {
		return nil, err
	}

	return s.repo.Goal.GetGoalByID(goalID)
}

// CloseGoal closes a goal to new contributions
func (s *GoalService) CloseGoal(goalID, userID uuid.UUID) (*models.Goal, error) {
	goal, err := s.repo.Goal.GetGoalByIDSimple(goalID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrGoalNotFound
		}
		return nil, err
	}

	// Check ownership
	if goal.OwnerID != userID {
		return nil, ErrUnauthorized
	}

	// Check status
	if goal.Status != models.GoalStatusOpen {
		return nil, ErrInvalidGoalStatus
	}

	goal.Status = models.GoalStatusClosed
	if err := s.repo.Goal.UpdateGoal(goal); err != nil {
		return nil, err
	}

	return goal, nil
}

// CancelGoal cancels a goal
func (s *GoalService) CancelGoal(goalID, userID uuid.UUID) (*models.Goal, error) {
	goal, err := s.repo.Goal.GetGoalByIDSimple(goalID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrGoalNotFound
		}
		return nil, err
	}

	// Check ownership
	if goal.OwnerID != userID {
		return nil, ErrUnauthorized
	}

	goal.Status = models.GoalStatusCancelled
	if err := s.repo.Goal.UpdateGoal(goal); err != nil {
		return nil, err
	}

	return goal, nil
}

// GetGoalProgress returns progress information for a goal
func (s *GoalService) GetGoalProgress(goalID uuid.UUID) (*GoalProgress, error) {
	goal, err := s.repo.Goal.GetGoalByID(goalID)
	if err != nil {
		return nil, err
	}

	totalContributions, err := s.repo.Goal.GetTotalConfirmedContributions(goalID)
	if err != nil {
		return nil, err
	}

	totalWithdrawals, err := s.repo.Goal.GetTotalCompletedWithdrawals(goalID)
	if err != nil {
		return nil, err
	}

	contributorCount, err := s.repo.Goal.GetContributorCount(goalID)
	if err != nil {
		return nil, err
	}

	milestones, err := s.repo.Milestone.GetMilestonesByGoalID(goalID)
	if err != nil {
		return nil, err
	}

	// Calculate milestone progress
	milestoneProgress := make([]MilestoneProgress, len(milestones))
	for i, milestone := range milestones {
		milestoneContributions, _ := s.repo.Milestone.GetTotalConfirmedContributionsByMilestone(milestone.ID)
		milestoneProgress[i] = MilestoneProgress{
			Milestone:       milestone,
			CurrentAmount:   milestoneContributions,
			ProgressPercent: calculatePercent(milestoneContributions, milestone.TargetAmount),
		}
	}

	return &GoalProgress{
		Goal:               *goal,
		TotalContributions: totalContributions,
		TotalWithdrawals:   totalWithdrawals,
		AvailableBalance:   totalContributions - totalWithdrawals,
		ProgressPercent:    calculatePercent(totalContributions, goal.TargetAmount),
		ContributorCount:   contributorCount,
		Milestones:         milestoneProgress,
	}, nil
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

// CreateMilestone creates a new milestone for a goal
func (s *GoalService) CreateMilestone(goalID, userID uuid.UUID, req CreateMilestoneRequest) (*models.Milestone, error) {
	goal, err := s.repo.Goal.GetGoalByIDSimple(goalID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrGoalNotFound
		}
		return nil, err
	}

	// Check ownership
	if goal.OwnerID != userID {
		return nil, ErrUnauthorized
	}

	// Validate
	if req.TargetAmount <= 0 {
		return nil, errors.New("target amount must be greater than 0")
	}

	if req.IsRecurring {
		if req.RecurrenceType == nil {
			return nil, errors.New("recurrence type required for recurring milestones")
		}
		if req.NextDueDate == nil {
			return nil, errors.New("next due date required for recurring milestones")
		}
	}

	// Get next order index if not provided
	orderIndex := req.OrderIndex
	if orderIndex == 0 {
		orderIndex, err = s.repo.Milestone.GetNextOrderIndex(goalID)
		if err != nil {
			return nil, err
		}
	}

	milestone := &models.Milestone{
		GoalID:             goalID,
		Title:              req.Title,
		Description:        req.Description,
		TargetAmount:       req.TargetAmount,
		OrderIndex:         orderIndex,
		IsRecurring:        req.IsRecurring,
		RecurrenceType:     req.RecurrenceType,
		RecurrenceInterval: req.RecurrenceInterval,
		NextDueDate:        req.NextDueDate,
		Status:             models.MilestoneStatusPending,
	}

	if err := s.repo.Milestone.CreateMilestone(milestone); err != nil {
		return nil, err
	}

	return milestone, nil
}

// CompleteMilestone marks a milestone as completed and creates next if recurring
func (s *GoalService) CompleteMilestone(milestoneID, userID uuid.UUID) (*models.Milestone, *models.Milestone, error) {
	milestone, err := s.repo.Milestone.GetMilestoneByID(milestoneID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil, ErrMilestoneNotFound
		}
		return nil, nil, err
	}

	// Check goal ownership
	goal, err := s.repo.Goal.GetGoalByIDSimple(milestone.GoalID)
	if err != nil {
		return nil, nil, err
	}
	if goal.OwnerID != userID {
		return nil, nil, ErrUnauthorized
	}

	// Mark as completed
	now := time.Now()
	milestone.Status = models.MilestoneStatusCompleted
	milestone.CompletedAt = &now

	if err := s.repo.Milestone.UpdateMilestone(milestone); err != nil {
		return nil, nil, err
	}

	// If recurring, create next milestone
	var nextMilestone *models.Milestone
	if milestone.IsRecurring && milestone.RecurrenceType != nil {
		nextOrderIndex, _ := s.repo.Milestone.GetNextOrderIndex(milestone.GoalID)
		nextDueDate := calculateNextDueDate(*milestone.NextDueDate, *milestone.RecurrenceType, milestone.RecurrenceInterval)

		nextMilestone = &models.Milestone{
			GoalID:             milestone.GoalID,
			Title:              generateNextTitle(milestone.Title),
			Description:        milestone.Description,
			TargetAmount:       milestone.TargetAmount,
			OrderIndex:         nextOrderIndex,
			IsRecurring:        true,
			RecurrenceType:     milestone.RecurrenceType,
			RecurrenceInterval: milestone.RecurrenceInterval,
			NextDueDate:        &nextDueDate,
			Status:             models.MilestoneStatusPending,
		}

		if err := s.repo.Milestone.CreateMilestone(nextMilestone); err != nil {
			return milestone, nil, err
		}
	}

	return milestone, nextMilestone, nil
}

// Helper functions

func calculatePercent(current, target int64) float64 {
	if target == 0 {
		return 0
	}
	return (float64(current) / float64(target)) * 100
}

func calculateNextDueDate(current time.Time, recurrenceType models.RecurrenceType, interval int) time.Time {
	switch recurrenceType {
	case models.RecurrenceWeekly:
		return current.AddDate(0, 0, 7*interval)
	case models.RecurrenceMonthly:
		return current.AddDate(0, interval, 0)
	case models.RecurrenceSemester:
		return current.AddDate(0, 6*interval, 0) // 6 months
	case models.RecurrenceYearly:
		return current.AddDate(interval, 0, 0)
	default:
		return current
	}
}

func generateNextTitle(currentTitle string) string {
	// Simple implementation - can be enhanced
	// TODO: Extract number and increment (e.g., "Semester 1" -> "Semester 2")
	return currentTitle + " (Next)"
}

// ValidateBankDetails validates bank account details
func ValidateBankDetails(bankName, accountNumber, accountName string) error {
	if bankName == "" {
		return errors.New("bank name is required")
	}
	if accountNumber == "" {
		return errors.New("account number is required")
	}
	if len(accountNumber) != 10 {
		return errors.New("account number must be 10 digits")
	}
	if accountName == "" {
		return errors.New("account name is required")
	}
	return nil
}
