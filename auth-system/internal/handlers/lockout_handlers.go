package handlers

import (
	"auth-system/internal/models"
	"auth-system/internal/service"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"net/http"
)

type LockoutHandler struct {
	authService *service.AuthService
}

func NewLockoutHandler(authService *service.AuthService) *LockoutHandler {
	return &LockoutHandler{
		authService: authService,
	}
}

// UnlockAccount allows users to request account unlock (usually via email)
func (h *LockoutHandler) UnlockAccount(c *gin.Context) {
	var req models.UnlockAccountRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := h.authService.UnlockAccount(req.Email)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Account unlock request processed. Check your email for confirmation.",
	})
}

// AdminUnlockAccount allows admins to unlock any account
func (h *LockoutHandler) AdminUnlockAccount(c *gin.Context) {
	var req models.AdminUnlockRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userID, err := uuid.Parse(req.UserID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	err = h.authService.AdminUnlockAccount(userID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Account unlocked successfully",
	})
}

// GetLockedAccounts returns all locked accounts (admin only)
func (h *LockoutHandler) GetLockedAccounts(c *gin.Context) {
	lockedAccounts, err := h.authService.GetLockedAccounts()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"locked_accounts": lockedAccounts,
		"count":           len(lockedAccounts),
	})
}

// GetAccountLockHistory returns lock history for the authenticated user
func (h *LockoutHandler) GetAccountLockHistory(c *gin.Context) {
	userIDStr, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	userID, err := uuid.Parse(userIDStr.(string))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	lockHistory, err := h.authService.GetAccountLockHistory(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"lock_history": lockHistory,
	})
}
