package service

import (
	"errors"
	"time"

	"github.com/gofund/services/goals-service/internal/repository"
	"github.com/gofund/shared/models"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

// ContributionService handles business logic for contributions
type ContributionService struct {
	repo *repository.Repository
}

// NewContributionService creates a new contribution service
func NewContributionService(repo *repository.Repository) *ContributionService {
	return &ContributionService{repo: repo}
}

// CreateContributionRequest represents a request to create a contribution
type CreateContributionRequest struct {
	GoalID      uuid.UUID
	MilestoneID *uuid.UUID
	Amount      int64
}

// CreateContribution creates a new contribution intent
func (s *ContributionService) CreateContribution(userID uuid.UUID, req CreateContributionRequest) (*models.Contribution, error) {
	// Validate amount
	if req.Amount <= 0 {
		return nil, errors.New("amount must be greater than 0")
	}

	// Check if goal exists and is open
	goal, err := s.repo.Goal.GetGoalByIDSimple(req.GoalID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrGoalNotFound
		}
		return nil, err
	}

	if goal.Status != models.GoalStatusOpen {
		return nil, errors.New("goal is not accepting contributions")
	}

	// Validate milestone if provided
	if req.MilestoneID != nil {
		milestone, err := s.repo.Milestone.GetMilestoneByID(*req.MilestoneID)
		if err != nil {
			return nil, errors.New("milestone not found")
		}
		if milestone.GoalID != req.GoalID {
			return nil, errors.New("milestone does not belong to this goal")
		}
	}

	contribution := &models.Contribution{
		GoalID:      req.GoalID,
		MilestoneID: req.MilestoneID,
		UserID:      userID,
		Amount:      req.Amount,
		Currency:    goal.Currency,
		Status:      models.ContributionStatusPending,
	}

	if err := s.repo.Contribution.CreateContribution(contribution); err != nil {
		return nil, err
	}

	return contribution, nil
}

// ConfirmContribution confirms a contribution after payment verification
func (s *ContributionService) ConfirmContribution(contributionID, paymentID uuid.UUID) error {
	contribution, err := s.repo.Contribution.GetContributionByID(contributionID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ErrContributionNotFound
		}
		return err
	}

	contribution.PaymentID = &paymentID
	contribution.Status = models.ContributionStatusConfirmed

	return s.repo.Contribution.UpdateContribution(contribution)
}

// GetContributionsByGoal retrieves all contributions for a goal
func (s *ContributionService) GetContributionsByGoal(goalID uuid.UUID) ([]models.Contribution, error) {
	return s.repo.Contribution.GetContributionsByGoalID(goalID)
}

// GetContributionsByUser retrieves all contributions by a user
func (s *ContributionService) GetContributionsByUser(userID uuid.UUID) ([]models.Contribution, error) {
	return s.repo.Contribution.GetContributionsByUserID(userID)
}

// WithdrawalService handles business logic for withdrawals
type WithdrawalService struct {
	repo *repository.Repository
}

// NewWithdrawalService creates a new withdrawal service
func NewWithdrawalService(repo *repository.Repository) *WithdrawalService {
	return &WithdrawalService{repo: repo}
}

// CreateWithdrawalRequest represents a request to create a withdrawal
type CreateWithdrawalRequest struct {
	GoalID        uuid.UUID
	MilestoneID   *uuid.UUID
	Amount        int64
	BankName      string
	AccountNumber string
	AccountName   string
}

