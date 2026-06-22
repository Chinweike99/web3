package repository

import (
	"auth-system/internal/config"
	"auth-system/internal/database"
	"auth-system/internal/models"
	"errors"
	"time"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type UserRepository struct {
	db *gorm.DB
}

func NewUserRepository() *UserRepository {
	return &UserRepository{
		db: database.DB,
	}
}

func (r *UserRepository) Create(user *models.User) error {
	return r.db.Create(user).Error
}

func (r *UserRepository) FindByEmail(email string) (*models.User, error) {
	var user models.User
	err := r.db.Where("email =  ?", email).First(&user).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &user, nil
}

func (r *UserRepository) FindByID(id uuid.UUID) (*models.User, error) {
	var user models.User
	err := r.db.First(&user, "id = ?", id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &user, nil
}

func (r *UserRepository) Update(user *models.User) error {
	return r.db.Save(user).Error
}

func (r *UserRepository) Delete(user *models.User) error {
	return r.db.Delete(user).Error
}

func (r *UserRepository) List() ([]models.User, error) {
	var users []models.User
	err := r.db.Find(&users).Error
	return users, err
}

func (r *UserRepository) CreateRefreshToken(token *models.RefreshToken) error {
	return r.db.Create(token).Error
}

func (r *UserRepository) FindRefreshToken(token string) (*models.RefreshToken, error) {
	var refreshToken models.RefreshToken
	err := r.db.Where("token = ? AND revoked = ?", token, false).First(&refreshToken).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &refreshToken, nil
}

func (r *UserRepository) RevokeRefreshToken(tokenID uuid.UUID) error {
	return r.db.Model(&models.RefreshToken{}).
		Where("id = ?", tokenID).
		Update("revoked", true).Error
}

func (r *UserRepository) RevokeAllUserRefreshTokens(userID uuid.UUID) error {
	return r.db.Model(&models.RefreshToken{}).
		Where("user_id = ?", userID).
		Update("revoked", true).Error
}

func (r *UserRepository) CreateEmailVerification(verification *models.EmailVerification) error {
	return r.db.Create(verification).Error
}

func (r *UserRepository) FindEmailVerificationByToken(token string) (*models.EmailVerification, error) {
	var verification models.EmailVerification
	err := r.db.Where("token = ? AND used_at IS NULL", token).First(&verification).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &verification, nil
}

func (r *UserRepository) MarkEmailAsVerified(userID uuid.UUID) error {
	now := time.Now()
	return r.db.Model(&models.User{}).
		Where("id = ?", userID).
		Updates(map[string]interface{}{
			"email_verified": true,
			"verified_at":    now,
		}).Error
}

func (r *UserRepository) MarkVerificatioTokenUsed(tokenID uuid.UUID) error {
	now := time.Now()
	return r.db.Model(&models.EmailVerification{}).Where("id = ?", tokenID).Update("used_at", now).Error
}

func (r *UserRepository) DeleteExpiredVerificationTokens() error {
	return r.db.Where("expires_at < ?", time.Now()).Delete(&models.EmailVerification{}).Error
}

func (r *UserRepository) CreatePasswordReset(reset *models.PasswordReset) error {
	return r.db.Create(reset).Error
}

func (r *UserRepository) FindPasswordResetByToken(token string) (*models.PasswordReset, error) {
	var reset models.PasswordReset
	err := r.db.Where("token = ? AND used_at IS NULL", token).First(&reset).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &reset, nil
}

func (r *UserRepository) MarkPasswordResetUsed(tokenID uuid.UUID) error {
	now := time.Now()
	return r.db.Model(&models.PasswordReset{}).
		Where("id = ?", tokenID).
		Update("used_at", now).Error
}

func (r *UserRepository) UpdatePassword(userID uuid.UUID, newPassword string) error {
	return r.db.Model(&models.User{}).Where("id = ?", userID).Update("password", newPassword).Error
}

func (r *UserRepository) DeleteExpiredPasswordResets() error {
	return r.db.Where("expires_at < ?", time.Now()).Delete(&models.PasswordReset{}).Error
}

func (r *UserRepository) EnableTwoFactor(userID uuid.UUID, secret string) error {
	return r.db.Model(&models.User{}).
		Where("id = ?", userID).
		Updates(map[string]interface{}{
			"two_factor_enabled": true,
			"two_factor_secret":  secret,
		}).Error
}

func (r *UserRepository) DisableTwoFactor(userID uuid.UUID) error {
	return r.db.Model(&models.User{}).
		Where("id = ?", userID).
		Updates(map[string]interface{}{
			"two_factor_enabled": false,
			"two_factor_secret":  nil,
		}).Error
}

func (r *UserRepository) VerifyPassword(user *models.User, password string) error {
	return bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
}

func (r *UserRepository) SaveBackupCodes(userID uuid.UUID, codes []models.BackupCode) error {
	return r.db.Create(&codes).Error
}

func (r *UserRepository) FindBackupCode(userID uuid.UUID, code string) (*models.BackupCode, error) {
	var backupCode models.BackupCode
	err := r.db.Where("user_id = ? AND code = ? AND used = ?", userID, code, false).
		First(&backupCode).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &backupCode, nil
}

func (r *UserRepository) MarkBackupCodeUsed(id uuid.UUID) error {
	return r.db.Model(&models.BackupCode{}).
		Where("id = ?", id).
		Update("used", true).Error
}

func (r *UserRepository) DeleteBackupCodes(userID uuid.UUID) error {
	return r.db.Where("user_id = ?", userID).Delete(&models.BackupCode{}).Error
}

func (r *UserRepository) GetUserByIDWith2FA(userID uuid.UUID) (*models.User, error) {
	var user models.User
	err := r.db.Preload("BackupCodes").First(&user, "id = ?", userID).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &user, nil
}

func (r *UserRepository) IsAccountLocked(userID uuid.UUID) (bool, *time.Time, error) {
	var user models.User
	err := r.db.Select("locked_until").First(&user, "id = ?", userID).Error
	if err != nil {
		return false, nil, err
	}

	if user.LockedUntil != nil && user.LockedUntil.After(time.Now()) {
		return true, user.LockedUntil, nil
	}

	return false, nil, nil
}

func (r *UserRepository) ResetLoginAttempts(userID uuid.UUID) error {
	return r.db.Model(&models.User{}).
		Where("id = ?", userID).
		Updates(map[string]interface{}{
			"login_attempts": 0,
			"locked_until":   nil,
		}).Error
}

func (r *UserRepository) CreateAccountLock(lock *models.AccountLock) error {
	return r.db.Create(lock).Error
}

func (r *UserRepository) UnlockAccount(userID uuid.UUID) error {
	return r.db.Model(&models.User{}).
		Where("id = ?", userID).
		Updates(map[string]interface{}{
			"login_attempts": 0,
			"locked_until":   nil,
		}).Error
}

func (r *UserRepository) GetAccountLocks(userID uuid.UUID, limit int) ([]models.AccountLock, error) {
	var locks []models.AccountLock
	err := r.db.Where("user_id = ?", userID).
		Order("created_at DESC").
		Limit(limit).
		Find(&locks).Error
	return locks, err
}

func (r *UserRepository) GetLockedAccounts() ([]models.User, error) {
	var users []models.User
	err := r.db.Where("locked_until > ?", time.Now()).
		Find(&users).Error
	return users, err
}

func (r *UserRepository) RecordLoginAttempt(userID uuid.UUID, success bool, ip string) error {
	user, err := r.FindByID(userID)
	if err != nil {
		return err
	}
	if user == nil {
		return nil
	}

	now := time.Now()

	if success {
		// Reset attempts on successful login
		return r.db.Model(&models.User{}).
			Where("id = ?", userID).
			Updates(map[string]interface{}{
				"login_attempts": 0,
				"locked_until":   nil,
				"last_login_at":  now,
				"last_login_ip":  ip,
			}).Error
	}

	// Increment failed attempts
	newAttempts := user.LoginAttempts + 1

	updateData := map[string]interface{}{
		"login_attempts":         newAttempts,
		"last_failed_attempt_at": now,
	}

	// Lock account if max attempts reached
	if newAttempts >= config.AppConfig.MaxLoginAttempts {
		lockedUntil := now.Add(time.Duration(config.AppConfig.LockoutDurationMinutes) * time.Minute)
		updateData["locked_until"] = lockedUntil
	}

	return r.db.Model(&models.User{}).
		Where("id = ?", userID).
		Updates(updateData).Error
}
