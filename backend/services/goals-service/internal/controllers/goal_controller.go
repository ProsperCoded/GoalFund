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

	c.JSON(http.StatusOK, gin.H{
		"data":  goals,
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

	c.JSON(http.StatusOK, goal)
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
