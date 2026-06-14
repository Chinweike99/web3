package service

import (
	"auth-system/internal/config"
	"auth-system/internal/models"
	"auth-system/internal/repository"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"log"
	"time"

	"gopkg.in/gomail.v2"
)

type EmailService struct {
	repo   *repository.UserRepository
	dialer *gomail.Dialer
}

func NewEmailService(repo *repository.UserRepository) *EmailService {
	dialer := gomail.NewDialer(
		config.AppConfig.SMTPHost,
		config.AppConfig.SMTPPort,
		config.AppConfig.SMTPUser,
		config.AppConfig.SMTPPassword,
	)
	return &EmailService{
		repo:   repo,
		dialer: dialer,
	}
}

func (s *EmailService) GenerateVerificationToken() (string, error) {
	bytes := make([]byte, 32)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes), nil
}

func (s *EmailService) SendVerificationEmail(user *models.User) error {
	token, err := s.GenerateVerificationToken()
	if err != nil {
		return err
	}

	// store verification token
	verification := &models.EmailVerification{
		UserID:    user.ID,
		Token:     token,
		ExpiresAt: time.Now().Add(24 * time.Hour),
	}

	if err := s.repo.CreateEmailVerification(verification); err != nil {
		return err
	}

	// Send Email
	m := gomail.NewMessage()
	m.SetHeader("From", config.AppConfig.SMTPFrom)
	m.SetHeader("To", user.Email)
	m.SetHeader("Subject", "Verify your email Address")

	verificationURL := fmt.Sprintf("%s/api/v1/auth/verify-email?token=%s", config.AppConfig.AppURL, token)
	body := fmt.Sprintf(`
		<html>
		<body>
			 <h2>Welcome to My Learning golang journey!</h2>
            <p>Please verify your email address by clicking the link below:</p>
            <p><a href="%s">Verify Email</a></p>
            <p>This link will expire in 24 hours.</p>
            <p>If you didn't create an account, please ignore this email.</p>
        </body>
        </html>
	`, verificationURL)
	m.SetBody("text/html", body)
	if err := s.dialer.DialAndSend(m); err != nil {
		log.Printf("Failed to send email: %v", err)
		return err
	}
	return nil
}

func (s *EmailService) SendWecomeEmail(user *models.User) error {
	m := gomail.NewMessage()
	m.SetHeader("From", config.AppConfig.SMTPFrom)
	m.SetHeader("To", user.Email)
	m.SetHeader("Subject", "Welcome to My Learning golang journey!")

	body := fmt.Sprintf(`
        <html>
        <body>
            <h2>Welcome %s!</h2>
            <p>Your email has been successfully verified.</p>
            <p>You can now log in to your account.</p>
        </body>
        </html>
    `, user.FirstName)

	m.SetBody("text/html", body)
	return s.dialer.DialAndSend(m)
}

func (s *EmailService) SendPasswordResetEmail(user *models.User, token string) error {
	m := gomail.NewMessage()
	m.SetHeader("From", config.AppConfig.SMTPFrom)
	m.SetHeader("To", user.Email)
	m.SetHeader("Subject", "Password Reset Request")

	resetURL := fmt.Sprintf("%s/reset-password?token=%s",
		config.AppConfig.AppURL, token)

	body := fmt.Sprintf(`
        <html>
        <body>
            <h2>Password Reset Request</h2>
            <p>We received a request to reset your password. Click the link below to create a new password:</p>
            <p><a href="%s">Reset Password</a></p>
            <p>This link will expire in 1 hour.</p>
            <p>If you didn't request this, please ignore this email.</p>
            <hr>
            <p>For security reasons, never share this link with anyone.</p>
        </body>
        </html>
    `, resetURL)

	m.SetBody("text/html", body)

	if err := s.dialer.DialAndSend(m); err != nil {
		log.Printf("Failed to send password reset email: %v", err)
		return err
	}

	return nil
}

func (s *EmailService) SendPasswordChangedEmail(user *models.User) error {
	m := gomail.NewMessage()
	m.SetHeader("From", config.AppConfig.SMTPFrom)
	m.SetHeader("To", user.Email)
	m.SetHeader("Subject", "Your Password Has Been Changed")

	body := fmt.Sprintf(`
        <html>
        <body>
            <h2>Password Changed Successfully</h2>
            <p>Hello %s,</p>
            <p>Your password has been successfully changed.</p>
            <p>If you did not perform this action, please contact our support team immediately.</p>
            <br>
            <p>Best regards,</p>
            <p>Security Team</p>
        </body>
        </html>
    `, user.FirstName)

	m.SetBody("text/html", body)

	return s.dialer.DialAndSend(m)
}
