package migrations

import (
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
