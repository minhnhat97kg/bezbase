package rbac

import (
	"fmt"
	"log"
	"strconv"

	"bezbase/internal/dto"
	"bezbase/internal/models"

	"github.com/casbin/casbin/v2"
	"github.com/casbin/casbin/v2/model"
	gormadapter "github.com/casbin/gorm-adapter/v3"
	"gorm.io/gorm"
)

type RBACService struct {
	enforcer *casbin.Enforcer
	db       *gorm.DB
}

func NewRBACService(db *gorm.DB) (*RBACService, error) {
	adapter, err := gormadapter.NewAdapterByDB(db)
	if err != nil {
		return nil, fmt.Errorf("failed to create casbin adapter: %w", err)
	}

	m := model.NewModel()
	err = m.LoadModelFromText(getRBACModel())
	if err != nil {
		return nil, fmt.Errorf("failed to load model from text: %w", err)
	}

	enforcer, err := casbin.NewEnforcer(m, adapter)
	if err != nil {
		return nil, fmt.Errorf("failed to create casbin enforcer: %w", err)
	}

	enforcer.EnableLog(true)
	err = enforcer.LoadPolicy()
	if err != nil {
		return nil, fmt.Errorf("failed to load policy: %w", err)
	}

	rbacService := &RBACService{
		enforcer: enforcer,
		db:       db,
	}

	if err := rbacService.initializeDefaultRoles(); err != nil {
		log.Printf("Warning: failed to initialize default roles: %v", err)
	}

	return rbacService, nil
}

func getRBACModel() string {
	return `
[request_definition]
r = sub, obj, act

[policy_definition]
p = sub, obj, act

[role_definition]
g = _, _

[policy_effect]
e = some(where (p.eft == allow))

[matchers]
m = g(r.sub, p.sub) && (r.obj == p.obj || p.obj == "*") && (r.act == p.act || p.act == "*")
`
}

func (r *RBACService) initializeDefaultRoles() error {
	// Default roles are now created via migration
	// Just ensure the permissions are set up for existing roles
	var roles []models.Role
	if err := r.db.Find(&roles).Error; err != nil {
		return err
	}

	for _, role := range roles {
		if !role.IsActive {
			continue
		}

		// Check if permissions already exist for this role
		permissions, err := r.enforcer.GetPermissionsForUser(role.Name)
		if err != nil {
			return err
		}

		// If no permissions exist, add default ones
		if len(permissions) == 0 {
			if err := r.addDefaultPermissionsForRole(role.Name); err != nil {
				return err
			}
		}
	}

	return nil
}

func (r *RBACService) addDefaultPermissionsForRole(roleName string) error {
	var permissions [][]string

	switch roleName {
	case "admin":
		permissions = [][]string{
			{"admin", "*", "*"},
		}
	case "moderator":
		permissions = [][]string{
			{"moderator", "users", "read"},
			{"moderator", "users", "update"},
			{"moderator", "posts", "*"},
		}
	case "user":
		permissions = [][]string{
			{"user", "profile", "read"},
			{"user", "profile", "update"},
			{"user", "posts", "create"},
			{"user", "posts", "read"},
		}
	}

	for _, policy := range permissions {
		exists, err := r.enforcer.HasPolicy(policy[0], policy[1], policy[2])
		if err != nil {
			return err
		}
		if !exists {
			if err := r.AddPermission(policy[0], policy[1], policy[2]); err != nil {
				return err
			}
		}
	}

	return nil
}

func (r *RBACService) CheckPermission(userID uint, resource, action string) (bool, error) {
	user := fmt.Sprintf("user:%d", userID)
	return r.enforcer.Enforce(user, resource, action)
}

func (r *RBACService) AddRole(role string) error {
	// For Casbin, roles are created implicitly when policies are added
	// We just need to ensure the role exists in the system
	return nil
}

