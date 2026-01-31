package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gofund/goals-service/internal/dto"
	"github.com/gofund/goals-service/internal/service"
	"github.com/google/uuid"
)

// RefundController handles refund-related endpoints
type RefundController struct {
	refundService *service.RefundService
}

// NewRefundController creates a new refund controller instance
func NewRefundController(refundService *service.RefundService) *RefundController {
	return &RefundController{
		refundService: refundService,
	}
}

// InitiateRefund handles refund initiation by goal owner
func (rc *RefundController) InitiateRefund(c *gin.Context) {
	// Extract user ID from header (set by Nginx after auth verification)
	userIDStr := c.GetHeader("X-User-ID")
	if userIDStr == "" {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "User not authenticated",
		})
		return
	}

	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid user ID",
		})
		return
	}

	var req dto.InitiateRefundRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid request format",
			"details": err.Error(),
		})
		return
	}

	// Initiate refund
	refund, err := rc.refundService.InitiateRefund(userID, &req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"refund":  refund,
		"message": "Refund initiated successfully",
	})
}

// GetRefund retrieves a refund by ID
func (rc *RefundController) GetRefund(c *gin.Context) {
	refundIDStr := c.Param("id")
	refundID, err := uuid.Parse(refundIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid refund ID",
		})
		return
	}

	refund, err := rc.refundService.GetRefund(refundID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"refund": refund,
	})
}

// GetGoalRefunds retrieves all refunds for a goal
func (rc *RefundController) GetGoalRefunds(c *gin.Context) {
	goalIDStr := c.Param("goalId")
	goalID, err := uuid.Parse(goalIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid goal ID",
		})
		return
	}

	refunds, err := rc.refundService.GetGoalRefunds(goalID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"refunds": refunds,
		"count":   len(refunds),
	})
}
