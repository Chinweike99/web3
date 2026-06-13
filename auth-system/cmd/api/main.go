package main

import (
	"auth-system/internal/config"
	"auth-system/internal/database"
	"auth-system/internal/handlers"
	"auth-system/internal/middlewares"
	"auth-system/internal/models"
	"auth-system/internal/repository"
	"auth-system/internal/service"
	"log"

	"github.com/gin-gonic/gin"
)

func main() {
	config.LoadConfig()

	// database.ConnectDB()

	database.ConnectDB()

	database.DB.AutoMigrate(
		&models.User{},
		&models.RefreshToken{},
	)

	// initialise repo
	userRepo := repository.NewUserRepository()
	emailService := service.NewEmailService(userRepo)
	// Initialize services
	tokenService := service.NewTokenService(userRepo)
	authService := service.NewAuthService(userRepo, tokenService, emailService)

	// Initialize handlers
	authHandler := handlers.NewAuthHandler(authService)

	// Initialse middleware
	authMiddleware := middlewares.NewAuthMiddleware(tokenService)

	// Setup  router
	router := gin.Default()

	// Public routes
	auth := router.Group("/api/v1/auth")
	{
		auth.POST("/register", authHandler.Register)
		auth.POST("/login", authHandler.Login)
		auth.POST("/refresh", authHandler.RefreshToken)
		auth.GET("/verify-email", authHandler.VerifyEmail)  // Add this
    	auth.POST("/resend-verification", authHandler.ResendVerification)  // Add this
	}

	// Protected routes
	protected := router.Group("/api/v1")
	protected.Use(authMiddleware.Authenticate())
	{
		protected.POST("/logout", authHandler.Logout)
		protected.GET("/profile", authHandler.GetProfile)

		// Admin only routes
		admin := protected.Group("/admin")
		admin.Use(authMiddleware.RequireRole("admin"))
		{
			admin.GET("/dashboard", func(c *gin.Context) {
				c.JSON(200, gin.H{"message": "Welcme to admin dashboard"})
			})
		}

		// Moderator and admin routes
		mod := protected.Group("/mod")
		mod.Use(authMiddleware.RequireRole("moderator", "admin"))
		{
			mod.GET("/reports", func(c *gin.Context) {
				c.JSON(200, gin.H{"message": "Reports access granted"})
			})
		}

	}

	// Start server
	log.Printf("Server starting on port %s", config.AppConfig.ServerPort)
	if err := router.Run(":" + config.AppConfig.ServerPort); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}
