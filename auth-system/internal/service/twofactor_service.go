package service

import (
	"auth-system/internal/models"
	"auth-system/internal/repository"
	"crypto/rand"
	"errors"

	"github.com/google/uuid"
	"github.com/pquerna/otp"
	"github.com/pquerna/otp/totp"
	"github.com/skip2/go-qrcode"
)

type TwoFactorService struct {
	userRepo *repository.UserRepository
}

func NewTwoFactorService(userRepo *repository.UserRepository) *TwoFactorService {
	return &TwoFactorService{
		userRepo: userRepo,
	}
}

// GenerateTOTPSecret generates a new TOTP secret for a user
func (s *TwoFactorService) GenerateTOTPSecret(email string) (string, string, error) {
	key, err := totp.Generate(totp.GenerateOpts{
		Issuer:      "AuthSystem",
		AccountName: email,
		Period:      30,
		Digits:      6,
		Algorithm:   otp.AlgorithmSHA1,
	})
	if err != nil {
		return "", "", err
	}

	return key.Secret(), key.URL(), nil
}

func (s *TwoFactorService) GenerateQRCode(issuer, accountName, secret string) ([]byte, error) {
	key, err := totp.Generate(totp.GenerateOpts{
		Issuer:      issuer,
		AccountName: accountName,
		Secret:      []byte(secret),
	})
	if err != nil {
		return nil, err
	}

	png, err := qrcode.Encode(key.URL(), qrcode.Medium, 256)
	if err != nil {
		return nil, err
	}

	return png, nil
}

func (s *TwoFactorService) EnableTwoFactor(userID uuid.UUID, token, secret string) error {
	// Validate TOTP token
	valid := totp.Validate(token, secret)
	if !valid {
		return errors.New("invalid TOTP token")
	}

	// Update user
	return s.userRepo.EnableTwoFactor(userID, secret)
}

func (s *TwoFactorService) DisableTwoFactor(userID uuid.UUID, password, token string) error {
	// Get user
	user, err := s.userRepo.FindByID(userID)
	if err != nil {
		return err
	}
	if user == nil {
		return errors.New("user not found")
	}

	// Verify password
	if err := s.userRepo.VerifyPassword(user, password); err != nil {
		return errors.New("invalid password")
	}

	// Validate TOTP if 2FA is enabled
	if user.TwoFactorEnabled {
		valid := totp.Validate(token, user.TwoFactorSecret)
		if !valid {
			return errors.New("invalid TOTP token")
		}
	}

	// Disable 2FA
	return s.userRepo.DisableTwoFactor(userID)
}

func (s *TwoFactorService) VerifyTwoFactorCode(secret, token string) bool {
	return totp.Validate(token, secret)
}

func (s *TwoFactorService) GenerateBackupCodes(userID uuid.UUID) ([]string, error) {
	// Generate 10 backup codes
	codes := make([]string, 10)
	backupCodeModels := make([]models.BackupCode, 10)

	for i := 0; i < 10; i++ {
		code, err := s.generateRandomCode()
		if err != nil {
			return nil, err
		}
		codes[i] = code
		backupCodeModels[i] = models.BackupCode{
			UserID: userID,
			Code:   code,
			Used:   false,
		}
	}

	// Save backup codes to database
	if err := s.userRepo.SaveBackupCodes(userID, backupCodeModels); err != nil {
		return nil, err
	}

	return codes, nil
}

func (s *TwoFactorService) VerifyBackupCode(userID uuid.UUID, code string) (bool, error) {
	backupCode, err := s.userRepo.FindBackupCode(userID, code)
	if err != nil {
		return false, err
	}
	if backupCode == nil {
		return false, errors.New("invalid backup code")
	}
	if backupCode.Used {
		return false, errors.New("backup code already used")
	}
	// Mark code as used
	if err := s.userRepo.MarkBackupCodeUsed(backupCode.ID); err != nil {
		return false, err
	}
	return true, nil
}

func (s *TwoFactorService) RegenerateBackupCodes(userID uuid.UUID) ([]string, error) {
	// Delete old backup codes
	if err := s.userRepo.DeleteBackupCodes(userID); err != nil {
		return nil, err
	}
	return s.GenerateBackupCodes(userID)
}

func (s *TwoFactorService) generateRandomCode() (string, error) {
	const chars = "ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	bytes := make([]byte, 8)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}

	for i, b := range bytes {
		bytes[i] = chars[b%byte(len(chars))]
	}
	code := string(bytes)
	return code[:4] + "-" + code[4:], nil
}

// GetTwoFactorStatus returns the 2FA status for a user
func (s *TwoFactorService) GetTwoFactorStatus(userID uuid.UUID) (bool, error) {
	user, err := s.userRepo.FindByID(userID)
	if err != nil {
		return false, err
	}
	if user == nil {
		return false, errors.New("user not found")
	}

	return user.TwoFactorEnabled, nil
}
