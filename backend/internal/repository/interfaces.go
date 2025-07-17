package repository

import (
	"bezbase/internal/models"
)

// UserRepository defines the interface for user data access
type UserRepository interface {
	GetByID(userID uint) (*models.User, error)
	GetByIDWithPreload(userID uint, preloads ...string) (*models.User, error)
	GetAll() ([]models.User, error)
	Search(searchTerm string) ([]models.User, error)
	Create(user *models.User) error
	Update(user *models.User) error
	Delete(userID uint) error
	UpdateStatus(userID uint, status models.UserStatus) error
	VerifyEmail(userID uint) error
}

// UserInfoRepository defines the interface for user info data access
type UserInfoRepository interface {
	GetByUserID(userID uint) (*models.UserInfo, error)
	GetByEmail(email string) (*models.UserInfo, error)
	GetByUsername(username string) (*models.UserInfo, error)
	Create(userInfo *models.UserInfo) error
	Update(userInfo *models.UserInfo) error
	Delete(userID uint) error
	IsEmailTaken(email string, excludeUserID uint) (bool, error)
	IsUsernameTaken(username string, excludeUserID uint) (bool, error)
}

// AuthProviderRepository defines the interface for auth provider data access
type AuthProviderRepository interface {
	GetByUserID(userID uint) ([]models.AuthProvider, error)
	GetByUsernameAndProvider(username string, provider models.AuthProviderType) (*models.AuthProvider, error)
	GetByProviderIDAndType(providerID string, provider models.AuthProviderType) (*models.AuthProvider, error)
	Create(authProvider *models.AuthProvider) error
	Update(authProvider *models.AuthProvider) error
	Delete(userID uint) error
	UpdateEmail(userID uint, provider models.AuthProviderType, newEmail string) error
}

// RoleRepository defines the interface for role data access
type RoleRepository interface {
	GetByID(id uint) (*models.Role, error)
	GetByName(name string) (*models.Role, error)
	GetAll() ([]models.Role, error)
	GetActive() ([]models.Role, error)
	GetWithPagination(page, pageSize int, searchFilter, statusFilter string, isSystemFilter *bool, sortField, sortOrder string) ([]models.Role, int, error)
	Create(role *models.Role) error
	Update(role *models.Role) error
	Delete(role *models.Role) error
}

// RuleRepository defines the interface for RBAC rule data access
type RuleRepository interface {
	GetPermissions(page, pageSize int, roleFilter, resourceFilter, actionFilter, sortField, sortOrder string) ([]models.Rule, int, error)
}