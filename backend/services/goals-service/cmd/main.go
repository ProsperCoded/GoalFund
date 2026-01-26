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
		port = "8083"
	}

	log.Printf("Goals Service starting on port %s", port)
	if err := r.Run(":" + port); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}

// setupRoutes configures all Goals Service routes
func setupRoutes(r *gin.Engine) {
	// Routes will be configured here
}
