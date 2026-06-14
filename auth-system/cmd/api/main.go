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
	"time"

	"github.com/gin-gonic/gin"
)

func main() {
	config.LoadConfig()

	// database.ConnectDB()
	database.ConnectDB()

	database.DB.AutoMigrate(
		&models.User{},
		&models.RefreshToken{},
		&models.EmailVerification{},
		&models.PasswordReset{},
	)
	log.Println("Running migrations...")

	// initialise repo
	userRepo := repository.NewUserRepository()
	emailService := service.NewEmailService(userRepo)
	rateLimitService := service.NewRateLimitService() 
	
	// Clean Expired token
	startCleanupJob(userRepo)
	// Initialize services
	tokenService := service.NewTokenService(userRepo)
	authService := service.NewAuthService(userRepo, tokenService, emailService, rateLimitService)

	// Initialize handlers
	authHandler := handlers.NewAuthHandler(authService)

	// Initialse middleware
	authMiddleware := middlewares.NewAuthMiddleware(tokenService)
	// rateLimiterMiddleware := middlewares.NewRateLimiterMiddleware()

	// Setup  router
	router := gin.Default()

	// Add global middlewares
    router.Use(rateLimitService.RateLimitMiddleware()) 
	// middlewares.SetupRateLimiting(router, rateLimiterMiddleware)

	// Public routes
	auth := router.Group("/api/v1/auth")
	{
		auth.POST("/register", authHandler.Register)
		auth.POST("/login", authHandler.Login)
		auth.POST("/refresh", authHandler.RefreshToken)
		auth.GET("/verify-email", authHandler.VerifyEmail)
		auth.POST("/resend-verification", authHandler.ResendVerification)
		auth.POST("/forgot-password", authHandler.ForgotPassword)
		auth.POST("/reset-password", authHandler.ResetPassword)
	}

	// Protected routes
	protected := router.Group("/api/v1")
	protected.Use(authMiddleware.Authenticate())
	{
		protected.POST("/logout", authHandler.Logout)
		protected.GET("/profile", authHandler.GetProfile)
		protected.POST("/chaneg-password", authHandler.ChangePassword)

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

func startCleanupJob(repo *repository.UserRepository) {
	ticker := time.NewTicker(24 * time.Hour)
	go func() {
		for range ticker.C {
			log.Println("Running cleanup job for expired tokens...")
			if err := repo.DeleteExpiredVerificationTokens(); err != nil {
				log.Printf("Failed to delete expired verification tokens: %v", err)
			}
			if err := repo.DeleteExpiredPasswordResets(); err != nil {
				log.Printf("Failed to delete expired password resets: %v", err)
			}
		}
	}()
}

// In main() function after initializing repository:
