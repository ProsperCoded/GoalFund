package main

import (
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/gofund/users-service/internal/router"
)

func main() {
	// Initialize Gin router
	r := gin.Default()

	// Add middleware for logging and recovery
	r.Use(gin.Logger())
	r.Use(gin.Recovery())

	// Setup all routes
	router.SetupRoutes(r)

	// Start server
	port := os.Getenv("PORT")
	if port == "" {
		port = "8084"
	}

	log.Printf("Users Service starting on port %s", port)
	log.Printf("Internal auth endpoint: /internal/verify")
	log.Printf("Public auth endpoints: /auth/*")
	log.Printf("Protected user endpoints: /users/*")
	
	if err := r.Run(":" + port); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}
