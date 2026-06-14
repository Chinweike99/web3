package middlewares

import (
	"auth-system/internal/config"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/ulule/limiter/v3"
	"github.com/ulule/limiter/v3/drivers/store/memory"
)

type RateLimiterMiddleware struct {
	store limiter.Store
}

func NewRateLimiterMiddleware() *RateLimiterMiddleware {
	store := memory.NewStore()
	return &RateLimiterMiddleware{
		store: store,
	}
}

func (rl *RateLimiterMiddleware) RateLimit(rate limiter.Rate, keyFunc func(c *gin.Context) string) gin.HandlerFunc {
	return func(c *gin.Context) {
		if !config.AppConfig.RateLimitEnabled {
			c.Next()
			return
		}

		key := keyFunc(c)
		limiterInstance := limiter.New(rl.store, rate)

		ctx := c.Request.Context()
		result, err := limiterInstance.Get(ctx, key)
		if err != nil {
			c.AbortWithStatusJSON(500, gin.H{"error": "Internal server error"})
			return
		}
		resetTime := time.Unix(result.Reset, 0)
		c.Header("X-RateLimit-Limit", strconv.FormatInt(rate.Limit, 10))
		c.Header("X-RateLimit-Remaining", strconv.FormatInt(result.Remaining, 10))
		c.Header("X-RateLimit-Reset", resetTime.Format(time.RFC3339))

		if result.Reached {
			retryAfter := time.Until(resetTime)

			c.Header(
				"Retry-After",
				strconv.FormatInt(int64(retryAfter.Seconds()), 10),
			)

			c.JSON(http.StatusTooManyRequests, gin.H{
				"error":       "Too many requests. Please try again later.",
				"retry_after": retryAfter.Seconds(),
			})

			c.Abort()
			return
		}
		c.Next()
	}
}

// IP-based rate limiting
func (m *RateLimiterMiddleware) RateLimitByIP(rate limiter.Rate) gin.HandlerFunc {
	return m.RateLimit(rate, func(c *gin.Context) string {
		return c.ClientIP()
	})
}

// Email-based rate limiting (for login/register)
func (m *RateLimiterMiddleware) RateLimitByEmail(rate limiter.Rate) gin.HandlerFunc {
	return m.RateLimit(rate, func(c *gin.Context) string {
		var email string
		if c.Request.Method == "POST" {
			var body map[string]interface{}
			if err := c.ShouldBindBodyWith(&body, binding.JSON); err == nil {
				if e, ok := body["email"]; ok {
					email = e.(string)
				}
			}
		}
		return email
	})
}

// User ID-based rate limiting
func (m *RateLimiterMiddleware) RateLimitByUserID(rate limiter.Rate) gin.HandlerFunc {
	return m.RateLimit(rate, func(c *gin.Context) string {
		userID, exists := c.Get("user_id")
		if exists {
			return userID.(string)
		}
		return c.ClientIP()
	})
}

// Combined rate limiting (IP + endpoint)
func (m *RateLimiterMiddleware) RateLimitByIPAndEndpoint(rate limiter.Rate) gin.HandlerFunc {
	return m.RateLimit(rate, func(c *gin.Context) string {
		return c.ClientIP() + ":" + c.FullPath()
	})
}
