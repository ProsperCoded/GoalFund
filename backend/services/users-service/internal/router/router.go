package router

import (
	"github.com/gin-gonic/gin"
	"github.com/gofund/users-service/internal/controllers"
	"github.com/gofund/users-service/internal/service"
)

// SetupRoutes configures all routes for the Users Service
func SetupRoutes(r *gin.Engine, authService *service.AuthService) {
	// Initialize controllers
	authController := controllers.NewAuthController(authService)

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
		
		// TODO: Add more user management endpoints
		// users.GET("/", authController.ListUsers)           // Admin only
		// users.GET("/:id", authController.GetUser)          // Admin or self
		// users.PUT("/:id", authController.UpdateUser)       // Admin or self
		// users.DELETE("/:id", authController.DeleteUser)    // Admin only
		// users.POST("/:id/roles", authController.AssignRole) // Admin only
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