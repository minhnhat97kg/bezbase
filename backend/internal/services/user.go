package services

import (
	"errors"

	"bezbase/internal/dto"
	"bezbase/internal/models"
	"bezbase/internal/pkg/contextx"
	"bezbase/internal/repository"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type UserService struct {
	userRepo         repository.UserRepository
	userInfoRepo     repository.UserInfoRepository
	authProviderRepo repository.AuthProviderRepository
	rbacService      *RBACService
	db               *gorm.DB
}

func NewUserService(
	userRepo repository.UserRepository,
	userInfoRepo repository.UserInfoRepository,
	authProviderRepo repository.AuthProviderRepository,
	rbacService *RBACService,
	db *gorm.DB,
) *UserService {
	return &UserService{
		userRepo:         userRepo,
		userInfoRepo:     userInfoRepo,
		authProviderRepo: authProviderRepo,
		rbacService:      rbacService,
		db:               db,
	}
}

// GetProfile retrieves user profile with all user information
func (s *UserService) GetProfile(ctx contextx.Contextx, userID uint) (*dto.UserResponse, error) {
	user, err := s.userRepo.GetByIDWithPreload(ctx, userID, "UserInfo")
	if err != nil {
		return nil, err
	}

	var roles []string
	if s.rbacService != nil {
		userRoles, err := s.rbacService.GetUserRoles(user.ID)
		if err == nil {
			roles = userRoles
		}
	}

	response := dto.ToUserResponseWithRoles(user, roles)
	return &response, nil
}

// UpdateProfile updates user information in UserInfo table
func (s *UserService) UpdateProfile(ctx contextx.Contextx, userID uint, req dto.UpdateProfileRequest) (*dto.UserResponse, error) {
	// Check if user exists
	_, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		return nil, err
	}

	// Get UserInfo
	userInfo, err := s.userInfoRepo.GetByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}

	// Check if username is being changed and if it's already taken
	if req.Username != "" && req.Username != userInfo.Username {
		taken, err := s.userInfoRepo.IsUsernameTaken(ctx, req.Username, userID)
		if err != nil {
			return nil, errors.New("failed to check username availability")
		}
		if taken {
			return nil, errors.New("username already taken")
		}
	}

	// Check if email is being changed and if it's already taken
	if req.Email != "" && req.Email != userInfo.Email {
		taken, err := s.userInfoRepo.IsEmailTaken(ctx, req.Email, userID)
		if err != nil {
			return nil, errors.New("failed to check email availability")
		}
		if taken {
			return nil, errors.New("email already taken")
		}
	}

	// Update fields that are provided
	if req.Username != "" {
		userInfo.Username = req.Username
	}
	if req.FirstName != "" {
		userInfo.FirstName = req.FirstName
	}
	if req.LastName != "" {
		userInfo.LastName = req.LastName
	}
	if req.Email != "" {
		userInfo.Email = req.Email
	}
	if req.AvatarURL != "" {
		userInfo.AvatarURL = req.AvatarURL
	}
	if req.Language != "" {
		userInfo.Language = req.Language
	}
	if req.Timezone != "" {
		userInfo.Timezone = req.Timezone
	}
	if req.Bio != "" {
		userInfo.Bio = req.Bio
	}
	if req.Location != "" {
		userInfo.Location = req.Location
	}
	if req.Website != "" {
		userInfo.Website = req.Website
	}
	if req.Phone != "" {
		userInfo.Phone = req.Phone
	}

	if err := s.userInfoRepo.Update(ctx, userInfo); err != nil {
		return nil, err
	}

	// Return updated profile
	return s.GetProfile(ctx, userID)
}

// ChangePassword changes user's password after verifying current password
func (s *UserService) ChangePassword(ctx contextx.Contextx, userID uint, currentPassword, newPassword string) error {
	// Get user's auth providers
	providers, err := s.authProviderRepo.GetByUserID(ctx, userID)
	if err != nil {
		return errors.New("user not found or no password set")
	}

	// Find email provider
	var authProvider *models.AuthProvider
	for _, p := range providers {
		if p.Provider == models.ProviderEmail {
			authProvider = &p
			break
		}
	}
	if authProvider == nil {
		return errors.New("user not found or no password set")
	}

	// Verify current password
	if err := bcrypt.CompareHashAndPassword([]byte(authProvider.Password), []byte(currentPassword)); err != nil {
		return errors.New("invalid current password")
	}

	// Hash new password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
	if err != nil {
		return errors.New("failed to hash new password")
	}

	// Update password in auth provider
	authProvider.Password = string(hashedPassword)
	if err := s.authProviderRepo.Update(ctx, authProvider); err != nil {
		return err
	}

	return nil
}

