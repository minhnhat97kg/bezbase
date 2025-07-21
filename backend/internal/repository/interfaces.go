package repository

import (
	"bezbase/internal/models"
	"bezbase/internal/pkg/contextx"
)

// UserRepository defines the interface for user data access
type UserRepository interface {
	GetByID(ctx contextx.Contextx, userID uint) (*models.User, error)
	GetByIDWithPreload(ctx contextx.Contextx, userID uint, preloads ...string) (*models.User, error)
	GetAll(ctx contextx.Contextx) ([]models.User, error)
	Search(ctx contextx.Contextx, searchTerm string) ([]models.User, error)
	Create(ctx contextx.Contextx, user *models.User) error
	Update(ctx contextx.Contextx, user *models.User) error
	Delete(ctx contextx.Contextx, userID uint) error
	UpdateStatus(ctx contextx.Contextx, userID uint, status models.UserStatus) error
	VerifyEmail(ctx contextx.Contextx, userID uint) error
	GetByIDDetailed(ctx contextx.Contextx, userID uint) (*models.User, error)
}

// UserInfoRepository defines the interface for user info data access
type UserInfoRepository interface {
	GetByUserID(ctx contextx.Contextx, userID uint) (*models.UserInfo, error)
	GetByEmail(ctx contextx.Contextx, email string) (*models.UserInfo, error)
	GetByUsername(ctx contextx.Contextx, username string) (*models.UserInfo, error)
	Create(ctx contextx.Contextx, userInfo *models.UserInfo) error
	Update(ctx contextx.Contextx, userInfo *models.UserInfo) error
	Delete(ctx contextx.Contextx, userID uint) error
	IsEmailTaken(ctx contextx.Contextx, email string, excludeUserID uint) (bool, error)
	IsUsernameTaken(ctx contextx.Contextx, username string, excludeUserID uint) (bool, error)
}

// AuthProviderRepository defines the interface for auth provider data access
type AuthProviderRepository interface {
	GetByUserID(ctx contextx.Contextx, userID uint) ([]models.AuthProvider, error)
	GetByUsernameAndProvider(ctx contextx.Contextx, username string, provider models.AuthProviderType) (*models.AuthProvider, error)
	GetByProviderIDAndType(ctx contextx.Contextx, providerID string, provider models.AuthProviderType) (*models.AuthProvider, error)
	GetByUserIDAndProvider(ctx contextx.Contextx, userID uint, provider models.AuthProviderType) (*models.AuthProvider, error)
	Create(ctx contextx.Contextx, authProvider *models.AuthProvider) error
	Update(ctx contextx.Contextx, authProvider *models.AuthProvider) error
	Delete(ctx contextx.Contextx, userID uint) error
	UpdateEmail(ctx contextx.Contextx, userID uint, provider models.AuthProviderType, newEmail string) error
}

// RoleRepository defines the interface for role data access
type RoleRepository interface {
	GetByID(ctx contextx.Contextx, id uint) (*models.Role, error)
	GetByName(ctx contextx.Contextx, name string) (*models.Role, error)
	GetAll(ctx contextx.Contextx) ([]models.Role, error)
	GetActive(ctx contextx.Contextx) ([]models.Role, error)
	GetWithPagination(ctx contextx.Contextx, page, pageSize int, searchFilter, statusFilter string, isSystemFilter *bool, sortField, sortOrder string) ([]models.Role, int, error)
	Create(ctx contextx.Contextx, role *models.Role) error
	Update(ctx contextx.Contextx, role *models.Role) error
	Delete(ctx contextx.Contextx, role *models.Role) error
}

// RuleRepository defines the interface for RBAC rule data access
type RuleRepository interface {
	GetPermissions(ctx contextx.Contextx, page, pageSize int, roleFilter, resourceFilter, actionFilter, sortField, sortOrder string) ([]models.Rule, int, error)
}


// RoleTemplateRepository defines the interface for role template data access
type RoleTemplateRepository interface {
	GetByID(ctx contextx.Contextx, id uint) (*models.RoleTemplate, error)
	GetByName(ctx contextx.Contextx, name string) (*models.RoleTemplate, error)
	GetByCategory(ctx contextx.Contextx, category string) ([]models.RoleTemplate, error)
	GetActive(ctx contextx.Contextx) ([]models.RoleTemplate, error)
	GetAll(ctx contextx.Contextx) ([]models.RoleTemplate, error)
	Create(ctx contextx.Contextx, template *models.RoleTemplate) error
	Update(ctx contextx.Contextx, template *models.RoleTemplate) error
	Delete(ctx contextx.Contextx, id uint) error
}

// ContextualPermissionRepository defines the interface for contextual permission data access
type ContextualPermissionRepository interface {
	GetByID(ctx contextx.Contextx, id uint) (*models.ContextualPermission, error)
	GetByRoleID(ctx contextx.Contextx, roleID uint) ([]models.ContextualPermission, error)
	GetByRoleIDAndContext(ctx contextx.Contextx, roleID uint, contextType string, contextValue string) ([]models.ContextualPermission, error)
	GetEffectivePermissions(ctx contextx.Contextx, roleID uint) ([]models.ContextualPermission, error)
	Create(ctx contextx.Contextx, permission *models.ContextualPermission) error
	Update(ctx contextx.Contextx, permission *models.ContextualPermission) error
	Delete(ctx contextx.Contextx, id uint) error
	DeleteByRoleID(ctx contextx.Contextx, roleID uint) error
}

// RoleInheritanceRepository defines the interface for role inheritance data access
type RoleInheritanceRepository interface {
	Create(ctx contextx.Contextx, inheritance *models.RoleInheritance) error
	GetByChildRole(ctx contextx.Contextx, childRoleID uint) (*models.RoleInheritance, error)
	DeleteByChildRole(ctx contextx.Contextx, childRoleID uint) error
	GetParentRoles(ctx contextx.Contextx, roleID uint) ([]models.Role, error)
	GetChildRoles(ctx contextx.Contextx, roleID uint) ([]models.Role, error)
}
