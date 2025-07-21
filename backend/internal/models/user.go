package models

import (
	"time"

	"gorm.io/gorm"
)

type UserStatus string

const (
	UserStatusActive    UserStatus = "active"
	UserStatusInactive  UserStatus = "inactive"
	UserStatusSuspended UserStatus = "suspended"
	UserStatusPending   UserStatus = "pending"
)

// User is the central entity that represents a user account
type User struct {
	ID           uint           `json:"id" gorm:"primaryKey"`
	Status       UserStatus     `json:"status" gorm:"not null;default:'pending'"`
	EmailVerified bool          `json:"email_verified" gorm:"default:false"`
	LastLoginAt  *time.Time     `json:"last_login_at"`
	CreatedAt    time.Time      `json:"created_at"`
	UpdatedAt    time.Time      `json:"updated_at"`
	DeletedAt    gorm.DeletedAt `json:"-" gorm:"index"`

	// Relationships
	UserInfo       *UserInfo           `json:"user_info,omitempty" gorm:"foreignKey:UserID"`
	AuthProviders  []AuthProvider      `json:"auth_providers,omitempty" gorm:"foreignKey:UserID"`
}

// GetPrimaryEmail returns the primary email from UserInfo
func (u *User) GetPrimaryEmail() string {
	if u.UserInfo != nil {
		return u.UserInfo.Email
	}
	return ""
}

// GetFullName returns the full name from UserInfo
func (u *User) GetFullName() string {
	if u.UserInfo != nil {
		return u.UserInfo.FirstName + " " + u.UserInfo.LastName
	}
	return ""
}

