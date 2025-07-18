package repository

import (
	"errors"

	"bezbase/internal/models"
	"bezbase/internal/pkg/contextx"

	"gorm.io/gorm"
)

type userRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) UserRepository {
	return &userRepository{db: db}
}

func (r *userRepository) GetByID(ctx contextx.Contextx, userID uint) (*models.User, error) {
	var user models.User
	if err := ctx.GetTxn(r.db).WithContext(ctx).First(&user, userID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("user not found")
		}
		return nil, err
	}
	return &user, nil
}

func (r *userRepository) GetByIDWithPreload(ctx contextx.Contextx, userID uint, preloads ...string) (*models.User, error) {
	var user models.User
	query := ctx.GetTxn(r.db).WithContext(ctx)
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

func (r *userRepository) GetAll(ctx contextx.Contextx) ([]models.User, error) {
	var users []models.User
	if err := ctx.GetTxn(r.db).Preload("UserInfo").Find(&users).Error; err != nil {
		return nil, errors.New("failed to get users")
	}
	return users, nil
}

func (r *userRepository) Search(ctx contextx.Contextx, searchTerm string) ([]models.User, error) {
	var users []models.User
	searchPattern := "%" + searchTerm + "%"

	if err := ctx.GetTxn(r.db).Preload("UserInfo").
		Joins("LEFT JOIN user_info ON users.id = user_info.user_id").
		Where("user_info.first_name ILIKE ? OR user_info.last_name ILIKE ? OR user_info.email ILIKE ?",
			searchPattern, searchPattern, searchPattern).
		Find(&users).Error; err != nil {
		return nil, errors.New("failed to search users")
	}
	return users, nil
}

func (r *userRepository) Create(ctx contextx.Contextx, user *models.User) error {
	if err := ctx.GetTxn(r.db).WithContext(ctx).Create(user).Error; err != nil {
		return errors.New("failed to create user")
	}
	return nil
}

func (r *userRepository) Update(ctx contextx.Contextx, user *models.User) error {
	if err := ctx.GetTxn(r.db).WithContext(ctx).Save(user).Error; err != nil {
		return errors.New("failed to update user")
	}
	return nil
}

func (r *userRepository) Delete(ctx contextx.Contextx, userID uint) error {
	if err := ctx.GetTxn(r.db).WithContext(ctx).Delete(&models.User{}, userID).Error; err != nil {
		return errors.New("failed to delete user")
	}
	return nil
}

func (r *userRepository) UpdateStatus(ctx contextx.Contextx, userID uint, status models.UserStatus) error {
	var user models.User
	if err := ctx.GetTxn(r.db).First(&user, userID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("user not found")
		}
		return err
	}

	user.Status = status
	if err := ctx.GetTxn(r.db).Save(&user).Error; err != nil {
		return errors.New("failed to update user status")
	}
	return nil
}

func (r *userRepository) VerifyEmail(ctx contextx.Contextx, userID uint) error {
	var user models.User
	if err := ctx.GetTxn(r.db).First(&user, userID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("user not found")
		}
		return err
	}

	user.EmailVerified = true
	if user.Status == models.UserStatusPending {
		user.Status = models.UserStatusActive
	}

	if err := ctx.GetTxn(r.db).Save(&user).Error; err != nil {
		return errors.New("failed to verify email")
	}
	return nil
}

func (r *userRepository) GetByIDDetailed(ctx contextx.Contextx, userID uint) (*models.User, error) {
	var user models.User
	if err := ctx.GetTxn(r.db).WithContext(ctx).Preload("UserInfo").Preload("AuthProviders").First(&user, userID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("user not found")
		}
		return nil, err
	}
	return &user, nil
}
