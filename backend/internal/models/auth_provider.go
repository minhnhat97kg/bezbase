package models

import (
	"time"

	"gorm.io/gorm"
)

type AuthProviderType string

const (
	ProviderEmail    AuthProviderType = "email"
	ProviderGoogle   AuthProviderType = "google"
	ProviderFacebook AuthProviderType = "facebook"
	ProviderGithub   AuthProviderType = "github"
	ProviderApple    AuthProviderType = "apple"
)

type AuthProvider struct {
	ID         uint             `json:"id" gorm:"primaryKey"`
	UserID     uint             `json:"user_id" gorm:"not null;index"`
	Provider   AuthProviderType `json:"provider" gorm:"not null;index"`
	ProviderID string           `json:"provider_id" gorm:"not null"`     // External provider's user ID
	UserName   string           `json:"user_name" gorm:"not null;index"` // Username or email for authentication
	Password   string           `json:"-" gorm:""`                       // Only for email provider
	Verified   bool             `json:"verified" gorm:"default:false"`
	CreatedAt  time.Time        `json:"created_at"`
	UpdatedAt  time.Time        `json:"updated_at"`
	DeletedAt  gorm.DeletedAt   `json:"-" gorm:"index"`

	// Relationships
	User User `json:"-" gorm:"foreignKey:UserID"`
}

// For composite unique constraint: one provider per user per provider type
func (AuthProvider) TableName() string {
	return "auth_providers"
}

// Index for unique constraint: user_id + provider should be unique
// This allows one account per provider per user
