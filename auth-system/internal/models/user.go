package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Role string

const (
	RoleUser      Role = "user"
	RoleAdmin     Role = "admin"
	RoleModerator Role = "moderator"
)

type User struct {
	ID            uuid.UUID      `gorm:"type:uuid;primaryKey;default:gen_random_uuid()" json:"id"`
	Email         string         `gorm:"uniqueIndex;not null" json:"email"`
	Password      string         `gorm:"not null" json:"-"`
	FirstName     string         `gorm:"not null" json:"first_name"`
	LastName      string         `gorm:"not null" json:"last_name"`
	Role          Role           `gorm:"type:varchar(20);default:'user'" json:"role"`
	IsActive      bool           `gorm:"default:true" json:"is_active"`
	CreatedAt     time.Time      `json:"created_at"`
	UpdatedAt     time.Time      `json:"updated_at"`
	DeletedAt     gorm.DeletedAt `gorm:"index" json:"-"`
	EmailVerified bool           `gorm:"default:false" json:"email_verified"`
	VerifiedAt    *time.Time     `json:"verified_at,omitempty"`
	// VerificationToken string         `gorm:"uniqueIndex" json:"_"`
	TwoFactorEnabled bool   `gorm:"default:false" json:"two_factor_enabled"`
    TwoFactorSecret  string `gorm:"type:varchar(255)" json:"-"`
    BackupCodes      []BackupCode `gorm:"foreignKey:UserID" json:"-"`
}

type BackupCode struct {
    ID        uuid.UUID `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
    UserID    uuid.UUID `gorm:"type:uuid;not null" json:"user_id"`
    Code      string    `gorm:"uniqueIndex;not null" json:"code"`
    Used      bool      `gorm:"default:false" json:"used"`
    CreatedAt time.Time `json:"created_at"`
}

// // Error implements [error].
// func (u *User) Error() string {
// 	panic("unimplemented")
// }

type EmailVerification struct {
	ID        uuid.UUID  `gorm:"type:uuid;primaryKey;default:gen_random_uuid()" json:"id"`
	UserID    uuid.UUID  `gorm:"type:uuid;not null" json:"user_id"`
	Token     string     `gorm:"uniqueIndex;not null" json:"token"`
	ExpiresAt time.Time  `gorm:"not null" json:"expires_at"`
	CreatedAt time.Time  `json:"created_at"`
	UsedAt    *time.Time `json:"used_at,omitempty"`
}

type RefreshToken struct {
	ID        uuid.UUID `gorm:"type:uuid;primaryKey;default:gen_random_uuid()" json:"id"`
	UserID    uuid.UUID `gorm:"type:uuid;not null" json:"user_id"`
	Token     string    `gorm:"uniqueIndex;not null" json:"token"`
	ExpiresAt time.Time `gorm:"not null" json:"expires_at"`
	CreatedAt time.Time `json:"created_at"`
	Revoked   bool      `gorm:"default:false" json:"revoked"`
}

type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6,max=20"`
}

type RegisterRequest struct {
	Email     string `json:"email" binding:"required,email"`
	Password  string `json:"password" binding:"required,min=6,max=20"`
	FirstName string `json:"first_name" binding:"required"`
	LastName  string `json:"last_name" binding:"required"`
}

type TokenResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	TokenType    string `json:"token_type"`
	ExpiresIn    int    `json:"expires_in"` // Access token expiration in seconds
}

type RefreshRequest struct {
	RefreshToken string `json:"refresh_token" binding:"required"`
}

type PasswordReset struct {
	ID        uuid.UUID  `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	UserID    uuid.UUID  `gorm:"type:uuid;not null" json:"user_id"`
	ExpiresAt time.Time  `gorm:"not null" json:"expires_at"`
	Token     string     `gorm:"uniqueIndex;not null" json:"token"`
	UsedAt    *time.Time `json:"used_at,omitempty"`
}

type ForgotPasswordRequest struct {
	Email string `json:"email" binding:"required,email"`
}

type ResetPasswordRequest struct {
	Token           string `json:"token" binding:"required"`
	NewPassword     string `json:"new_password" binding:"required,min=6"`
	ConfirmPassword string `json:"confirm_password" binding:"required,min=6"`
}

type ChangePasswordRequest struct {
	CurrentPassword string `json:"current_password" binding:"required"`
	NewPassword     string `json:"new_password" binding:"required,min=6"`
	ConfirmPassword string `json:"confirm_password" binding:"required,min=6"`
}

type TwoFactorSetupRequest struct {
    Password string `json:"password" binding:"required"`
}

type TwoFactorVerifyRequest struct {
    Token string `json:"token" binding:"required"`
}

type TwoFactorLoginRequest struct {
    Email    string `json:"email" binding:"required,email"`
    Password string `json:"password" binding:"required"`
    Token    string `json:"token" binding:"required"`
}

type BackupCodeRequest struct {
    Code string `json:"code" binding:"required"`
}

type TwoFactorDisableRequest struct {
    Password string `json:"password" binding:"required"`
    Token    string `json:"token" binding:"required"`
}