// GetUserByID retrieves user with basic info (used by middleware)
func (s *UserService) GetUserByID(ctx contextx.Contextx, userID uint) (*models.User, error) {
	return s.userRepo.GetByIDWithPreload(ctx, userID, "UserInfo")
}

// UpdateUserStatus updates user status (active, inactive, suspended, pending)
func (s *UserService) UpdateUserStatus(ctx contextx.Contextx, userID uint, status models.UserStatus) error {
	return s.userRepo.UpdateStatus(ctx, userID, status)
}

// VerifyEmail marks user as email verified
func (s *UserService) VerifyEmail(ctx contextx.Contextx, userID uint) error {
	return s.userRepo.VerifyEmail(ctx, userID)
}

// GetUserAuthProviders returns all auth providers for a user
func (s *UserService) GetUserAuthProviders(ctx contextx.Contextx, userID uint) ([]models.AuthProvider, error) {
	return s.authProviderRepo.GetByUserID(ctx, userID)
}

// GetAllUsers returns a list of all users with their basic information
func (s *UserService) GetAllUsers(ctx contextx.Contextx) ([]dto.UserResponse, error) {
	users, err := s.userRepo.GetAll(ctx)
	if err != nil {
		return nil, err
	}

	var userResponses []dto.UserResponse
	for _, user := range users {
		var roles []string
		if s.rbacService != nil {
			userRoles, err := s.rbacService.GetUserRoles(user.ID)
			if err == nil {
				roles = userRoles
			}
		}
		userResponse := dto.ToUserResponseWithRoles(&user, roles)
		userResponses = append(userResponses, userResponse)
	}

	return userResponses, nil
}

// SearchUsers searches for users by name or email
func (s *UserService) SearchUsers(ctx contextx.Contextx, searchTerm string) ([]dto.UserResponse, error) {
	users, err := s.userRepo.Search(ctx, searchTerm)
	if err != nil {
		return nil, err
	}

	var userResponses []dto.UserResponse
	for _, user := range users {
		var roles []string
		if s.rbacService != nil {
			userRoles, err := s.rbacService.GetUserRoles(user.ID)
			if err == nil {
				roles = userRoles
			}
		}
		userResponse := dto.ToUserResponseWithRoles(&user, roles)
		userResponses = append(userResponses, userResponse)
	}

	return userResponses, nil
}

// DeleteUser soft deletes a user and related data
func (s *UserService) DeleteUser(ctx contextx.Contextx, userID uint) error {
	tx := s.db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// Soft delete user (cascades to related tables due to GORM relationships)
	if err := tx.Delete(&models.User{}, userID).Error; err != nil {
		tx.Rollback()
		return errors.New("failed to delete user")
	}

	// Soft delete UserInfo
	if err := tx.Where("user_id = ?", userID).Delete(&models.UserInfo{}).Error; err != nil {
		tx.Rollback()
		return errors.New("failed to delete user info")
	}

	// Soft delete AuthProviders
	if err := tx.Where("user_id = ?", userID).Delete(&models.AuthProvider{}).Error; err != nil {
		tx.Rollback()
		return errors.New("failed to delete auth providers")
	}

	if err := tx.Commit().Error; err != nil {
		return errors.New("failed to delete user account")
	}

	return nil
}