// CreateWithdrawal creates a new withdrawal request
func (s *WithdrawalService) CreateWithdrawal(userID uuid.UUID, req CreateWithdrawalRequest) (*models.Withdrawal, error) {
	// Validate amount
	if req.Amount <= 0 {
		return nil, errors.New("amount must be greater than 0")
	}

	// Get goal
	goal, err := s.repo.Goal.GetGoalByIDSimple(req.GoalID)
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

	// Check goal status
	if goal.Status == models.GoalStatusCancelled {
		return nil, ErrInvalidGoalStatus
	}

	// Determine bank details (use provided or fall back to goal's bank details)
	bankName := req.BankName
	accountNumber := req.AccountNumber
	accountName := req.AccountName

	if bankName == "" {
		bankName = goal.BankName
	}
	if accountNumber == "" {
		accountNumber = goal.AccountNumber
	}
	if accountName == "" {
		accountName = goal.AccountName
	}

	// Validate bank details
	if err := ValidateBankDetails(bankName, accountNumber, accountName); err != nil {
		return nil, ErrBankDetailsRequired
	}

	// Calculate available balance
	totalContributions, err := s.repo.Goal.GetTotalConfirmedContributions(req.GoalID)
	if err != nil {
		return nil, err
	}

	totalWithdrawals, err := s.repo.Goal.GetTotalCompletedWithdrawals(req.GoalID)
	if err != nil {
		return nil, err
	}

	availableBalance := totalContributions - totalWithdrawals

	if req.Amount > availableBalance {
		return nil, ErrInsufficientBalance
	}

	// Validate milestone if provided
	if req.MilestoneID != nil {
		milestone, err := s.repo.Milestone.GetMilestoneByID(*req.MilestoneID)
		if err != nil {
			return nil, errors.New("milestone not found")
		}
		if milestone.GoalID != req.GoalID {
			return nil, errors.New("milestone does not belong to this goal")
		}
	}

	withdrawal := &models.Withdrawal{
		GoalID:        req.GoalID,
		MilestoneID:   req.MilestoneID,
		OwnerID:       userID,
		Amount:        req.Amount,
		Currency:      goal.Currency,
		BankName:      bankName,
		AccountNumber: accountNumber,
		AccountName:   accountName,
		Status:        models.WithdrawalStatusPending,
		RequestedAt:   time.Now(),
	}

	if err := s.repo.Withdrawal.CreateWithdrawal(withdrawal); err != nil {
		return nil, err
	}

	return withdrawal, nil
}

// CompleteWithdrawal marks a withdrawal as completed
func (s *WithdrawalService) CompleteWithdrawal(withdrawalID, ledgerTransactionID uuid.UUID) error {
	withdrawal, err := s.repo.Withdrawal.GetWithdrawalByID(withdrawalID)
	if err != nil {
		return err
	}

	now := time.Now()
	withdrawal.Status = models.WithdrawalStatusCompleted
	withdrawal.LedgerTransactionID = &ledgerTransactionID
	withdrawal.CompletedAt = &now

	return s.repo.Withdrawal.UpdateWithdrawal(withdrawal)
}

// GetWithdrawalsByGoal retrieves all withdrawals for a goal
func (s *WithdrawalService) GetWithdrawalsByGoal(goalID uuid.UUID) ([]models.Withdrawal, error) {
	return s.repo.Withdrawal.GetWithdrawalsByGoalID(goalID)
}

// ProofService handles business logic for proofs
type ProofService struct {
	repo *repository.Repository
}

// NewProofService creates a new proof service
func NewProofService(repo *repository.Repository) *ProofService {
	return &ProofService{repo: repo}
}

// CreateProofRequest represents a request to create a proof
type CreateProofRequest struct {
	GoalID      uuid.UUID
	MilestoneID *uuid.UUID
	Title       string
	Description string
	MediaURLs   []string
}

