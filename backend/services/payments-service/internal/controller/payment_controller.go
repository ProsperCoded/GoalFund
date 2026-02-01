package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gofund/payments-service/internal/dto"
	"github.com/gofund/payments-service/internal/service"
	"github.com/gofund/shared/logger"
)

// PaymentController handles payment-related HTTP requests
type PaymentController struct {
	paymentService *service.PaymentService
}

// NewPaymentController creates a new payment controller
func NewPaymentController(paymentService *service.PaymentService) *PaymentController {
	return &PaymentController{
		paymentService: paymentService,
	}
}

// InitializePayment handles POST /api/v1/payments/initialize
func (pc *PaymentController) InitializePayment(c *gin.Context) {
	var req dto.InitializePaymentRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		logger.Error("Invalid payment initialization request", map[string]interface{}{
			"error": err.Error(),
		})
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "Invalid request body",
			"error":   err.Error(),
		})
		return
	}

	// Initialize payment
	resp, err := pc.paymentService.InitializePayment(c.Request.Context(), &req)
	if err != nil {
		logger.Error("Failed to initialize payment", map[string]interface{}{
			"error":   err.Error(),
			"user_id": req.UserID.String(),
			"goal_id": req.GoalID.String(),
			"amount":  req.Amount,
		})
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "error",
			"message": "Failed to initialize payment",
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status": "success",
		"data":   resp,
	})
}

// VerifyPayment handles GET /api/v1/payments/verify/:reference
func (pc *PaymentController) VerifyPayment(c *gin.Context) {
	reference := c.Param("reference")
	if reference == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "Reference is required",
		})
		return
	}

	// Verify payment
	resp, err := pc.paymentService.VerifyPayment(c.Request.Context(), reference)
	if err != nil {
		logger.Error("Failed to verify payment", map[string]interface{}{
			"error":     err.Error(),
			"reference": reference,
		})
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "error",
			"message": "Failed to verify payment",
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status": "success",
		"data":   resp,
	})
}

// GetPaymentStatus handles GET /api/v1/payments/:paymentId/status
func (pc *PaymentController) GetPaymentStatus(c *gin.Context) {
	paymentID := c.Param("paymentId")
	if paymentID == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "Payment ID is required",
		})
		return
	}

	// Get payment status
	resp, err := pc.paymentService.GetPaymentStatus(c.Request.Context(), paymentID)
	if err != nil {
		logger.Error("Failed to get payment status", map[string]interface{}{
			"error":      err.Error(),
			"payment_id": paymentID,
		})
		c.JSON(http.StatusNotFound, gin.H{
			"status":  "error",
			"message": "Payment not found",
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status": "success",
		"data":   resp,
	})
}

// ListBanks handles GET /api/v1/payments/banks
func (pc *PaymentController) ListBanks(c *gin.Context) {
	country := c.DefaultQuery("country", "nigeria")

	// Get bank list
	banks, err := pc.paymentService.ListBanks(c.Request.Context(), country)
	if err != nil {
		logger.Error("Failed to list banks", map[string]interface{}{
			"error":   err.Error(),
			"country": country,
		})
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "error",
			"message": "Failed to list banks",
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status": "success",
		"data":   banks,
	})
}

// ResolveAccount handles GET /api/v1/payments/resolve-account
func (pc *PaymentController) ResolveAccount(c *gin.Context) {
	accountNumber := c.Query("account_number")
	bankCode := c.Query("bank_code")

	if accountNumber == "" || bankCode == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "account_number and bank_code are required",
		})
		return
	}

	req := &dto.ResolveAccountRequest{
		AccountNumber: accountNumber,
		BankCode:      bankCode,
	}

	// Resolve account
	resp, err := pc.paymentService.ResolveAccount(c.Request.Context(), req)
	if err != nil {
		logger.Error("Failed to resolve account", map[string]interface{}{
			"error":          err.Error(),
			"account_number": accountNumber,
			"bank_code":      bankCode,
		})
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "error",
			"message": "Failed to resolve account",
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status": "success",
		"data":   resp,
	})
}