// CreateUser creates a new user with UserInfo
func (s *UserService) CreateUser(ctx contextx.Contextx, req dto.CreateUserRequest) (*dto.UserResponse, error) {
	// Check if user with this email already exists
	var existingUserInfo models.UserInfo
	if err := s.db.Where("email = ?", req.Email).First(&existingUserInfo).Error; err == nil {
		return nil, errors.New("user with this email already exists")
	}

	// Check if username is already taken
	if err := s.db.Where("username = ?", req.Username).First(&existingUserInfo).Error; err == nil {
		return nil, errors.New("username already taken")
	}

	// Hash the password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, errors.New("failed to hash password")
	}

	// Begin transaction
	tx := s.db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// Create user record
	user := models.User{
		Status:        models.UserStatus(req.Status),
		EmailVerified: false,
	}

	if err := tx.Create(&user).Error; err != nil {
		tx.Rollback()
		return nil, errors.New("failed to create user")
	}

	// Create user info record
	userInfo := models.UserInfo{
		UserID:    user.ID,
		Username:  req.Username,
		FirstName: req.FirstName,
		LastName:  req.LastName,
		Email:     req.Email,
		Language:  req.Language,
		Timezone:  req.Timezone,
		Bio:       req.Bio,
		Location:  req.Location,
		Website:   req.Website,
		Phone:     req.Phone,
	}

	if err := tx.Create(&userInfo).Error; err != nil {
		tx.Rollback()
		return nil, errors.New("failed to create user info")
	}

	// Create auth provider (password-based)
	authProvider := models.AuthProvider{
		UserID:     user.ID,
		Provider:   models.ProviderEmail,
		ProviderID: req.Email,
		UserName:   req.Username,
		Password:   string(hashedPassword),
	}

	if err := tx.Create(&authProvider).Error; err != nil {
		tx.Rollback()
		return nil, errors.New("failed to create auth provider")
	}

	// Commit transaction
	if err := tx.Commit().Error; err != nil {
		return nil, errors.New("failed to create user account")
	}

	// Return created user
	user.UserInfo = &userInfo
	var roles []string
	if s.rbacService != nil {
		userRoles, err := s.rbacService.GetUserRoles(user.ID)
		if err == nil {
			roles = userRoles
		}
	}
	response := dto.ToUserResponseWithRoles(&user, roles)
	return &response, nil
}

// GetUserByID retrieves a user by ID with all information
func (s *UserService) GetUserByIDDetailed(userID uint) (*dto.UserResponse, error) {
	var user models.User
	if err := s.db.Preload("UserInfo").First(&user, userID).Error; err != nil {
		return nil, errors.New("user not found")
	}

	var roles []string
	if s.rbacService != nil {
		userRoles, err := s.rbacService.GetUserRoles(user.ID)
		if err == nil {
			roles = userRoles
		}
	}

	response := dto.ToUserResponseWithRoles(&user, roles)
	return &response, nil
}

// UpdateUser updates user information
func (s *UserService) UpdateUser(ctx contextx.Contextx, userID uint, req dto.UpdateUserRequest) (*dto.UserResponse, error) {
	// Check if user exists
	var user models.User
	if err := s.db.Preload("UserInfo").First(&user, userID).Error; err != nil {
		return nil, errors.New("user not found")
	}

	// Check if email is being changed and if it's already taken
	if req.Email != "" && req.Email != user.UserInfo.Email {
		var existingUser models.User
		if err := s.db.Joins("UserInfo").Where("user_info.email = ? AND users.id != ?", req.Email, userID).First(&existingUser).Error; err == nil {
			return nil, errors.New("email already taken by another user")
		}
	}

	// Begin transaction
	tx := s.db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// Update user status if provided
	if req.Status != "" {
		user.Status = models.UserStatus(req.Status)
		if err := tx.Save(&user).Error; err != nil {
			tx.Rollback()
			return nil, errors.New("failed to update user status")
		}
	}

	// Update user info
	if req.FirstName != "" {
		user.UserInfo.FirstName = req.FirstName
	}
	if req.LastName != "" {
		user.UserInfo.LastName = req.LastName
	}
	if req.Email != "" {
		user.UserInfo.Email = req.Email
	}
	if req.Language != "" {
		user.UserInfo.Language = req.Language
	}
	if req.Timezone != "" {
		user.UserInfo.Timezone = req.Timezone
	}
	if req.Bio != "" {
		user.UserInfo.Bio = req.Bio
	}
	if req.Location != "" {
		user.UserInfo.Location = req.Location
	}
	if req.Website != "" {
		user.UserInfo.Website = req.Website
	}
	if req.Phone != "" {
		user.UserInfo.Phone = req.Phone
	}

	if err := tx.Save(user.UserInfo).Error; err != nil {
		tx.Rollback()
		return nil, errors.New("failed to update user info")
	}

	// Update auth provider email if email was changed
	if req.Email != "" && req.Email != user.UserInfo.Email {
		if err := tx.Model(&models.AuthProvider{}).Where("user_id = ? AND provider_type = ?", userID, models.ProviderEmail).Update("provider_id", req.Email).Error; err != nil {
			tx.Rollback()
			return nil, errors.New("failed to update auth provider")
		}
	}

	// Commit transaction
	if err := tx.Commit().Error; err != nil {
		return nil, errors.New("failed to save user changes")
	}

	// Return updated user
	var roles []string
	if s.rbacService != nil {
		userRoles, err := s.rbacService.GetUserRoles(user.ID)
		if err == nil {
			roles = userRoles
		}
	}
	response := dto.ToUserResponseWithRoles(&user, roles)
	return &response, nil
}
