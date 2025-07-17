package dto

type SendVerificationEmailRequest struct {
	UserID uint `json:"user_id" validate:"required"`
}

type VerifyEmailRequest struct {
	Token string `json:"token" validate:"required"`
}

type ResendVerificationEmailRequest struct {
	UserID uint `json:"user_id" validate:"required"`
}

type EmailVerificationResponse struct {
	Message string `json:"message"`
	Success bool   `json:"success"`
}