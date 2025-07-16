package services

import (
	"errors"

	"bezbase/internal/dto"
	"bezbase/internal/models"

	"gorm.io/gorm"
)

type UserService struct {
	db *gorm.DB
}

func NewUserService(db *gorm.DB) *UserService {
	return &UserService{
		db: db,
	}
}

// GetProfile retrieves user profile with all user information
func (s *UserService) GetProfile(userID uint) (*dto.UserResponse, error) {
	var user models.User
	if err := s.db.Preload("UserInfo").First(&user, userID).Error; err != nil {
		return nil, errors.New("user not found")
	}

	response := dto.ToUserResponse(&user)
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
		userResponse := dto.ToUserResponse(&user)
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
		userResponse := dto.ToUserResponse(&user)
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
