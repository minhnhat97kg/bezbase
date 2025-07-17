package repository

import (
	"fmt"

	"bezbase/internal/models"

	"gorm.io/gorm"
)

type ruleRepository struct {
	db *gorm.DB
}

func NewRuleRepository(db *gorm.DB) RuleRepository {
	return &ruleRepository{db: db}
}

func (r *ruleRepository) GetPermissions(page, pageSize int, roleFilter, resourceFilter, actionFilter, sortField, sortOrder string) ([]models.Rule, int, error) {
	// Query casbin_rule table directly for permissions (ptype = 'p')
	query := r.db.Model(&models.Rule{}).Where("ptype = 'p'")

	// Apply filters
	if roleFilter != "" {
		query = query.Where("v0 LIKE ?", "%"+roleFilter+"%")
	}
	if resourceFilter != "" {
		query = query.Where("v1 LIKE ?", "%"+resourceFilter+"%")
	}
	if actionFilter != "" {
		query = query.Where("v2 LIKE ?", "%"+actionFilter+"%")
	}

	// Count total records
	var total int64
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, fmt.Errorf("failed to count permissions: %w", err)
	}

	// Apply sorting
	orderClause := "id ASC" // default sort
	if sortField != "" {
		switch sortField {
		case "role":
			orderClause = fmt.Sprintf("v0 %s", sortOrder)
		case "resource":
			orderClause = fmt.Sprintf("v1 %s", sortOrder)
		case "action":
			orderClause = fmt.Sprintf("v2 %s", sortOrder)
		default:
			orderClause = fmt.Sprintf("id %s", sortOrder)
		}
	}

	rules := make([]models.Rule, 0)
	offset := (page - 1) * pageSize
	if err := query.Order(orderClause).Offset(offset).Limit(pageSize).
		Select("id, v0, v1, v2").Find(&rules).Error; err != nil {
		return nil, 0, fmt.Errorf("failed to get permissions: %w", err)
	}

	return rules, int(total), nil
}