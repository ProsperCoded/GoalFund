package router

import (
	"github.com/gin-gonic/gin"
	"github.com/gofund/users-service/internal/controllers"
	"github.com/gofund/users-service/internal/service"
)

// SetupRoutes configures all routes for the Users Service
func SetupRoutes(r *gin.Engine, authService *service.AuthService, userService *service.UserService, kycService *service.KYCService) {
	// Initialize controllers
	authController := controllers.NewAuthController(authService, userService)
	userController := controllers.NewUserController(authService, userService)
	kycController := controllers.NewKYCController(kycService)


	// Health check endpoint
	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "healthy", "service": "users-service"})
	})

	// Internal routes (called by Nginx, not exposed externally)
	internal := r.Group("/internal")
	{
		// Auth verification endpoint for Nginx auth_request
		internal.GET("/verify", authController.VerifyToken)
	}

	// Public authentication routes (no auth required)
	auth := r.Group("/auth")
	{
		auth.POST("/login", authController.Login)
		auth.POST("/register", authController.Register)
		auth.POST("/refresh", authController.RefreshToken)
		auth.POST("/logout", authController.Logout)
		auth.POST("/forgot-password", authController.ForgotPassword)
		auth.POST("/reset-password", authController.ResetPassword)
	}

	// Protected user routes (auth required - handled by Nginx)
	users := r.Group("/users")
	{
		users.GET("/profile", authController.GetProfile)
		users.PUT("/profile", authController.UpdateProfile)
		users.PUT("/settlement-account", userController.UpdateSettlementAccount)
		
		// KYC verification routes
		kyc := users.Group("/kyc")
		{
			kyc.POST("/submit-nin", kycController.SubmitNIN)
			kyc.GET("/status", kycController.GetKYCStatus)
		}
	}

	// Public user routes (for guest contributions/onboarding)
	publicUsers := r.Group("/public/users")
	{
		publicUsers.POST("/contribution-signup", userController.CreateLightweightUser)
		publicUsers.POST("/set-password", userController.SetPassword)
	}


	// TODO: Add role management routes
	// roles := r.Group("/roles")
	// {
	//     roles.GET("/", roleController.ListRoles)
	//     roles.POST("/", roleController.CreateRole)
	//     roles.GET("/:id", roleController.GetRole)
	//     roles.PUT("/:id", roleController.UpdateRole)
	//     roles.DELETE("/:id", roleController.DeleteRole)
	// }
}