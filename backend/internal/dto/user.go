package dto

import (
	"time"

	"bezbase/internal/models"
)

type UpdateProfileRequest struct {
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	AvatarURL string `json:"avatar_url"`
	Language  string `json:"language"`
	Timezone  string `json:"timezone"`
}

type CreateUserRequest struct {
	FirstName string `json:"first_name" binding:"required"`
	LastName  string `json:"last_name" binding:"required"`
	Email     string `json:"email" binding:"required,email"`
	Password  string `json:"password" binding:"required,min=8"`
	Status    string `json:"status" binding:"required,oneof=active inactive suspended pending"`
	Language  string `json:"language"`
	Timezone  string `json:"timezone"`
	Bio       string `json:"bio"`
	Location  string `json:"location"`
	Website   string `json:"website"`
	Phone     string `json:"phone"`
}

type UpdateUserRequest struct {
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Email     string `json:"email"`
	Status    string `json:"status" binding:"omitempty,oneof=active inactive suspended pending"`
	Language  string `json:"language"`
	Timezone  string `json:"timezone"`
	Bio       string `json:"bio"`
	Location  string `json:"location"`
	Website   string `json:"website"`
	Phone     string `json:"phone"`
}

type UserResponse struct {
	ID            uint       `json:"id"`
	Status        string     `json:"status"`
	EmailVerified bool       `json:"email_verified"`
	LastLoginAt   *time.Time `json:"last_login_at,omitempty"`
	CreatedAt     time.Time  `json:"created_at"`
	UpdatedAt     time.Time  `json:"updated_at"`
	FirstName     string     `json:"first_name"`
	LastName      string     `json:"last_name"`
	Email         string     `json:"email"`
	AvatarURL     string     `json:"avatar_url,omitempty"`
	Language      string     `json:"language,omitempty"`
	Timezone      string     `json:"timezone,omitempty"`
	Bio           string     `json:"bio,omitempty"`
	Location      string     `json:"location,omitempty"`
	Website       string     `json:"website,omitempty"`
	Phone         string     `json:"phone,omitempty"`
	Roles         []string   `json:"roles,omitempty"`
}

// ToUserResponse converts a User model to a UserResponse DTO
func ToUserResponse(user *models.User) UserResponse {
	resp := UserResponse{
		ID:            user.ID,
		Status:        string(user.Status),
		EmailVerified: user.EmailVerified,
		LastLoginAt:   user.LastLoginAt,
		CreatedAt:     user.CreatedAt,
		UpdatedAt:     user.UpdatedAt,
	}
	
	if user.UserInfo != nil {
		resp.FirstName = user.UserInfo.FirstName
		resp.LastName = user.UserInfo.LastName
		resp.Email = user.UserInfo.Email
		resp.AvatarURL = user.UserInfo.AvatarURL
		resp.Language = user.UserInfo.Language
		resp.Timezone = user.UserInfo.Timezone
		resp.Bio = user.UserInfo.Bio
		resp.Location = user.UserInfo.Location
		resp.Website = user.UserInfo.Website
		resp.Phone = user.UserInfo.Phone
	}
	
	return resp
}

// ToUserResponseWithRoles converts a User model to a UserResponse DTO with roles
func ToUserResponseWithRoles(user *models.User, roles []string) UserResponse {
	resp := ToUserResponse(user)
	resp.Roles = roles
	return resp
}
