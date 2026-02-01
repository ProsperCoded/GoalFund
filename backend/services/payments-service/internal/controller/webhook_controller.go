package controller

import (
	"encoding/json"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gofund/payments-service/internal/dto"
	"github.com/gofund/payments-service/internal/service"
	"github.com/gofund/shared/logger"
)

// WebhookController handles webhook-related HTTP requests
type WebhookController struct {
	webhookService *service.WebhookService
}

// NewWebhookController creates a new webhook controller
func NewWebhookController(webhookService *service.WebhookService) *WebhookController {
	return &WebhookController{
		webhookService: webhookService,
	}
}

// HandleWebhook handles POST /api/v1/payments/webhook
func (wc *WebhookController) HandleWebhook(c *gin.Context) {
	// Get webhook body from context (set by middleware)
	bodyInterface, exists := c.Get("webhook_body")
	if !exists {
		logger.Error("Webhook body not found in context", nil)
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "Invalid webhook request",
		})
		return
	}

	body, ok := bodyInterface.([]byte)
	if !ok {
		logger.Error("Invalid webhook body type", nil)
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "Invalid webhook request",
		})
		return
	}

	// Get signature from header
	signature := c.GetHeader("x-paystack-signature")

	// Parse webhook payload
	var payload dto.WebhookPayload
	if err := json.Unmarshal(body, &payload); err != nil {
		logger.Error("Failed to parse webhook payload", map[string]interface{}{
			"error": err.Error(),
		})
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "Invalid webhook payload",
		})
		return
	}

	logger.Info("Received webhook", map[string]interface{}{
		"event": payload.Event,
	})

	// Process webhook asynchronously (return 200 immediately)
	go func() {
		if err := wc.webhookService.ProcessWebhook(c.Request.Context(), &payload, signature); err != nil {
			logger.Error("Failed to process webhook", map[string]interface{}{
				"error": err.Error(),
				"event": payload.Event,
			})
		}
	}()

	// Return 200 OK immediately to Paystack
	c.JSON(http.StatusOK, gin.H{
		"status": "received",
	})
}
