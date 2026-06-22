package service

import (
	"auth-system/internal/config"
	"auth-system/internal/models"
	"auth-system/internal/repository"
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type AuthService struct {
	userRepo         *repository.UserRepository
	tokenService     *TokenService
	emailService     *EmailService
	rateLimitService *RateLimitService
	twoFactorService *TwoFactorService
}

func NewAuthService(userRepo *repository.UserRepository,
	tokenService *TokenService,
	emailService *EmailService,
	rateLimitService *RateLimitService,
	twoFactorService *TwoFactorService,
) *AuthService {
	return &AuthService{
		userRepo:         userRepo,
		tokenService:     tokenService,
		emailService:     emailService,
		rateLimitService: rateLimitService,
		twoFactorService: twoFactorService,
	}
}

func (s *AuthService) Register(req *models.RegisterRequest, ip string) (*models.User, error) {
	// Track registration attempt
	defer func() {
		s.rateLimitService.RecordAttempt(ip, req.Email, true)
	}()

	existingUser, err := s.userRepo.FindByEmail(req.Email)
	if err != nil {
		return nil, err
	}
	if existingUser != nil {
		return nil, errors.New("user with this email already exists")
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}
	user := &models.User{
		ID:        uuid.New(),
		Email:     req.Email,
		Password:  string(hashedPassword),
		FirstName: req.FirstName,
		LastName:  req.LastName,
		Role:      models.RoleUser,
		IsActive:  true,
	}

	err = s.userRepo.Create(user)
	if err != nil {
		return nil, err
	}

	if err := s.emailService.SendVerificationEmail(user); err != nil {
		log.Printf("Failed to send verification email: %v", err)
	}

	return user, nil
}

func (s *AuthService) Login(req *models.LoginRequest, ip string) (*models.User, *models.TokenResponse, error) {

	// Track attempt
	success := false
	defer func() {
		s.rateLimitService.RecordAttempt(ip, req.Email, success)
	}()

	user, err := s.userRepo.FindByEmail(req.Email)
	if err != nil {
		return nil, nil, err
	}
	if user == nil {
		return nil, nil, errors.New("Invalid credentials")
	}
	if !user.IsActive {
		return nil, nil, errors.New("Account is inactive")
	}

	// Check if account is locked
	isLocked, lockedUntil, err := s.userRepo.IsAccountLocked(user.ID)
	if err != nil {
		return nil, nil, err
	}
	if isLocked {
		remainingMinutes := int(time.Until(*lockedUntil).Minutes())
		return nil, nil, fmt.Errorf("account is locked. Please try again in %d minutes", remainingMinutes)
	}

	if !user.EmailVerified {
		return nil, nil, errors.New("please verify your email address before logging in")
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password))
	if err != nil {
		// Record failed attempt
		s.userRepo.RecordLoginAttempt(user.ID, false, ip)
		newAttempts := user.LoginAttempts + 1

		if newAttempts >= config.AppConfig.MaxLoginAttempts-2 && config.AppConfig.NotifyOnLockout {
			go s.emailService.SendSuspiciousActivityAlert(user, ip, newAttempts)
		}
		// Send lockout notification if account is now locked
		if newAttempts >= config.AppConfig.MaxLoginAttempts {
			if config.AppConfig.NotifyOnLockout {
				lockedUntil := time.Now().Add(time.Duration(config.AppConfig.LockoutDurationMinutes) * time.Minute)
				go s.emailService.SendAccountLockedEmail(user, lockedUntil, ip)
			}

			// Log the lock event
			lock := &models.AccountLock{
				UserID:       user.ID,
				IPAddress:    ip,
				AttemptCount: newAttempts,
				LockedAt:     time.Now(),
			}
			s.userRepo.CreateAccountLock(lock)
		}
		return nil, nil, errors.New("Invalid credentials")
	}

	success = true
	// Record successful login
	s.userRepo.RecordLoginAttempt(user.ID, true, ip)

	// Check if 2FA is enabled
	if user.TwoFactorEnabled {
		return user, nil, errors.New("2FA_REQUIRED")
	}

	// Generate tokens
	accessToken, err := s.tokenService.GenerateAccessToken(user)
	if err != nil {
		return nil, nil, err
	}
	refreshToken, _, err := s.tokenService.GenerateRefreshToken(user.ID)
	if err != nil {
		return nil, nil, err
	}
	tokenResponse := &models.TokenResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		TokenType:    "Bearer",
		ExpiresIn:    config.AppConfig.JwtAccessExpiration * 60,
	}
	return user, tokenResponse, nil
}

