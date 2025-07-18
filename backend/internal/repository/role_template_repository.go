package repository

import (
	"bezbase/internal/models"
	"bezbase/internal/pkg/contextx"
	"errors"

	"gorm.io/gorm"
)

type roleTemplateRepository struct {
	db *gorm.DB
}

func NewRoleTemplateRepository(db *gorm.DB) RoleTemplateRepository {
	return &roleTemplateRepository{db: db}
}

func (r *roleTemplateRepository) GetByID(ctx contextx.Contextx, id uint) (*models.RoleTemplate, error) {
	var template models.RoleTemplate
	err := ctx.GetTxn(r.db).Where("id = ?", id).First(&template).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, gorm.ErrRecordNotFound
		}
		return nil, err
	}
	return &template, nil
}

func (r *roleTemplateRepository) GetByName(ctx contextx.Contextx, name string) (*models.RoleTemplate, error) {
	var template models.RoleTemplate
	err := ctx.GetTxn(r.db).Where("name = ?", name).First(&template).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, gorm.ErrRecordNotFound
		}
		return nil, err
	}
	return &template, nil
}

func (r *roleTemplateRepository) GetByCategory(ctx contextx.Contextx, category string) ([]models.RoleTemplate, error) {
	var templates []models.RoleTemplate
	err := ctx.GetTxn(r.db).Where("category = ? AND is_active = ?", category, true).
		Order("display_name").Find(&templates).Error
	return templates, err
}

func (r *roleTemplateRepository) GetActive(ctx contextx.Contextx) ([]models.RoleTemplate, error) {
	var templates []models.RoleTemplate
	err := ctx.GetTxn(r.db).Where("is_active = ?", true).Order("category, display_name").Find(&templates).Error
	return templates, err
}

func (r *roleTemplateRepository) GetAll(ctx contextx.Contextx) ([]models.RoleTemplate, error) {
	var templates []models.RoleTemplate
	err := ctx.GetTxn(r.db).Order("category, display_name").Find(&templates).Error
	return templates, err
}

func (r *roleTemplateRepository) Create(ctx contextx.Contextx, template *models.RoleTemplate) error {
	return ctx.GetTxn(r.db).Create(template).Error
}

func (r *roleTemplateRepository) Update(ctx contextx.Contextx, template *models.RoleTemplate) error {
	return ctx.GetTxn(r.db).Save(template).Error
}

func (r *roleTemplateRepository) Delete(ctx contextx.Contextx, id uint) error {
	return ctx.GetTxn(r.db).Delete(&models.RoleTemplate{}, id).Error
}