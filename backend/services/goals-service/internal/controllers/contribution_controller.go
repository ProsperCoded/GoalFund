package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gofund/goals-service/internal/dto"
	"github.com/gofund/goals-service/internal/service"
	"github.com/google/uuid"
)

// ContributionController handles contribution and withdrawal related endpoints
type ContributionController struct {
	contributionService *service.ContributionService
	withdrawalService   *service.WithdrawalService
	proofService         *service.ProofService
	voteService          *service.VoteService
}

// NewContributionController creates a new contribution controller instance
func NewContributionController(
	contributionService *service.ContributionService,
	withdrawalService *service.WithdrawalService,
	proofService *service.ProofService,
	voteService *service.VoteService,
) *ContributionController {
	return &ContributionController{
		contributionService: contributionService,
		withdrawalService:   withdrawalService,
		proofService:        proofService,
		voteService:         voteService,
	}
}

// CreateContribution handles contribution creation
func (cc *ContributionController) CreateContribution(c *gin.Context) {
	userIDStr := c.GetHeader("X-User-ID")
	userID, _ := uuid.Parse(userIDStr)

	var req dto.CreateContributionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	contribution, err := cc.contributionService.CreateContribution(userID, req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, contribution)
}

// CreateWithdrawal handles withdrawal request creation
func (cc *ContributionController) CreateWithdrawal(c *gin.Context) {
	userIDStr := c.GetHeader("X-User-ID")
	userID, _ := uuid.Parse(userIDStr)

	var req dto.CreateWithdrawalRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	withdrawal, err := cc.withdrawalService.CreateWithdrawal(userID, req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, withdrawal)
}

// CreateProof handles proof submission
func (cc *ContributionController) CreateProof(c *gin.Context) {
	userIDStr := c.GetHeader("X-User-ID")
	userID, _ := uuid.Parse(userIDStr)

	var req dto.CreateProofRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	proof, err := cc.proofService.CreateProof(userID, req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, proof)
}

// CreateVote handles voting on a proof
func (cc *ContributionController) CreateVote(c *gin.Context) {
	userIDStr := c.GetHeader("X-User-ID")
	userID, _ := uuid.Parse(userIDStr)

	var req dto.CreateVoteRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	vote, err := cc.voteService.CreateVote(userID, req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, vote)
}

// GetVoteStats retrieves vote statistics for a proof
func (cc *ContributionController) GetVoteStats(c *gin.Context) {
	proofIDStr := c.Param("proofId")
	proofID, err := uuid.Parse(proofIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid proof ID"})
		return
	}

	stats, err := cc.voteService.GetVoteStats(proofID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, stats)
}

// GetProofs retrieves all proofs for a goal
func (cc *ContributionController) GetProofs(c *gin.Context) {
	goalIDStr := c.Query("goalId")
	goalID, err := uuid.Parse(goalIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid goal ID"})
		return
	}

	proofs, err := cc.proofService.GetProofsByGoal(goalID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, proofs)
}
