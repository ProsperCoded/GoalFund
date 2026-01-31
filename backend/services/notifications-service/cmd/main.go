package main

import (
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/gofund/shared/metrics"
	"github.com/joho/godotenv"
	gintrace "gopkg.in/DataDog/dd-trace-go.v1/contrib/gin-gonic/gin"
)

func main() {
	// Load .env file if it exists (for local development)
	if err := godotenv.Load(); err != nil {
		log.Printf("No .env file found, using system environment variables")
	}

	// Initialize Datadog tracing and metrics
	serviceName := getEnv("DD_SERVICE", "notifications-service")
	env := getEnv("DD_ENV", "dev")
	version := getEnv("DD_VERSION", "1.0.0")
	
	if err := metrics.InitDatadog(serviceName, env, version); err != nil {
		log.Printf("Warning: Failed to initialize Datadog: %v", err)
	} else {
		log.Printf("Datadog initialized successfully for %s", serviceName)
	}
	defer metrics.StopDatadog()

	// Initialize router
	r := gin.Default()

	// Add Datadog APM middleware
	r.Use(gintrace.Middleware(serviceName))

	// Setup routes
	setupRoutes(r)

	// Start server
	port := os.Getenv("PORT")
	if port == "" {
		port = "8085"
	}

	log.Printf("Notifications Service starting on port %s", port)
	if err := r.Run(":" + port); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}

// getEnv gets environment variable with fallback
func getEnv(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}

// setupRoutes configures all Notifications Service routes
func setupRoutes(r *gin.Engine) {
	// Routes will be configured here
}
