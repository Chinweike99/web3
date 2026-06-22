package middlewares

import (
	"auth-system/internal/config"
	"github.com/gin-gonic/gin"
	"github.com/ulule/limiter/v3"
	"time"
)

func SetupRateLimiting(router *gin.Engine, rateLimiter *RateLimiterMiddleware) {
	// Strict limits for authentication endpoints
	loginRate := limiter.Rate{
		Period: 1 * time.Minute,
		Limit:  int64(config.AppConfig.RateLimitLoginPerMin),
	}

	registerRate := limiter.Rate{
		Period: 1 * time.Minute,
		Limit:  int64(config.AppConfig.RateLimitRegisterPerMin),
	}

	forgotPasswordRate := limiter.Rate{
		Period: 1 * time.Hour,
		Limit:  3,
	}

	resetPasswordRate := limiter.Rate{
		Period: 1 * time.Hour,
		Limit:  5,
	}

	verifyEmailRate := limiter.Rate{
		Period: 1 * time.Hour,
		Limit:  10,
	}

	// General API rate limit
	apiRate := limiter.Rate{
		Period: 1 * time.Minute,
		Limit:  int64(config.AppConfig.RateLimitPerMinute),
	}

	// Apply to specific routes
	authGroup := router.Group("/api/v1/auth")
	{
		authGroup.POST("/login",
			rateLimiter.RateLimitByEmail(loginRate),
		)
		authGroup.POST("/register",
			rateLimiter.RateLimitByEmail(registerRate),
		)
		authGroup.POST("/forgot-password",
			rateLimiter.RateLimitByEmail(forgotPasswordRate),
		)
		authGroup.POST("/reset-password",
			rateLimiter.RateLimitByIP(resetPasswordRate),
		)
		authGroup.GET("/verify-email",
			rateLimiter.RateLimitByIP(verifyEmailRate),
		)
		authGroup.POST("/resend-verification",
			rateLimiter.RateLimitByEmail(forgotPasswordRate),
		)
		authGroup.POST("/refresh",
			rateLimiter.RateLimitByIP(apiRate),
		)
	}

	// Protected routes with rate limiting
	protectedGroup := router.Group("/api/v1")
	protectedGroup.Use(rateLimiter.RateLimitByIPAndEndpoint(apiRate))
	{
		protectedGroup.POST("/logout")
		protectedGroup.GET("/profile")
		protectedGroup.POST("/change-password")
	}
}
