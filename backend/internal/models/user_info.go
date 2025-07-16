package models

import (
	"time"

	"gorm.io/gorm"
)

type UserInfo struct {
	ID          uint           `json:"id" gorm:"primaryKey"`
	UserID      uint           `json:"user_id" gorm:"not null;uniqueIndex"` // One-to-one with User
	FirstName   string         `json:"first_name" gorm:"not null"`
	LastName    string         `json:"last_name" gorm:"not null"`
	Email       string         `json:"email" gorm:"not null;uniqueIndex"` // Primary email for display
	Avatar      string         `json:"avatar" gorm:""`                     // Profile picture URL
	Bio         string         `json:"bio" gorm:""`                        // User biography
	Location    string         `json:"location" gorm:""`                   // User location
	Website     string         `json:"website" gorm:""`                    // Personal website
	Phone       string         `json:"phone" gorm:""`                      // Phone number
	DateOfBirth *time.Time     `json:"date_of_birth" gorm:""`             // Birth date
	Gender      string         `json:"gender" gorm:""`                     // Gender
	Timezone    string         `json:"timezone" gorm:"default:'UTC'"`      // User timezone
	Language    string         `json:"language" gorm:"default:'en'"`       // Preferred language
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `json:"-" gorm:"index"`
	
	// Relationships
	User User `json:"-" gorm:"foreignKey:UserID"`
}

func (UserInfo) TableName() string {
	return "user_info"
}