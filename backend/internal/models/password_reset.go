package models

import (
	"time"

	"gorm.io/gorm"
)

type PasswordResetToken struct {
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

func (PasswordResetToken) TableName() string {
	return "password_reset_tokens"
}

// IsExpired checks if the token has expired
func (prt *PasswordResetToken) IsExpired() bool {
	return time.Now().After(prt.ExpiresAt)
}

// IsUsed checks if the token has been used
func (prt *PasswordResetToken) IsUsed() bool {
	return prt.UsedAt != nil
}

// IsValid checks if the token is valid (not expired and not used)
func (prt *PasswordResetToken) IsValid() bool {
	return !prt.IsExpired() && !prt.IsUsed()
}