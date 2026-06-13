package service

import (
	"auth-system/internal/config"
	"auth-system/internal/models"
	"auth-system/internal/repository"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

type TokenService struct {
	repo *repository.UserRepository
}

// func (s *TokenService) GenerateRefreshToken(id uuid.UUID) (string, error) {
// 	panic("unimplemented")
// }

type Claims struct {
	UserID string `json:"user_id"`
	Email  string `json:"email"`
	Role   string `json:"role"`
	jwt.RegisteredClaims
}

func NewTokenService(repo *repository.UserRepository) *TokenService {
	return &TokenService{
		repo: repo,
	}
}

func (s *TokenService) GenerateAccessToken(user *models.User) (string, error) {
	expirationTime := time.Now().Add(time.Duration(config.AppConfig.JwtAccessExpiration) * time.Minute)
	claims := &Claims{
		UserID: user.ID.String(),
		Email:  user.Email,
		Role:   string(user.Role),
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Issuer:    "auth-system",
			Subject:   user.ID.String(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(config.AppConfig.JWTSecret))
	if err != nil {
		return "", err
	}
	return tokenString, nil
}

func (s *TokenService) GenerateRefreshToken(userID uuid.UUID) (string, *models.RefreshToken, error) {
	tokenString := uuid.New().String()
	expiresAt := time.Now().Add(time.Duration(config.AppConfig.JwtRefreshExpiration) * time.Hour)

	refreshToken := &models.RefreshToken{
		UserID:    userID,
		Token:     tokenString,
		ExpiresAt: expiresAt,
		Revoked:   false,
	}

	err := s.repo.CreateRefreshToken(refreshToken)
	if err != nil {
		return "", nil, err
	}

	return tokenString, refreshToken, nil
}

func (s *TokenService) ValidateAccessToken(tokenString string) (*Claims, error) {
	claims := &Claims{}
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(config.AppConfig.JWTSecret), nil
	})
	if err != nil {
		return nil, err
	}
	if !token.Valid {
		return nil, fmt.Errorf("invalid token")
	}
	return claims, nil
}

func (s *TokenService) ValidateRefreshToken(tokenString string) (*models.RefreshToken, error) {
	refreshToken, err := s.repo.FindRefreshToken(tokenString)
	if err != nil {
		return nil, err
	}
	if refreshToken == nil || refreshToken.Revoked || refreshToken.ExpiresAt.Before(time.Now()) {
		return nil, fmt.Errorf("invalid or expired refresh token")
	}
	return refreshToken, nil
}
