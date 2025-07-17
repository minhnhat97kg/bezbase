package repository

import (
	"errors"

	"bezbase/internal/models"

	"gorm.io/gorm"
)

type authProviderRepository struct {
	db *gorm.DB
}

func NewAuthProviderRepository(db *gorm.DB) AuthProviderRepository {
	return &authProviderRepository{db: db}
}

func (r *authProviderRepository) GetByUserID(userID uint) ([]models.AuthProvider, error) {
	var providers []models.AuthProvider
	if err := r.db.Where("user_id = ?", userID).Find(&providers).Error; err != nil {
		return nil, errors.New("failed to get auth providers")
	}
	return providers, nil
}

func (r *authProviderRepository) GetByUsernameAndProvider(username string, provider models.AuthProviderType) (*models.AuthProvider, error) {
	var authProvider models.AuthProvider
	if err := r.db.Where("user_name = ? AND provider = ?", username, provider).First(&authProvider).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("auth provider not found")
		}
		return nil, err
	}
	return &authProvider, nil
}

func (r *authProviderRepository) GetByProviderIDAndType(providerID string, provider models.AuthProviderType) (*models.AuthProvider, error) {
	var authProvider models.AuthProvider
	if err := r.db.Where("provider_id = ? AND provider = ?", providerID, provider).First(&authProvider).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("auth provider not found")
		}
		return nil, err
	}
	return &authProvider, nil
}

func (r *authProviderRepository) Create(authProvider *models.AuthProvider) error {
	if err := r.db.Create(authProvider).Error; err != nil {
		return errors.New("failed to create auth provider")
	}
	return nil
}

func (r *authProviderRepository) Update(authProvider *models.AuthProvider) error {
	if err := r.db.Save(authProvider).Error; err != nil {
		return errors.New("failed to update auth provider")
	}
	return nil
}

func (r *authProviderRepository) Delete(userID uint) error {
	if err := r.db.Where("user_id = ?", userID).Delete(&models.AuthProvider{}).Error; err != nil {
		return errors.New("failed to delete auth providers")
	}
	return nil
}

func (r *authProviderRepository) UpdateEmail(userID uint, provider models.AuthProviderType, newEmail string) error {
	if err := r.db.Model(&models.AuthProvider{}).
		Where("user_id = ? AND provider = ?", userID, provider).
		Update("provider_id", newEmail).Error; err != nil {
		return errors.New("failed to update auth provider email")
	}
	return nil
}