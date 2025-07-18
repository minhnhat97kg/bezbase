package repository

import (
	"bezbase/internal/models"
	"bezbase/internal/pkg/contextx"
	"errors"

	"gorm.io/gorm"
)

type organizationUserRepository struct {
	db *gorm.DB
}

func NewOrganizationUserRepository(db *gorm.DB) OrganizationUserRepository {
	return &organizationUserRepository{db: db}
}

func (r *organizationUserRepository) GetByID(ctx contextx.Contextx, id uint) (*models.OrganizationUser, error) {
	var orgUser models.OrganizationUser
	err := ctx.GetTxn(r.db).Preload("Organization").Preload("User").Preload("User.UserInfo").
		Where("id = ?", id).First(&orgUser).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, gorm.ErrRecordNotFound
		}
		return nil, err
	}
	return &orgUser, nil
}

func (r *organizationUserRepository) GetByOrgID(ctx contextx.Contextx, orgID uint) ([]models.OrganizationUser, error) {
	var orgUsers []models.OrganizationUser
	err := ctx.GetTxn(r.db).Preload("User").Preload("User.UserInfo").
		Where("org_id = ?", orgID).Find(&orgUsers).Error
	return orgUsers, err
}

func (r *organizationUserRepository) GetByUserID(ctx contextx.Contextx, userID uint) ([]models.OrganizationUser, error) {
	var orgUsers []models.OrganizationUser
	err := ctx.GetTxn(r.db).Preload("Organization").
		Where("user_id = ?", userID).Find(&orgUsers).Error
	return orgUsers, err
}

func (r *organizationUserRepository) GetByOrgIDAndUserID(ctx contextx.Contextx, orgID uint, userID uint) (*models.OrganizationUser, error) {
	var orgUser models.OrganizationUser
	err := ctx.GetTxn(r.db).Where("org_id = ? AND user_id = ?", orgID, userID).First(&orgUser).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, gorm.ErrRecordNotFound
		}
		return nil, err
	}
	return &orgUser, nil
}

func (r *organizationUserRepository) GetByOrgIDAndEmail(ctx contextx.Contextx, orgID uint, email string) (*models.OrganizationUser, error) {
	var orgUser models.OrganizationUser
	err := ctx.GetTxn(r.db).Joins("JOIN users ON organization_users.user_id = users.id").
		Joins("JOIN user_info ON users.id = user_info.user_id").
		Where("organization_users.org_id = ? AND user_info.email = ?", orgID, email).
		First(&orgUser).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, gorm.ErrRecordNotFound
		}
		return nil, err
	}
	return &orgUser, nil
}

func (r *organizationUserRepository) GetByOrgIDAndRole(ctx contextx.Contextx, orgID uint, role string) ([]models.OrganizationUser, error) {
	var orgUsers []models.OrganizationUser
	err := ctx.GetTxn(r.db).Preload("User").Preload("User.UserInfo").
		Where("org_id = ? AND role = ?", orgID, role).Find(&orgUsers).Error
	return orgUsers, err
}

func (r *organizationUserRepository) Create(ctx contextx.Contextx, orgUser *models.OrganizationUser) error {
	return ctx.GetTxn(r.db).Create(orgUser).Error
}

func (r *organizationUserRepository) Update(ctx contextx.Contextx, orgUser *models.OrganizationUser) error {
	return ctx.GetTxn(r.db).Save(orgUser).Error
}

func (r *organizationUserRepository) Delete(ctx contextx.Contextx, id uint) error {
	return ctx.GetTxn(r.db).Delete(&models.OrganizationUser{}, id).Error
}

func (r *organizationUserRepository) DeleteByOrgIDAndUserID(ctx contextx.Contextx, orgID uint, userID uint) error {
	return ctx.GetTxn(r.db).Where("org_id = ? AND user_id = ?", orgID, userID).Delete(&models.OrganizationUser{}).Error
}