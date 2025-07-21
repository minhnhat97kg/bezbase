package services

import (
	"context"
	"fmt"
	"log"
	"strconv"
	"strings"

	"bezbase/internal/dto"
	"bezbase/internal/models"
	"bezbase/internal/pkg/contextx"
	"bezbase/internal/repository"

	"github.com/casbin/casbin/v2"
	"github.com/casbin/casbin/v2/model"
	gormadapter "github.com/casbin/gorm-adapter/v3"
	"gorm.io/gorm"
)

type RBACService struct {
	enforcer *casbin.Enforcer
	roleRepo repository.RoleRepository
	ruleRepo repository.RuleRepository
	db       *gorm.DB
}

// AssignDefaultRoleToUser assigns the default 'user' role to a user if they have no roles
func (r *RBACService) AssignDefaultRoleToUser(ctx contextx.Contextx, userID uint) error {
	subject := fmt.Sprintf("user:%d", userID)
	roles, err := r.enforcer.GetRolesForUser(subject)
	if err != nil {
		return err
	}
	if len(roles) == 0 {
		_, err := r.enforcer.AddRoleForUser(subject, "user")
		if err != nil {
			return err
		}
		return r.enforcer.SavePolicy()
	}
	return nil
}

// GetPermissionsForUser returns all permissions for a user (resource, action)
func (r *RBACService) GetPermissionsForUser(ctx contextx.Contextx, userID uint) ([]string, error) {
	subject := fmt.Sprintf("user:%d", userID)
	var result []string

	// Get direct permissions for user
	perms, err := r.enforcer.GetPermissionsForUser(subject)
	if err != nil {
		return nil, err
	}
	for _, perm := range perms {
		if len(perm) >= 3 {
			result = append(result, fmt.Sprintf("%s:%s:%s", subject, perm[1], perm[2]))
		}
	}

	// Get roles assigned to user
	roles, err := r.enforcer.GetRolesForUser(subject)
	if err != nil {
		return nil, err
	}
	// If user has no roles, treat as if they have 'user' role
	if len(roles) == 0 {
		roles = append(roles, "user")
	}
	for _, role := range roles {
		rolePerms, err := r.enforcer.GetPermissionsForUser(role)
		if err != nil {
			continue
		}
		for _, perm := range rolePerms {
			if len(perm) >= 3 {
				result = append(result, fmt.Sprintf("%s:%s:%s", role, perm[1], perm[2]))
			}
		}
	}

	return result, nil
}

