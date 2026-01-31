package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gofund/users-service/internal/dto"
	"github.com/gofund/users-service/internal/service"
)

// UserController handles user-related endpoints
type UserController struct {
	authService *service.AuthService
	userService *service.UserService
}

// NewUserController creates a new user controller instance
func NewUserController(authService *service.AuthService, userService *service.UserService) *UserController {
	return &UserController{
		authService: authService,
		userService: userService,
	}
}


// CreateLightweightUser handles email-only user creation for contributions
// This is a public endpoint that doesn't require authentication
func (uc *UserController) CreateLightweightUser(c *gin.Context) {
	var req dto.RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid request format",
			"details": err.Error(),
		})
		return
	}

	// Register user (lightweight - no password)
	response, err := uc.authService.Register(&req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"user":    response.User,
		"message": "User account created/updated successfully",
	})
}


// SetPassword handles first-time password setup for lightweight users
// This is a public endpoint
func (uc *UserController) SetPassword(c *gin.Context) {
	var req dto.SetPasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid request format",
			"details": err.Error(),
		})
		return
	}

	// Set password and return auth tokens
	response, err := uc.authService.SetPassword(&req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, response)
}

// UpdateSettlementAccount handles updating user's settlement account details
// Requires authentication
func (uc *UserController) UpdateSettlementAccount(c *gin.Context) {
	// Extract user ID from header (set by Nginx after auth verification)
	userID := c.GetHeader("X-User-ID")
	if userID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "User not authenticated",
		})
		return
	}

	var req struct {
		BankName      string `json:"bank_name" binding:"required"`
		AccountNumber string `json:"account_number" binding:"required"`
		AccountName   string `json:"account_name" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid request format",
			"details": err.Error(),
		})
		return
	}

	// Update settlement account
	if err := uc.userService.UpdateSettlementAccount(userID, req.BankName, req.AccountNumber, req.AccountName); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}


	c.JSON(http.StatusOK, gin.H{
		"message": "Settlement account updated successfully",
	})
}
