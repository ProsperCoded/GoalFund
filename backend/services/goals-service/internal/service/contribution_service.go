package service

import (
	"errors"
	"time"

	"github.com/gofund/goals-service/internal/dto"
	"github.com/gofund/goals-service/internal/repository"
	"github.com/gofund/shared/events"
	"github.com/gofund/shared/messaging"
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

// CreateContribution creates a new contribution intent
func (s *ContributionService) CreateContribution(userID uuid.UUID, req dto.CreateContributionRequest) (*models.Contribution, error) {
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

// GetContributionByID retrieves a single contribution by ID
func (s *ContributionService) GetContributionByID(contributionID uuid.UUID) (*models.Contribution, error) {
	contribution, err := s.repo.Contribution.GetContributionByID(contributionID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrContributionNotFound
		}
		return nil, err
	}
	return contribution, nil
}

// WithdrawalService handles business logic for withdrawals
type WithdrawalService struct {
	repo *repository.Repository
}

// NewWithdrawalService creates a new withdrawal service
func NewWithdrawalService(repo *repository.Repository) *WithdrawalService {
	return &WithdrawalService{repo: repo}
}

// CreateWithdrawal creates a new withdrawal request
func (s *WithdrawalService) CreateWithdrawal(userID uuid.UUID, req dto.CreateWithdrawalRequest) (*models.Withdrawal, error) {
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
		bankName = goal.DepositBankName
	}
	if accountNumber == "" {
		accountNumber = goal.DepositAccountNumber
	}
	if accountName == "" {
		accountName = goal.DepositAccountName
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
	repo      *repository.Repository
	publisher messaging.Publisher
}

// NewProofService creates a new proof service
func NewProofService(repo *repository.Repository, publisher messaging.Publisher) *ProofService {
	return &ProofService{repo: repo, publisher: publisher}
}

// CreateProof creates a new proof
func (s *ProofService) CreateProof(userID uuid.UUID, req dto.CreateProofRequest) (*models.Proof, error) {
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

	// Emit event
	if s.publisher != nil {
		event := events.ProofSubmitted{
			ID:        uuid.New().String(),
			GoalID:    proof.GoalID.String(),
			ProofID:   proof.ID.String(),
			CreatedAt: time.Now().Unix(),
		}
		s.publisher.Publish("ProofSubmitted", event)
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
	repo      *repository.Repository
	publisher messaging.Publisher
}

// NewVoteService creates a new vote service
func NewVoteService(repo *repository.Repository, publisher messaging.Publisher) *VoteService {
	return &VoteService{repo: repo, publisher: publisher}
}

// CreateVote creates a new vote or updates existing
func (s *VoteService) CreateVote(userID uuid.UUID, req dto.CreateVoteRequest) (*models.Vote, error) {
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

	// Check if proof is now verified
	go s.checkProofVerification(proof.GoalID, req.ProofID)

	return vote, nil
}

func (s *VoteService) checkProofVerification(goalID, proofID uuid.UUID) {
	_, satisfied, err := s.repo.Vote.GetVoteStats(proofID)
	if err != nil {
		return
	}

	contributorCount, err := s.repo.Goal.GetContributorCount(goalID)
	if err != nil {
		return
	}

	// Threshold: max(3, 5% of contributors)
	threshold := int64(3)
	fivePercent := int64(float64(contributorCount) * 0.05)
	if fivePercent > threshold {
		threshold = fivePercent
	}

	if satisfied >= threshold {
		// Proof is verified!
		// Update goal status if needed? Or just emit event.
		// Usually Goal status updates to VERIFIED.
		
		if s.publisher != nil {
			event := events.ProofVerified{
				ID:        uuid.New().String(),
				GoalID:    goalID.String(),
				ProofID:   proofID.String(),
				CreatedAt: time.Now().Unix(),
			}
			s.publisher.Publish("ProofVerified", event)
		}
	}
}

// GetVotesByProof retrieves all votes for a proof
func (s *VoteService) GetVotesByProof(proofID uuid.UUID) ([]models.Vote, error) {
	return s.repo.Vote.GetVotesByProofID(proofID)
}

// GetVoteStats retrieves vote statistics for a proof
func (s *VoteService) GetVoteStats(proofID uuid.UUID) (*dto.VoteStats, error) {
	total, satisfied, err := s.repo.Vote.GetVoteStats(proofID)
	if err != nil {
		return nil, err
	}

	satisfactionRate := float64(0)
	if total > 0 {
		satisfactionRate = (float64(satisfied) / float64(total)) * 100
	}

	return &dto.VoteStats{
		TotalVotes:       total,
		SatisfiedVotes:   satisfied,
		UnsatisfiedVotes: total - satisfied,
		SatisfactionRate: satisfactionRate,
	}, nil
}