// CreateProof creates a new proof
func (s *ProofService) CreateProof(userID uuid.UUID, req CreateProofRequest) (*models.Proof, error) {
	// Get goal
	goal, err := s.repo.Goal.GetGoalByIDSimple(req.GoalID)
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

	// Validate milestone if provided
	if req.MilestoneID != nil {
		milestone, err := s.repo.Milestone.GetMilestoneByID(*req.MilestoneID)
		if err != nil {
			return nil, errors.New("milestone not found")
		}
		if milestone.GoalID != req.GoalID {
			return nil, errors.New("milestone does not belong to this goal")
		}
	}

	proof := &models.Proof{
		GoalID:      req.GoalID,
		MilestoneID: req.MilestoneID,
		SubmittedBy: userID,
		Title:       req.Title,
		Description: req.Description,
		MediaURLs:   req.MediaURLs,
		SubmittedAt: time.Now(),
	}

	if err := s.repo.Proof.CreateProof(proof); err != nil {
		return nil, err
	}

	return proof, nil
}

// GetProof retrieves a proof by ID
func (s *ProofService) GetProof(proofID uuid.UUID) (*models.Proof, error) {
	proof, err := s.repo.Proof.GetProofByID(proofID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrProofNotFound
		}
		return nil, err
	}
	return proof, nil
}

// GetProofsByGoal retrieves all proofs for a goal
func (s *ProofService) GetProofsByGoal(goalID uuid.UUID) ([]models.Proof, error) {
	return s.repo.Proof.GetProofsByGoalID(goalID)
}

// VoteService handles business logic for votes
type VoteService struct {
	repo *repository.Repository
}

// NewVoteService creates a new vote service
func NewVoteService(repo *repository.Repository) *VoteService {
	return &VoteService{repo: repo}
}

// CreateVoteRequest represents a request to create a vote
type CreateVoteRequest struct {
	ProofID     uuid.UUID
	IsSatisfied bool
	Comment     string
}

// CreateVote creates a new vote or updates existing
func (s *VoteService) CreateVote(userID uuid.UUID, req CreateVoteRequest) (*models.Vote, error) {
	// Get proof
	proof, err := s.repo.Proof.GetProofByID(req.ProofID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrProofNotFound
		}
		return nil, err
	}

	// Check if user is a contributor
	isContributor, err := s.repo.Goal.IsUserContributor(proof.GoalID, userID)
	if err != nil {
		return nil, err
	}
	if !isContributor {
		return nil, ErrNotContributor
	}

	// Check if user already voted
	existingVote, err := s.repo.Vote.GetVoteByProofAndVoter(req.ProofID, userID)
	if err == nil {
		// Update existing vote
		existingVote.IsSatisfied = req.IsSatisfied
		existingVote.Comment = req.Comment
		existingVote.VotedAt = time.Now()

		if err := s.repo.Vote.UpdateVote(existingVote); err != nil {
			return nil, err
		}
		return existingVote, nil
	}

	// Create new vote
	vote := &models.Vote{
		ProofID:     req.ProofID,
		VoterID:     userID,
		IsSatisfied: req.IsSatisfied,
		Comment:     req.Comment,
		VotedAt:     time.Now(),
	}

	if err := s.repo.Vote.CreateVote(vote); err != nil {
		return nil, err
	}

	return vote, nil
}

// GetVotesByProof retrieves all votes for a proof
func (s *VoteService) GetVotesByProof(proofID uuid.UUID) ([]models.Vote, error) {
	return s.repo.Vote.GetVotesByProofID(proofID)
}

// GetVoteStats retrieves vote statistics for a proof
func (s *VoteService) GetVoteStats(proofID uuid.UUID) (*VoteStats, error) {
	total, satisfied, err := s.repo.Vote.GetVoteStats(proofID)
	if err != nil {
		return nil, err
	}

	satisfactionRate := float64(0)
	if total > 0 {
		satisfactionRate = (float64(satisfied) / float64(total)) * 100
	}

	return &VoteStats{
		TotalVotes:       total,
		SatisfiedVotes:   satisfied,
		UnsatisfiedVotes: total - satisfied,
		SatisfactionRate: satisfactionRate,
	}, nil
}

// VoteStats represents vote statistics
type VoteStats struct {
	TotalVotes       int64
	SatisfiedVotes   int64
	UnsatisfiedVotes int64
	SatisfactionRate float64
}
