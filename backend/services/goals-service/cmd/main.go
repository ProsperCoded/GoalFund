package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gofund/goals-service/internal/config"
	"github.com/gofund/goals-service/internal/controllers"
	"github.com/gofund/goals-service/internal/events"
	"github.com/gofund/goals-service/internal/middleware"
	"github.com/gofund/goals-service/internal/repository"
	"github.com/gofund/goals-service/internal/service"
	"github.com/gofund/shared/database"
	"github.com/gofund/shared/messaging"
	"github.com/gofund/shared/metrics"
	"github.com/joho/godotenv"
	gintrace "gopkg.in/DataDog/dd-trace-go.v1/contrib/gin-gonic/gin"
	"gorm.io/gorm/logger"
)

func main() {
	// Load .env file
	if err := godotenv.Load(); err != nil {
		log.Printf("No .env file found, using system environment variables")
	}

	// Load config
	cfg := config.LoadConfig()

	// Initialize Datadog
	if err := metrics.InitDatadog(cfg.Datadog.Service, cfg.Datadog.Env, cfg.Datadog.Version); err != nil {
		log.Printf("Warning: Failed to initialize Datadog: %v", err)
	}
	defer metrics.StopDatadog()

	// Initialize Database
	db, err := database.SetupDatabase(database.Config{
		Host:     cfg.Database.Host,
		Port:     stringToInt(cfg.Database.Port),
		User:     cfg.Database.User,
		Password: cfg.Database.Password,
		DBName:   cfg.Database.DBName,
		SSLMode:  cfg.Database.SSLMode,
		LogLevel: logger.Info,
	})
	if err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}

	// Initialize RabbitMQ
	rabbitConn, err := messaging.NewRabbitMQConnection(cfg.RabbitMQ.URL)
	if err != nil {
		log.Printf("Warning: Failed to connect to RabbitMQ: %v", err)
	} else {
		defer rabbitConn.Close()
	}

	// Initialize Messaging
	var publisher messaging.Publisher
	if rabbitConn != nil {
		var err error
		publisher, err = messaging.NewRabbitMQPublisher(rabbitConn, cfg.RabbitMQ.Exchange)
		if err != nil {
			log.Printf("Warning: Failed to create RabbitMQ publisher: %v", err)
		}
	}

	// Initialize Repositories
	repo := repository.NewRepository(db)

	// Initialize Services
	goalService := service.NewGoalService(repo)
	contributionService := service.NewContributionService(repo)
	withdrawalService := service.NewWithdrawalService(repo)
	proofService := service.NewProofService(repo, publisher)
	voteService := service.NewVoteService(repo, publisher)
	refundService := service.NewRefundService(db, publisher)

	// Initialize Event Handlers
	eventHandler := events.NewEventHandler(contributionService, goalService, publisher)

	// Start consuming events if RabbitMQ is connected
	if rabbitConn != nil {
		consumer, err := messaging.NewRabbitMQConsumer(rabbitConn, cfg.RabbitMQ.Exchange, cfg.RabbitMQ.QueueName)
		if err != nil {
			log.Printf("Failed to create RabbitMQ consumer: %v", err)
		} else {
			err = consumer.Consume("PaymentVerified", eventHandler.HandlePaymentVerified)
			if err != nil {
				log.Printf("Failed to start consuming PaymentVerified: %v", err)
			}
		}
	}

	// Initialize Controllers
	goalController := controllers.NewGoalController(goalService)
	contributionController := controllers.NewContributionController(contributionService, withdrawalService, proofService, voteService)
	refundController := controllers.NewRefundController(refundService)

	// Setup Router
	if cfg.Server.Env == "production" {
		gin.SetMode(gin.ReleaseMode)
	}
	r := gin.Default()

	// Datadog tracing middleware
	r.Use(gintrace.Middleware(cfg.Datadog.Service))

	// Routes
	api := r.Group("/api/v1/goals")
	{
		// Public routes (or read-only)
		api.GET("", goalController.ListPublicGoals)
		api.GET("/list", goalController.ListPublicGoals) // Alias for frontend compatibility
		api.GET("/:id", goalController.GetGoal)
		api.GET("/view/:id", goalController.GetGoal) // Alias for frontend compatibility
		api.GET("/:id/progress", goalController.GetGoalProgress)
		api.GET("/proofs", contributionController.GetProofs)
		api.GET("/proofs/:proofId/stats", contributionController.GetVoteStats)

		// Protected routes
		protected := api.Group("")
		protected.Use(middleware.AuthMiddleware())
		{
			protected.GET("/my", goalController.GetMyGoals)
			protected.POST("", goalController.CreateGoal)
			protected.PATCH("/:id", goalController.UpdateGoal)
			protected.POST("/:id/milestones", goalController.CreateMilestone)
			protected.GET("/:goalId/milestones", goalController.GetGoalMilestones)
			protected.POST("/milestones/:milestoneId/complete", goalController.CompleteMilestone)
			
			protected.POST("/contribute", contributionController.CreateContribution)
			protected.POST("/withdraw", contributionController.CreateWithdrawal)
			protected.POST("/proofs", contributionController.CreateProof)
			protected.POST("/votes", contributionController.CreateVote)

			protected.POST("/refunds", refundController.InitiateRefund)
			protected.GET("/refunds/:id", refundController.GetRefund)
			protected.GET("/goals/:goalId/refunds", refundController.GetGoalRefunds)
		}
	}

	// Contributions routes
	contributions := r.Group("/api/v1/contributions")
	contributions.Use(middleware.AuthMiddleware())
	{
		contributions.GET("/my", contributionController.GetMyContributions)
		contributions.GET("/:id", contributionController.GetContribution)
		contributions.POST("", contributionController.CreateContribution)
	}

	// Health check
	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "UP", "timestamp": time.Now()})
	})

	// Start Server with Graceful Shutdown
	port := cfg.Server.Port
	srv := &http.Server{
		Addr:    ":" + port,
		Handler: r,
	}

	go func() {
		log.Printf("Goals Service starting on port %s", port)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Failed to start server: %v", err)
		}
	}()

	// Wait for interrupt signal to gracefully shut down the server with
	// a timeout of 5 seconds.
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal("Server forced to shutdown:", err)
	}

	log.Println("Server exiting")
}

func stringToInt(s string) int {
	var n int
	fmt.Sscanf(s, "%d", &n)
	return n
}
