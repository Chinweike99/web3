package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)



type Role string

const (
	RoleUser  Role = "user"
	RoleAdmin Role = "admin"
	RoleModerator Role = "moderator"
)

type User struct {
	ID uuid.UUID  `gorm: "type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	Email string `gorm:"uniqueIndex;not null" json:"email"`
	Password string `gorm:"not null" json:"-"` // - means exclude from JSON
	FirstName string `gorm: "not null": json:"first_name"`
	LastName string `gorm: "not null": json:"last_name"`
	Role Role`gorm:"type:varchar(20);default:'user'" json:"role"`
	IsActive bool `gorm:"default:true" json:"is_active"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

type RefreshToken struct {
	ID uuid.UUID `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	UserID uuid.UUID `gorm:"type:uuid;not null" json:"user_id"`
	Token string `gorm:uniqueIndex;not null" json:"token"`
	ExpiresAt time.Time `gorm:"not null" json:"expires_at"`
	CreatedAt time.Time `json:"created_at"`
	Revoked bool `gorm:"default:false" json:"revoked"`
}

type LoginRequest struct {
	Email string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6,max=20"`
}

type RegisterRequest struct {
	Email string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6,max=20"`
	FirstName string `json:"first_name" binding:"required"`
	LastName string `json:"last_name" binding:"required"`
}

type TokenResponse struct {
	AccessToken string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	TokenType string `json:"token_type"`
	ExpiresIn int `json:"expires_in"` // Access token expiration in seconds
}