package models

import (
	"time"

	"gorm.io/gorm"
)

type RoleInheritance struct {
	ID           uint      `json:"id" gorm:"primaryKey"`
	ParentRoleID uint      `json:"parent_role_id" gorm:"not null;index"`
	ChildRoleID  uint      `json:"child_role_id" gorm:"not null;index"`
	Depth        int       `json:"depth" gorm:"not null;default:1"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`

	// Relationships
	ParentRole Role `json:"parent_role,omitempty" gorm:"foreignKey:ParentRoleID"`
	ChildRole  Role `json:"child_role,omitempty" gorm:"foreignKey:ChildRoleID"`
}

type ContextualPermission struct {
	ID           uint           `json:"id" gorm:"primaryKey"`
	RoleID       uint           `json:"role_id" gorm:"not null;index"`
	Resource     string         `json:"resource" gorm:"not null;size:100" validate:"required,min=1,max=100"`
	Action       string         `json:"action" gorm:"not null;size:50" validate:"required,min=1,max=50"`
	ContextType  string         `json:"context_type" gorm:"size:50" validate:"omitempty,min=1,max=50"`
	ContextValue string         `json:"context_value" gorm:"size:255" validate:"omitempty,min=1,max=255"`
	IsGranted    bool           `json:"is_granted" gorm:"default:true"`
	CreatedAt    time.Time      `json:"created_at"`
	UpdatedAt    time.Time      `json:"updated_at"`
	DeletedAt    gorm.DeletedAt `json:"-" gorm:"index"`

	// Relationships
	Role Role `json:"role,omitempty" gorm:"foreignKey:RoleID"`
}

