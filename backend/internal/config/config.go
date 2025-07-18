package config

import (
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	DatabaseURL string
	JWTSecret   string
	Port        string
}

func Load() *Config {
	// Try to load .env file, ignore error if not found
	_ = godotenv.Load()
	return &Config{
		DatabaseURL: getEnvOrDefault("DATABASE_URL", "postgres://user:password@localhost/bezbase?sslmode=disable"),
		JWTSecret:   getEnvOrDefault("JWT_SECRET", "your-secret-key-change-this-in-production"),
		Port:        getEnvOrDefault("PORT", "8080"),
	}
}

func getEnvOrDefault(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
