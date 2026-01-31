package main

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/gofund/notifications-service/internal/config"
	"github.com/gofund/notifications-service/internal/handlers"
	"github.com/gofund/notifications-service/internal/repository"
	"github.com/gofund/notifications-service/internal/service"
	"github.com/gofund/shared/database"
	"github.com/gofund/shared/messaging"
	"github.com/gofund/shared/metrics"
	"github.com/joho/godotenv"
	gintrace "gopkg.in/DataDog/dd-trace-go.v1/contrib/gin-gonic/gin"
)

func main() {
	// Load .env file if it exists (for local development)
	if err := godotenv.Load(); err != nil {
		log.Printf("No .env file found, using system environment variables")
	}

	// Load configuration
	cfg := config.LoadConfig()

	// Initialize Datadog tracing and metrics
	if err := metrics.InitDatadog(cfg.DDService, cfg.DDEnv, cfg.DDVersion); err != nil {
		log.Printf("Warning: Failed to initialize Datadog: %v", err)
	} else {
		log.Printf("Datadog initialized successfully for %s", cfg.DDService)
	}
	defer metrics.StopDatadog()

	// Initialize database connection
	db, err := database.NewPostgresConnection(cfg.GetDSN())
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()
	log.Println("Database connection established")

	// Initialize repositories
	notificationRepo := repository.NewNotificationRepository(db)
	preferenceRepo := repository.NewPreferenceRepository(db)

	// Initialize email service
	emailService, err := service.NewEmailService(cfg)
	if err != nil {
		log.Printf("Warning: Failed to initialize email service: %v", err)
	}

	// Initialize notification service
	notificationService := service.NewNotificationService(
		notificationRepo,
		preferenceRepo,
		emailService,
	)

	// Initialize event handler
	eventHandler := handlers.NewEventHandler(notificationService)

	// Initialize RabbitMQ connection
	rabbitConn, err := messaging.NewRabbitMQConnection(cfg.RabbitMQURL)
	if err != nil {
		log.Fatalf("Failed to connect to RabbitMQ: %v", err)
	}
	defer rabbitConn.Close()
	log.Println("RabbitMQ connection established")

	// Initialize RabbitMQ consumer
	consumer, err := messaging.NewRabbitMQConsumer(rabbitConn, cfg.RabbitMQExchange, cfg.RabbitMQQueue)
	if err != nil {
		log.Fatalf("Failed to create RabbitMQ consumer: %v", err)
	}

	// Start consuming events
	go startEventConsumers(consumer, eventHandler)

	// Initialize HTTP router
	r := gin.Default()

	// Add Datadog APM middleware
	r.Use(gintrace.Middleware(cfg.DDService))

	// Setup HTTP routes
	setupRoutes(r, notificationService)

	// Start server
	log.Printf("Notifications Service starting on port %s", cfg.Port)
	if err := r.Run(":" + cfg.Port); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}

// startEventConsumers starts consuming events from RabbitMQ
func startEventConsumers(consumer *messaging.RabbitMQConsumer, eventHandler *handlers.EventHandler) {
	log.Println("Starting event consumers...")

	// Payment events
	if err := consumer.Consume("PaymentVerified", eventHandler.HandlePaymentVerified); err != nil {
		log.Printf("Failed to consume PaymentVerified events: %v", err)
	}

	// Contribution events
	if err := consumer.Consume("ContributionConfirmed", eventHandler.HandleContributionConfirmed); err != nil {
		log.Printf("Failed to consume ContributionConfirmed events: %v", err)
	}

	// Withdrawal events
	if err := consumer.Consume("WithdrawalRequested", eventHandler.HandleWithdrawalRequested); err != nil {
		log.Printf("Failed to consume WithdrawalRequested events: %v", err)
	}

	if err := consumer.Consume("WithdrawalCompleted", eventHandler.HandleWithdrawalCompleted); err != nil {
		log.Printf("Failed to consume WithdrawalCompleted events: %v", err)
	}

	// Proof events
	if err := consumer.Consume("ProofSubmitted", eventHandler.HandleProofSubmitted); err != nil {
		log.Printf("Failed to consume ProofSubmitted events: %v", err)
	}

	if err := consumer.Consume("ProofVoted", eventHandler.HandleProofVoted); err != nil {
		log.Printf("Failed to consume ProofVoted events: %v", err)
	}

	// Goal events
	if err := consumer.Consume("GoalFunded", eventHandler.HandleGoalFunded); err != nil {
		log.Printf("Failed to consume GoalFunded events: %v", err)
	}

	// User events
	if err := consumer.Consume("UserSignedUp", eventHandler.HandleUserSignedUp); err != nil {
		log.Printf("Failed to consume UserSignedUp events: %v", err)
	}

	if err := consumer.Consume("PasswordResetRequested", eventHandler.HandlePasswordResetRequested); err != nil {
		log.Printf("Failed to consume PasswordResetRequested events: %v", err)
	}

	if err := consumer.Consume("EmailVerificationRequested", eventHandler.HandleEmailVerificationRequested); err != nil {
		log.Printf("Failed to consume EmailVerificationRequested events: %v", err)
	}

	if err := consumer.Consume("KYCVerified", eventHandler.HandleKYCVerified); err != nil {
		log.Printf("Failed to consume KYCVerified events: %v", err)
	}

	log.Println("Event consumers started successfully")
}

// setupRoutes configures all HTTP routes
func setupRoutes(r *gin.Engine, notificationService service.NotificationService) {
	// Initialize HTTP handler
	notificationHandler := handlers.NewNotificationHandler(notificationService)

	// Health check
	r.GET("/api/v1/notifications/health", notificationHandler.HealthCheck)

	// API routes (require authentication in production)
	api := r.Group("/api/v1/notifications")
	{
		// Notification endpoints
		api.GET("", notificationHandler.GetNotifications)
		api.GET("/:id", notificationHandler.GetNotification)
		api.PUT("/:id/read", notificationHandler.MarkAsRead)
		api.DELETE("/:id", notificationHandler.DeleteNotification)
		api.GET("/unread/count", notificationHandler.GetUnreadCount)

		// Preference endpoints
		api.GET("/preferences", notificationHandler.GetPreferences)
		api.PUT("/preferences", notificationHandler.UpdatePreferences)
	}

	log.Println("HTTP routes configured")
}
