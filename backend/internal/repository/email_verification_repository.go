package repository

import (
	"bezbase/internal/models"
	"time"

	"gorm.io/gorm"
)

type EmailVerificationRepository interface {
	Create(token *models.EmailVerificationToken) error
	GetByToken(token string) (*models.EmailVerificationToken, error)
	GetByUserID(userID uint) (*models.EmailVerificationToken, error)
	MarkAsUsed(token string) error
	DeleteExpiredTokens() error
	DeleteByUserID(userID uint) error
}

type emailVerificationRepository struct {
	db *gorm.DB
}

func NewEmailVerificationRepository(db *gorm.DB) EmailVerificationRepository {
	return &emailVerificationRepository{db: db}
}

func (r *emailVerificationRepository) Create(token *models.EmailVerificationToken) error {
	return r.db.Create(token).Error
}

func (r *emailVerificationRepository) GetByToken(token string) (*models.EmailVerificationToken, error) {
	var verificationToken models.EmailVerificationToken
	err := r.db.Where("token = ? AND used_at IS NULL", token).First(&verificationToken).Error
	if err != nil {
		return nil, err
	}
	return &verificationToken, nil
}

func (r *emailVerificationRepository) GetByUserID(userID uint) (*models.EmailVerificationToken, error) {
	var verificationToken models.EmailVerificationToken
	err := r.db.Where("user_id = ? AND used_at IS NULL", userID).
		Order("created_at DESC").
		First(&verificationToken).Error
	if err != nil {
		return nil, err
	}
	return &verificationToken, nil
}

func (r *emailVerificationRepository) MarkAsUsed(token string) error {
	now := time.Now()
	return r.db.Model(&models.EmailVerificationToken{}).
		Where("token = ?", token).
		Update("used_at", now).Error
}

func (r *emailVerificationRepository) DeleteExpiredTokens() error {
	return r.db.Where("expires_at < ?", time.Now()).
		Delete(&models.EmailVerificationToken{}).Error
}

func (r *emailVerificationRepository) DeleteByUserID(userID uint) error {
	return r.db.Where("user_id = ?", userID).
		Delete(&models.EmailVerificationToken{}).Error
}