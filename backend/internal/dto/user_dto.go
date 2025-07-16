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

type UserResponse struct {
	ID            uint              `json:"id"`
	Status        string            `json:"status"`
	EmailVerified bool              `json:"email_verified"`
	LastLoginAt   *time.Time        `json:"last_login_at,omitempty"`
	CreatedAt     time.Time         `json:"created_at"`
	UpdatedAt     time.Time         `json:"updated_at"`
	UserInfo      *UserInfoResponse `json:"user_info,omitempty"`
}

type UserInfoResponse struct {
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Email     string `json:"email"`
	AvatarURL string `json:"avatar_url,omitempty"`
	Language  string `json:"language,omitempty"`
	Timezone  string `json:"timezone,omitempty"`
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
		resp.UserInfo = &UserInfoResponse{
			FirstName: user.UserInfo.FirstName,
			LastName:  user.UserInfo.LastName,
			Email:     user.UserInfo.Email,
			AvatarURL: user.UserInfo.AvatarURL,
			Language:  user.UserInfo.Language,
			Timezone:  user.UserInfo.Timezone,
		}
	}
	return resp
}
