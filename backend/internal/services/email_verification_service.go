package services

import (
	"fmt"
	"time"

	"bezbase/internal/models"
	"bezbase/internal/pkg/contextx"
	"bezbase/internal/repository"

	"gorm.io/gorm"
)

type EmailVerificationService struct {
	userRepo         repository.UserRepository
	verificationRepo repository.EmailVerificationRepository
	emailService     *EmailService
}

func NewEmailVerificationService(
	userRepo repository.UserRepository,
	verificationRepo repository.EmailVerificationRepository,
	emailService *EmailService,
) *EmailVerificationService {
	return &EmailVerificationService{
		userRepo:         userRepo,
		verificationRepo: verificationRepo,
		emailService:     emailService,
	}
}

func (s *EmailVerificationService) SendVerificationEmail(ctx contextx.Contextx, userID uint) error {
	// Get user
	user, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		return fmt.Errorf("user not found: %w", err)
	}

	// Check if user already verified
	if user.EmailVerified {
		return fmt.Errorf("email already verified")
	}

	// Get user's email
	email := user.GetPrimaryEmail()
	if email == "" {
		return fmt.Errorf("user has no email address")
	}

	// Send verification email
	return s.emailService.SendVerificationEmail(ctx, user, email)
}

func (s *EmailVerificationService) VerifyEmail(ctx contextx.Contextx, token string) error {
	// Get verification token
	verificationToken, err := s.verificationRepo.GetByToken(token)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return fmt.Errorf("invalid or expired verification token")
		}
		return fmt.Errorf("failed to get verification token: %w", err)
	}

	// Check if token is valid
	if !verificationToken.IsValid() {
		return fmt.Errorf("verification token is expired or already used")
	}

	// Get user
	user, err := s.userRepo.GetByID(ctx, verificationToken.UserID)
	if err != nil {
		return fmt.Errorf("user not found: %w", err)
	}

	// Mark token as used
	if err := s.verificationRepo.MarkAsUsed(token); err != nil {
		return fmt.Errorf("failed to mark token as used: %w", err)
	}

	// Update user's email verification status
	user.EmailVerified = true
	if user.Status == models.UserStatusPending {
		user.Status = models.UserStatusActive
	}

	if err := s.userRepo.Update(ctx, user); err != nil {
		return fmt.Errorf("failed to update user: %w", err)
	}

	return nil
}

func (s *EmailVerificationService) ResendVerificationEmail(ctx contextx.Contextx, userID uint) error {
	// Get user
	user, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		return fmt.Errorf("user not found: %w", err)
	}

	// Check if user already verified
	if user.EmailVerified {
		return fmt.Errorf("email already verified")
	}

	// Check if there's a recent verification token (prevent spam)
	existingToken, err := s.verificationRepo.GetByUserID(userID)
	if err == nil && existingToken != nil {
		// Check if the token was created less than 1 minute ago
		if time.Since(existingToken.CreatedAt) < time.Minute {
			return fmt.Errorf("verification email was sent recently, please wait before requesting again")
		}
	}

	// Send new verification email
	return s.SendVerificationEmail(ctx, userID)
}

