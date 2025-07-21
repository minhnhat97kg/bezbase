package migrations

import (
	"strings"
	"bezbase/internal/models"

	"github.com/go-gormigrate/gormigrate/v2"
	"gorm.io/gorm"
)

// GetMigrations returns all migrations in chronological order
func GetMigrations() []*gormigrate.Migration {
	return []*gormigrate.Migration{
		{
			ID: "20250716_001_create_complete_schema",
			Migrate: func(tx *gorm.DB) error {
				// Create Users table
				type User struct {
					ID            uint         `gorm:"primaryKey"`
					Status        string       `gorm:"not null;default:'pending'"`
					EmailVerified bool         `gorm:"default:false"`
					LastLoginAt   *interface{} `gorm:"type:timestamp"`
					CreatedAt     interface{}  `gorm:"type:timestamp;not null;default:CURRENT_TIMESTAMP"`
					UpdatedAt     interface{}  `gorm:"type:timestamp;not null;default:CURRENT_TIMESTAMP"`
					DeletedAt     *interface{} `gorm:"type:timestamp;index"`
				}

				if err := tx.AutoMigrate(&User{}); err != nil {
					return err
				}

				// Create UserInfo table
				type UserInfo struct {
					ID          uint         `gorm:"primaryKey"`
					UserID      uint         `gorm:"not null;uniqueIndex"`
					FirstName   string       `gorm:"not null"`
					LastName    string       `gorm:"not null"`
					Email       string       `gorm:"not null;uniqueIndex"`
					Avatar      string       `gorm:""`
					Bio         string       `gorm:""`
					Location    string       `gorm:""`
					Website     string       `gorm:""`
					Phone       string       `gorm:""`
					DateOfBirth *interface{} `gorm:"type:timestamp"`
					Gender      string       `gorm:""`
					Timezone    string       `gorm:"default:'UTC'"`
					Language    string       `gorm:"default:'en'"`
					CreatedAt   interface{}  `gorm:"type:timestamp;not null;default:CURRENT_TIMESTAMP"`
					UpdatedAt   interface{}  `gorm:"type:timestamp;not null;default:CURRENT_TIMESTAMP"`
					DeletedAt   *interface{} `gorm:"type:timestamp;index"`
				}

				if err := tx.AutoMigrate(&UserInfo{}); err != nil {
					return err
				}

				// Create AuthProviders table with user_name (not email)
				type AuthProvider struct {
					ID         uint         `gorm:"primaryKey"`
					UserID     uint         `gorm:"not null;index"`
					Provider   string       `gorm:"not null;index"`
					ProviderID string       `gorm:"not null"`
					UserName   string       `gorm:"not null;index"` // Username or email for authentication
					Password   string       `gorm:""`
					Verified   bool         `gorm:"default:false"`
					CreatedAt  interface{}  `gorm:"type:timestamp;not null;default:CURRENT_TIMESTAMP"`
					UpdatedAt  interface{}  `gorm:"type:timestamp;not null;default:CURRENT_TIMESTAMP"`
					DeletedAt  *interface{} `gorm:"type:timestamp;index"`
				}

				if err := tx.AutoMigrate(&AuthProvider{}); err != nil {
					return err
				}

				// Add foreign key constraints
				if err := tx.Exec("ALTER TABLE user_info ADD CONSTRAINT fk_user_info_user_id FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE").Error; err != nil {
					return err
				}

				if err := tx.Exec("ALTER TABLE auth_providers ADD CONSTRAINT fk_auth_providers_user_id FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE").Error; err != nil {
					return err
				}

				// Add composite unique constraint for user_id + provider
				if err := tx.Exec("CREATE UNIQUE INDEX IF NOT EXISTS idx_auth_providers_user_provider ON auth_providers(user_id, provider)").Error; err != nil {
					return err
				}

				// Add additional indexes for performance
				indexes := []string{
					"CREATE INDEX IF NOT EXISTS idx_users_status ON users(status)",
					"CREATE INDEX IF NOT EXISTS idx_users_email_verified ON users(email_verified)",
					"CREATE INDEX IF NOT EXISTS idx_user_info_email ON user_info(email)",
					"CREATE INDEX IF NOT EXISTS idx_auth_providers_user_name ON auth_providers(user_name)",
					"CREATE INDEX IF NOT EXISTS idx_auth_providers_provider ON auth_providers(provider)",
					"CREATE INDEX IF NOT EXISTS idx_auth_providers_provider_id ON auth_providers(provider_id)",
				}

				for _, query := range indexes {
					if err := tx.Exec(query).Error; err != nil {
						return err
					}
				}

				return nil
			},
			Rollback: func(tx *gorm.DB) error {
				// Drop tables in reverse order due to foreign key constraints
				return tx.Migrator().DropTable("auth_providers", "user_info", "users")
			},
		},
		{
			ID: "20250716_002_add_casbin_policies",
			Migrate: func(tx *gorm.DB) error {
				// Create roles table
				type Role struct {
					ID          uint         `gorm:"primaryKey"`
					Name        string       `gorm:"uniqueIndex;not null;size:100"`
					DisplayName string       `gorm:"not null;size:255"`
					Description string       `gorm:"size:500"`
					IsSystem    bool         `gorm:"default:false"`
					IsActive    bool         `gorm:"default:true"`
					CreatedAt   interface{}  `gorm:"type:timestamp;not null;default:CURRENT_TIMESTAMP"`
					UpdatedAt   interface{}  `gorm:"type:timestamp;not null;default:CURRENT_TIMESTAMP"`
					DeletedAt   *interface{} `gorm:"type:timestamp;index"`
				}

				if err := tx.AutoMigrate(&Role{}); err != nil {
					return err
				}

				// Add indexes for better performance
				indexes := []string{
					"CREATE INDEX IF NOT EXISTS idx_roles_name ON roles(name)",
					"CREATE INDEX IF NOT EXISTS idx_roles_is_active ON roles(is_active)",
					"CREATE INDEX IF NOT EXISTS idx_roles_is_system ON roles(is_system)",
				}

				for _, query := range indexes {
					if err := tx.Exec(query).Error; err != nil {
						return err
					}
				}

				// Insert default roles
				defaultRoles := []Role{
					{
						Name:        "admin",
						DisplayName: "Administrator",
						Description: "Full system access with all permissions",
						IsSystem:    true,
						IsActive:    true,
					},
					{
						Name:        "user",
						DisplayName: "Standard User",
						Description: "Basic user access with limited permissions",
						IsSystem:    true,
						IsActive:    true,
					},
					{
						Name:        "editor",
						DisplayName: "Content Editor",
						Description: "Can create and edit content",
						IsSystem:    false,
						IsActive:    true,
					},
				}

				for _, role := range defaultRoles {
					if err := tx.Create(&role).Error; err != nil {
						return err
					}
				}

				return nil
			},
			Rollback: func(tx *gorm.DB) error {
				return tx.Migrator().DropTable("roles")
			},
		},
		// {
		// 	ID: "20250717_001_add_username_to_user_info",
		// 	Migrate: func(tx *gorm.DB) error {
		// 		// Add username column to user_info table
		// 		if err := tx.Exec("ALTER TABLE user_info ADD COLUMN username VARCHAR(255) NOT NULL DEFAULT ''").Error; err != nil {
		// 			return err
		// 		}
		//
		// 		// Create unique index on username
		// 		if err := tx.Exec("CREATE UNIQUE INDEX idx_user_info_username ON user_info(username)").Error; err != nil {
		// 			return err
		// 		}
		//
		// 		return nil
		// 	},
		// 	Rollback: func(tx *gorm.DB) error {
		// 		// Drop the username column
		// 		return tx.Exec("ALTER TABLE user_info DROP COLUMN username").Error
		// 	},
		// },
		{
			ID: "20250717_002_add_email_verification_tokens",
			Migrate: func(tx *gorm.DB) error {
				// Create EmailVerificationToken table
				type EmailVerificationToken struct {
					ID        uint         `gorm:"primaryKey"`
					UserID    uint         `gorm:"not null;index"`
					Token     string       `gorm:"not null;uniqueIndex;size:255"`
					Email     string       `gorm:"not null;size:255"`
					ExpiresAt interface{}  `gorm:"type:timestamp;not null"`
					UsedAt    *interface{} `gorm:"type:timestamp"`
					CreatedAt interface{}  `gorm:"type:timestamp;not null;default:CURRENT_TIMESTAMP"`
					UpdatedAt interface{}  `gorm:"type:timestamp;not null;default:CURRENT_TIMESTAMP"`
					DeletedAt *interface{} `gorm:"type:timestamp;index"`
				}

				if err := tx.AutoMigrate(&EmailVerificationToken{}); err != nil {
					return err
				}

				// Add foreign key constraint
				if err := tx.Exec("ALTER TABLE email_verification_tokens ADD CONSTRAINT fk_email_verification_tokens_user_id FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE").Error; err != nil {
					return err
				}

				// Add indexes for performance
				indexes := []string{
					"CREATE INDEX IF NOT EXISTS idx_email_verification_tokens_user_id ON email_verification_tokens(user_id)",
					"CREATE INDEX IF NOT EXISTS idx_email_verification_tokens_token ON email_verification_tokens(token)",
					"CREATE INDEX IF NOT EXISTS idx_email_verification_tokens_expires_at ON email_verification_tokens(expires_at)",
					"CREATE INDEX IF NOT EXISTS idx_email_verification_tokens_used_at ON email_verification_tokens(used_at)",
				}

				for _, query := range indexes {
					if err := tx.Exec(query).Error; err != nil {
						return err
					}
				}

				return nil
			},
			Rollback: func(tx *gorm.DB) error {
				return tx.Migrator().DropTable("email_verification_tokens")
			},
		},
		{
			ID: "20250717_003_add_password_reset_tokens",
			Migrate: func(tx *gorm.DB) error {
				// Create PasswordResetToken table
				type PasswordResetToken struct {
					ID        uint         `gorm:"primaryKey"`
					UserID    uint         `gorm:"not null;index"`
					Token     string       `gorm:"not null;uniqueIndex;size:255"`
					Email     string       `gorm:"not null;size:255"`
					ExpiresAt interface{}  `gorm:"type:timestamp;not null"`
					UsedAt    *interface{} `gorm:"type:timestamp"`
					CreatedAt interface{}  `gorm:"type:timestamp;not null;default:CURRENT_TIMESTAMP"`
					UpdatedAt interface{}  `gorm:"type:timestamp;not null;default:CURRENT_TIMESTAMP"`
					DeletedAt *interface{} `gorm:"type:timestamp;index"`
				}

				if err := tx.AutoMigrate(&PasswordResetToken{}); err != nil {
					return err
				}

				// Add foreign key constraint
				if err := tx.Exec("ALTER TABLE password_reset_tokens ADD CONSTRAINT fk_password_reset_tokens_user_id FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE").Error; err != nil {
					return err
				}

				// Add indexes for performance
				indexes := []string{
					"CREATE INDEX IF NOT EXISTS idx_password_reset_tokens_user_id ON password_reset_tokens(user_id)",
					"CREATE INDEX IF NOT EXISTS idx_password_reset_tokens_token ON password_reset_tokens(token)",
					"CREATE INDEX IF NOT EXISTS idx_password_reset_tokens_expires_at ON password_reset_tokens(expires_at)",
					"CREATE INDEX IF NOT EXISTS idx_password_reset_tokens_used_at ON password_reset_tokens(used_at)",
				}

				for _, query := range indexes {
					if err := tx.Exec(query).Error; err != nil {
						return err
					}
				}

				return nil
			},
			Rollback: func(tx *gorm.DB) error {
				return tx.Migrator().DropTable("password_reset_tokens")
			},
		},
		{
			ID: "20250717_004_add_multi_tenancy",
			Migrate: func(tx *gorm.DB) error {
				// Create Organizations table
				type Organization struct {
					ID          uint         `gorm:"primaryKey"`
					Name        string       `gorm:"not null;size:255"`
					Slug        string       `gorm:"uniqueIndex;not null;size:100"`
					Domain      string       `gorm:"size:255"`
					Settings    string       `gorm:"type:jsonb"`
					IsActive    bool         `gorm:"default:true"`
					PlanType    string       `gorm:"default:'free';size:50"`
					CreatedAt   interface{}  `gorm:"type:timestamp;not null;default:CURRENT_TIMESTAMP"`
					UpdatedAt   interface{}  `gorm:"type:timestamp;not null;default:CURRENT_TIMESTAMP"`
					DeletedAt   *interface{} `gorm:"type:timestamp;index"`
				}

				if err := tx.AutoMigrate(&Organization{}); err != nil {
					return err
				}

				// Create OrganizationUsers table for many-to-many relationship
				type OrganizationUser struct {
					ID           uint         `gorm:"primaryKey"`
					OrgID        uint         `gorm:"not null;index"`
					UserID       uint         `gorm:"not null;index"`
					Role         string       `gorm:"not null;default:'member';size:50"`
					IsPrimary    bool         `gorm:"default:false"`
					JoinedAt     interface{}  `gorm:"type:timestamp;not null;default:CURRENT_TIMESTAMP"`
					CreatedAt    interface{}  `gorm:"type:timestamp;not null;default:CURRENT_TIMESTAMP"`
					UpdatedAt    interface{}  `gorm:"type:timestamp;not null;default:CURRENT_TIMESTAMP"`
					DeletedAt    *interface{} `gorm:"type:timestamp;index"`
				}

				if err := tx.AutoMigrate(&OrganizationUser{}); err != nil {
					return err
				}

				// Create OrganizationInvitations table
				type OrganizationInvitation struct {
					ID          uint         `gorm:"primaryKey"`
					OrgID       uint         `gorm:"not null;index"`
					Email       string       `gorm:"not null;size:255"`
					Role        string       `gorm:"not null;default:'member';size:50"`
					Token       string       `gorm:"not null;uniqueIndex;size:255"`
					ExpiresAt   interface{}  `gorm:"type:timestamp;not null"`
					InvitedBy   uint         `gorm:"not null"`
					AcceptedAt  *interface{} `gorm:"type:timestamp"`
					CreatedAt   interface{}  `gorm:"type:timestamp;not null;default:CURRENT_TIMESTAMP"`
					UpdatedAt   interface{}  `gorm:"type:timestamp;not null;default:CURRENT_TIMESTAMP"`
					DeletedAt   *interface{} `gorm:"type:timestamp;index"`
				}

				if err := tx.AutoMigrate(&OrganizationInvitation{}); err != nil {
					return err
				}

				// Add current_org_id to users table if it doesn't exist
				var columnExists bool
				err := tx.Raw("SELECT EXISTS (SELECT 1 FROM information_schema.columns WHERE table_name = 'users' AND column_name = 'current_org_id')").Scan(&columnExists).Error
				if err != nil {
					return err
				}
				if !columnExists {
					if err := tx.Exec("ALTER TABLE users ADD COLUMN current_org_id INTEGER").Error; err != nil {
						return err
					}
				}

				// Add org_id to roles table for organization-specific roles if it doesn't exist
				err = tx.Raw("SELECT EXISTS (SELECT 1 FROM information_schema.columns WHERE table_name = 'roles' AND column_name = 'org_id')").Scan(&columnExists).Error
				if err != nil {
					return err
				}
				if !columnExists {
					if err := tx.Exec("ALTER TABLE roles ADD COLUMN org_id INTEGER").Error; err != nil {
						return err
					}
				}

				// Add foreign key constraints
				constraints := []string{
					"ALTER TABLE organization_users ADD CONSTRAINT fk_organization_users_org_id FOREIGN KEY (org_id) REFERENCES organizations(id) ON DELETE CASCADE",
					"ALTER TABLE organization_users ADD CONSTRAINT fk_organization_users_user_id FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE",
					"ALTER TABLE organization_invitations ADD CONSTRAINT fk_organization_invitations_org_id FOREIGN KEY (org_id) REFERENCES organizations(id) ON DELETE CASCADE",
					"ALTER TABLE organization_invitations ADD CONSTRAINT fk_organization_invitations_invited_by FOREIGN KEY (invited_by) REFERENCES users(id) ON DELETE CASCADE",
					"ALTER TABLE users ADD CONSTRAINT fk_users_current_org_id FOREIGN KEY (current_org_id) REFERENCES organizations(id) ON DELETE SET NULL",
					"ALTER TABLE roles ADD CONSTRAINT fk_roles_org_id FOREIGN KEY (org_id) REFERENCES organizations(id) ON DELETE CASCADE",
				}

				for _, constraint := range constraints {
					if err := tx.Exec(constraint).Error; err != nil {
						// Ignore constraint already exists errors
						if !strings.Contains(err.Error(), "already exists") {
							return err
						}
					}
				}

				// Add unique constraint for user-organization relationship
				if err := tx.Exec("CREATE UNIQUE INDEX IF NOT EXISTS idx_organization_users_org_user ON organization_users(org_id, user_id) WHERE deleted_at IS NULL").Error; err != nil {
					return err
				}

				// Add indexes for performance
				indexes := []string{
					"CREATE INDEX IF NOT EXISTS idx_organizations_slug ON organizations(slug)",
					"CREATE INDEX IF NOT EXISTS idx_organizations_domain ON organizations(domain)",
					"CREATE INDEX IF NOT EXISTS idx_organizations_is_active ON organizations(is_active)",
					"CREATE INDEX IF NOT EXISTS idx_organization_users_org_id ON organization_users(org_id)",
					"CREATE INDEX IF NOT EXISTS idx_organization_users_user_id ON organization_users(user_id)",
					"CREATE INDEX IF NOT EXISTS idx_organization_users_role ON organization_users(role)",
					"CREATE INDEX IF NOT EXISTS idx_organization_invitations_email ON organization_invitations(email)",
					"CREATE INDEX IF NOT EXISTS idx_organization_invitations_token ON organization_invitations(token)",
					"CREATE INDEX IF NOT EXISTS idx_organization_invitations_expires_at ON organization_invitations(expires_at)",
					"CREATE INDEX IF NOT EXISTS idx_users_current_org_id ON users(current_org_id)",
					"CREATE INDEX IF NOT EXISTS idx_roles_org_id ON roles(org_id)",
				}

				for _, query := range indexes {
					if err := tx.Exec(query).Error; err != nil {
						return err
					}
				}

				// Create default personal organization for existing users
				// First, create a default organization
				defaultOrg := Organization{
					Name:     "Default Organization",
					Slug:     "default",
					Settings: "{}",
					IsActive: true,
					PlanType: "free",
				}

				if err := tx.Create(&defaultOrg).Error; err != nil {
					return err
				}

				// Add all existing users to this default organization as owners
				if err := tx.Exec(`
					INSERT INTO organization_users (org_id, user_id, role, is_primary, joined_at, created_at, updated_at)
					SELECT ?, id, 'owner', true, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP
					FROM users
					WHERE deleted_at IS NULL
				`, defaultOrg.ID).Error; err != nil {
					return err
				}

				// Set current_org_id for all existing users
				if err := tx.Exec("UPDATE users SET current_org_id = ? WHERE deleted_at IS NULL", defaultOrg.ID).Error; err != nil {
					return err
				}

				return nil
			},
			Rollback: func(tx *gorm.DB) error {
				// Remove foreign key constraints first
				constraints := []string{
					"ALTER TABLE roles DROP CONSTRAINT IF EXISTS fk_roles_org_id",
					"ALTER TABLE users DROP CONSTRAINT IF EXISTS fk_users_current_org_id",
					"ALTER TABLE organization_invitations DROP CONSTRAINT IF EXISTS fk_organization_invitations_invited_by",
					"ALTER TABLE organization_invitations DROP CONSTRAINT IF EXISTS fk_organization_invitations_org_id",
					"ALTER TABLE organization_users DROP CONSTRAINT IF EXISTS fk_organization_users_user_id",
					"ALTER TABLE organization_users DROP CONSTRAINT IF EXISTS fk_organization_users_org_id",
				}

				for _, constraint := range constraints {
					tx.Exec(constraint)
				}

				// Drop columns
				tx.Exec("ALTER TABLE roles DROP COLUMN IF EXISTS org_id")
				tx.Exec("ALTER TABLE users DROP COLUMN IF EXISTS current_org_id")

				// Drop tables in reverse order
				return tx.Migrator().DropTable("organization_invitations", "organization_users", "organizations")
			},
		},
		{
			ID: "20250717_005_add_advanced_rbac",
			Migrate: func(tx *gorm.DB) error {
				// Add parent_role_id and hierarchy_level to roles table
				if err := tx.Exec("ALTER TABLE roles ADD COLUMN parent_role_id INTEGER").Error; err != nil {
					return err
				}
				if err := tx.Exec("ALTER TABLE roles ADD COLUMN hierarchy_level INTEGER DEFAULT 0").Error; err != nil {
					return err
				}

				// Create RoleInheritance table for efficient role hierarchy queries
				type RoleInheritance struct {
					ID            uint `gorm:"primaryKey"`
					ParentRoleID  uint `gorm:"not null;index"`
					ChildRoleID   uint `gorm:"not null;index"`
					Depth         int  `gorm:"not null;default:1"`
					CreatedAt     interface{} `gorm:"type:timestamp;not null;default:CURRENT_TIMESTAMP"`
					UpdatedAt     interface{} `gorm:"type:timestamp;not null;default:CURRENT_TIMESTAMP"`
				}

				if err := tx.AutoMigrate(&RoleInheritance{}); err != nil {
					return err
				}

				// Create ContextualPermissions table for context-aware permissions
				type ContextualPermission struct {
					ID           uint         `gorm:"primaryKey"`
					RoleID       uint         `gorm:"not null;index"`
					Resource     string       `gorm:"not null;size:100"`
					Action       string       `gorm:"not null;size:50"`
					ContextType  string       `gorm:"size:50"` // e.g., "organization", "department", "project"
					ContextValue string       `gorm:"size:255"` // specific ID or pattern
					IsGranted    bool         `gorm:"default:true"`
					CreatedAt    interface{}  `gorm:"type:timestamp;not null;default:CURRENT_TIMESTAMP"`
					UpdatedAt    interface{}  `gorm:"type:timestamp;not null;default:CURRENT_TIMESTAMP"`
					DeletedAt    *interface{} `gorm:"type:timestamp;index"`
				}

				if err := tx.AutoMigrate(&ContextualPermission{}); err != nil {
					return err
				}

				// Create RoleTemplates table for predefined role configurations
				type RoleTemplate struct {
					ID          uint         `gorm:"primaryKey"`
					Name        string       `gorm:"not null;size:100"`
					DisplayName string       `gorm:"not null;size:255"`
					Description string       `gorm:"size:500"`
					Category    string       `gorm:"size:100"` // e.g., "system", "business", "department"
					Config      string       `gorm:"type:jsonb"` // JSON configuration for permissions
					IsActive    bool         `gorm:"default:true"`
					CreatedAt   interface{}  `gorm:"type:timestamp;not null;default:CURRENT_TIMESTAMP"`
					UpdatedAt   interface{}  `gorm:"type:timestamp;not null;default:CURRENT_TIMESTAMP"`
					DeletedAt   *interface{} `gorm:"type:timestamp;index"`
				}

				if err := tx.AutoMigrate(&RoleTemplate{}); err != nil {
					return err
				}

				// Add foreign key constraints
				constraints := []string{
					"ALTER TABLE roles ADD CONSTRAINT fk_roles_parent_role_id FOREIGN KEY (parent_role_id) REFERENCES roles(id) ON DELETE SET NULL",
					"ALTER TABLE role_inheritances ADD CONSTRAINT fk_role_inheritances_parent_role_id FOREIGN KEY (parent_role_id) REFERENCES roles(id) ON DELETE CASCADE",
					"ALTER TABLE role_inheritances ADD CONSTRAINT fk_role_inheritances_child_role_id FOREIGN KEY (child_role_id) REFERENCES roles(id) ON DELETE CASCADE",
					"ALTER TABLE contextual_permissions ADD CONSTRAINT fk_contextual_permissions_role_id FOREIGN KEY (role_id) REFERENCES roles(id) ON DELETE CASCADE",
				}

				for _, constraint := range constraints {
					if err := tx.Exec(constraint).Error; err != nil {
						return err
					}
				}

				// Add unique constraints
				uniqueConstraints := []string{
					"CREATE UNIQUE INDEX IF NOT EXISTS idx_role_inheritances_parent_child ON role_inheritances(parent_role_id, child_role_id)",
					"CREATE UNIQUE INDEX IF NOT EXISTS idx_contextual_permissions_role_resource_action_context ON contextual_permissions(role_id, resource, action, context_type, context_value) WHERE deleted_at IS NULL",
					"CREATE UNIQUE INDEX IF NOT EXISTS idx_role_templates_name ON role_templates(name) WHERE deleted_at IS NULL",
				}

				for _, constraint := range uniqueConstraints {
					if err := tx.Exec(constraint).Error; err != nil {
						return err
					}
				}

				// Add indexes for performance
				indexes := []string{
					"CREATE INDEX IF NOT EXISTS idx_roles_parent_role_id ON roles(parent_role_id)",
					"CREATE INDEX IF NOT EXISTS idx_roles_hierarchy_level ON roles(hierarchy_level)",
					"CREATE INDEX IF NOT EXISTS idx_role_inheritances_parent_role_id ON role_inheritances(parent_role_id)",
					"CREATE INDEX IF NOT EXISTS idx_role_inheritances_child_role_id ON role_inheritances(child_role_id)",
					"CREATE INDEX IF NOT EXISTS idx_role_inheritances_depth ON role_inheritances(depth)",
					"CREATE INDEX IF NOT EXISTS idx_contextual_permissions_role_id ON contextual_permissions(role_id)",
					"CREATE INDEX IF NOT EXISTS idx_contextual_permissions_resource ON contextual_permissions(resource)",
					"CREATE INDEX IF NOT EXISTS idx_contextual_permissions_action ON contextual_permissions(action)",
					"CREATE INDEX IF NOT EXISTS idx_contextual_permissions_context_type ON contextual_permissions(context_type)",
					"CREATE INDEX IF NOT EXISTS idx_role_templates_category ON role_templates(category)",
					"CREATE INDEX IF NOT EXISTS idx_role_templates_is_active ON role_templates(is_active)",
				}

				for _, query := range indexes {
					if err := tx.Exec(query).Error; err != nil {
						return err
					}
				}

				// Create default role templates
				defaultTemplates := []RoleTemplate{
					{
						Name:        "organization_admin",
						DisplayName: "Organization Administrator",
						Description: "Full administrative access within an organization",
						Category:    "system",
						Config:      `{"permissions": ["*:*"], "context": "organization", "inheritable": true}`,
						IsActive:    true,
					},
					{
						Name:        "team_lead",
						DisplayName: "Team Lead",
						Description: "Team management and project oversight",
						Category:    "business",
						Config:      `{"permissions": ["users:read", "users:update", "projects:*"], "context": "team", "inheritable": true}`,
						IsActive:    true,
					},
					{
						Name:        "project_manager",
						DisplayName: "Project Manager",
						Description: "Project management and coordination",
						Category:    "business",
						Config:      `{"permissions": ["projects:*", "reports:read"], "context": "project", "inheritable": false}`,
						IsActive:    true,
					},
					{
						Name:        "viewer",
						DisplayName: "Viewer",
						Description: "Read-only access to assigned resources",
						Category:    "basic",
						Config:      `{"permissions": ["*:read"], "context": "resource", "inheritable": false}`,
						IsActive:    true,
					},
				}

				for _, template := range defaultTemplates {
					if err := tx.Create(&template).Error; err != nil {
						return err
					}
				}

				// Update existing roles with hierarchy levels
				// Admin roles get level 0 (highest)
				if err := tx.Exec("UPDATE roles SET hierarchy_level = 0 WHERE name = 'admin'").Error; err != nil {
					return err
				}

				// Create organization-specific admin roles that inherit from global admin
				if err := tx.Exec(`
					INSERT INTO roles (name, display_name, description, is_system, is_active, org_id, parent_role_id, hierarchy_level, created_at, updated_at)
					SELECT 
						CONCAT('org_admin_', o.id),
						'Organization Administrator',
						CONCAT('Administrator for ', o.name),
						false,
						true,
						o.id,
						r.id,
						1,
						CURRENT_TIMESTAMP,
						CURRENT_TIMESTAMP
					FROM organizations o
					CROSS JOIN roles r
					WHERE r.name = 'admin' AND r.is_system = true
				`).Error; err != nil {
					return err
				}

				return nil
			},
			Rollback: func(tx *gorm.DB) error {
				// Remove foreign key constraints first
				constraints := []string{
					"ALTER TABLE contextual_permissions DROP CONSTRAINT IF EXISTS fk_contextual_permissions_role_id",
					"ALTER TABLE role_inheritances DROP CONSTRAINT IF EXISTS fk_role_inheritances_child_role_id",
					"ALTER TABLE role_inheritances DROP CONSTRAINT IF EXISTS fk_role_inheritances_parent_role_id",
					"ALTER TABLE roles DROP CONSTRAINT IF EXISTS fk_roles_parent_role_id",
				}

				for _, constraint := range constraints {
					tx.Exec(constraint)
				}

				// Drop columns from roles table
				tx.Exec("ALTER TABLE roles DROP COLUMN IF EXISTS hierarchy_level")
				tx.Exec("ALTER TABLE roles DROP COLUMN IF EXISTS parent_role_id")

				// Drop tables in reverse order
				return tx.Migrator().DropTable("role_templates", "contextual_permissions", "role_inheritances")
			},
		},
		{
			ID: "20250720_001_remove_organizations",
			Migrate: func(tx *gorm.DB) error {
				// Remove foreign key constraints first
				constraints := []string{
					"ALTER TABLE roles DROP CONSTRAINT IF EXISTS fk_roles_org_id",
					"ALTER TABLE users DROP CONSTRAINT IF EXISTS fk_users_current_org_id",
					"ALTER TABLE organization_invitations DROP CONSTRAINT IF EXISTS fk_organization_invitations_invited_by",
					"ALTER TABLE organization_invitations DROP CONSTRAINT IF EXISTS fk_organization_invitations_org_id",
					"ALTER TABLE organization_users DROP CONSTRAINT IF EXISTS fk_organization_users_user_id",
					"ALTER TABLE organization_users DROP CONSTRAINT IF EXISTS fk_organization_users_org_id",
				}

				for _, constraint := range constraints {
					tx.Exec(constraint)
				}

				// Drop columns from users and roles tables
				tx.Exec("ALTER TABLE users DROP COLUMN IF EXISTS current_org_id")
				tx.Exec("ALTER TABLE roles DROP COLUMN IF EXISTS org_id")

				// Drop organization-related tables
				tx.Migrator().DropTable("organization_invitations", "organization_users", "organizations")

				return nil
			},
			Rollback: func(tx *gorm.DB) error {
				// This rollback would recreate the organization system - complex operation
				// For now, just return nil as this is meant to be a permanent removal
				return nil
			},
		},
		{
			ID: "20250718_001_seed_roles",
			Migrate: func(tx *gorm.DB) error {
				// Check if roles already exist, if not create them
				var count int64
				if err := tx.Model(&models.Role{}).Count(&count).Error; err != nil {
					return err
				}
				
				// Only seed if no roles exist
				if count == 0 {
					roles := []models.Role{
						{
							Name:           "admin",
							DisplayName:    "System Administrator",
							Description:    "Full system access with all permissions",
							IsSystem:       true,
							IsActive:       true,
							HierarchyLevel: 0,
						},
						{
							Name:           "user",
							DisplayName:    "Standard User",
							Description:    "Basic user access with limited permissions",
							IsSystem:       true,
							IsActive:       true,
							HierarchyLevel: 2,
						},
						{
							Name:           "editor",
							DisplayName:    "Content Editor",
							Description:    "Can create and edit content",
							IsSystem:       false,
							IsActive:       true,
							HierarchyLevel: 1,
						},
						{
							Name:           "viewer",
							DisplayName:    "Viewer",
							Description:    "Read-only access to resources",
							IsSystem:       false,
							IsActive:       true,
							HierarchyLevel: 3,
						},
					}
					
					for _, role := range roles {
						if err := tx.Create(&role).Error; err != nil {
							return err
						}
					}
				}
				
				return nil
			},
			Rollback: func(tx *gorm.DB) error {
				// Don't delete system roles, just leave them
				return nil
			},
		},
	}
}

// GetInitialMigration returns the initial migration that creates all tables at once
// This is used for fresh installations
func GetInitialMigration() *gormigrate.Migration {
	return &gormigrate.Migration{
		ID: "SCHEMA_INIT",
		Migrate: func(tx *gorm.DB) error {
			// Create all tables using the actual models
			return tx.AutoMigrate(
				&models.User{},
				&models.UserInfo{},
				&models.AuthProvider{},
				&models.EmailVerificationToken{},
				&models.PasswordResetToken{},
			)
		},
		Rollback: func(tx *gorm.DB) error {
			return tx.Migrator().DropTable(
				&models.PasswordResetToken{},
				&models.EmailVerificationToken{},
				&models.AuthProvider{},
				&models.UserInfo{},
				&models.User{},
			)
		},
	}
}
