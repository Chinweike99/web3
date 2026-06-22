package service

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"sync"
	"time"
)

type RateLimitService struct {
	attempts   map[string][]time.Time
	mu         sync.RWMutex
	blockedIPs map[string]time.Time
	blockedMu  sync.RWMutex
}

type AttemptInfo struct {
	IP        string
	Email     string
	Timestamp time.Time
	Success   bool
}

func NewRateLimitService() *RateLimitService {
	rls := &RateLimitService{
		attempts:   make(map[string][]time.Time),
		blockedIPs: make(map[string]time.Time),
	}

	// Start cleanup goroutine
	go rls.cleanupOldAttempts()

	return rls
}

func (rls *RateLimitService) RecordAttempt(ip, email string, success bool) {
	rls.mu.Lock()
	defer rls.mu.Unlock()

	key := ip + ":" + email
	if _, exists := rls.attempts[key]; !exists {
		rls.attempts[key] = []time.Time{}
	}

	rls.attempts[key] = append(rls.attempts[key], time.Now())

	// Keep only last 100 attempts
	if len(rls.attempts[key]) > 100 {
		rls.attempts[key] = rls.attempts[key][len(rls.attempts[key])-100:]
	}

	// Check for brute force pattern (multiple failed attempts)
	if !success {
		rls.checkBruteForce(ip, email)
	}
}

func (rls *RateLimitService) checkBruteForce(ip, email string) {
	key := ip + ":" + email
	rls.mu.RLock()
	attempts := rls.attempts[key]
	rls.mu.RUnlock()

	// Count failed attempts in last 15 minutes
	fifteenMinutesAgo := time.Now().Add(-15 * time.Minute)
	failedCount := 0
	for _, attempt := range attempts {
		if attempt.After(fifteenMinutesAgo) {
			failedCount++
		}
	}

	// Block IP after 10 failed attempts in 15 minutes
	if failedCount >= 10 {
		rls.blockedMu.Lock()
		rls.blockedIPs[ip] = time.Now().Add(30 * time.Minute)
		rls.blockedMu.Unlock()
	}
}

func (rls *RateLimitService) IsIPBlocked(ip string) bool {
	rls.blockedMu.RLock()
	defer rls.blockedMu.RUnlock()

	if blockedUntil, exists := rls.blockedIPs[ip]; exists {
		if time.Now().Before(blockedUntil) {
			return true
		}
		// Remove expired block
		delete(rls.blockedIPs, ip)
	}
	return false
}

func (rls *RateLimitService) cleanupOldAttempts() {
	ticker := time.NewTicker(1 * time.Hour)
	for range ticker.C {
		rls.mu.Lock()
		oneHourAgo := time.Now().Add(-1 * time.Hour)
		for key, attempts := range rls.attempts {
			newAttempts := []time.Time{}
			for _, attempt := range attempts {
				if attempt.After(oneHourAgo) {
					newAttempts = append(newAttempts, attempt)
				}
			}
			if len(newAttempts) == 0 {
				delete(rls.attempts, key)
			} else {
				rls.attempts[key] = newAttempts
			}
		}
		rls.mu.Unlock()
	}
}

func (rls *RateLimitService) RateLimitMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		clientIP := c.ClientIP()

		// Check if IP is blocked
		if rls.IsIPBlocked(clientIP) {
			c.JSON(http.StatusForbidden, gin.H{
				"error": "Your IP has been temporarily blocked due to suspicious activity",
			})
			c.Abort()
			return
		}

		c.Next()
	}
}
