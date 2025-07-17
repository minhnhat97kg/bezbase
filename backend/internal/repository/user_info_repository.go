package repository

import (
	"errors"

	"bezbase/internal/models"

	"gorm.io/gorm"
)

type userInfoRepository struct {
	db *gorm.DB
}

func NewUserInfoRepository(db *gorm.DB) UserInfoRepository {
	return &userInfoRepository{db: db}
}

func (r *userInfoRepository) GetByUserID(userID uint) (*models.UserInfo, error) {
	var userInfo models.UserInfo
	if err := r.db.Where("user_id = ?", userID).First(&userInfo).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("user info not found")
		}
		return nil, err
	}
	return &userInfo, nil
}

func (r *userInfoRepository) GetByEmail(email string) (*models.UserInfo, error) {
	var userInfo models.UserInfo
	if err := r.db.Where("email = ?", email).First(&userInfo).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("user info not found")
		}
		return nil, err
	}
	return &userInfo, nil
}

func (r *userInfoRepository) GetByUsername(username string) (*models.UserInfo, error) {
	var userInfo models.UserInfo
	if err := r.db.Where("username = ?", username).First(&userInfo).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("user info not found")
		}
		return nil, err
	}
	return &userInfo, nil
}

func (r *userInfoRepository) Create(userInfo *models.UserInfo) error {
	if err := r.db.Create(userInfo).Error; err != nil {
		return errors.New("failed to create user info")
	}
	return nil
}

func (r *userInfoRepository) Update(userInfo *models.UserInfo) error {
	if err := r.db.Save(userInfo).Error; err != nil {
		return errors.New("failed to update user info")
	}
	return nil
}

func (r *userInfoRepository) Delete(userID uint) error {
	if err := r.db.Where("user_id = ?", userID).Delete(&models.UserInfo{}).Error; err != nil {
		return errors.New("failed to delete user info")
	}
	return nil
}

func (r *userInfoRepository) IsEmailTaken(email string, excludeUserID uint) (bool, error) {
	var count int64
	query := r.db.Model(&models.UserInfo{}).Where("email = ?", email)
	
	if excludeUserID > 0 {
		query = query.Where("user_id != ?", excludeUserID)
	}
	
	if err := query.Count(&count).Error; err != nil {
		return false, err
	}
	
	return count > 0, nil
}

func (r *userInfoRepository) IsUsernameTaken(username string, excludeUserID uint) (bool, error) {
	var count int64
	query := r.db.Model(&models.UserInfo{}).Where("username = ?", username)
	
	if excludeUserID > 0 {
		query = query.Where("user_id != ?", excludeUserID)
	}
	
	if err := query.Count(&count).Error; err != nil {
		return false, err
	}
	
	return count > 0, nil
}