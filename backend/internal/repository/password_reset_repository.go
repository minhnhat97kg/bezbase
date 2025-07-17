package repository

import (
	"bezbase/internal/models"
	"time"

	"gorm.io/gorm"
)

type PasswordResetRepository interface {
	Create(token *models.PasswordResetToken) error
	GetByToken(token string) (*models.PasswordResetToken, error)
	GetByUserID(userID uint) (*models.PasswordResetToken, error)
	MarkAsUsed(token string) error
	DeleteExpiredTokens() error
	DeleteByUserID(userID uint) error
}

type passwordResetRepository struct {
	db *gorm.DB
}

func NewPasswordResetRepository(db *gorm.DB) PasswordResetRepository {
	return &passwordResetRepository{db: db}
}

func (r *passwordResetRepository) Create(token *models.PasswordResetToken) error {
	return r.db.Create(token).Error
}

func (r *passwordResetRepository) GetByToken(token string) (*models.PasswordResetToken, error) {
	var resetToken models.PasswordResetToken
	err := r.db.Where("token = ? AND used_at IS NULL", token).First(&resetToken).Error
	if err != nil {
		return nil, err
	}
	return &resetToken, nil
}

func (r *passwordResetRepository) GetByUserID(userID uint) (*models.PasswordResetToken, error) {
	var resetToken models.PasswordResetToken
	err := r.db.Where("user_id = ? AND used_at IS NULL", userID).
		Order("created_at DESC").
		First(&resetToken).Error
	if err != nil {
		return nil, err
	}
	return &resetToken, nil
}

func (r *passwordResetRepository) MarkAsUsed(token string) error {
	now := time.Now()
	return r.db.Model(&models.PasswordResetToken{}).
		Where("token = ?", token).
		Update("used_at", now).Error
}

func (r *passwordResetRepository) DeleteExpiredTokens() error {
	return r.db.Where("expires_at < ?", time.Now()).
		Delete(&models.PasswordResetToken{}).Error
}

func (r *passwordResetRepository) DeleteByUserID(userID uint) error {
	return r.db.Where("user_id = ?", userID).
		Delete(&models.PasswordResetToken{}).Error
}