package services

import (
	"bytes"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"html/template"
	"log"
	"net/smtp"
	"os"
	"time"

	"bezbase/internal/models"
	"bezbase/internal/pkg/contextx"
	"bezbase/internal/repository"
)

type EmailService struct {
	verificationRepo repository.EmailVerificationRepository
	smtpHost         string
	smtpPort         string
	smtpUsername     string
	smtpPassword     string
	fromEmail        string
	baseURL          string
}

func NewEmailService(verificationRepo repository.EmailVerificationRepository) *EmailService {
	return &EmailService{
		verificationRepo: verificationRepo,
		smtpHost:         getEnv("SMTP_HOST", "smtp.gmail.com"),
		smtpPort:         getEnv("SMTP_PORT", "587"),
		smtpUsername:     getEnv("SMTP_USERNAME", ""),
		smtpPassword:     getEnv("SMTP_PASSWORD", ""),
		fromEmail:        getEnv("FROM_EMAIL", "noreply@bezbase.com"),
		baseURL:          getEnv("BASE_URL", "http://localhost:3000"),
	}
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
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

	return s.sendEmail(email, subject, body)
}

func (s *EmailService) SendPasswordResetEmail(ctx contextx.Contextx, user *models.User, token string) error {
	subject := "Reset your password"
	resetURL := fmt.Sprintf("%s/reset-password?token=%s", s.baseURL, token)

	body, err := s.generatePasswordResetHTML(user.GetFullName(), resetURL)
	if err != nil {
		return fmt.Errorf("failed to generate email body: %w", err)
	}

	return s.sendEmail(user.GetPrimaryEmail(), subject, body)
}

func (s *EmailService) sendEmail(to, subject, body string) error {
	// Skip sending email if SMTP is not configured
	if s.smtpUsername == "" || s.smtpPassword == "" {
		log.Printf("Email would be sent to %s with subject: %s", to, subject)
		log.Printf("Email body: %s", body)
		return nil
	}

	// Setup headers
	headers := map[string]string{
		"From":         s.fromEmail,
		"To":           to,
		"Subject":      subject,
		"MIME-Version": "1.0",
		"Content-Type": "text/html; charset=UTF-8",
	}

	// Build message
	message := ""
	for k, v := range headers {
		message += fmt.Sprintf("%s: %s\r\n", k, v)
	}
	message += "\r\n" + body

	// Setup authentication
	auth := smtp.PlainAuth("", s.smtpUsername, s.smtpPassword, s.smtpHost)

	// Send email
	addr := fmt.Sprintf("%s:%s", s.smtpHost, s.smtpPort)
	return smtp.SendMail(addr, auth, s.fromEmail, []string{to}, []byte(message))
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

