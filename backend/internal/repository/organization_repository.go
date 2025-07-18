package repository

import (
	"bezbase/internal/models"
	"bezbase/internal/pkg/contextx"
	"errors"

	"gorm.io/gorm"
)

type organizationRepository struct {
	db *gorm.DB
}

func NewOrganizationRepository(db *gorm.DB) OrganizationRepository {
	return &organizationRepository{db: db}
}

func (r *organizationRepository) GetByID(ctx contextx.Contextx, orgID uint) (*models.Organization, error) {
	var org models.Organization
	err := ctx.GetTxn(r.db).Where("id = ?", orgID).First(&org).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, gorm.ErrRecordNotFound
		}
		return nil, err
	}
	return &org, nil
}

func (r *organizationRepository) GetByIDWithRelations(ctx contextx.Contextx, orgID uint) (*models.Organization, error) {
	var org models.Organization
	err := ctx.GetTxn(r.db).Preload("Users").Preload("Users.User").Preload("Users.User.UserInfo").
		Preload("Invitations").Preload("Invitations.InvitedByUser").
		Preload("Roles").Where("id = ?", orgID).First(&org).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, gorm.ErrRecordNotFound
		}
		return nil, err
	}
	return &org, nil
}

func (r *organizationRepository) GetBySlug(ctx contextx.Contextx, slug string) (*models.Organization, error) {
	var org models.Organization
	err := ctx.GetTxn(r.db).Where("slug = ?", slug).First(&org).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, gorm.ErrRecordNotFound
		}
		return nil, err
	}
	return &org, nil
}

func (r *organizationRepository) GetAll(ctx contextx.Contextx) ([]models.Organization, error) {
	var orgs []models.Organization
	err := ctx.GetTxn(r.db).Where("is_active = ?", true).Find(&orgs).Error
	return orgs, err
}

func (r *organizationRepository) Create(ctx contextx.Contextx, org *models.Organization) error {
	return ctx.GetTxn(r.db).Create(org).Error
}

func (r *organizationRepository) Update(ctx contextx.Contextx, org *models.Organization) error {
	return ctx.GetTxn(r.db).Save(org).Error
}

func (r *organizationRepository) Delete(ctx contextx.Contextx, orgID uint) error {
	return ctx.GetTxn(r.db).Delete(&models.Organization{}, orgID).Error
}