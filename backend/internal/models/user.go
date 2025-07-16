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
	ID            uint           `json:"id" gorm:"primaryKey"`
	Status        UserStatus     `json:"status" gorm:"not null;default:'pending'"`
	EmailVerified bool           `json:"email_verified" gorm:"default:false"`
	LastLoginAt   *time.Time     `json:"last_login_at"`
	CreatedAt     time.Time      `json:"created_at"`
	UpdatedAt     time.Time      `json:"updated_at"`
	DeletedAt     gorm.DeletedAt `json:"-" gorm:"index"`
	
	// Relationships
	UserInfo      *UserInfo      `json:"user_info,omitempty" gorm:"foreignKey:UserID"`
	AuthProviders []AuthProvider `json:"auth_providers,omitempty" gorm:"foreignKey:UserID"`
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

// DTOs for API requests/responses
type LoginRequest struct {
	Username string `json:"username" validate:"required"` // Can be username or email
	Password string `json:"password" validate:"required,min=6"`
}

type RegisterRequest struct {
	Email     string `json:"email" validate:"required,email"`
	Password  string `json:"password" validate:"required,min=6"`
	FirstName string `json:"first_name" validate:"required"`
	LastName  string `json:"last_name" validate:"required"`
}

type UpdateProfileRequest struct {
	FirstName   string `json:"first_name"`
	LastName    string `json:"last_name"`
	Bio         string `json:"bio"`
	Location    string `json:"location"`
	Website     string `json:"website"`
	Phone       string `json:"phone"`
	Gender      string `json:"gender"`
	Timezone    string `json:"timezone"`
	Language    string `json:"language"`
}

// Response DTOs
type UserResponse struct {
	ID            uint       `json:"id"`
	Status        UserStatus `json:"status"`
	EmailVerified bool       `json:"email_verified"`
	Email         string     `json:"email"`
	FirstName     string     `json:"first_name"`
	LastName      string     `json:"last_name"`
	Avatar        string     `json:"avatar"`
	Bio           string     `json:"bio"`
	Location      string     `json:"location"`
	Website       string     `json:"website"`
	CreatedAt     time.Time  `json:"created_at"`
	UpdatedAt     time.Time  `json:"updated_at"`
}

type AuthResponse struct {
	Token string       `json:"token"`
	User  UserResponse `json:"user"`
}

// Helper function to convert User + UserInfo to UserResponse
func ToUserResponse(user *User) UserResponse {
	response := UserResponse{
		ID:            user.ID,
		Status:        user.Status,
		EmailVerified: user.EmailVerified,
		CreatedAt:     user.CreatedAt,
		UpdatedAt:     user.UpdatedAt,
	}
	
	if user.UserInfo != nil {
		response.Email = user.UserInfo.Email
		response.FirstName = user.UserInfo.FirstName
		response.LastName = user.UserInfo.LastName
		response.Avatar = user.UserInfo.Avatar
		response.Bio = user.UserInfo.Bio
		response.Location = user.UserInfo.Location
		response.Website = user.UserInfo.Website
	}
	
	return response
}

