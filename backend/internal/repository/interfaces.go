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
