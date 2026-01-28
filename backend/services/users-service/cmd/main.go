package main

import (
	"log"
	"os"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gofund/shared/database"
	"github.com/gofund/shared/jwt"
	"github.com/gofund/users-service/internal/repository"
	"github.com/gofund/users-service/internal/router"
	"github.com/gofund/users-service/internal/service"
	"gorm.io/gorm/logger"
)

func main() {
	// Initialize database
	dbConfig := database.Config{
		Host:     getEnv("USERS_DB_HOST", "localhost"),
		Port:     getEnvInt("USERS_DB_PORT", 5432),
		User:     getEnv("USERS_DB_USER", "postgres"),
		Password: getEnv("USERS_DB_PASSWORD", "postgres"),
		DBName:   getEnv("USERS_DB_NAME", "users_db"),
		SSLMode:  getEnv("USERS_DB_SSLMODE", "disable"),
		LogLevel: logger.Info,
	}

	db, err := database.SetupDatabase(dbConfig)
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	// Initialize repositories
	userRepo := repository.NewUserRepository(db)
	sessionRepo := repository.NewSessionRepository(db)

	// Initialize JWT service
	jwtSecret := getEnv("JWT_SECRET", "your-super-secret-jwt-key-change-in-production")
	jwtService := jwt.NewJWTService(jwtSecret, time.Hour, 30*24*time.Hour) // 1 hour access, 30 days refresh

	// Initialize services
	authService := service.NewAuthService(userRepo, sessionRepo, jwtService)

	// Initialize Gin router
	r := gin.Default()

	// Add middleware for logging and recovery
	r.Use(gin.Logger())
	r.Use(gin.Recovery())

	// Setup all routes
	router.SetupRoutes(r, authService)

	// Start server
	port := getEnv("PORT", "8084")

	log.Printf("Users Service starting on port %s", port)
	log.Printf("Internal auth endpoint: /internal/verify")
	log.Printf("Public auth endpoints: /auth/*")
	log.Printf("Protected user endpoints: /users/*")
	log.Printf("Database: %s:%d/%s", dbConfig.Host, dbConfig.Port, dbConfig.DBName)
	
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

// getEnvInt gets environment variable as integer with fallback
func getEnvInt(key string, fallback int) int {
	if value := os.Getenv(key); value != "" {
		if intValue, err := strconv.Atoi(value); err == nil {
			return intValue
		}
	}
	return fallback
}
