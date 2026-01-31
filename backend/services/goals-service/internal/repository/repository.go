package repository

import (

	"github.com/gofund/shared/models"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

// GoalRepository handles database operations for goals
type GoalRepository struct {
	db *gorm.DB
}

// NewGoalRepository creates a new goal repository
func NewGoalRepository(db *gorm.DB) *GoalRepository {
	return &GoalRepository{db: db}
}

// CreateGoal creates a new goal
func (r *GoalRepository) CreateGoal(goal *models.Goal) error {
	return r.db.Create(goal).Error
}

// GetGoalByID retrieves a goal by ID
func (r *GoalRepository) GetGoalByID(id uuid.UUID) (*models.Goal, error) {
	var goal models.Goal
	err := r.db.Preload("Milestones").
		Preload("Contributions").
		Preload("Withdrawals").
		Preload("Proofs").
		First(&goal, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &goal, nil
}

// GetGoalByIDSimple retrieves a goal without preloading relationships
func (r *GoalRepository) GetGoalByIDSimple(id uuid.UUID) (*models.Goal, error) {
	var goal models.Goal
	err := r.db.First(&goal, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &goal, nil
}

// GetGoalsByOwnerID retrieves all goals for a specific owner
func (r *GoalRepository) GetGoalsByOwnerID(ownerID uuid.UUID) ([]models.Goal, error) {
	var goals []models.Goal
	err := r.db.Where("owner_id = ?", ownerID).
		Order("created_at DESC").
		Find(&goals).Error
	return goals, err
}

// GetGoals retrieves goals with filters
func (r *GoalRepository) GetGoals(status *models.GoalStatus, limit, offset int) ([]models.Goal, int64, error) {
	var goals []models.Goal
	var total int64

	query := r.db.Model(&models.Goal{})

	if status != nil {
		query = query.Where("status = ?", *status)
	}

	// Get total count
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Get paginated results
	err := query.Limit(limit).Offset(offset).
		Order("created_at DESC").
		Find(&goals).Error

	return goals, total, err
}

// UpdateGoal updates a goal
func (r *GoalRepository) UpdateGoal(goal *models.Goal) error {
	return r.db.Save(goal).Error
}

// DeleteGoal deletes a goal
func (r *GoalRepository) DeleteGoal(id uuid.UUID) error {
	return r.db.Delete(&models.Goal{}, "id = ?", id).Error
}

// GetTotalConfirmedContributions calculates total confirmed contributions for a goal
func (r *GoalRepository) GetTotalConfirmedContributions(goalID uuid.UUID) (int64, error) {
	var total int64
	err := r.db.Model(&models.Contribution{}).
		Where("goal_id = ? AND status = ?", goalID, models.ContributionStatusConfirmed).
		Select("COALESCE(SUM(amount), 0)").
		Scan(&total).Error
	return total, err
}

// GetTotalCompletedWithdrawals calculates total completed withdrawals for a goal
func (r *GoalRepository) GetTotalCompletedWithdrawals(goalID uuid.UUID) (int64, error) {
	var total int64
	err := r.db.Model(&models.Withdrawal{}).
		Where("goal_id = ? AND status = ?", goalID, models.WithdrawalStatusCompleted).
		Select("COALESCE(SUM(amount), 0)").
		Scan(&total).Error
	return total, err
}

// GetContributorCount returns the number of unique contributors for a goal
func (r *GoalRepository) GetContributorCount(goalID uuid.UUID) (int64, error) {
	var count int64
	err := r.db.Model(&models.Contribution{}).
		Where("goal_id = ? AND status = ?", goalID, models.ContributionStatusConfirmed).
		Distinct("user_id").
		Count(&count).Error
	return count, err
}

// IsUserContributor checks if a user has contributed to a goal
func (r *GoalRepository) IsUserContributor(goalID, userID uuid.UUID) (bool, error) {
	var count int64
	err := r.db.Model(&models.Contribution{}).
		Where("goal_id = ? AND user_id = ? AND status = ?", goalID, userID, models.ContributionStatusConfirmed).
		Count(&count).Error
	return count > 0, err
}

// MilestoneRepository handles database operations for milestones
type MilestoneRepository struct {
	db *gorm.DB
}

// NewMilestoneRepository creates a new milestone repository
func NewMilestoneRepository(db *gorm.DB) *MilestoneRepository {
	return &MilestoneRepository{db: db}
}

// CreateMilestone creates a new milestone
func (r *MilestoneRepository) CreateMilestone(milestone *models.Milestone) error {
	return r.db.Create(milestone).Error
}

// GetMilestoneByID retrieves a milestone by ID
func (r *MilestoneRepository) GetMilestoneByID(id uuid.UUID) (*models.Milestone, error) {
	var milestone models.Milestone
	err := r.db.First(&milestone, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &milestone, nil
}

// GetMilestonesByGoalID retrieves all milestones for a goal
func (r *MilestoneRepository) GetMilestonesByGoalID(goalID uuid.UUID) ([]models.Milestone, error) {
	var milestones []models.Milestone
	err := r.db.Where("goal_id = ?", goalID).
		Order("order_index ASC").
		Find(&milestones).Error
	return milestones, err
}

// UpdateMilestone updates a milestone
func (r *MilestoneRepository) UpdateMilestone(milestone *models.Milestone) error {
	return r.db.Save(milestone).Error
}

// DeleteMilestone deletes a milestone
func (r *MilestoneRepository) DeleteMilestone(id uuid.UUID) error {
	return r.db.Delete(&models.Milestone{}, "id = ?", id).Error
}

// GetTotalConfirmedContributionsByMilestone calculates total confirmed contributions for a milestone
func (r *MilestoneRepository) GetTotalConfirmedContributionsByMilestone(milestoneID uuid.UUID) (int64, error) {
	var total int64
	err := r.db.Model(&models.Contribution{}).
		Where("milestone_id = ? AND status = ?", milestoneID, models.ContributionStatusConfirmed).
		Select("COALESCE(SUM(amount), 0)").
		Scan(&total).Error
	return total, err
}

// GetNextOrderIndex gets the next available order index for a goal's milestones
func (r *MilestoneRepository) GetNextOrderIndex(goalID uuid.UUID) (int, error) {
	var maxOrder int
	err := r.db.Model(&models.Milestone{}).
		Where("goal_id = ?", goalID).
		Select("COALESCE(MAX(order_index), 0)").
		Scan(&maxOrder).Error
	return maxOrder + 1, err
}

// ContributionRepository handles database operations for contributions
type ContributionRepository struct {
	db *gorm.DB
}

// NewContributionRepository creates a new contribution repository
func NewContributionRepository(db *gorm.DB) *ContributionRepository {
	return &ContributionRepository{db: db}
}

// CreateContribution creates a new contribution
func (r *ContributionRepository) CreateContribution(contribution *models.Contribution) error {
	return r.db.Create(contribution).Error
}

// GetContributionByID retrieves a contribution by ID
func (r *ContributionRepository) GetContributionByID(id uuid.UUID) (*models.Contribution, error) {
	var contribution models.Contribution
	err := r.db.First(&contribution, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &contribution, nil
}

// GetContributionByPaymentID retrieves a contribution by payment ID
func (r *ContributionRepository) GetContributionByPaymentID(paymentID uuid.UUID) (*models.Contribution, error) {
	var contribution models.Contribution
	err := r.db.First(&contribution, "payment_id = ?", paymentID).Error
	if err != nil {
		return nil, err
	}
	return &contribution, nil
}

// GetContributionsByGoalID retrieves all contributions for a goal
func (r *ContributionRepository) GetContributionsByGoalID(goalID uuid.UUID) ([]models.Contribution, error) {
	var contributions []models.Contribution
	err := r.db.Where("goal_id = ?", goalID).
		Order("created_at DESC").
		Find(&contributions).Error
	return contributions, err
}

// GetContributionsByUserID retrieves all contributions by a user
func (r *ContributionRepository) GetContributionsByUserID(userID uuid.UUID) ([]models.Contribution, error) {
	var contributions []models.Contribution
	err := r.db.Where("user_id = ?", userID).
		Order("created_at DESC").
		Find(&contributions).Error
	return contributions, err
}

// UpdateContribution updates a contribution
func (r *ContributionRepository) UpdateContribution(contribution *models.Contribution) error {
	return r.db.Save(contribution).Error
}

// WithdrawalRepository handles database operations for withdrawals
type WithdrawalRepository struct {
	db *gorm.DB
}

// NewWithdrawalRepository creates a new withdrawal repository
func NewWithdrawalRepository(db *gorm.DB) *WithdrawalRepository {
	return &WithdrawalRepository{db: db}
}

// CreateWithdrawal creates a new withdrawal
func (r *WithdrawalRepository) CreateWithdrawal(withdrawal *models.Withdrawal) error {
	return r.db.Create(withdrawal).Error
}

// GetWithdrawalByID retrieves a withdrawal by ID
func (r *WithdrawalRepository) GetWithdrawalByID(id uuid.UUID) (*models.Withdrawal, error) {
	var withdrawal models.Withdrawal
	err := r.db.First(&withdrawal, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &withdrawal, nil
}

// GetWithdrawalsByGoalID retrieves all withdrawals for a goal
func (r *WithdrawalRepository) GetWithdrawalsByGoalID(goalID uuid.UUID) ([]models.Withdrawal, error) {
	var withdrawals []models.Withdrawal
	err := r.db.Where("goal_id = ?", goalID).
		Order("requested_at DESC").
		Find(&withdrawals).Error
	return withdrawals, err
}

// UpdateWithdrawal updates a withdrawal
func (r *WithdrawalRepository) UpdateWithdrawal(withdrawal *models.Withdrawal) error {
	return r.db.Save(withdrawal).Error
}

// ProofRepository handles database operations for proofs
type ProofRepository struct {
	db *gorm.DB
}

// NewProofRepository creates a new proof repository
func NewProofRepository(db *gorm.DB) *ProofRepository {
	return &ProofRepository{db: db}
}

// CreateProof creates a new proof
func (r *ProofRepository) CreateProof(proof *models.Proof) error {
	return r.db.Create(proof).Error
}

// GetProofByID retrieves a proof by ID with votes
func (r *ProofRepository) GetProofByID(id uuid.UUID) (*models.Proof, error) {
	var proof models.Proof
	err := r.db.Preload("Votes").First(&proof, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &proof, nil
}

// GetProofsByGoalID retrieves all proofs for a goal
func (r *ProofRepository) GetProofsByGoalID(goalID uuid.UUID) ([]models.Proof, error) {
	var proofs []models.Proof
	err := r.db.Where("goal_id = ?", goalID).
		Order("submitted_at DESC").
		Find(&proofs).Error
	return proofs, err
}

// UpdateProof updates a proof
func (r *ProofRepository) UpdateProof(proof *models.Proof) error {
	return r.db.Save(proof).Error
}

// DeleteProof deletes a proof
func (r *ProofRepository) DeleteProof(id uuid.UUID) error {
	return r.db.Delete(&models.Proof{}, "id = ?", id).Error
}

// VoteRepository handles database operations for votes
type VoteRepository struct {
	db *gorm.DB
}

// NewVoteRepository creates a new vote repository
func NewVoteRepository(db *gorm.DB) *VoteRepository {
	return &VoteRepository{db: db}
}

// CreateVote creates a new vote
func (r *VoteRepository) CreateVote(vote *models.Vote) error {
	return r.db.Create(vote).Error
}

// GetVoteByProofAndVoter retrieves a vote by proof ID and voter ID
func (r *VoteRepository) GetVoteByProofAndVoter(proofID, voterID uuid.UUID) (*models.Vote, error) {
	var vote models.Vote
	err := r.db.First(&vote, "proof_id = ? AND voter_id = ?", proofID, voterID).Error
	if err != nil {
		return nil, err
	}
	return &vote, nil
}

// GetVotesByProofID retrieves all votes for a proof
func (r *VoteRepository) GetVotesByProofID(proofID uuid.UUID) ([]models.Vote, error) {
	var votes []models.Vote
	err := r.db.Where("proof_id = ?", proofID).
		Order("voted_at DESC").
		Find(&votes).Error
	return votes, err
}

// UpdateVote updates a vote
func (r *VoteRepository) UpdateVote(vote *models.Vote) error {
	return r.db.Save(vote).Error
}

// GetVoteStats returns vote statistics for a proof
func (r *VoteRepository) GetVoteStats(proofID uuid.UUID) (total, satisfied int64, err error) {
	err = r.db.Model(&models.Vote{}).
		Where("proof_id = ?", proofID).
		Count(&total).Error
	if err != nil {
		return 0, 0, err
	}

	err = r.db.Model(&models.Vote{}).
		Where("proof_id = ? AND is_satisfied = ?", proofID, true).
		Count(&satisfied).Error

	return total, satisfied, err
}

// DeleteVote deletes a vote
func (r *VoteRepository) DeleteVote(id uuid.UUID) error {
	return r.db.Delete(&models.Vote{}, "id = ?", id).Error
}

// Repository aggregates all repositories
type Repository struct {
	Goal         *GoalRepository
	Milestone    *MilestoneRepository
	Contribution *ContributionRepository
	Withdrawal   *WithdrawalRepository
	Proof        *ProofRepository
	Vote         *VoteRepository
}

// NewRepository creates a new repository instance
func NewRepository(db *gorm.DB) *Repository {
	return &Repository{
		Goal:         NewGoalRepository(db),
		Milestone:    NewMilestoneRepository(db),
		Contribution: NewContributionRepository(db),
		Withdrawal:   NewWithdrawalRepository(db),
		Proof:        NewProofRepository(db),
		Vote:         NewVoteRepository(db),
	}
}
