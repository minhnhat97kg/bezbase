package dto

type PasswordResetRequest struct {
	Email string `json:"email" validate:"required,email"`
}

type ResetPasswordRequest struct {
	Token       string `json:"token" validate:"required"`
	NewPassword string `json:"new_password" validate:"required,min=8"`
}

type ValidateResetTokenRequest struct {
	Token string `json:"token" validate:"required"`
}

type PasswordResetResponse struct {
	Message string `json:"message"`
	Success bool   `json:"success"`
}