func (r *RBACService) CreateRole(req dto.CreateRoleRequest) (*models.Role, error) {
	// Check if role already exists
	var existingRole models.Role
	if err := r.db.Where("name = ?", req.Name).First(&existingRole).Error; err == nil {
		return nil, fmt.Errorf("role with name '%s' already exists", req.Name)
	}

	role := models.Role{
		Name:        req.Name,
		DisplayName: req.DisplayName,
		Description: req.Description,
		IsSystem:    false,
		IsActive:    true,
	}

	if req.IsActive != nil {
		role.IsActive = *req.IsActive
	}

	if err := role.Validate(); err != nil {
		return nil, fmt.Errorf("validation failed: %w", err)
	}

	if err := r.db.Create(&role).Error; err != nil {
		return nil, fmt.Errorf("failed to create role: %w", err)
	}

	return &role, nil
}

func (r *RBACService) UpdateRole(id uint, req dto.UpdateRoleRequest) (*models.Role, error) {
	var role models.Role
	if err := r.db.First(&role, id).Error; err != nil {
		return nil, fmt.Errorf("role not found: %w", err)
	}

	if role.IsSystem {
		return nil, fmt.Errorf("cannot update system role")
	}

	if req.DisplayName != "" {
		role.DisplayName = req.DisplayName
	}
	if req.Description != "" {
		role.Description = req.Description
	}
	if req.IsActive != nil {
		role.IsActive = *req.IsActive
	}

	if err := role.Validate(); err != nil {
		return nil, fmt.Errorf("validation failed: %w", err)
	}

	if err := r.db.Save(&role).Error; err != nil {
		return nil, fmt.Errorf("failed to update role: %w", err)
	}

	return &role, nil
}

func (r *RBACService) GetRoleByName(name string) (*models.Role, error) {
	var role models.Role
	if err := r.db.Where("name = ?", name).First(&role).Error; err != nil {
		return nil, fmt.Errorf("role not found: %w", err)
	}
	return &role, nil
}

func (r *RBACService) GetRoleByID(id uint) (*models.Role, error) {
	var role models.Role
	if err := r.db.First(&role, id).Error; err != nil {
		return nil, fmt.Errorf("role not found: %w", err)
	}
	return &role, nil
}

func (r *RBACService) GetAllRolesWithMetadata() ([]models.Role, error) {
	var roles []models.Role
	if err := r.db.Find(&roles).Error; err != nil {
		return nil, fmt.Errorf("failed to get roles: %w", err)
	}
	return roles, nil
}

func (r *RBACService) GetActiveRoles() ([]models.Role, error) {
	var roles []models.Role
	if err := r.db.Where("is_active = ?", true).Find(&roles).Error; err != nil {
		return nil, fmt.Errorf("failed to get active roles: %w", err)
	}
	return roles, nil
}

func (r *RBACService) AddPermission(role, resource, action string) error {
	// Validate role exists
	roleModel, err := r.GetRoleByName(role)
	if err != nil {
		return fmt.Errorf("role validation failed: %w", err)
	}

	if !roleModel.IsActive {
		return fmt.Errorf("cannot add permission to inactive role: %s", role)
	}

	_, err = r.enforcer.AddPolicy(role, resource, action)
	if err != nil {
		return fmt.Errorf("failed to add permission: %w", err)
	}
	return r.enforcer.SavePolicy()
}

func (r *RBACService) RemovePermission(role, resource, action string) error {
	_, err := r.enforcer.RemovePolicy(role, resource, action)
	if err != nil {
		return fmt.Errorf("failed to remove permission: %w", err)
	}
	return r.enforcer.SavePolicy()
}

func (r *RBACService) AssignRoleToUser(userID uint, role string) error {
	// Validate role exists and is active
	roleModel, err := r.GetRoleByName(role)
	if err != nil {
		return fmt.Errorf("role validation failed: %w", err)
	}

	if !roleModel.IsActive {
		return fmt.Errorf("cannot assign inactive role: %s", role)
	}

	user := fmt.Sprintf("user:%d", userID)
	_, err = r.enforcer.AddRoleForUser(user, role)
	if err != nil {
		return fmt.Errorf("failed to assign role to user: %w", err)
	}
	return r.enforcer.SavePolicy()
}

