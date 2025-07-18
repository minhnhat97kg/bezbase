package repository

import (
	"bezbase/internal/models"
	"bezbase/internal/pkg/contextx"
	"errors"
	"time"

	"gorm.io/gorm"
)

type organizationInvitationRepository struct {
	db *gorm.DB
}

func NewOrganizationInvitationRepository(db *gorm.DB) OrganizationInvitationRepository {
	return &organizationInvitationRepository{db: db}
}

func (r *organizationInvitationRepository) GetByID(ctx contextx.Contextx, id uint) (*models.OrganizationInvitation, error) {
	var invitation models.OrganizationInvitation
	err := ctx.GetTxn(r.db).Preload("Organization").Preload("InvitedByUser").Preload("InvitedByUser.UserInfo").
		Where("id = ?", id).First(&invitation).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, gorm.ErrRecordNotFound
		}
		return nil, err
	}
	return &invitation, nil
}

func (r *organizationInvitationRepository) GetByToken(ctx contextx.Contextx, token string) (*models.OrganizationInvitation, error) {
	var invitation models.OrganizationInvitation
	err := ctx.GetTxn(r.db).Preload("Organization").Where("token = ?", token).First(&invitation).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, gorm.ErrRecordNotFound
		}
		return nil, err
	}
	return &invitation, nil
}

func (r *organizationInvitationRepository) GetByOrgID(ctx contextx.Contextx, orgID uint) ([]models.OrganizationInvitation, error) {
	var invitations []models.OrganizationInvitation
	err := ctx.GetTxn(r.db).Preload("InvitedByUser").Preload("InvitedByUser.UserInfo").
		Where("org_id = ?", orgID).Order("created_at DESC").Find(&invitations).Error
	return invitations, err
}

func (r *organizationInvitationRepository) GetByOrgIDAndEmail(ctx contextx.Contextx, orgID uint, email string) (*models.OrganizationInvitation, error) {
	var invitation models.OrganizationInvitation
	err := ctx.GetTxn(r.db).Where("org_id = ? AND email = ? AND accepted_at IS NULL AND expires_at > ?",
		orgID, email, time.Now()).First(&invitation).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, gorm.ErrRecordNotFound
		}
		return nil, err
	}
	return &invitation, nil
}

func (r *organizationInvitationRepository) GetPendingByEmail(ctx contextx.Contextx, email string) ([]models.OrganizationInvitation, error) {
	var invitations []models.OrganizationInvitation
	err := ctx.GetTxn(r.db).Preload("Organization").
		Where("email = ? AND accepted_at IS NULL AND expires_at > ?", email, time.Now()).
		Order("created_at DESC").Find(&invitations).Error
	return invitations, err
}

func (r *organizationInvitationRepository) Create(ctx contextx.Contextx, invitation *models.OrganizationInvitation) error {
	return ctx.GetTxn(r.db).Create(invitation).Error
}

func (r *organizationInvitationRepository) Update(ctx contextx.Contextx, invitation *models.OrganizationInvitation) error {
	return ctx.GetTxn(r.db).Save(invitation).Error
}

func (r *organizationInvitationRepository) Delete(ctx contextx.Contextx, id uint) error {
	return ctx.GetTxn(r.db).Delete(&models.OrganizationInvitation{}, id).Error
}

func (r *organizationInvitationRepository) DeleteExpired(ctx contextx.Contextx) error {
	return ctx.GetTxn(r.db).Where("expires_at < ? AND accepted_at IS NULL", time.Now()).
		Delete(&models.OrganizationInvitation{}).Error
}
