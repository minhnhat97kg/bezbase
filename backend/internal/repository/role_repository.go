package repository

import (
	"errors"
	"fmt"

	"bezbase/internal/models"
	"bezbase/internal/pkg/contextx"

	"gorm.io/gorm"
)

type roleRepository struct {
	db *gorm.DB
}

func NewRoleRepository(db *gorm.DB) RoleRepository {
	return &roleRepository{db: db}
}

func (r *roleRepository) GetByID(ctx contextx.Contextx, id uint) (*models.Role, error) {
	var role models.Role
	if err := ctx.GetTxn(r.db).First(&role, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("role not found")
		}
		return nil, err
	}
	return &role, nil
}

func (r *roleRepository) GetByName(ctx contextx.Contextx, name string) (*models.Role, error) {
	var role models.Role
	if err := ctx.GetTxn(r.db).Where("name = ?", name).First(&role).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("role not found")
		}
		return nil, err
	}
	return &role, nil
}

func (r *roleRepository) GetAll(ctx contextx.Contextx) ([]models.Role, error) {
	var roles []models.Role
	if err := ctx.GetTxn(r.db).Find(&roles).Error; err != nil {
		return nil, errors.New("failed to get roles")
	}
	return roles, nil
}

func (r *roleRepository) GetActive(ctx contextx.Contextx) ([]models.Role, error) {
	var roles []models.Role
	if err := ctx.GetTxn(r.db).Where("is_active = ?", true).Find(&roles).Error; err != nil {
		return nil, errors.New("failed to get active roles")
	}
	return roles, nil
}

func (r *roleRepository) GetWithPagination(ctx contextx.Contextx, page, pageSize int, searchFilter, statusFilter string, isSystemFilter *bool, sortField, sortOrder string) ([]models.Role, int, error) {
	var roles []models.Role
	var total int64

	query := ctx.GetTxn(r.db).Model(&models.Role{})

	// Apply search filter
	if searchFilter != "" {
		query = query.Where("name LIKE ? OR display_name LIKE ?", "%"+searchFilter+"%", "%"+searchFilter+"%")
	}

	// Apply status filter
	if statusFilter != "" {
		switch statusFilter {
		case "active":
			query = query.Where("is_active = ?", true)
		case "inactive":
			query = query.Where("is_active = ?", false)
		}
	}

	// Apply is_system filter
	if isSystemFilter != nil {
		query = query.Where("is_system = ?", *isSystemFilter)
	}

	// Count total records
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, fmt.Errorf("failed to count roles: %w", err)
	}

	// Apply sorting
	if sortField != "" {
		allowedSortFields := map[string]bool{
			"name":         true,
			"display_name": true,
			"created_at":   true,
		}
		if allowedSortFields[sortField] {
			query = query.Order(sortField + " " + sortOrder)
		} else {
			query = query.Order("created_at " + sortOrder)
		}
	} else {
		query = query.Order("created_at " + sortOrder)
	}

	// Apply pagination
	offset := (page - 1) * pageSize
	query = query.Offset(offset).Limit(pageSize)

	if err := query.Find(&roles).Error; err != nil {
		return nil, 0, fmt.Errorf("failed to get roles: %w", err)
	}

	return roles, int(total), nil
}

func (r *roleRepository) Create(ctx contextx.Contextx, role *models.Role) error {
	if err := ctx.GetTxn(r.db).Create(role).Error; err != nil {
		return errors.New("failed to create role")
	}
	return nil
}

func (r *roleRepository) Update(ctx contextx.Contextx, role *models.Role) error {
	if err := ctx.GetTxn(r.db).Save(role).Error; err != nil {
		return errors.New("failed to update role")
	}
	return nil
}

func (r *roleRepository) Delete(ctx contextx.Contextx, role *models.Role) error {
	if err := ctx.GetTxn(r.db).Delete(role).Error; err != nil {
		return errors.New("failed to delete role")
	}
	return nil
}
