package service

import (
	"errors"
	"time"

	"github.com/gofund/goals-service/internal/dto"
	"github.com/gofund/shared/events"
	"github.com/gofund/shared/messaging"
	"github.com/gofund/shared/models"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

// RefundService handles refund business logic
type RefundService struct {
	db        *gorm.DB
	publisher messaging.Publisher
}

// NewRefundService creates a new refund service instance
func NewRefundService(db *gorm.DB, publisher messaging.Publisher) *RefundService {
	return &RefundService{
		db:        db,
		publisher: publisher,
	}
}

// InitiateRefund initiates a refund for a goal
func (rs *RefundService) InitiateRefund(initiatedBy uuid.UUID, req *dto.InitiateRefundRequest) (*models.Refund, error) {
	goalID, err := uuid.Parse(req.GoalID)
	if err != nil {
		return nil, errors.New("invalid goal ID")
	}

	// Start transaction
	tx := rs.db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// Get goal
	var goal models.Goal
	if err := tx.First(&goal, "id = ?", goalID).Error; err != nil {
		tx.Rollback()
		return nil, errors.New("goal not found")
	}

	// Verify initiator is goal owner
	if goal.OwnerID != initiatedBy {
		tx.Rollback()
		return nil, errors.New("only goal owner can initiate refunds")
	}

	// Verify goal is cancelled or closed
	if goal.Status != models.GoalStatusCancelled && goal.Status != models.GoalStatusClosed {
		tx.Rollback()
		return nil, errors.New("can only refund cancelled or closed goals")
	}

	// Check if refund already exists for this goal
	var existingRefund models.Refund
	err = tx.Where("goal_id = ? AND status IN ?", goalID, []models.RefundStatus{
		models.RefundStatusPending,
		models.RefundStatusProcessing,
	}).First(&existingRefund).Error
	
	if err == nil {
		tx.Rollback()
		return nil, errors.New("refund already in progress for this goal")
	}

	// Get all confirmed contributions
	var contributions []models.Contribution
	if err := tx.Where("goal_id = ? AND status = ?", goalID, models.ContributionStatusConfirmed).Find(&contributions).Error; err != nil {
		tx.Rollback()
		return nil, errors.New("failed to fetch contributions")
	}

	if len(contributions) == 0 {
		tx.Rollback()
		return nil, errors.New("no confirmed contributions to refund")
	}

	// Calculate total refund amount
	var totalContributed int64
	for _, contrib := range contributions {
		totalContributed += contrib.Amount
	}

	totalRefundAmount := int64(float64(totalContributed) * (req.RefundPercentage / 100.0))

	// Create refund record
	refund := &models.Refund{
		GoalID:            goalID,
		InitiatedBy:       initiatedBy,
		RefundPercentage:  req.RefundPercentage,
		TotalRefundAmount: totalRefundAmount,
		Currency:          goal.Currency,
		Reason:            req.Reason,
		Status:            models.RefundStatusPending,
	}

	if err := tx.Create(refund).Error; err != nil {
		tx.Rollback()
		return nil, errors.New("failed to create refund")
	}

	// Create refund disbursements for each contributor
	// First, get settlement account details for all users
	userIDs := make([]uuid.UUID, len(contributions))
	for i, contrib := range contributions {
		userIDs[i] = contrib.UserID
	}

	var users []models.User
	if err := tx.Where("id IN ?", userIDs).Find(&users).Error; err != nil {
		tx.Rollback()
		return nil, errors.New("failed to fetch user settlement accounts")
	}

	// Map users by ID for easy lookup
	userMap := make(map[uuid.UUID]*models.User)
	for i := range users {
		userMap[users[i].ID] = &users[i]
	}

	// Create disbursements
	for _, contrib := range contributions {
		refundAmount := int64(float64(contrib.Amount) * (req.RefundPercentage / 100.0))
		
		user := userMap[contrib.UserID]
		disbursement := &models.RefundDisbursement{
			RefundID:       refund.ID,
			ContributionID: contrib.ID,
			UserID:         contrib.UserID,
			Amount:         refundAmount,
			Currency:       contrib.Currency,
			Status:         models.RefundStatusPending,
		}

		// Include settlement account if available
		if user != nil {
			disbursement.SettlementBankName = user.SettlementBankName
			disbursement.SettlementAccountNumber = user.SettlementAccountNumber
			disbursement.SettlementAccountName = user.SettlementAccountName
		}

		if err := tx.Create(disbursement).Error; err != nil {
			tx.Rollback()
			return nil, errors.New("failed to create refund disbursement")
		}
	}

	// Commit transaction
	if err := tx.Commit().Error; err != nil {
		return nil, errors.New("failed to commit refund transaction")
	}

	// Load disbursements for response
	if err := rs.db.Preload("Disbursements").First(refund, refund.ID).Error; err != nil {
		return nil, errors.New("failed to load refund details")
	}

	// Emit RefundInitiated event
	if rs.publisher != nil {
		event := events.RefundInitiated{
			ID:                uuid.New().String(),
			RefundID:          refund.ID.String(),
			GoalID:            refund.GoalID.String(),
			InitiatedBy:       refund.InitiatedBy.String(),
			RefundPercentage:  refund.RefundPercentage,
			TotalRefundAmount: refund.TotalRefundAmount,
			CreatedAt:         time.Now().Unix(),
		}
		rs.publisher.Publish("RefundInitiated", event)
	}

	return refund, nil
}

// GetRefund retrieves a refund by ID
func (rs *RefundService) GetRefund(refundID uuid.UUID) (*models.Refund, error) {
	var refund models.Refund
	if err := rs.db.Preload("Disbursements").First(&refund, "id = ?", refundID).Error; err != nil {
		return nil, errors.New("refund not found")
	}
	return &refund, nil
}

// GetGoalRefunds retrieves all refunds for a goal
func (rs *RefundService) GetGoalRefunds(goalID uuid.UUID) ([]models.Refund, error) {
	var refunds []models.Refund
	if err := rs.db.Preload("Disbursements").Where("goal_id = ?", goalID).Find(&refunds).Error; err != nil {
		return nil, errors.New("failed to fetch refunds")
	}
	return refunds, nil
}

// UpdateRefundStatus updates the status of a refund
func (rs *RefundService) UpdateRefundStatus(refundID uuid.UUID, status models.RefundStatus) error {
	updates := map[string]interface{}{
		"status": status,
	}
	
	if status == models.RefundStatusCompleted {
		now := time.Now()
		updates["completed_at"] = &now
	}

	if err := rs.db.Model(&models.Refund{}).Where("id = ?", refundID).Updates(updates).Error; err != nil {
		return err
	}

	// Emit event if completed
	if status == models.RefundStatusCompleted && rs.publisher != nil {
		var refund models.Refund
		if err := rs.db.First(&refund, "id = ?", refundID).Error; err == nil {
			event := events.RefundCompleted{
				ID:                uuid.New().String(),
				RefundID:          refund.ID.String(),
				GoalID:            refund.GoalID.String(),
				TotalRefundAmount: refund.TotalRefundAmount,
				CompletedAt:       time.Now().Unix(),
			}
			rs.publisher.Publish("RefundCompleted", event)
		}
	}

	return nil
}

// UpdateDisbursementStatus updates the status of a refund disbursement
func (rs *RefundService) UpdateDisbursementStatus(disbursementID uuid.UUID, status models.RefundStatus, ledgerTxID *uuid.UUID) error {
	updates := map[string]interface{}{
		"status": status,
	}
	
	if ledgerTxID != nil {
		updates["ledger_transaction_id"] = ledgerTxID
	}
	
	if status == models.RefundStatusCompleted {
		now := time.Now()
		updates["completed_at"] = &now
	}

	if err := rs.db.Model(&models.RefundDisbursement{}).Where("id = ?", disbursementID).Updates(updates).Error; err != nil {
		return err
	}

	// Emit event if completed
	if status == models.RefundStatusCompleted && rs.publisher != nil {
		var disbursement models.RefundDisbursement
		if err := rs.db.Preload("Refund").First(&disbursement, "id = ?", disbursementID).Error; err == nil {
			event := events.ContributionRefunded{
				ID:             uuid.New().String(),
				ContributionID: disbursement.ContributionID.String(),
				UserID:         disbursement.UserID.String(),
				GoalID:         disbursement.Refund.GoalID.String(),
				RefundAmount:   disbursement.Amount,
				CreatedAt:      time.Now().Unix(),
			}
			rs.publisher.Publish("ContributionRefunded", event)
		}
	}

	return nil
}
