package main

import (
	"context"
	"log"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gofund/payments-service/internal/config"
	"github.com/gofund/payments-service/internal/controller"
	"github.com/gofund/payments-service/internal/middleware"
	"github.com/gofund/payments-service/internal/repository"
	"github.com/gofund/payments-service/internal/service"
	"github.com/gofund/shared/messaging"
	"github.com/gofund/shared/metrics"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	gintrace "gopkg.in/DataDog/dd-trace-go.v1/contrib/gin-gonic/gin"
)

func main() {
	// Load configuration
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	// Initialize Datadog tracing and metrics
	if err := metrics.InitDatadog(cfg.ServiceName, cfg.DatadogEnv, cfg.DatadogVersion); err != nil {
		log.Printf("Warning: Failed to initialize Datadog: %v", err)
	} else {
		log.Printf("Datadog initialized successfully for %s", cfg.ServiceName)
	}
	defer metrics.StopDatadog()

	// Connect to MongoDB
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	mongoClient, err := mongo.Connect(ctx, options.Client().ApplyURI(cfg.MongoDBURI))
	if err != nil {
		log.Fatalf("Failed to connect to MongoDB: %v", err)
	}
	defer mongoClient.Disconnect(context.Background())

	// Ping MongoDB to verify connection
	if err := mongoClient.Ping(ctx, nil); err != nil {
		log.Fatalf("Failed to ping MongoDB: %v", err)
	}
	log.Printf("Connected to MongoDB successfully")

	// Get database
	db := mongoClient.Database(cfg.MongoDBDatabase)

	// Initialize repositories
	paymentRepo := repository.NewPaymentRepository(db)
	webhookRepo := repository.NewWebhookRepository(db)
	idempotencyRepo := repository.NewIdempotencyRepository(db)

	// Ensure indexes
	if err := paymentRepo.EnsureIndexes(context.Background()); err != nil {
		log.Printf("Warning: Failed to create payment indexes: %v", err)
	}
	if err := webhookRepo.EnsureIndexes(context.Background()); err != nil {
		log.Printf("Warning: Failed to create webhook indexes: %v", err)
	}
	if err := idempotencyRepo.EnsureIndexes(context.Background()); err != nil {
		log.Printf("Warning: Failed to create idempotency indexes: %v", err)
	}

	// Initialize RabbitMQ connection
	rabbitConn, err := messaging.NewRabbitMQConnection(cfg.RabbitMQURL)
	if err != nil {
		log.Fatalf("Failed to connect to RabbitMQ: %v", err)
	}
	defer rabbitConn.Close()

	// Initialize RabbitMQ publisher
	eventPublisher, err := messaging.NewRabbitMQPublisher(rabbitConn, cfg.RabbitMQExchange)
	if err != nil {
		log.Fatalf("Failed to initialize event publisher: %v", err)
	}
	log.Printf("Connected to RabbitMQ successfully")

	// Initialize Paystack client
	paystackClient := service.NewPaystackClient(cfg.PaystackSecretKey, cfg.PaystackBaseURL)

	// Initialize services
	paymentService := service.NewPaymentService(
		paymentRepo,
		idempotencyRepo,
		paystackClient,
		eventPublisher,
	)

	webhookService := service.NewWebhookService(
		webhookRepo,
		paymentRepo,
		eventPublisher,
	)

	// Initialize controllers
	paymentController := controller.NewPaymentController(paymentService)
	webhookController := controller.NewWebhookController(webhookService)

	// Initialize router
	r := gin.Default()

	// Add Datadog APM middleware
	r.Use(gintrace.Middleware(cfg.ServiceName))

	// Health check endpoint
	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status":  "healthy",
			"service": cfg.ServiceName,
		})
	})

	// Setup routes
	setupRoutes(r, paymentController, webhookController, cfg)

	// Start server
	log.Printf("Payments Service starting on port %s", cfg.ServicePort)
	if err := r.Run(":" + cfg.ServicePort); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}

// setupRoutes configures all Payments Service routes
func setupRoutes(
	r *gin.Engine,
	paymentController *controller.PaymentController,
	webhookController *controller.WebhookController,
	cfg *config.Config,
) {
	// API v1 routes
	v1 := r.Group("/api/v1/payments")
	{
		// Payment routes
		v1.POST("/initialize", paymentController.InitializePayment)
		v1.GET("/verify/:reference", paymentController.VerifyPayment)
		v1.GET("/:paymentId/status", paymentController.GetPaymentStatus)
		v1.GET("/banks", paymentController.ListBanks)
		v1.GET("/resolve-account", paymentController.ResolveAccount)

		// Webhook route (with signature verification middleware)
		webhookSecret := cfg.PaystackSecretKey
		if cfg.PaystackWebhookSecret != "" {
			webhookSecret = cfg.PaystackWebhookSecret
		}
		v1.POST("/webhook", middleware.WebhookAuthMiddleware(webhookSecret), webhookController.HandleWebhook)
	}

	log.Printf("Routes configured successfully")
}
