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
					ID            uint           `gorm:"primaryKey"`
					Status        string         `gorm:"not null;default:'pending'"`
					EmailVerified bool           `gorm:"default:false"`
					LastLoginAt   *interface{}   `gorm:"type:timestamp"`
					CreatedAt     interface{}    `gorm:"type:timestamp;not null;default:CURRENT_TIMESTAMP"`
					UpdatedAt     interface{}    `gorm:"type:timestamp;not null;default:CURRENT_TIMESTAMP"`
					DeletedAt     *interface{}   `gorm:"type:timestamp;index"`
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
					ID         uint        `gorm:"primaryKey"`
					UserID     uint        `gorm:"not null;index"`
					Provider   string      `gorm:"not null;index"`
					ProviderID string      `gorm:"not null"`
					UserName   string      `gorm:"not null;index"` // Username or email for authentication
					Password   string      `gorm:""`
					Verified   bool        `gorm:"default:false"`
					CreatedAt  interface{} `gorm:"type:timestamp;not null;default:CURRENT_TIMESTAMP"`
					UpdatedAt  interface{} `gorm:"type:timestamp;not null;default:CURRENT_TIMESTAMP"`
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
			)
		},
		Rollback: func(tx *gorm.DB) error {
			return tx.Migrator().DropTable(
				&models.AuthProvider{},
				&models.UserInfo{},
				&models.User{},
			)
		},
	}
}