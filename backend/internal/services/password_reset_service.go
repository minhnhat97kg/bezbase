package services

import (
	"fmt"
	"time"

	"bezbase/internal/models"
	"bezbase/internal/pkg/auth"
	"bezbase/internal/pkg/contextx"
	"bezbase/internal/repository"

	"gorm.io/gorm"
)

type PasswordResetService struct {
	userRepo          repository.UserRepository
	userInfoRepo      repository.UserInfoRepository
	authProviderRepo  repository.AuthProviderRepository
	passwordResetRepo repository.PasswordResetRepository
	emailService      *EmailService
}

func NewPasswordResetService(
	userRepo repository.UserRepository,
	userInfoRepo repository.UserInfoRepository,
	authProviderRepo repository.AuthProviderRepository,
	passwordResetRepo repository.PasswordResetRepository,
	emailService *EmailService,
) *PasswordResetService {
	return &PasswordResetService{
		userRepo:          userRepo,
		userInfoRepo:      userInfoRepo,
		authProviderRepo:  authProviderRepo,
		passwordResetRepo: passwordResetRepo,
		emailService:      emailService,
	}
}

func (s *PasswordResetService) RequestPasswordReset(ctx contextx.Contextx, email string) error {
	// Find user by email
	userInfo, err := s.userInfoRepo.GetByEmail(ctx, email)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			// Don't reveal if email exists or not for security
			return nil
		}
		return fmt.Errorf("failed to find user: %w", err)
	}

	// Get user details
	user, err := s.userRepo.GetByID(ctx, userInfo.UserID)
	if err != nil {
		return fmt.Errorf("failed to get user: %w", err)
	}

	// Check if there's a recent password reset token (prevent spam)
	existingToken, err := s.passwordResetRepo.GetByUserID(user.ID)
	if err == nil && existingToken != nil {
		// Check if the token was created less than 5 minutes ago
		if time.Since(existingToken.CreatedAt) < 5*time.Minute {
			return fmt.Errorf("password reset was requested recently, please wait before requesting again")
		}
	}

	// Generate password reset token
	token, err := generateSecureToken()
	if err != nil {
		return fmt.Errorf("failed to generate reset token: %w", err)
	}

	// Create password reset token record
	resetToken := &models.PasswordResetToken{
		UserID:    user.ID,
		Token:     token,
		Email:     email,
		ExpiresAt: time.Now().Add(1 * time.Hour), // 1 hour expiry
	}

	// Delete any existing tokens for this user
	s.passwordResetRepo.DeleteByUserID(user.ID)

	// Save new token
	if err := s.passwordResetRepo.Create(resetToken); err != nil {
		return fmt.Errorf("failed to save reset token: %w", err)
	}

	// Send password reset email
	if err := s.emailService.SendPasswordResetEmail(ctx, user, token); err != nil {
		return fmt.Errorf("failed to send password reset email: %w", err)
	}

	return nil
}

func (s *PasswordResetService) ResetPassword(ctx contextx.Contextx, token, newPassword string) error {
	// Get password reset token
	resetToken, err := s.passwordResetRepo.GetByToken(token)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return fmt.Errorf("invalid or expired reset token")
		}
		return fmt.Errorf("failed to get reset token: %w", err)
	}

	// Check if token is valid
	if !resetToken.IsValid() {
		return fmt.Errorf("password reset token is expired or already used")
	}

	// Get user
	user, err := s.userRepo.GetByID(ctx, resetToken.UserID)
	if err != nil {
		return fmt.Errorf("user not found: %w", err)
	}

	// Hash the new password
	hashedPassword, err := auth.HashPassword(newPassword)
	if err != nil {
		return fmt.Errorf("failed to hash password: %w", err)
	}

	// Update user's password in auth provider
	authProvider, err := s.authProviderRepo.GetByUserIDAndProvider(ctx, user.ID, models.ProviderEmail)
	if err != nil {
		return fmt.Errorf("failed to get auth provider: %w", err)
	}

	authProvider.Password = hashedPassword
	if err := s.authProviderRepo.Update(ctx, authProvider); err != nil {
		return fmt.Errorf("failed to update password: %w", err)
	}

	// Mark token as used
	if err := s.passwordResetRepo.MarkAsUsed(token); err != nil {
		return fmt.Errorf("failed to mark token as used: %w", err)
	}

	return nil
}

func (s *PasswordResetService) ValidateResetToken(ctx contextx.Contextx, token string) error {
	// Get password reset token
	resetToken, err := s.passwordResetRepo.GetByToken(token)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return fmt.Errorf("invalid or expired reset token")
		}
		return fmt.Errorf("failed to get reset token: %w", err)
	}

	// Check if token is valid
	if !resetToken.IsValid() {
		return fmt.Errorf("password reset token is expired or already used")
	}

	return nil
}