func (s *AuthService) RefreshAccessToken(refreshTokenString string) (*models.TokenResponse, error) {
	refreshToken, err := s.tokenService.ValidateRefreshToken(refreshTokenString)
	if err != nil {
		return nil, err
	}
	user, err := s.userRepo.FindByID(refreshToken.UserID)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, errors.New("userr not found")
	}
	// Generate accessToken
	accessToken, err := s.tokenService.GenerateAccessToken(user)
	if err != nil {
		return nil, err
	}
	tokenResponse := &models.TokenResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshTokenString,
		TokenType:    "Bearer",
		ExpiresIn:    config.AppConfig.JwtAccessExpiration,
	}
	return tokenResponse, nil
}

func (s *AuthService) Logout(userID uuid.UUID) error {
	return s.userRepo.RevokeAllUserRefreshTokens(userID)
}

func (s *AuthService) VerifyEmail(token string) error {
	verification, err := s.userRepo.FindEmailVerificationByToken(token)
	if err != nil {
		return err
	}

	if verification == nil {
		return errors.New("Invalid or expired token")
	}

	if verification.ExpiresAt.Before(time.Now()) {
		return errors.New("Verification token expired")
	}

	if err := s.userRepo.MarkEmailAsVerified(verification.UserID); err != nil {
		return err
	}

	if err := s.userRepo.MarkVerificatioTokenUsed(verification.ID); err != nil {
		return err
	}

	user, err := s.userRepo.FindByID(verification.UserID)
	if err != nil {
		return err
	}
	go s.emailService.SendWecomeEmail(user)

	return nil
}

func (s *AuthService) ResendVerificationEmail(email string) error {
	user, err := s.userRepo.FindByEmail(email)
	if err != nil {
		return err
	}
	if user == nil {
		return errors.New("User not found")
	}
	if user.EmailVerified {
		return errors.New("email already verified")
	}
	return s.emailService.SendVerificationEmail(user)
}

func (s *AuthService) ForgotPassword(email string) error {
	user, err := s.userRepo.FindByEmail(email)
	if err != nil {
		return err
	}
	if user == nil {
		return errors.New("User not found")
	}

	token, err := s.generatePasswordResetToken()
	if err != nil {
		return err
	}

	passwordReset := &models.PasswordReset{
		UserID:    user.ID,
		Token:     token,
		ExpiresAt: time.Now().Add(1 * time.Hour),
	}
	if err := s.userRepo.CreatePasswordReset(passwordReset); err != nil {
		return err
	}
	// Send reset email
	go s.emailService.SendPasswordResetEmail(user, token)
	return nil
}

func (s *AuthService) ResetPassword(req *models.ResetPasswordRequest) error {
	// Validate passwords match
	if req.NewPassword != req.ConfirmPassword {
		return errors.New("passwords do not match")
	}
	reset, err := s.userRepo.FindPasswordResetByToken(req.Token)
	if err != nil {
		return err
	}
	if reset == nil {
		return errors.New("invalid or expired reset token")
	}

	// Check if token is expired
	if reset.ExpiresAt.Before(time.Now()) {
		return errors.New("reset token has expired")
	}

	// Hash new password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.NewPassword), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	// Update password
	if err := s.userRepo.UpdatePassword(reset.UserID, string(hashedPassword)); err != nil {
		return err
	}

	// Mark reset token as used
	if err := s.userRepo.MarkPasswordResetUsed(reset.ID); err != nil {
		return err
	}
	if err := s.userRepo.RevokeAllUserRefreshTokens(reset.UserID); err != nil {
		log.Printf("Failed to revoke refresh tokens: %v", err)
	}
	user, err := s.userRepo.FindByID(reset.UserID)
	if err == nil && user != nil {
		go s.emailService.SendPasswordChangedEmail(user)
	}

	return nil
}

