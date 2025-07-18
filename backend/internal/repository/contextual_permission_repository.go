package repository

import (
	"bezbase/internal/models"
	"bezbase/internal/pkg/contextx"
	"errors"
	"fmt"

	"gorm.io/gorm"
)

type contextualPermissionRepository struct {
	db *gorm.DB
}

func NewContextualPermissionRepository(db *gorm.DB) ContextualPermissionRepository {
	return &contextualPermissionRepository{db: db}
}

func (r *contextualPermissionRepository) GetByID(ctx contextx.Contextx, id uint) (*models.ContextualPermission, error) {
	var permission models.ContextualPermission
	err := ctx.GetTxn(r.db).Preload("Role").Where("id = ?", id).First(&permission).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, gorm.ErrRecordNotFound
		}
		return nil, err
	}
	return &permission, nil
}

func (r *contextualPermissionRepository) GetByRoleID(ctx contextx.Contextx, roleID uint) ([]models.ContextualPermission, error) {
	var permissions []models.ContextualPermission
	err := ctx.GetTxn(r.db).Where("role_id = ?", roleID).Order("resource, action").Find(&permissions).Error
	return permissions, err
}

func (r *contextualPermissionRepository) GetByRoleIDAndContext(ctx contextx.Contextx, roleID uint, contextType string, contextValue string) ([]models.ContextualPermission, error) {
	var permissions []models.ContextualPermission
	query := ctx.GetTxn(r.db).Where("role_id = ?", roleID)

	if contextType != "" {
		query = query.Where("context_type = ?", contextType)
	}
	if contextValue != "" {
		query = query.Where("context_value = ?", contextValue)
	}

	err := query.Order("resource, action").Find(&permissions).Error
	return permissions, err
}

func (r *contextualPermissionRepository) GetEffectivePermissions(ctx contextx.Contextx, roleID uint, orgID *uint) ([]models.ContextualPermission, error) {
	var permissions []models.ContextualPermission

	query := `
		WITH RECURSIVE role_hierarchy AS (
			SELECT r.id, r.parent_role_id, 0 as level
			FROM roles r
			WHERE r.id = ?
			
			UNION ALL
			
			SELECT r.id, r.parent_role_id, rh.level + 1
			FROM roles r
			JOIN role_hierarchy rh ON r.id = rh.parent_role_id
			WHERE rh.level < 10  -- Prevent infinite loops
		)
		SELECT DISTINCT cp.*
		FROM contextual_permissions cp
		JOIN role_hierarchy rh ON cp.role_id = rh.id
		WHERE cp.deleted_at IS NULL
	`

	args := []interface{}{roleID}

	if orgID != nil {
		query += ` AND (cp.context_type = 'organization' AND cp.context_value = ? OR cp.context_type IS NULL OR cp.context_type = '')`
		args = append(args, fmt.Sprintf("%d", *orgID))
	}

	query += ` ORDER BY rh.level, cp.resource, cp.action`

	err := ctx.GetTxn(r.db).Raw(query, args...).Scan(&permissions).Error
	return permissions, err
}

func (r *contextualPermissionRepository) Create(ctx contextx.Contextx, permission *models.ContextualPermission) error {
	return ctx.GetTxn(r.db).Create(permission).Error
}

func (r *contextualPermissionRepository) Update(ctx contextx.Contextx, permission *models.ContextualPermission) error {
	return ctx.GetTxn(r.db).Save(permission).Error
}

func (r *contextualPermissionRepository) Delete(ctx contextx.Contextx, id uint) error {
	return ctx.GetTxn(r.db).Delete(&models.ContextualPermission{}, id).Error
}

func (r *contextualPermissionRepository) DeleteByRoleID(ctx contextx.Contextx, roleID uint) error {
	return ctx.GetTxn(r.db).Where("role_id = ?", roleID).Delete(&models.ContextualPermission{}).Error
}
