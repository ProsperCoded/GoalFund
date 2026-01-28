package controllers

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/gofund/users-service/internal/service"
)

// AuthController handles authentication-related endpoints
type AuthController struct {
	authService *service.AuthService
}

// NewAuthController creates a new auth controller instance
func NewAuthController(authService *service.AuthService) *AuthController {
	return &AuthController{
		authService: authService,
	}
}

// VerifyToken handles internal token verification for Nginx auth_request
// This endpoint is called by Nginx for every protected route
func (ac *AuthController) VerifyToken(c *gin.Context) {
	// Extract Authorization header
	authHeader := c.GetHeader("Authorization")
	if authHeader == "" {
		c.Status(http.StatusUnauthorized)
		return
	}

	// Extract Bearer token
	tokenParts := strings.Split(authHeader, " ")
	if len(tokenParts) != 2 || tokenParts[0] != "Bearer" {
		c.Status(http.StatusUnauthorized)
		return
	}

	token := tokenParts[1]
	if token == "" {
		c.Status(http.StatusUnauthorized)
		return
	}

	// Validate access token
	claims, err := ac.authService.ValidateAccessToken(token)
	if err != nil {
		c.Status(http.StatusUnauthorized)
		return
	}

	// Set user context headers for downstream services
	c.Header("X-User-ID", claims.UserID)
	c.Header("X-User-Email", claims.Email)
	c.Header("X-User-Role", strings.Join(claims.Roles, ","))
	
	c.Status(http.StatusOK)
}

// Login handles user authentication and JWT token generation
func (ac *AuthController) Login(c *gin.Context) {
	var req service.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid request format",
			"details": err.Error(),
		})
		return
	}

	// Authenticate user
	response, err := ac.authService.Login(&req)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, response)
}

// Register handles user registration
func (ac *AuthController) Register(c *gin.Context) {
	var req service.RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid request format",
			"details": err.Error(),
		})
		return
	}

	// Register user
	response, err := ac.authService.Register(&req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, response)
}

// RefreshToken handles JWT token refresh
func (ac *AuthController) RefreshToken(c *gin.Context) {
	var req service.RefreshRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid request format",
			"details": err.Error(),
		})
		return
	}

	// Refresh token
	tokenPair, err := ac.authService.RefreshToken(&req)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, tokenPair)
}

// Logout handles user logout (token invalidation)
func (ac *AuthController) Logout(c *gin.Context) {
	var req service.RefreshRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid request format",
			"details": err.Error(),
		})
		return
	}

	// Logout (invalidate refresh token session)
	if err := ac.authService.Logout(req.RefreshToken); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Logged out successfully",
	})
}

// ForgotPassword handles password reset requests
func (ac *AuthController) ForgotPassword(c *gin.Context) {
	var req service.ForgotPasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid request format",
			"details": err.Error(),
		})
		return
	}

	// Initiate password reset
	if err := ac.authService.ForgotPassword(&req); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to process password reset request",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "If the email exists, a password reset link has been sent",
	})
}

// ResetPassword handles password reset with token
func (ac *AuthController) ResetPassword(c *gin.Context) {
	var req service.ResetPasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid request format",
			"details": err.Error(),
		})
		return
	}

	// Reset password
	if err := ac.authService.ResetPassword(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Password reset successfully",
	})
}

// GetProfile returns the current user's profile
func (ac *AuthController) GetProfile(c *gin.Context) {
	// Extract user ID from header (set by Nginx after auth verification)
	userID := c.GetHeader("X-User-ID")
	if userID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "User not authenticated",
		})
		return
	}

	// Get user profile
	profile, err := ac.authService.GetProfile(userID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"user": profile,
	})
}

// UpdateProfile handles user profile updates
func (ac *AuthController) UpdateProfile(c *gin.Context) {
	// Extract user ID from header (set by Nginx after auth verification)
	userID := c.GetHeader("X-User-ID")
	if userID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "User not authenticated",
		})
		return
	}

	var req service.UpdateProfileRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid request format",
			"details": err.Error(),
		})
		return
	}

	// Update user profile
	profile, err := ac.authService.UpdateProfile(userID, &req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"user": profile,
		"message": "Profile updated successfully",
	})
}