func (r *RBACService) RemoveRoleFromUser(userID uint, role string) error {
	user := fmt.Sprintf("user:%d", userID)
	_, err := r.enforcer.DeleteRoleForUser(user, role)
	if err != nil {
		return fmt.Errorf("failed to remove role from user: %w", err)
	}
	return r.enforcer.SavePolicy()
}

func (r *RBACService) GetUserRoles(userID uint) ([]string, error) {
	user := fmt.Sprintf("user:%d", userID)
	return r.enforcer.GetRolesForUser(user)
}

func (r *RBACService) GetUsersWithRole(role string) ([]uint, error) {
	subjects, err := r.enforcer.GetUsersForRole(role)
	if err != nil {
		return nil, fmt.Errorf("failed to get users for role: %w", err)
	}

	var userIDs []uint
	for _, subject := range subjects {
		if len(subject) > 5 && subject[:5] == "user:" {
			userIDStr := subject[5:]
			userID, err := strconv.ParseUint(userIDStr, 10, 32)
			if err != nil {
				continue
			}
			userIDs = append(userIDs, uint(userID))
		}
	}

	return userIDs, nil
}

func (r *RBACService) GetAllRoles() ([]string, error) {
	return r.enforcer.GetAllRoles()
}

func (r *RBACService) GetPermissionsForRole(role string) ([][]string, error) {
	return r.enforcer.GetPermissionsForUser(role)
}

type PermissionResponse struct {
	ID       int    `json:"id"`
	Role     string `json:"role"`
	Resource string `json:"resource"`
	Action   string `json:"action"`
}

func (r *RBACService) GetAllPermissions(page, pageSize int, roleFilter, resourceFilter, actionFilter, sortField, sortOrder string) ([]dto.PermissionResponse, int, error) {
	// Query casbin_rule table directly for permissions (ptype = 'p')
	query := r.db.Table("casbin_rule").Where("ptype = 'p'")

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

	// Apply pagination and get results
	var casbinRules []struct {
		ID int    `json:"id"`
		V0 string `json:"v0"`
		V1 string `json:"v1"`
		V2 string `json:"v2"`
	}

	offset := (page - 1) * pageSize
	if err := query.Order(orderClause).Offset(offset).Limit(pageSize).
		Select("id, v0, v1, v2").Find(&casbinRules).Error; err != nil {
		return nil, 0, fmt.Errorf("failed to get permissions: %w", err)
	}

	// Convert to PermissionResponse format
	permissions := make([]dto.PermissionResponse, len(casbinRules))
	for i, rule := range casbinRules {
		permissions[i] = dto.PermissionResponse{
			ID:       rule.ID,
			Role:     rule.V0,
			Resource: rule.V1,
			Action:   rule.V2,
		}
	}

	return permissions, int(total), nil
}

func (r *RBACService) DeleteRole(role string) error {
	// Get role from database
	var roleModel models.Role
	if err := r.db.Where("name = ?", role).First(&roleModel).Error; err != nil {
		return fmt.Errorf("role not found: %w", err)
	}

	if roleModel.IsSystem {
		return fmt.Errorf("cannot delete system role")
	}

	// Delete from Casbin
	_, err := r.enforcer.DeleteRole(role)
	if err != nil {
		return fmt.Errorf("failed to delete role from enforcer: %w", err)
	}

	// Delete from database
	if err := r.db.Delete(&roleModel).Error; err != nil {
		return fmt.Errorf("failed to delete role from database: %w", err)
	}

	return r.enforcer.SavePolicy()
}

func (r *RBACService) DeleteRoleByID(id uint) error {
	var roleModel models.Role
	if err := r.db.First(&roleModel, id).Error; err != nil {
		return fmt.Errorf("role not found: %w", err)
	}

	return r.DeleteRole(roleModel.Name)
}

func (r *RBACService) ReloadPolicy() error {
	return r.enforcer.LoadPolicy()
}
