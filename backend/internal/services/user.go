package services

import (
	"errors"

	"bezbase/internal/dto"
	"bezbase/internal/models"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type UserService struct {
	db          *gorm.DB
	rbacService *RBACService
}

func NewUserService(db *gorm.DB, rbacService *RBACService) *UserService {
	return &UserService{
		db:          db,
		rbacService: rbacService,
	}
}

// GetProfile retrieves user profile with all user information
func (s *UserService) GetProfile(userID uint) (*dto.UserResponse, error) {
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

// UpdateProfile updates user information in UserInfo table
func (s *UserService) UpdateProfile(userID uint, req dto.UpdateProfileRequest) (*dto.UserResponse, error) {
	// Start transaction for atomic update
	tx := s.db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// Check if user exists
	var user models.User
	if err := tx.First(&user, userID).Error; err != nil {
		tx.Rollback()
		return nil, errors.New("user not found")
	}

	// Update UserInfo
	var userInfo models.UserInfo
	if err := tx.Where("user_id = ?", userID).First(&userInfo).Error; err != nil {
		tx.Rollback()
		return nil, errors.New("user info not found")
	}

	// Update fields that are provided
	if req.FirstName != "" {
		userInfo.FirstName = req.FirstName
	}
	if req.LastName != "" {
		userInfo.LastName = req.LastName
	}
	// if req.AvatarURL != "" {
	// 	userInfo.AvatarURL = req.AvatarURL
	// }
	if req.Language != "" {
		userInfo.Language = req.Language
	}
	if req.Timezone != "" {
		userInfo.Timezone = req.Timezone
	}

	if err := tx.Save(&userInfo).Error; err != nil {
		tx.Rollback()
		return nil, errors.New("failed to update profile")
	}

	// Commit transaction
	if err := tx.Commit().Error; err != nil {
		return nil, errors.New("failed to save profile changes")
	}

	// Return updated profile
	return s.GetProfile(userID)
}

// GetUserByID retrieves user with basic info (used by middleware)
func (s *UserService) GetUserByID(userID uint) (*models.User, error) {
	var user models.User
	if err := s.db.Preload("UserInfo").First(&user, userID).Error; err != nil {
		return nil, errors.New("user not found")
	}

	return &user, nil
}

// UpdateUserStatus updates user status (active, inactive, suspended, pending)
func (s *UserService) UpdateUserStatus(userID uint, status models.UserStatus) error {
	var user models.User
	if err := s.db.First(&user, userID).Error; err != nil {
		return errors.New("user not found")
	}

	user.Status = status
	if err := s.db.Save(&user).Error; err != nil {
		return errors.New("failed to update user status")
	}

	return nil
}

// VerifyEmail marks user as email verified
func (s *UserService) VerifyEmail(userID uint) error {
	var user models.User
	if err := s.db.First(&user, userID).Error; err != nil {
		return errors.New("user not found")
	}

	user.EmailVerified = true
	if status := user.Status; status == models.UserStatusPending {
		user.Status = models.UserStatusActive
	}

	if err := s.db.Save(&user).Error; err != nil {
		return errors.New("failed to verify email")
	}

	return nil
}

// GetUserAuthProviders returns all auth providers for a user
func (s *UserService) GetUserAuthProviders(userID uint) ([]models.AuthProvider, error) {
	var providers []models.AuthProvider
	if err := s.db.Where("user_id = ?", userID).Find(&providers).Error; err != nil {
		return nil, errors.New("failed to get auth providers")
	}

	return providers, nil
}

// GetAllUsers returns a list of all users with their basic information
func (s *UserService) GetAllUsers() ([]dto.UserResponse, error) {
	var users []models.User
	if err := s.db.Preload("UserInfo").Find(&users).Error; err != nil {
		return nil, errors.New("failed to get users")
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
func (s *UserService) SearchUsers(searchTerm string) ([]dto.UserResponse, error) {
	var users []models.User

	// Search in both UserInfo (first_name, last_name, email) and Users table
	searchPattern := "%" + searchTerm + "%"

	if err := s.db.Preload("UserInfo").
		Joins("LEFT JOIN user_info ON users.id = user_info.user_id").
		Where("user_info.first_name ILIKE ? OR user_info.last_name ILIKE ? OR user_info.email ILIKE ?",
			searchPattern, searchPattern, searchPattern).
		Find(&users).Error; err != nil {
		return nil, errors.New("failed to search users")
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
func (s *UserService) DeleteUser(userID uint) error {
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
func (s *UserService) CreateUser(req dto.CreateUserRequest) (*dto.UserResponse, error) {
	// Check if user with this email already exists
	var existingUser models.User
	if err := s.db.Joins("UserInfo").Where("user_info.email = ?", req.Email).First(&existingUser).Error; err == nil {
		return nil, errors.New("user with this email already exists")
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
func (s *UserService) UpdateUser(userID uint, req dto.UpdateUserRequest) (*dto.UserResponse, error) {
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
