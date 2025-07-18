package services

import (
	"bytes"
	"context"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"html/template"
	"time"

	"bezbase/internal/config"
	"bezbase/internal/models"
	"bezbase/internal/pkg/contextx"
	"bezbase/internal/repository"
	"bezbase/internal/thirdparty/email"
)

type EmailService struct {
	verificationRepo repository.EmailVerificationRepository
	emailProvider    email.EmailProvider
	baseURL          string
}

func NewEmailService(verificationRepo repository.EmailVerificationRepository, cfg *config.EmailConfig, baseURL string) *EmailService {
	var emailProvider email.EmailProvider
	
	// Create provider based on configuration
	switch cfg.Provider {
	case "smtp":
		smtpConfig := &email.SMTPConfig{
			Host:      cfg.SMTPHost,
			Port:      cfg.SMTPPort,
			Username:  cfg.SMTPUsername,
			Password:  cfg.SMTPPassword,
			FromEmail: cfg.FromEmail,
		}
		emailProvider = email.NewSMTPProvider(smtpConfig)
	default:
		// Default to SMTP if provider not specified
		smtpConfig := &email.SMTPConfig{
			Host:      cfg.SMTPHost,
			Port:      cfg.SMTPPort,
			Username:  cfg.SMTPUsername,
			Password:  cfg.SMTPPassword,
			FromEmail: cfg.FromEmail,
		}
		emailProvider = email.NewSMTPProvider(smtpConfig)
	}
	
	return &EmailService{
		verificationRepo: verificationRepo,
		emailProvider:    emailProvider,
		baseURL:          baseURL,
	}
}


func (s *EmailService) SendVerificationEmail(ctx contextx.Contextx, user *models.User, email string) error {
	// Generate verification token
	token, err := generateSecureToken()
	if err != nil {
		return fmt.Errorf("failed to generate verification token: %w", err)
	}

	// Create verification token record
	verificationToken := &models.EmailVerificationToken{
		UserID:    user.ID,
		Token:     token,
		Email:     email,
		ExpiresAt: time.Now().Add(24 * time.Hour), // 24 hours expiry
	}

	// Delete any existing tokens for this user
	if err := s.verificationRepo.DeleteByUserID(user.ID); err != nil {
		return fmt.Errorf("failed to delete existing tokens: %w", err)
	}

	// Save new token
	if err := s.verificationRepo.Create(verificationToken); err != nil {
		return fmt.Errorf("failed to save verification token: %w", err)
	}

	// Send email
	subject := "Verify your email address"
	verificationURL := fmt.Sprintf("%s/verify-email?token=%s", s.baseURL, token)

	body, err := s.generateEmailVerificationHTML(user.GetFullName(), verificationURL)
	if err != nil {
		return fmt.Errorf("failed to generate email body: %w", err)
	}

	return s.emailProvider.SendEmail(context.Background(), email, subject, body)
}

func (s *EmailService) SendPasswordResetEmail(ctx contextx.Contextx, user *models.User, token string) error {
	subject := "Reset your password"
	resetURL := fmt.Sprintf("%s/reset-password?token=%s", s.baseURL, token)

	body, err := s.generatePasswordResetHTML(user.GetFullName(), resetURL)
	if err != nil {
		return fmt.Errorf("failed to generate email body: %w", err)
	}

	return s.emailProvider.SendEmail(context.Background(), user.GetPrimaryEmail(), subject, body)
}


func (s *EmailService) generateEmailVerificationHTML(name, verificationURL string) (string, error) {
	tmpl := `
<!DOCTYPE html>
<html>
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Verify Your Email</title>
    <style>
        body { font-family: Arial, sans-serif; line-height: 1.6; color: #333; }
        .container { max-width: 600px; margin: 0 auto; padding: 20px; }
        .header { background-color: #007bff; color: white; text-align: center; padding: 20px; }
        .content { padding: 30px; background-color: #f9f9f9; }
        .button { display: inline-block; padding: 12px 30px; background-color: #007bff; color: white; text-decoration: none; border-radius: 5px; margin: 20px 0; }
        .footer { text-align: center; margin-top: 30px; color: #666; font-size: 12px; }
    </style>
</head>
<body>
    <div class="container">
        <div class="header">
            <h1>BezBase</h1>
        </div>
        <div class="content">
            <h2>Welcome{{if .Name}}, {{.Name}}{{end}}!</h2>
            <p>Thank you for registering with BezBase. To complete your registration, please verify your email address by clicking the button below:</p>
            <a href="{{.VerificationURL}}" class="button">Verify Email Address</a>
            <p>If the button doesn't work, you can also copy and paste the following link into your browser:</p>
            <p style="word-break: break-all; color: #007bff;">{{.VerificationURL}}</p>
            <p>This verification link will expire in 24 hours.</p>
            <p>If you didn't create an account with BezBase, please ignore this email.</p>
        </div>
        <div class="footer">
            <p>© 2024 BezBase. All rights reserved.</p>
        </div>
    </div>
</body>
</html>`

	t, err := template.New("email").Parse(tmpl)
	if err != nil {
		return "", err
	}

	var buf bytes.Buffer
	data := struct {
		Name            string
		VerificationURL string
	}{
		Name:            name,
		VerificationURL: verificationURL,
	}

	if err := t.Execute(&buf, data); err != nil {
		return "", err
	}

	return buf.String(), nil
}

func (s *EmailService) generatePasswordResetHTML(name, resetURL string) (string, error) {
	tmpl := `
<!DOCTYPE html>
<html>
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Reset Your Password</title>
    <style>
        body { font-family: Arial, sans-serif; line-height: 1.6; color: #333; }
        .container { max-width: 600px; margin: 0 auto; padding: 20px; }
        .header { background-color: #007bff; color: white; text-align: center; padding: 20px; }
        .content { padding: 30px; background-color: #f9f9f9; }
        .button { display: inline-block; padding: 12px 30px; background-color: #007bff; color: white; text-decoration: none; border-radius: 5px; margin: 20px 0; }
        .footer { text-align: center; margin-top: 30px; color: #666; font-size: 12px; }
    </style>
</head>
<body>
    <div class="container">
        <div class="header">
            <h1>BezBase</h1>
        </div>
        <div class="content">
            <h2>Password Reset Request{{if .Name}} for {{.Name}}{{end}}</h2>
            <p>We received a request to reset the password for your BezBase account. Click the button below to reset your password:</p>
            <a href="{{.ResetURL}}" class="button">Reset Password</a>
            <p>If the button doesn't work, you can also copy and paste the following link into your browser:</p>
            <p style="word-break: break-all; color: #007bff;">{{.ResetURL}}</p>
            <p>This password reset link will expire in 1 hour.</p>
            <p>If you didn't request a password reset, please ignore this email. Your password will remain unchanged.</p>
        </div>
        <div class="footer">
            <p>© 2024 BezBase. All rights reserved.</p>
        </div>
    </div>
</body>
</html>`

	t, err := template.New("email").Parse(tmpl)
	if err != nil {
		return "", err
	}

	var buf bytes.Buffer
	data := struct {
		Name     string
		ResetURL string
	}{
		Name:     name,
		ResetURL: resetURL,
	}

	if err := t.Execute(&buf, data); err != nil {
		return "", err
	}

	return buf.String(), nil
}

func generateSecureToken() (string, error) {
	bytes := make([]byte, 32)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes), nil
}

