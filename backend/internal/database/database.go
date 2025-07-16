package database

import (
	"fmt"
	"log"

	"bezbase/internal/database/migrations"

	"github.com/go-gormigrate/gormigrate/v2"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func Connect(databaseURL string) (*gorm.DB, error) {
	db, err := gorm.Open(postgres.Open(databaseURL), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	sqlDB, err := db.DB()
	if err != nil {
		return nil, fmt.Errorf("failed to get underlying sql.DB: %w", err)
	}

	if err := sqlDB.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	return db, nil
}

func Migrate(db *gorm.DB) error {
	m := gormigrate.New(db, gormigrate.DefaultOptions, migrations.GetMigrations())

	// Set the initial migration for fresh databases
	m.InitSchema(func(tx *gorm.DB) error {
		log.Println("Running initial schema migration...")
		return migrations.GetInitialMigration().Migrate(tx)
	})

	if err := m.Migrate(); err != nil {
		return fmt.Errorf("failed to run migrations: %w", err)
	}

	log.Println("Migrations completed successfully")
	return nil
}

// RollbackMigration rolls back the last migration
func RollbackMigration(db *gorm.DB) error {
	m := gormigrate.New(db, gormigrate.DefaultOptions, migrations.GetMigrations())
	
	if err := m.RollbackLast(); err != nil {
		return fmt.Errorf("failed to rollback migration: %w", err)
	}

	log.Println("Rollback completed successfully")
	return nil
}

// MigrateTo runs migrations up to a specific migration ID
func MigrateTo(db *gorm.DB, migrationID string) error {
	m := gormigrate.New(db, gormigrate.DefaultOptions, migrations.GetMigrations())
	
	if err := m.MigrateTo(migrationID); err != nil {
		return fmt.Errorf("failed to migrate to %s: %w", migrationID, err)
	}

	log.Printf("Migration to %s completed successfully", migrationID)
	return nil
}

// GetMigrationStatus returns the current migration status
func GetMigrationStatus(db *gorm.DB) error {
	// Check if migrations table exists
	if !db.Migrator().HasTable("migrations") {
		log.Println("No migrations table found - database appears to be empty")
		return nil
	}

	// Get all applied migrations
	var migrations []struct {
		ID string `gorm:"column:id"`
	}
	
	if err := db.Table("migrations").Select("id").Find(&migrations).Error; err != nil {
		return fmt.Errorf("failed to get migration status: %w", err)
	}

	log.Printf("Applied migrations (%d):", len(migrations))
	for _, migration := range migrations {
		log.Printf("  - %s", migration.ID)
	}

	return nil
}