func NewRBACService(
	roleRepo repository.RoleRepository,
	ruleRepo repository.RuleRepository,
	db *gorm.DB,
) (*RBACService, error) {
	adapter, err := gormadapter.NewAdapterByDBUseTableName(db, "", "rules")
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
		roleRepo: roleRepo,
		ruleRepo: ruleRepo,
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
	roles, err := r.roleRepo.GetAll(contextx.Background())
	if err != nil {
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
			{"admin", models.ResourceTypeAll.String(), models.ActionTypeAll.String()},
		}
	case "moderator":
		permissions = [][]string{
			{"moderator", models.ResourceTypeUser.String(), models.ActionTypeRead.String()},
			{"moderator", models.ResourceTypeUser.String(), models.ActionTypeUpdate.String()},
			{"moderator", models.ResourceTypePost.String(), models.ActionTypeAll.String()},
		}
	case "user":
		permissions = [][]string{
			{"user", models.ResourceTypeProfile.String(), models.ActionTypeRead.String()},
			{"user", models.ResourceTypeProfile.String(), models.ActionTypeUpdate.String()},
			{"user", models.ResourceTypePost.String(), models.ActionTypeCreate.String()},
			{"user", models.ResourceTypePost.String(), models.ActionTypeRead.String()},
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
	// Check roles for user
	roles, err := r.enforcer.GetRolesForUser(user)
	if err != nil {
		return false, err
	}
	
	// If user has no roles, assign default 'user' role automatically
	if len(roles) == 0 {
		if err := r.AssignDefaultRoleToUser(contextx.Background(), userID); err != nil {
			// If assignment fails, fall back to checking with default role
			roles = append(roles, "user")
		} else {
			// Re-fetch roles after assignment
			roles, err = r.enforcer.GetRolesForUser(user)
			if err != nil {
				return false, err
			}
		}
	}
	
	// Check as each role (including inherited permissions)
	for _, roleName := range roles {
		// Check direct permission
		ok, err := r.enforcer.Enforce(roleName, resource, action)
		if err != nil {
			continue
		}
		if ok {
			return true, nil
		}
		
		// Check inherited permissions from parent roles
		// First get the role by name to get its ID
		roleObj, err := r.roleRepo.GetByName(contextx.Background(), roleName)
		if err == nil {
			if hasInheritedPermission, err := r.checkInheritedPermissions(roleObj.ID, resource, action); err == nil && hasInheritedPermission {
				return true, nil
			}
		}
	}
	
	return false, nil
}

func (r *RBACService) AddRole(role string) error {
	// For Casbin, roles are created implicitly when policies are added
	// We just need to ensure the role exists in the system
	return nil
}

func (r *RBACService) CreateRole(ctx contextx.Contextx, req dto.CreateRoleRequest) (*models.Role, error) {
	// Check if role already exists
	if _, err := r.roleRepo.GetByName(ctx, req.Name); err == nil {
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

	if err := r.roleRepo.Create(ctx, &role); err != nil {
		return nil, fmt.Errorf("failed to create role: %w", err)
	}

	return &role, nil
}

func (r *RBACService) UpdateRole(ctx contextx.Contextx, id uint, req dto.UpdateRoleRequest) (*models.Role, error) {
	role, err := r.roleRepo.GetByID(ctx, id)
	if err != nil {
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

	if err := r.roleRepo.Update(ctx, role); err != nil {
		return nil, fmt.Errorf("failed to update role: %w", err)
	}

	return role, nil
}

func (r *RBACService) GetRoleByName(ctx contextx.Contextx, name string) (*models.Role, error) {
	var role models.Role
	if err := r.db.WithContext(ctx).Where("name = ?", name).First(&role).Error; err != nil {
		return nil, fmt.Errorf("role not found: %w", err)
	}
	return &role, nil
}

func (r *RBACService) GetRoleByID(ctx contextx.Contextx, id uint) (*models.Role, error) {
	var role models.Role
	if err := r.db.WithContext(ctx).First(&role, id).Error; err != nil {
		return nil, fmt.Errorf("role not found: %w", err)
	}
	return &role, nil
}

func (r *RBACService) GetAllRolesWithMetadata(ctx contextx.Contextx) ([]models.Role, error) {
	var roles []models.Role
	if err := r.db.WithContext(ctx).Find(&roles).Error; err != nil {
		return nil, fmt.Errorf("failed to get roles: %w", err)
	}
	return roles, nil
}

func (r *RBACService) GetAllRolesWithPagination(ctx contextx.Contextx, page, pageSize int, searchFilter, statusFilter string, isSystemFilter *bool, sortField, sortOrder string) ([]models.Role, int, error) {
	var roles []models.Role
	var total int64

	query := r.db.WithContext(ctx).Model(&models.Role{})

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

func (r *RBACService) GetActiveRoles(ctx contextx.Contextx) ([]models.Role, error) {
	var roles []models.Role
	if err := r.db.WithContext(ctx).Where("is_active = ?", true).Find(&roles).Error; err != nil {
		return nil, fmt.Errorf("failed to get active roles: %w", err)
	}
	return roles, nil
}

func (r *RBACService) AddPermission(role, resource, action string) error {
	// Validate role exists
	roleModel, err := r.GetRoleByName(contextx.Background(), role)
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
	roleModel, err := r.GetRoleByName(contextx.Background(), role)
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
	roles, err := r.enforcer.GetRolesForUser(user)
	if err != nil {
		return nil, err
	}
	if len(roles) == 0 {
		return []string{"user"}, nil
	}
	return roles, nil
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

func (r *RBACService) GetAllPermissions(page, pageSize int, roleFilter, resourceFilter, actionFilter, permissionFilter, sortField, sortOrder string) ([]dto.PermissionResponse, int, error) {
	// Get hardcoded permissions to create a mapping
	hardcodedPermissions := models.GetHardcodedPermissions()
	permissionMap := make(map[string]models.Permission)
	for _, perm := range hardcodedPermissions {
		key := fmt.Sprintf("%s:%s", perm.Resource.String(), perm.Action.String())
		permissionMap[key] = perm
	}

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
	if permissionFilter != "" {
		// Filter by the permission field from hardcoded permissions
		// First, find all hardcoded permissions that match the filter
		var matchingResourceActions []string
		for _, perm := range hardcodedPermissions {
			if strings.Contains(strings.ToLower(perm.Permission), strings.ToLower(permissionFilter)) {
				matchingResourceActions = append(matchingResourceActions, fmt.Sprintf("%s:%s", perm.Resource.String(), perm.Action.String()))
			}
		}

		if len(matchingResourceActions) > 0 {
			// Create OR conditions for each matching resource:action
			conditions := make([]string, len(matchingResourceActions))
			args := make([]interface{}, len(matchingResourceActions))
			for i, resourceAction := range matchingResourceActions {
				conditions[i] = "CONCAT(v1, ':', v2) = ?"
				args[i] = resourceAction
			}
			query = query.Where(fmt.Sprintf("(%s)", strings.Join(conditions, " OR ")), args...)
		} else {
			// If no hardcoded permissions match, also check resource:action format
			query = query.Where("CONCAT(v1, ':', v2) LIKE ?", "%"+permissionFilter+"%")
		}
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
		case "permission":
			// For permission sorting, we'll need to sort by resource:action format
			// since we can't easily sort by the hardcoded permission values in SQL
			orderClause = fmt.Sprintf("CONCAT(v1, ':', v2) %s", sortOrder)
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

	// Convert to PermissionResponse format using hardcoded permissions
	permissions := make([]dto.PermissionResponse, len(rules))
	for i, rule := range rules {
		key := fmt.Sprintf("%s:%s", rule.V1, rule.V2)
		permission := key // default to resource:action format

		// Use the hardcoded permission if available
		if hardcodedPerm, exists := permissionMap[key]; exists {
			permission = hardcodedPerm.Permission
		}

		permissions[i] = dto.PermissionResponse{
			ID:         rule.ID,
			Role:       rule.V0,
			Resource:   rule.V1,
			Action:     rule.V2,
			Permission: permission,
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

// Check permission with contextual information (simplified without organization context)
func (r *RBACService) CheckPermissionWithContext(userID uint, resource, action string) (bool, error) {
	// Use the standard permission check since we removed organization context
	return r.CheckPermission(userID, resource, action)
}


// Check contextual permissions
func (r *RBACService) checkContextualPermissions(roleID uint, resource, action string) (bool, error) {
	var count int64
	query := r.db.Model(&models.ContextualPermission{}).
		Where("role_id = ? AND resource = ? AND action = ? AND is_granted = ?",
			roleID, resource, action, true)


	err := query.Count(&count).Error
	return count > 0, err
}

// Check inherited permissions from parent roles
func (r *RBACService) checkInheritedPermissions(roleID uint, resource, action string) (bool, error) {
	// Get all parent roles in the hierarchy
	parentRoles, err := models.GetAllParentRoles(r.db, roleID)
	if err != nil {
		return false, err
	}

	for _, parentRole := range parentRoles {
		// Check Casbin permissions for parent role
		ok, err := r.enforcer.Enforce(parentRole.Name, resource, action)
		if err != nil {
			continue
		}
		if ok {
			return true, nil
		}

		// Check contextual permissions for parent role
		hasContextualPermission, err := r.checkContextualPermissions(parentRole.ID, resource, action)
		if err != nil {
			continue
		}
		if hasContextualPermission {
			return true, nil
		}
	}

	return false, nil
}

// Create role from template
func (r *RBACService) CreateRoleFromTemplate(templateID uint, customName string) (*models.Role, error) {
	var template models.RoleTemplate
	if err := r.db.First(&template, templateID).Error; err != nil {
		return nil, fmt.Errorf("template not found: %w", err)
	}

	roleName := template.Name
	if customName != "" {
		roleName = customName
	}

	role := &models.Role{
		Name:           roleName,
		DisplayName:    template.DisplayName,
		Description:    template.Description,
		IsSystem:       false,
		IsActive:       true,
		HierarchyLevel: 2, // Templates create level 2 roles by default
	}

	ctx := contextx.NewContextx(context.Background())
	roleReq := dto.CreateRoleRequest{
		Name:        role.Name,
		DisplayName: role.DisplayName,
		Description: role.Description,
	}

	createdRole, err := r.CreateRole(ctx, roleReq)
	if err != nil {
		return nil, err
	}

	// TODO: Parse template.Config and create contextual permissions
	// This would involve parsing the JSON config and creating the appropriate permissions

	return createdRole, nil
}


// Update role hierarchy
func (r *RBACService) SetRoleParent(childRoleID uint, parentRoleID *uint) error {
	childRole, err := r.GetRoleByID(contextx.Background(), childRoleID)
	if err != nil {
		return err
	}

	if childRole.IsSystem {
		return fmt.Errorf("cannot modify system role hierarchy")
	}

	if parentRoleID != nil {
		parentRole, err := r.GetRoleByID(contextx.Background(), *parentRoleID)
		if err != nil {
			return fmt.Errorf("parent role not found: %w", err)
		}

		// Comprehensive circular dependency check
		if err := r.validateRoleHierarchy(childRoleID, *parentRoleID); err != nil {
			return err
		}

		// Update hierarchy level
		childRole.HierarchyLevel = parentRole.HierarchyLevel + 1
	} else {
		childRole.HierarchyLevel = 0
	}

	// Simply update the parent_role_id field
	childRole.ParentRoleID = parentRoleID
	return r.db.Save(childRole).Error
}

// GetEligibleParentRoles returns roles that can be safely set as parent for the given role
func (r *RBACService) GetEligibleParentRoles(ctx contextx.Contextx, roleID uint) ([]models.Role, error) {
	// Get all active roles
	var allRoles []models.Role
	if err := r.db.Where("is_active = ? AND id != ?", true, roleID).Find(&allRoles).Error; err != nil {
		return nil, err
	}

	// Get the current role to understand its position in hierarchy (for future enhancements)
	_, err := r.GetRoleByID(ctx, roleID)
	if err != nil {
		return nil, err
	}

	// Build a map of role relationships for easier traversal
	childMap := make(map[uint][]uint) // parent_id -> []child_ids
	for _, role := range allRoles {
		if role.ParentRoleID != nil {
			childMap[*role.ParentRoleID] = append(childMap[*role.ParentRoleID], role.ID)
		}
	}

	// Function to check if a role is a descendant of another role
	isDescendant := func(potentialDescendant, ancestor uint) bool {
		visited := make(map[uint]bool)
		var checkDescendants func(uint) bool
		checkDescendants = func(currentRole uint) bool {
			if visited[currentRole] {
				return false // Prevent infinite loops
			}
			visited[currentRole] = true
			
			children := childMap[currentRole]
			for _, child := range children {
				if child == potentialDescendant {
					return true
				}
				if checkDescendants(child) {
					return true
				}
			}
			return false
		}
		return checkDescendants(ancestor)
	}

	// Filter out ineligible roles
	var eligibleRoles []models.Role
	for _, role := range allRoles {
		// Cannot assign system roles as parents
		if role.IsSystem {
			continue
		}
		
		// Cannot assign self as parent
		if role.ID == roleID {
			continue
		}
		
		// Cannot assign descendants as parents (would create circular dependency)
		if isDescendant(role.ID, roleID) {
			continue
		}
		
		// Cannot assign roles that would create a hierarchy level > 10 (arbitrary limit)
		if role.HierarchyLevel >= 9 {
			continue
		}
		
		eligibleRoles = append(eligibleRoles, role)
	}

	return eligibleRoles, nil
}

// validateRoleHierarchy checks if setting parentRoleID as parent of childRoleID would create circular dependency
func (r *RBACService) validateRoleHierarchy(childRoleID, parentRoleID uint) error {
	// Cannot set self as parent
	if childRoleID == parentRoleID {
		return fmt.Errorf("cannot set role as its own parent")
	}

	// Check if parent role would become descendant of child role
	visited := make(map[uint]bool)
	var checkAncestors func(uint) bool
	checkAncestors = func(currentRoleID uint) bool {
		if visited[currentRoleID] {
			return false // Prevent infinite loops
		}
		visited[currentRoleID] = true

		var role models.Role
		if err := r.db.First(&role, currentRoleID).Error; err != nil {
			return false
		}

		if role.ParentRoleID == nil {
			return false
		}

		if *role.ParentRoleID == childRoleID {
			return true // Found circular dependency
		}

		return checkAncestors(*role.ParentRoleID)
	}

	if checkAncestors(parentRoleID) {
		return fmt.Errorf("circular dependency detected: parent role %d is a descendant of child role %d", parentRoleID, childRoleID)
	}

	// Check depth limit
	depth := 0
	currentRoleID := parentRoleID
	visited = make(map[uint]bool)
	
	for currentRoleID != 0 && depth < 10 {
		if visited[currentRoleID] {
			return fmt.Errorf("circular dependency detected in role hierarchy")
		}
		visited[currentRoleID] = true

		var role models.Role
		if err := r.db.First(&role, currentRoleID).Error; err != nil {
			break
		}

		if role.ParentRoleID == nil {
			break
		}

		currentRoleID = *role.ParentRoleID
		depth++
	}

	if depth >= 10 {
		return fmt.Errorf("role hierarchy depth limit exceeded (max 10 levels)")
	}

	return nil
}
