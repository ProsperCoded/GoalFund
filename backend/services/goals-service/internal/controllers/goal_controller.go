package controllers

import (
	"net/http"

	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/gofund/goals-service/internal/dto"
	"github.com/gofund/goals-service/internal/service"
	"github.com/google/uuid"
)

// GoalController handles goal-related endpoints
type GoalController struct {
	goalService *service.GoalService
}

// NewGoalController creates a new goal controller instance
func NewGoalController(goalService *service.GoalService) *GoalController {
	return &GoalController{
		goalService: goalService,
	}
}

// ListPublicGoals handles retrieving public goals with pagination
func (gc *GoalController) ListPublicGoals(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", "10"))

	goals, total, err := gc.goalService.ListPublicGoals(page, pageSize)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Enrich goals with computed fields
	enrichedGoals := make([]gin.H, len(goals))
	for i, goal := range goals {
		progress, _ := gc.goalService.GetGoalProgress(goal.ID)
		currentAmount := int64(0)
		contributorCount := int64(0)
		if progress != nil {
			currentAmount = progress.TotalContributions
			contributorCount = progress.ContributorCount
		}
		
		enrichedGoals[i] = gin.H{
			"id":                        goal.ID,
			"owner_id":                  goal.OwnerID,
			"title":                     goal.Title,
			"description":               goal.Description,
			"target_amount":             goal.TargetAmount,
			"fixed_contribution_amount": goal.FixedContributionAmount,
			"current_amount":            currentAmount,
			"currency":                  goal.Currency,
			"deadline":                  goal.Deadline,
			"status":                    goal.Status,
			"is_public":                 goal.IsPublic,
			"created_at":                goal.CreatedAt,
			"updated_at":                goal.UpdatedAt,
			"contributor_count":         contributorCount,
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"data":  enrichedGoals,
		"total": total,
		"page":  page,
		"size":  pageSize,
	})
}

// CreateGoal handles goal creation
func (gc *GoalController) CreateGoal(c *gin.Context) {
	userIDStr := c.GetHeader("X-User-ID")
	if userIDStr == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	var req dto.CreateGoalRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	goal, err := gc.goalService.CreateGoal(userID, req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, goal)
}

// GetGoal retrieves a goal by ID
func (gc *GoalController) GetGoal(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid goal ID"})
		return
	}

	goal, err := gc.goalService.GetGoal(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	// Get progress info for current_amount and contributor_count
	progress, err := gc.goalService.GetGoalProgress(id)
	if err != nil {
		c.JSON(http.StatusOK, goal)
		return
	}

	// Return goal with computed fields
	c.JSON(http.StatusOK, gin.H{
		"id":                        goal.ID,
		"owner_id":                  goal.OwnerID,
		"title":                     goal.Title,
		"description":               goal.Description,
		"target_amount":             goal.TargetAmount,
		"fixed_contribution_amount": goal.FixedContributionAmount,
		"current_amount":            progress.TotalContributions,
		"currency":                  goal.Currency,
		"deadline":                  goal.Deadline,
		"status":                    goal.Status,
		"is_public":                 goal.IsPublic,
		"deposit_bank_name":         goal.DepositBankName,
		"deposit_account_number":    goal.DepositAccountNumber,
		"deposit_account_name":      goal.DepositAccountName,
		"created_at":                goal.CreatedAt,
		"updated_at":                goal.UpdatedAt,
		"contributor_count":         progress.ContributorCount,
		"milestones":                goal.Milestones,
		"contributions":             goal.Contributions,
		"withdrawals":               goal.Withdrawals,
		"proofs":                    goal.Proofs,
	})
}

// UpdateGoal updates a goal
func (gc *GoalController) UpdateGoal(c *gin.Context) {
	userIDStr := c.GetHeader("X-User-ID")
	userID, _ := uuid.Parse(userIDStr)

	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid goal ID"})
		return
	}

	var req dto.UpdateGoalRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	goal, err := gc.goalService.UpdateGoal(id, userID, req)
	if err != nil {
		status := http.StatusBadRequest
		if err == service.ErrUnauthorized {
			status = http.StatusForbidden
		} else if err == service.ErrGoalNotFound {
			status = http.StatusNotFound
		}
		c.JSON(status, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, goal)
}

// GetGoalProgress returns progress information for a goal
func (gc *GoalController) GetGoalProgress(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid goal ID"})
		return
	}

	progress, err := gc.goalService.GetGoalProgress(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, progress)
}

// CreateMilestone creates a new milestone for a goal
func (gc *GoalController) CreateMilestone(c *gin.Context) {
	userIDStr := c.GetHeader("X-User-ID")
	userID, _ := uuid.Parse(userIDStr)

	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid goal ID"})
		return
	}

	var req dto.CreateMilestoneRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	milestone, err := gc.goalService.CreateMilestone(id, userID, req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, milestone)
}

// CompleteMilestone marks a milestone as completed
func (gc *GoalController) CompleteMilestone(c *gin.Context) {
	userIDStr := c.GetHeader("X-User-ID")
	userID, _ := uuid.Parse(userIDStr)

	milestoneIDStr := c.Param("milestoneId")
	milestoneID, err := uuid.Parse(milestoneIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid milestone ID"})
		return
	}

	milestone, nextMilestone, err := gc.goalService.CompleteMilestone(milestoneID, userID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"completed": milestone,
		"next":      nextMilestone,
	})
}

// GetMyGoals retrieves all goals created by the authenticated user
func (gc *GoalController) GetMyGoals(c *gin.Context) {
	userIDStr := c.GetHeader("X-User-ID")
	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid user ID"})
		return
	}

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))

	goals, total, err := gc.goalService.ListUserGoals(userID, page, pageSize)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Enrich goals with computed fields
	enrichedGoals := make([]gin.H, len(goals))
	for i, goal := range goals {
		progress, _ := gc.goalService.GetGoalProgress(goal.ID)
		currentAmount := int64(0)
		contributorCount := int64(0)
		if progress != nil {
			currentAmount = progress.TotalContributions
			contributorCount = progress.ContributorCount
		}
		
		enrichedGoals[i] = gin.H{
			"id":                        goal.ID,
			"owner_id":                  goal.OwnerID,
			"title":                     goal.Title,
			"description":               goal.Description,
			"target_amount":             goal.TargetAmount,
			"fixed_contribution_amount": goal.FixedContributionAmount,
			"current_amount":            currentAmount,
			"currency":                  goal.Currency,
			"deadline":                  goal.Deadline,
			"status":                    goal.Status,
			"is_public":                 goal.IsPublic,
			"deposit_bank_name":         goal.DepositBankName,
			"deposit_account_number":    goal.DepositAccountNumber,
			"deposit_account_name":      goal.DepositAccountName,
			"created_at":                goal.CreatedAt,
			"updated_at":                goal.UpdatedAt,
			"contributor_count":         contributorCount,
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"goals": enrichedGoals,
		"total": total,
		"page":  page,
		"limit": pageSize,
	})
}

// GetGoalMilestones retrieves all milestones for a goal
func (gc *GoalController) GetGoalMilestones(c *gin.Context) {
	goalIDStr := c.Param("id")
	goalID, err := uuid.Parse(goalIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid goal ID"})
		return
	}

	milestones, err := gc.goalService.GetGoalMilestones(goalID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"milestones": milestones})
}
