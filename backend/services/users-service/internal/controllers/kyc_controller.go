package controllers

import (
	"net/http"

	"github.com/gofund/users-service/internal/service"
	"github.com/gin-gonic/gin"
)

// KYCController handles KYC verification endpoints
type KYCController struct {
	kycService *service.KYCService
}

// NewKYCController creates a new KYC controller instance
func NewKYCController(kycService *service.KYCService) *KYCController {
	return &KYCController{
		kycService: kycService,
	}
}

// SubmitNIN handles NIN submission for KYC verification
// @Summary Submit NIN for KYC verification
// @Description Submit National Identification Number for basic KYC verification (auto-approved in this dummy implementation)
// @Tags KYC
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body service.SubmitNINRequest true "NIN submission request"
// @Success 200 {object} service.KYCStatusResponse
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /api/v1/users/kyc/submit-nin [post]
func (c *KYCController) SubmitNIN(ctx *gin.Context) {
	// Get user ID from context (set by auth middleware)
	userID, exists := ctx.Get("user_id")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, ErrorResponse{
			Error: "unauthorized",
		})
		return
	}

	var req service.SubmitNINRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, ErrorResponse{
			Error: err.Error(),
		})
		return
	}

	response, err := c.kycService.SubmitNIN(userID.(string), &req)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, ErrorResponse{
			Error: err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, response)
}

// GetKYCStatus retrieves the KYC verification status
// @Summary Get KYC verification status
// @Description Retrieve the current KYC verification status for the authenticated user
// @Tags KYC
// @Produce json
// @Security BearerAuth
// @Success 200 {object} service.KYCStatusResponse
// @Failure 401 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /api/v1/users/kyc/status [get]
func (c *KYCController) GetKYCStatus(ctx *gin.Context) {
	// Get user ID from context (set by auth middleware)
	userID, exists := ctx.Get("user_id")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, ErrorResponse{
			Error: "unauthorized",
		})
		return
	}

	response, err := c.kycService.GetKYCStatus(userID.(string))
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, ErrorResponse{
			Error: err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, response)
}

// ErrorResponse represents an error response
type ErrorResponse struct {
	Error string `json:"error"`
}
