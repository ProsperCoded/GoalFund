package middleware

import (
	"crypto/hmac"
	"crypto/sha512"
	"encoding/hex"
	"io"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gofund/shared/metrics"
)

// WebhookAuthMiddleware verifies Paystack webhook signatures
func WebhookAuthMiddleware(webhookSecret string) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get signature from header
		signature := c.GetHeader("x-paystack-signature")
		if signature == "" {
			log.Printf("[ERROR] Missing webhook signature")
			metrics.IncrementCounter("webhook.signature.missing")
			c.JSON(http.StatusUnauthorized, gin.H{
				"status":  "error",
				"message": "Missing webhook signature",
			})
			c.Abort()
			return
		}

		// Read request body
		body, err := io.ReadAll(c.Request.Body)
		if err != nil {
			log.Printf("[ERROR] Failed to read webhook body: %v", err)
			c.JSON(http.StatusBadRequest, gin.H{
				"status":  "error",
				"message": "Failed to read request body",
			})
			c.Abort()
			return
		}

		// Verify signature
		if !verifySignature(body, signature, webhookSecret) {
			log.Printf("[ERROR] Invalid webhook signature: %s", signature)
			metrics.IncrementCounter("webhook.signature.invalid")
			c.JSON(http.StatusUnauthorized, gin.H{
				"status":  "error",
				"message": "Invalid webhook signature",
			})
			c.Abort()
			return
		}

		// Store body in context for controller to use
		c.Set("webhook_body", body)

		metrics.IncrementCounter("webhook.signature.valid")
		c.Next()
	}
}

// verifySignature verifies the webhook signature using HMAC SHA512
func verifySignature(payload []byte, signature, secret string) bool {
	mac := hmac.New(sha512.New, []byte(secret))
	mac.Write(payload)
	expectedSignature := hex.EncodeToString(mac.Sum(nil))

	return hmac.Equal([]byte(signature), []byte(expectedSignature))
}
