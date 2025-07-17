package models

import (
	"time"

	"gorm.io/gorm"
)

type EmailVerificationToken struct {
	ID        uint           `json:"id" gorm:"primaryKey"`
	UserID    uint           `json:"user_id" gorm:"not null;index"`
	Token     string         `json:"token" gorm:"not null;uniqueIndex"`
	Email     string         `json:"email" gorm:"not null"`
	ExpiresAt time.Time      `json:"expires_at" gorm:"not null"`
	UsedAt    *time.Time     `json:"used_at"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`

	// Relationships
	User User `json:"-" gorm:"foreignKey:UserID"`
}

func (EmailVerificationToken) TableName() string {
	return "email_verification_tokens"
}

// IsExpired checks if the token has expired
func (evt *EmailVerificationToken) IsExpired() bool {
	return time.Now().After(evt.ExpiresAt)
}

// IsUsed checks if the token has been used
func (evt *EmailVerificationToken) IsUsed() bool {
	return evt.UsedAt != nil
}

// IsValid checks if the token is valid (not expired and not used)
func (evt *EmailVerificationToken) IsValid() bool {
	return !evt.IsExpired() && !evt.IsUsed()
}