package config

import (
	"os"

	"github.com/joho/godotenv"
)

// EmailConfig contains email service configuration
type EmailConfig struct {
	SMTPHost     string
	SMTPPort     string
	SMTPUsername string
	SMTPPassword string
	FromEmail    string
	Provider     string // "smtp", "sendgrid", "mailgun", etc.
}

// DatabaseConfig contains database configuration
type DatabaseConfig struct {
	URL string
}

// AuthConfig contains authentication configuration
type AuthConfig struct {
	JWTSecret string
}

// ServerConfig contains server configuration
type ServerConfig struct {
	Port    string
	BaseURL string
}

// Config is the main configuration struct containing all service configs
type Config struct {
	Database DatabaseConfig
	Auth     AuthConfig
	Server   ServerConfig
	Email    EmailConfig
}

func Load() *Config {
	// Try to load .env file, ignore error if not found
	_ = godotenv.Load()
	return &Config{
		Database: DatabaseConfig{
			URL: getEnvOrDefault("DATABASE_URL", "postgres://user:password@localhost/bezbase?sslmode=disable"),
		},
		Auth: AuthConfig{
			JWTSecret: getEnvOrDefault("JWT_SECRET", "your-secret-key-change-this-in-production"),
		},
		Server: ServerConfig{
			Port:    getEnvOrDefault("PORT", "8080"),
			BaseURL: getEnvOrDefault("BASE_URL", "http://localhost:3000"),
		},
		Email: EmailConfig{
			SMTPHost:     getEnvOrDefault("SMTP_HOST", "smtp.gmail.com"),
			SMTPPort:     getEnvOrDefault("SMTP_PORT", "587"),
			SMTPUsername: getEnvOrDefault("SMTP_USERNAME", ""),
			SMTPPassword: getEnvOrDefault("SMTP_PASSWORD", ""),
			FromEmail:    getEnvOrDefault("FROM_EMAIL", "noreply@bezbase.com"),
			Provider:     getEnvOrDefault("EMAIL_PROVIDER", "smtp"),
		},
	}
}

func getEnvOrDefault(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
