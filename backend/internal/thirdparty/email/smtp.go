package email

import (
	"context"
	"fmt"
	"log"
	"net/smtp"
)

// SMTPConfig contains SMTP-specific configuration
type SMTPConfig struct {
	Host     string
	Port     string
	Username string
	Password string
	FromEmail string
}

// SMTPProvider implements EmailProvider using SMTP
type SMTPProvider struct {
	config *SMTPConfig
}

// NewSMTPProvider creates a new SMTP email provider
func NewSMTPProvider(config *SMTPConfig) *SMTPProvider {
	return &SMTPProvider{
		config: config,
	}
}

// SendEmail sends an email using SMTP
func (s *SMTPProvider) SendEmail(ctx context.Context, to, subject, body string) error {
	// Skip sending email if SMTP is not configured
	if !s.IsConfigured() {
		log.Printf("Email would be sent to %s with subject: %s", to, subject)
		log.Printf("Email body: %s", body)
		return nil
	}

	// Setup headers
	headers := map[string]string{
		"From":         s.config.FromEmail,
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
	auth := smtp.PlainAuth("", s.config.Username, s.config.Password, s.config.Host)

	// Send email
	addr := fmt.Sprintf("%s:%s", s.config.Host, s.config.Port)
	return smtp.SendMail(addr, auth, s.config.FromEmail, []string{to}, []byte(message))
}

// IsConfigured checks if SMTP is properly configured
func (s *SMTPProvider) IsConfigured() bool {
	return s.config.Username != "" && s.config.Password != ""
}