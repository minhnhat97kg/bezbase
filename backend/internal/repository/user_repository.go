package repository

import (
	"errors"

	"bezbase/internal/models"

	"gorm.io/gorm"
)

type userRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) UserRepository {
	return &userRepository{db: db}
}

func (r *userRepository) GetByID(userID uint) (*models.User, error) {
	var user models.User
	if err := r.db.First(&user, userID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("user not found")
		}
		return nil, err
	}
	return &user, nil
}

func (r *userRepository) GetByIDWithPreload(userID uint, preloads ...string) (*models.User, error) {
	var user models.User
	query := r.db
	
	for _, preload := range preloads {
		query = query.Preload(preload)
	}
	
	if err := query.First(&user, userID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("user not found")
		}
		return nil, err
	}
	return &user, nil
}

func (r *userRepository) GetAll() ([]models.User, error) {
	var users []models.User
	if err := r.db.Preload("UserInfo").Find(&users).Error; err != nil {
		return nil, errors.New("failed to get users")
	}
	return users, nil
}

func (r *userRepository) Search(searchTerm string) ([]models.User, error) {
	var users []models.User
	searchPattern := "%" + searchTerm + "%"

	if err := r.db.Preload("UserInfo").
		Joins("LEFT JOIN user_info ON users.id = user_info.user_id").
		Where("user_info.first_name ILIKE ? OR user_info.last_name ILIKE ? OR user_info.email ILIKE ?",
			searchPattern, searchPattern, searchPattern).
		Find(&users).Error; err != nil {
		return nil, errors.New("failed to search users")
	}
	return users, nil
}

func (r *userRepository) Create(user *models.User) error {
	if err := r.db.Create(user).Error; err != nil {
		return errors.New("failed to create user")
	}
	return nil
}

func (r *userRepository) Update(user *models.User) error {
	if err := r.db.Save(user).Error; err != nil {
		return errors.New("failed to update user")
	}
	return nil
}

func (r *userRepository) Delete(userID uint) error {
	if err := r.db.Delete(&models.User{}, userID).Error; err != nil {
		return errors.New("failed to delete user")
	}
	return nil
}

func (r *userRepository) UpdateStatus(userID uint, status models.UserStatus) error {
	var user models.User
	if err := r.db.First(&user, userID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("user not found")
		}
		return err
	}

	user.Status = status
	if err := r.db.Save(&user).Error; err != nil {
		return errors.New("failed to update user status")
	}
	return nil
}

func (r *userRepository) VerifyEmail(userID uint) error {
	var user models.User
	if err := r.db.First(&user, userID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("user not found")
		}
		return err
	}

	user.EmailVerified = true
	if user.Status == models.UserStatusPending {
		user.Status = models.UserStatusActive
	}

	if err := r.db.Save(&user).Error; err != nil {
		return errors.New("failed to verify email")
	}
	return nil
}