type RoleTemplate struct {
	ID          uint           `json:"id" gorm:"primaryKey"`
	Name        string         `json:"name" gorm:"not null;size:100" validate:"required,min=1,max=100"`
	DisplayName string         `json:"display_name" gorm:"not null;size:255" validate:"required,min=1,max=255"`
	Description string         `json:"description" gorm:"size:500" validate:"omitempty,max=500"`
	Category    string         `json:"category" gorm:"size:100" validate:"omitempty,oneof=system business department basic"`
	Config      string         `json:"config" gorm:"type:jsonb"`
	IsActive    bool           `json:"is_active" gorm:"default:true"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `json:"-" gorm:"index"`
}

// TableName returns the table name for RoleInheritance
func (RoleInheritance) TableName() string {
	return "role_inheritances"
}

// TableName returns the table name for ContextualPermission
func (ContextualPermission) TableName() string {
	return "contextual_permissions"
}

// TableName returns the table name for RoleTemplate
func (RoleTemplate) TableName() string {
	return "role_templates"
}

// BeforeCreate sets default values for ContextualPermission
func (cp *ContextualPermission) BeforeCreate(tx *gorm.DB) error {
	// Default is to grant permission
	return nil
}

// BeforeCreate sets default values for RoleTemplate
func (rt *RoleTemplate) BeforeCreate(tx *gorm.DB) error {
	if rt.Category == "" {
		rt.Category = "basic"
	}
	return nil
}

// PermissionKey returns a unique key for this contextual permission
func (cp *ContextualPermission) PermissionKey() string {
	key := cp.Resource + ":" + cp.Action
	if cp.ContextType != "" && cp.ContextValue != "" {
		key += ":" + cp.ContextType + ":" + cp.ContextValue
	}
	return key
}

// IsGlobal checks if this is a global permission (no context)
func (cp *ContextualPermission) IsGlobal() bool {
	return cp.ContextType == "" || cp.ContextValue == ""
}

// IsContextSpecific checks if this permission is context-specific
func (cp *ContextualPermission) IsContextSpecific() bool {
	return cp.ContextType != "" && cp.ContextValue != ""
}

// MatchesContext checks if this permission matches the given context
func (cp *ContextualPermission) MatchesContext(contextType, contextValue string) bool {
	if cp.IsGlobal() {
		return true
	}
	return cp.ContextType == contextType && (cp.ContextValue == contextValue || cp.ContextValue == "*")
}

// Extended Role model methods for hierarchy support
// These methods should be added to the existing Role model

// GetAllParentRoles returns all parent roles in the hierarchy
func GetAllParentRoles(db *gorm.DB, roleID uint) ([]Role, error) {
	var roles []Role
	err := db.Raw(`
		WITH RECURSIVE role_hierarchy AS (
			SELECT r.*, 0 as level
			FROM roles r
			WHERE r.id = ?
			
			UNION ALL
			
			SELECT r.*, rh.level + 1
			FROM roles r
			JOIN role_hierarchy rh ON r.id = rh.parent_role_id
		)
		SELECT * FROM role_hierarchy WHERE level > 0
		ORDER BY level
	`, roleID).Scan(&roles).Error
	
	return roles, err
}

// GetAllChildRoles returns all child roles in the hierarchy
func GetAllChildRoles(db *gorm.DB, roleID uint) ([]Role, error) {
	var roles []Role
	err := db.Raw(`
		WITH RECURSIVE role_hierarchy AS (
			SELECT r.*, 0 as level
			FROM roles r
			WHERE r.id = ?
			
			UNION ALL
			
			SELECT r.*, rh.level + 1
			FROM roles r
			JOIN role_hierarchy rh ON rh.id = r.parent_role_id
		)
		SELECT * FROM role_hierarchy WHERE level > 0
		ORDER BY level
	`, roleID).Scan(&roles).Error
	
	return roles, err
}

// GetEffectivePermissions returns all permissions for a role including inherited ones
func GetEffectivePermissions(db *gorm.DB, roleID uint, orgID *uint) ([]ContextualPermission, error) {
	var permissions []ContextualPermission
	
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
		query += ` AND (cp.context_type = 'organization' AND cp.context_value = ? OR cp.context_type IS NULL)`
		args = append(args, *orgID)
	}
	
	query += ` ORDER BY rh.level, cp.resource, cp.action`
	
	err := db.Raw(query, args...).Scan(&permissions).Error
	return permissions, err
}

// CreateRoleHierarchy creates the hierarchy relationship and updates inheritance table
func CreateRoleHierarchy(db *gorm.DB, parentRoleID, childRoleID uint) error {
	// Update the role's parent_role_id
	if err := db.Model(&Role{}).Where("id = ?", childRoleID).Update("parent_role_id", parentRoleID).Error; err != nil {
		return err
	}
	
	// Create or update role inheritance entries for efficient querying
	// This creates entries for all ancestor relationships
	err := db.Exec(`
		WITH RECURSIVE ancestors AS (
			SELECT id, parent_role_id, 1 as depth
			FROM roles
			WHERE id = ?
			
			UNION ALL
			
			SELECT r.id, r.parent_role_id, a.depth + 1
			FROM roles r
			JOIN ancestors a ON r.id = a.parent_role_id
			WHERE a.depth < 10
		)
		INSERT INTO role_inheritances (parent_role_id, child_role_id, depth, created_at, updated_at)
		SELECT a.parent_role_id, ?, a.depth, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP
		FROM ancestors
		WHERE a.parent_role_id IS NOT NULL
		ON CONFLICT (parent_role_id, child_role_id) DO UPDATE SET 
			depth = EXCLUDED.depth,
			updated_at = CURRENT_TIMESTAMP
	`, parentRoleID, childRoleID).Error
	
	return err
}

// RemoveRoleHierarchy removes the hierarchy relationship
func RemoveRoleHierarchy(db *gorm.DB, childRoleID uint) error {
	// Remove parent relationship
	if err := db.Model(&Role{}).Where("id = ?", childRoleID).Update("parent_role_id", nil).Error; err != nil {
		return err
	}
	
	// Remove all inheritance entries involving this role as a child
	return db.Where("child_role_id = ?", childRoleID).Delete(&RoleInheritance{}).Error
}