package repository

import (
	"bezbase/internal/models"
	"bezbase/internal/pkg/contextx"

	"gorm.io/gorm"
)

// Interface defined in interfaces.go

type roleInheritanceRepository struct {
	db *gorm.DB
}

func NewRoleInheritanceRepository(db *gorm.DB) RoleInheritanceRepository {
	return &roleInheritanceRepository{
		db: db,
	}
}

func (r *roleInheritanceRepository) Create(ctx contextx.Contextx, inheritance *models.RoleInheritance) error {
	return ctx.GetTxn(r.db).Create(inheritance).Error
}

func (r *roleInheritanceRepository) GetByChildRole(ctx contextx.Contextx, childRoleID uint) (*models.RoleInheritance, error) {
	var inheritance models.RoleInheritance
	err := ctx.GetTxn(r.db).Where("child_role_id = ?", childRoleID).First(&inheritance).Error
	if err != nil {
		return nil, err
	}
	return &inheritance, nil
}

func (r *roleInheritanceRepository) DeleteByChildRole(ctx contextx.Contextx, childRoleID uint) error {
	return ctx.GetTxn(r.db).Where("child_role_id = ?", childRoleID).Delete(&models.RoleInheritance{}).Error
}

func (r *roleInheritanceRepository) GetParentRoles(ctx contextx.Contextx, roleID uint) ([]models.Role, error) {
	var roles []models.Role
	
	query := `
	WITH RECURSIVE role_hierarchy AS (
		SELECT r.id, r.name, r.display_name, r.description, r.is_system, r.is_active, 
			   r.org_id, r.parent_role_id, r.hierarchy_level, r.created_at, r.updated_at
		FROM roles r 
		WHERE r.id = ?
		
		UNION ALL
		
		SELECT p.id, p.name, p.display_name, p.description, p.is_system, p.is_active,
			   p.org_id, p.parent_role_id, p.hierarchy_level, p.created_at, p.updated_at
		FROM roles p
		INNER JOIN role_hierarchy rh ON p.id = rh.parent_role_id
	)
	SELECT DISTINCT * FROM role_hierarchy WHERE id != ?
	ORDER BY hierarchy_level ASC
	`
	
	err := ctx.GetTxn(r.db).Raw(query, roleID, roleID).Scan(&roles).Error
	return roles, err
}

func (r *roleInheritanceRepository) GetChildRoles(ctx contextx.Contextx, roleID uint) ([]models.Role, error) {
	var roles []models.Role
	
	query := `
	WITH RECURSIVE role_hierarchy AS (
		SELECT r.id, r.name, r.display_name, r.description, r.is_system, r.is_active,
			   r.org_id, r.parent_role_id, r.hierarchy_level, r.created_at, r.updated_at
		FROM roles r 
		WHERE r.id = ?
		
		UNION ALL
		
		SELECT c.id, c.name, c.display_name, c.description, c.is_system, c.is_active,
			   c.org_id, c.parent_role_id, c.hierarchy_level, c.created_at, c.updated_at
		FROM roles c
		INNER JOIN role_hierarchy rh ON c.parent_role_id = rh.id
	)
	SELECT DISTINCT * FROM role_hierarchy WHERE id != ?
	ORDER BY hierarchy_level DESC
	`
	
	err := ctx.GetTxn(r.db).Raw(query, roleID, roleID).Scan(&roles).Error
	return roles, err
}