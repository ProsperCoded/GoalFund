package main

import (
	"log"
	"os"

	"github.com/gin-gonic/gin"
)

func main() {
	// Initialize router
	r := gin.Default()

	// Setup routes
	setupRoutes(r)

	// Start server
	port := os.Getenv("PORT")
	if port == "" {
		port = "8081"
	}

	log.Printf("Payments Service starting on port %s", port)
	if err := r.Run(":" + port); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}

// setupRoutes configures all Payments Service routes
func setupRoutes(r *gin.Engine) {
	// Routes will be configured here
}
