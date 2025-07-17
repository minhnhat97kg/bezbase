package services

import (
	"errors"
	"time"

	"bezbase/internal/dto"
	"bezbase/internal/models"
	"bezbase/internal/pkg/auth"
	"bezbase/internal/repository"

	"gorm.io/gorm"
)

type AuthService struct {
	userRepo         repository.UserRepository
	userInfoRepo     repository.UserInfoRepository
	authProviderRepo repository.AuthProviderRepository
	jwtSecret        string
	db               *gorm.DB
}

func NewAuthService(
	userRepo repository.UserRepository,
	userInfoRepo repository.UserInfoRepository,
	authProviderRepo repository.AuthProviderRepository,
	jwtSecret string,
	db *gorm.DB,
) *AuthService {
	return &AuthService{
		userRepo:         userRepo,
		userInfoRepo:     userInfoRepo,
		authProviderRepo: authProviderRepo,
		jwtSecret:        jwtSecret,
		db:               db,
	}
}

// Register creates a new user with email/password authentication
func (s *AuthService) Register(req dto.RegisterRequest) (*dto.AuthResponse, error) {
	tx := s.db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// Check if username is already taken
	var existingUserInfo models.UserInfo
	if err := tx.Where("username = ?", req.Username).First(&existingUserInfo).Error; err == nil {
		tx.Rollback()
		return nil, errors.New("username already taken")
	}

	// Check if email is already registered
	if err := tx.Where("email = ?", req.Email).First(&existingUserInfo).Error; err == nil {
		tx.Rollback()
		return nil, errors.New("email already registered")
	}

	// Check if username is already registered in auth providers
	var existingProvider models.AuthProvider
	if err := tx.Where("user_name = ? AND provider = ?", req.Username, models.ProviderEmail).First(&existingProvider).Error; err == nil {
		tx.Rollback()
		return nil, errors.New("username already registered")
	}

	// Hash password
	hashedPassword, err := auth.HashPassword(req.Password)
	if err != nil {
		tx.Rollback()
		return nil, errors.New("failed to hash password")
	}

	// Create user
	user := models.User{
		Status:        models.UserStatusPending,
		EmailVerified: false,
	}
	if err := tx.Create(&user).Error; err != nil {
		tx.Rollback()
		return nil, errors.New("failed to create user")
	}

	// Create user info
	userInfo := models.UserInfo{
		UserID:    user.ID,
		Username:  req.Username,
		FirstName: req.FirstName,
		LastName:  req.LastName,
		Email:     req.Email,
		Language:  "en",
		Timezone:  "UTC",
	}
	if err := tx.Create(&userInfo).Error; err != nil {
		tx.Rollback()
		return nil, errors.New("failed to create user info")
	}

	// Create auth provider
	authProvider := models.AuthProvider{
		UserID:     user.ID,
		Provider:   models.ProviderEmail,
		ProviderID: req.Email, // For email provider, provider_id is the email
		UserName:   req.Username, // Use username for login
		Password:   hashedPassword,
		Verified:   false,
	}
	if err := tx.Create(&authProvider).Error; err != nil {
		tx.Rollback()
		return nil, errors.New("failed to create auth provider")
	}

	// Commit transaction
	if err := tx.Commit().Error; err != nil {
		return nil, errors.New("failed to save user data")
	}

	// Load relationships for response
	user.UserInfo = &userInfo

	// Generate JWT token
	token, err := auth.GenerateToken(user.ID, user.UserInfo.Username, s.jwtSecret)
	if err != nil {
		return nil, errors.New("failed to generate token")
	}

	return &dto.AuthResponse{
		Token: token,
		User:  dto.ToUserResponse(&user),
	}, nil
}

// LoginWithUsername authenticates user with username and password
func (s *AuthService) LoginWithUsername(req dto.LoginRequest) (*dto.AuthResponse, error) {
	// Find auth provider by username
	var authProvider models.AuthProvider
	if err := s.db.Where("user_name = ? AND provider = ?", req.Username, models.ProviderEmail).First(&authProvider).Error; err != nil {
		return nil, errors.New("invalid credentials")
	}

	// Check password
	if !auth.CheckPasswordHash(req.Password, authProvider.Password) {
		return nil, errors.New("invalid credentials")
	}

	// Get user with info
	var user models.User
	if err := s.db.Preload("UserInfo").First(&user, authProvider.UserID).Error; err != nil {
		return nil, errors.New("user not found")
	}

	// Update last login time
	now := time.Now()
	user.LastLoginAt = &now
	s.db.Save(&user)

	// Generate JWT token
	token, err := auth.GenerateToken(user.ID, user.UserInfo.Username, s.jwtSecret)
	if err != nil {
		return nil, errors.New("failed to generate token")
	}

	return &dto.AuthResponse{
		Token: token,
		User:  dto.ToUserResponse(&user),
	}, nil
}