func (s *AuthService) ChangePassword(userID uuid.UUID, req *models.ChangePasswordRequest) error {
	if req.NewPassword != req.ConfirmPassword {
		return errors.New("passwords do not match")
	}
	user, err := s.userRepo.FindByID(userID)
	if err != nil {
		return err
	}
	if user == nil {
		return errors.New("user not found")
	}
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.CurrentPassword))
	if err != nil {
		return errors.New("current password is incorrect")
	}
	// Hash new password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.NewPassword), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	if err := s.userRepo.UpdatePassword(userID, string(hashedPassword)); err != nil {
		return err
	}
	if err := s.userRepo.RevokeAllUserRefreshTokens(userID); err != nil {
		log.Printf("Failed to revoke refresh tokens: %v", err)
	}
	go s.emailService.SendPasswordChangedEmail(user)
	return nil
}

func (s *AuthService) generatePasswordResetToken() (string, error) {
	bytes := make([]byte, 32)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes), nil
}

func (s *AuthService) VerifyTwoFactorLogin(req *models.TwoFactorLoginRequest, ip string) (*models.User, *models.TokenResponse, error) {
	user, err := s.userRepo.FindByEmail(req.Email)
	if err != nil {
		return nil, nil, err
	}
	if user == nil {
		return nil, nil, errors.New("Invalid credentials")
	}

	// Verify password
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password))
	if err != nil {
		s.rateLimitService.RecordAttempt(ip, req.Email, false)
		return nil, nil, errors.New("invalid credentials")
	}

	valid := s.twoFactorService.VerifyTwoFactorCode(user.TwoFactorSecret, req.Token)
	if !valid {
		// Check if it's a backup code
		backupValid, err := s.twoFactorService.VerifyBackupCode(user.ID, req.Token)
		if err != nil || !backupValid {
			s.rateLimitService.RecordAttempt(ip, req.Email, false)
			return nil, nil, errors.New("invalid 2FA code")
		}
	}
	s.rateLimitService.RecordAttempt(ip, req.Email, true)

	// Generate tokens
	accessToken, err := s.tokenService.GenerateAccessToken(user)
	if err != nil {
		return nil, nil, err
	}

	refreshToken, _, err := s.tokenService.GenerateRefreshToken(user.ID)
	if err != nil {
		return nil, nil, err
	}

	tokenResponse := &models.TokenResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		TokenType:    "Bearer",
		ExpiresIn:    config.AppConfig.JwtAccessExpiration * 60,
	}
	return user, tokenResponse, nil
}

func (s *AuthService) GetUserByID(userID uuid.UUID) (*models.User, error) {
	return s.userRepo.FindByID(userID)
}

func (s *AuthService) VerifyUserPassword(user *models.User, password string) error {
	return s.userRepo.VerifyPassword(user, password)
}

// UnlockAccount manually unlocks a locked account
func (s *AuthService) UnlockAccount(email string) error {
	user, err := s.userRepo.FindByEmail(email)
	if err != nil {
		return err
	}
	if user == nil {
		return errors.New("user not found")
	}

	// Check if account is actually locked
	isLocked, _, err := s.userRepo.IsAccountLocked(user.ID)
	if err != nil {
		return err
	}
	if !isLocked {
		return errors.New("account is not locked")
	}

	// Unlock account
	if err := s.userRepo.UnlockAccount(user.ID); err != nil {
		return err
	}

	// Send unlock notification
	if config.AppConfig.NotifyOnLockout {
		go s.emailService.SendAccountUnlockedEmail(user)
	}

	return nil
}

// AdminUnlockAccount unlocks an account by admin
func (s *AuthService) AdminUnlockAccount(userID uuid.UUID) error {
	user, err := s.userRepo.FindByID(userID)
	if err != nil {
		return err
	}
	if user == nil {
		return errors.New("user not found")
	}

	if err := s.userRepo.UnlockAccount(user.ID); err != nil {
		return err
	}

	// Send notification
	if config.AppConfig.NotifyOnLockout {
		go s.emailService.SendAccountUnlockedEmail(user)
	}

	return nil
}

// GetLockedAccounts returns all currently locked accounts
func (s *AuthService) GetLockedAccounts() ([]models.User, error) {
	return s.userRepo.GetLockedAccounts()
}

// GetAccountLockHistory returns lock history for a user
func (s *AuthService) GetAccountLockHistory(userID uuid.UUID) ([]models.AccountLock, error) {
	return s.userRepo.GetAccountLocks(userID, 10)
}
