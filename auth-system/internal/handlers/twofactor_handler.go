package handlers

import (
	"auth-system/internal/models"
	"auth-system/internal/service"
	"encoding/base64"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type TwoFactorHandler struct {
	twoFactorService *service.TwoFactorService
	authService      *service.AuthService
}

func NewTwoFactorHandler(twoFactorService *service.TwoFactorService, authService *service.AuthService) *TwoFactorHandler {
	return &TwoFactorHandler{
		twoFactorService: twoFactorService,
		authService:      authService,
	}
}

// Setup2FA initiates 2FA setup for a user
func (h *TwoFactorHandler) Setup2FA(c *gin.Context) {
	userIDStr, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	var req models.TwoFactorSetupRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userID, err := uuid.Parse(userIDStr.(string))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	// Get user to verify password
	user, err := h.authService.GetUserByID(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get user"})
		return
	}

	// Verify password
	if err := h.authService.VerifyUserPassword(user, req.Password); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid password"})
		return
	}

	// Generate TOTP secret
	secret, url, err := h.twoFactorService.GenerateTOTPSecret(user.Email)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate secret"})
		return
	}

	// Generate QR code
	qrCodeBytes, err := h.twoFactorService.GenerateQRCode("AuthSystem", user.Email, secret)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate QR code"})
		return
	}

	qrCodeBase64 := base64.StdEncoding.EncodeToString(qrCodeBytes)

	c.JSON(http.StatusOK, gin.H{
		"secret":  secret,
		"url":     url,
		"qr_code": qrCodeBase64,
		"message": "Scan the QR code with Google Authenticator or similar app",
	})
}

func (h *TwoFactorHandler) Enable2FA(c *gin.Context) {
	userIDStr, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	var req models.TwoFactorVerifyRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userID, err := uuid.Parse(userIDStr.(string))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	// Get the secret from temporary storage (you might want to store it temporarily)
	// For simplicity, we'll assume the secret is passed or stored in session
	// This is a simplified version - in production, store the secret temporarily

	// For this example, we'll need to retrieve the secret from a temporary store
	// Let's assume we have a way to get it
	secret := c.GetHeader("X-2FA-Secret")
	if secret == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "2FA secret not found"})
		return
	}

	err = h.twoFactorService.EnableTwoFactor(userID, req.Token, secret)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Generate backup codes
	backupCodes, err := h.twoFactorService.GenerateBackupCodes(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate backup codes"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":      "2FA enabled successfully",
		"backup_codes": backupCodes,
		"warning":      "Save these backup codes securely. You won't see them again.",
	})
}

// Disable2FA disables 2FA for the user
func (h *TwoFactorHandler) Disable2FA(c *gin.Context) {
	userIDStr, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	var req models.TwoFactorDisableRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userID, err := uuid.Parse(userIDStr.(string))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	err = h.twoFactorService.DisableTwoFactor(userID, req.Password, req.Token)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"message": "2FA disabled successfully",
	})
}

func (h *TwoFactorHandler) Verify2FALogin(c *gin.Context) {
	var req models.TwoFactorLoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ip := c.ClientIP()
	user, tokenResponse, err := h.authService.VerifyTwoFactorLogin(&req, ip)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"message": "Login successful",
		"user": gin.H{
			"id":         user.ID,
			"email":      user.Email,
			"first_name": user.FirstName,
			"last_name":  user.LastName,
			"role":       user.Role,
		},
		"tokens": tokenResponse,
	})
}

func (h *TwoFactorHandler) Get2FAStatus(c *gin.Context) {
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

	enabled, err := h.twoFactorService.GetTwoFactorStatus(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"two_factor_enabled": enabled,
	})
}

func (h *TwoFactorHandler) RegenerateBackupCodes(c *gin.Context) {
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
	backupCodes, err := h.twoFactorService.RegenerateBackupCodes(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"message":      "Backup codes regenerated successfully",
		"backup_codes": backupCodes,
		"warning":      "Save these backup codes securely. Old codes are invalid now.",
	})
}
