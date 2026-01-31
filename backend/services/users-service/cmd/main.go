package main

import (
	"log"
	"os"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gofund/shared/database"
	"github.com/gofund/shared/jwt"
	"github.com/gofund/shared/messaging"
	"github.com/gofund/shared/metrics"
	"github.com/gofund/users-service/internal/repository"
	"github.com/gofund/users-service/internal/router"
	"github.com/gofund/users-service/internal/service"
	"github.com/joho/godotenv"
	gintrace "gopkg.in/DataDog/dd-trace-go.v1/contrib/gin-gonic/gin"
	"gorm.io/gorm/logger"
)

func main() {
	// Load .env file if it exists (for local development)
	if err := godotenv.Load(); err != nil {
		log.Printf("No .env file found, using system environment variables")
	}

	// Initialize Datadog tracing and metrics
	serviceName := getEnv("DD_SERVICE", "users-service")
	env := getEnv("DD_ENV", "dev")
	version := getEnv("DD_VERSION", "1.0.0")
	
	if err := metrics.InitDatadog(serviceName, env, version); err != nil {
		log.Printf("Warning: Failed to initialize Datadog: %v", err)
	} else {
		log.Printf("Datadog initialized successfully for %s", serviceName)
	}
	defer metrics.StopDatadog()

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

	// Initialize RabbitMQ connection
	rabbitmqURL := getEnv("RABBITMQ_URL", "amqp://guest:guest@localhost:5672/")
	rabbitConn, err := messaging.NewRabbitMQConnection(rabbitmqURL)
	if err != nil {
		log.Printf("Warning: Failed to connect to RabbitMQ: %v", err)
		log.Printf("Continuing without messaging capabilities...")
	}
	defer func() {
		if rabbitConn != nil {
			rabbitConn.Close()
		}
	}()

	// Initialize messaging services
	var publisher messaging.Publisher
	var consumer messaging.Consumer
	if rabbitConn != nil {
		exchangeName := getEnv("RABBITMQ_EXCHANGE", "gofund.events")
		queueName := getEnv("RABBITMQ_QUEUE", "users.notifications")

		publisher, err = messaging.NewRabbitMQPublisher(rabbitConn, exchangeName)
		if err != nil {
			log.Printf("Warning: Failed to create publisher: %v", err)
		}

		consumer, err = messaging.NewRabbitMQConsumer(rabbitConn, exchangeName, queueName)
		if err != nil {
			log.Printf("Warning: Failed to create consumer: %v", err)
		}
	}

	// Initialize services
	eventService := service.NewEventService(publisher)
	authService := service.NewAuthService(userRepo, sessionRepo, jwtService, eventService)
	kycService := service.NewKYCService(userRepo, eventService)

	// Start notification consumers if available
	if consumer != nil {
		notificationService := service.NewNotificationService(consumer)
		if err := notificationService.StartConsumers(); err != nil {
			log.Printf("Warning: Failed to start notification consumers: %v", err)
		} else {
			log.Printf("Notification consumers started successfully")
		}
	}

	// Initialize Gin router
	r := gin.Default()

	// Add Datadog APM middleware (automatically tracks HTTP requests)
	r.Use(gintrace.Middleware(serviceName))
	
	// Add middleware for logging and recovery
	r.Use(gin.Logger())
	r.Use(gin.Recovery())

	// Setup all routes
	router.SetupRoutes(r, authService, kycService)

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
