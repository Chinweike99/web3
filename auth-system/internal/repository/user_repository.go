package repository

import (
	"auth-system/internal/database"
	"auth-system/internal/models"
	"errors"
	"time"

	"github.com/google/uuid"
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