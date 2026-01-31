package service

import (
	"errors"
	"fmt"
	"time"

	"github.com/gofund/ledger-service/internal/dto"
	"github.com/gofund/shared/models"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

// RefundLedgerService handles refund-related ledger operations
type RefundLedgerService struct {
	db *gorm.DB
}

// NewRefundLedgerService creates a new refund ledger service instance
func NewRefundLedgerService(db *gorm.DB) *RefundLedgerService {
	return &RefundLedgerService{db: db}
}

// ProcessRefund processes a refund by creating ledger entries
// This reverses the contribution ledger entry
func (rls *RefundLedgerService) ProcessRefund(req *dto.ProcessRefundRequest) (*uuid.UUID, error) {
	// Start transaction
	tx := rls.db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// Get or create user account
	var userAccount models.Account
	err := tx.Where("entity_id = ? AND account_type = ? AND currency = ?",
		req.UserID, models.AccountTypeUser, req.Currency).First(&userAccount).Error
	
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			tx.Rollback()
			return nil, errors.New("user account not found")
		}
		tx.Rollback()
		return nil, err
	}

	// Get goal account
	var goalAccount models.Account
	err = tx.Where("entity_id = ? AND account_type = ? AND currency = ?",
		req.GoalID, models.AccountTypeGoal, req.Currency).First(&goalAccount).Error
	
	if err != nil {
		tx.Rollback()
		return nil, errors.New("goal account not found")
	}

	// Create transaction record
	transaction := &models.Transaction{
		ID:          uuid.New(),
		Type:        models.TransactionTypeRefund,
		Description: fmt.Sprintf("Refund for contribution %s", req.ContributionID.String()),
		Amount:      req.Amount,
		Currency:    req.Currency,
		Metadata: map[string]interface{}{
			"refund_id":       req.RefundID.String(),
			"goal_id":         req.GoalID.String(),
			"user_id":         req.UserID.String(),
			"contribution_id": req.ContributionID.String(),
		},
		Status:          models.TransactionStatusCompleted,
		TransactionDate: time.Now(),
	}

	if err := tx.Create(transaction).Error; err != nil {
		tx.Rollback()
		return nil, errors.New("failed to create transaction record")
	}

	// Create ledger entries (reverse of contribution)
	// Debit goal account (decrease goal balance)
	debitEntry := &models.LedgerEntry{
		ID:            uuid.New(),
		AccountID:     goalAccount.ID,
		TransactionID: transaction.ID,
		EntryType:     models.EntryTypeDebit,
		Amount:        req.Amount,
		Currency:      req.Currency,
		Description:   fmt.Sprintf("Refund to user %s", req.UserID.String()),
		Metadata: map[string]interface{}{
			"refund_id":       req.RefundID.String(),
			"contribution_id": req.ContributionID.String(),
		},
		CreatedAt: time.Now(),
	}

	if err := tx.Create(debitEntry).Error; err != nil {
		tx.Rollback()
		return nil, errors.New("failed to create debit entry")
	}

	// Credit user account (increase user balance for refund)
	creditEntry := &models.LedgerEntry{
		ID:            uuid.New(),
		AccountID:     userAccount.ID,
		TransactionID: transaction.ID,
		EntryType:     models.EntryTypeCredit,
		Amount:        req.Amount,
		Currency:      req.Currency,
		Description:   fmt.Sprintf("Refund from goal %s", req.GoalID.String()),
		Metadata: map[string]interface{}{
			"refund_id":       req.RefundID.String(),
			"contribution_id": req.ContributionID.String(),
		},
		CreatedAt: time.Now(),
	}

	if err := tx.Create(creditEntry).Error; err != nil {
		tx.Rollback()
		return nil, errors.New("failed to create credit entry")
	}

	// Update balance snapshots
	if err := rls.updateBalanceSnapshot(tx, userAccount.ID); err != nil {
		tx.Rollback()
		return nil, errors.New("failed to update user balance snapshot")
	}

	if err := rls.updateBalanceSnapshot(tx, goalAccount.ID); err != nil {
		tx.Rollback()
		return nil, errors.New("failed to update goal balance snapshot")
	}

	// Commit transaction
	if err := tx.Commit().Error; err != nil {
		return nil, errors.New("failed to commit refund transaction")
	}

	return &transaction.ID, nil
}

// updateBalanceSnapshot updates the balance snapshot for an account
func (rls *RefundLedgerService) updateBalanceSnapshot(tx *gorm.DB, accountID uuid.UUID) error {
	// Calculate balance from ledger entries
	var balance int64
	
	// Sum credits
	var creditSum int64
	if err := tx.Model(&models.LedgerEntry{}).
		Where("account_id = ? AND entry_type = ?", accountID, models.EntryTypeCredit).
		Select("COALESCE(SUM(amount), 0)").
		Scan(&creditSum).Error; err != nil {
		return err
	}

	// Sum debits
	var debitSum int64
	if err := tx.Model(&models.LedgerEntry{}).
		Where("account_id = ? AND entry_type = ?", accountID, models.EntryTypeDebit).
		Select("COALESCE(SUM(amount), 0)").
		Scan(&debitSum).Error; err != nil {
		return err
	}

	balance = creditSum - debitSum

	// Get currency from account
	var account models.Account
	if err := tx.First(&account, "id = ?", accountID).Error; err != nil {
		return err
	}

	// Update or create balance snapshot
	var snapshot models.BalanceSnapshot
	err := tx.Where("account_id = ?", accountID).First(&snapshot).Error
	
	if errors.Is(err, gorm.ErrRecordNotFound) {
		// Create new snapshot
		snapshot = models.BalanceSnapshot{
			ID:        uuid.New(),
			AccountID: accountID,
			Balance:   balance,
			Currency:  account.Currency,
			UpdatedAt: time.Now(),
		}
		return tx.Create(&snapshot).Error
	}
	
	// Update existing snapshot
	snapshot.Balance = balance
	snapshot.UpdatedAt = time.Now()
	return tx.Save(&snapshot).Error
}

// GetAccountBalance retrieves the current balance for an account
func (rls *RefundLedgerService) GetAccountBalance(accountID uuid.UUID) (int64, error) {
	var snapshot models.BalanceSnapshot
	err := rls.db.Where("account_id = ?", accountID).First(&snapshot).Error
	
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return 0, nil
		}
		return 0, err
	}
	
	return snapshot.Balance, nil
}

// GetGoalBalance retrieves the balance for a goal
func (rls *RefundLedgerService) GetGoalBalance(goalID uuid.UUID, currency string) (int64, error) {
	var account models.Account
	err := rls.db.Where("entity_id = ? AND account_type = ? AND currency = ?",
		goalID, models.AccountTypeGoal, currency).First(&account).Error
	
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return 0, nil
		}
		return 0, err
	}
	
	return rls.GetAccountBalance(account.ID)
}

// GetUserBalance retrieves the balance for a user
func (rls *RefundLedgerService) GetUserBalance(userID uuid.UUID, currency string) (int64, error) {
	var account models.Account
	err := rls.db.Where("entity_id = ? AND account_type = ? AND currency = ?",
		userID, models.AccountTypeUser, currency).First(&account).Error
	
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return 0, nil
		}
		return 0, err
	}
	
	return rls.GetAccountBalance(account.ID)
}