// RegisterWithSocialProvider creates a user from social login (future implementation)
func (s *AuthService) RegisterWithSocialProvider(provider models.AuthProviderType, providerID, email, firstName, lastName string) (*dto.AuthResponse, error) {
	tx := s.db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// Check if this social account is already linked
	var existingProvider models.AuthProvider
	if err := tx.Where("provider_id = ? AND provider = ?", providerID, provider).First(&existingProvider).Error; err == nil {
		// User already exists, just login
		tx.Rollback()
		return s.loginWithSocialProvider(provider, providerID)
	}

	// Create new user
	user := models.User{
		Status:        models.UserStatusActive, // Social logins are typically pre-verified
		EmailVerified: true,
	}
	if err := tx.Create(&user).Error; err != nil {
		tx.Rollback()
		return nil, errors.New("failed to create user")
	}

	// Create user info
	userInfo := models.UserInfo{
		UserID:    user.ID,
		FirstName: firstName,
		LastName:  lastName,
		Email:     email,
		Language:  "en",
		Timezone:  "UTC",
	}
	if err := tx.Create(&userInfo).Error; err != nil {
		tx.Rollback()
		return nil, errors.New("failed to create user info")
	}

	// Create auth provider
	authProvider := models.AuthProvider{
		UserID:     user.ID,
		Provider:   provider,
		ProviderID: providerID,
		UserName:   email, // For social providers, use email as username
		Verified:   true,  // Social logins are typically pre-verified
	}
	if err := tx.Create(&authProvider).Error; err != nil {
		tx.Rollback()
		return nil, errors.New("failed to create auth provider")
	}

	// Commit transaction
	if err := tx.Commit().Error; err != nil {
		return nil, errors.New("failed to save user data")
	}

	// Load relationships for response
	user.UserInfo = &userInfo

	// Generate JWT token
	token, err := auth.GenerateToken(user.ID, user.UserInfo.Username, s.jwtSecret)
	if err != nil {
		return nil, errors.New("failed to generate token")
	}

	return &dto.AuthResponse{
		Token: token,
		User:  dto.ToUserResponse(&user),
	}, nil
}

// loginWithSocialProvider handles login for existing social accounts
func (s *AuthService) loginWithSocialProvider(provider models.AuthProviderType, providerID string) (*dto.AuthResponse, error) {
	// Find auth provider
	var authProvider models.AuthProvider
	if err := s.db.Where("provider_id = ? AND provider = ?", providerID, provider).First(&authProvider).Error; err != nil {
		return nil, errors.New("social account not found")
	}

	// Get user with info
	var user models.User
	if err := s.db.Preload("UserInfo").First(&user, authProvider.UserID).Error; err != nil {
		return nil, errors.New("user not found")
	}

	// Update last login time
	now := time.Now()
	user.LastLoginAt = &now
	s.db.Save(&user)

	// Generate JWT token
	token, err := auth.GenerateToken(user.ID, user.UserInfo.Email, s.jwtSecret)
	if err != nil {
		return nil, errors.New("failed to generate token")
	}

	return &dto.AuthResponse{
		Token: token,
		User:  dto.ToUserResponse(&user),
	}, nil
}

// LinkSocialProvider links a social account to existing user
func (s *AuthService) LinkSocialProvider(userID uint, provider models.AuthProviderType, providerID, email string) error {
	// Check if this social account is already linked to another user
	var existingProvider models.AuthProvider
	if err := s.db.Where("provider_id = ? AND provider = ?", providerID, provider).First(&existingProvider).Error; err == nil {
		return errors.New("social account already linked to another user")
	}

	// Create new auth provider link
	authProvider := models.AuthProvider{
		UserID:     userID,
		Provider:   provider,
		ProviderID: providerID,
		UserName:   email, // For social providers, use email as username
		Verified:   true,
	}

	if err := s.db.Create(&authProvider).Error; err != nil {
		return errors.New("failed to link social account")
	}

	return nil
}
