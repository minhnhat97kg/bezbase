package database

import (
	"fmt"
	"log"
	"os"
	"testing"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// TestDB creates a test database connection using SQLite in memory
func TestDB(t *testing.T) *gorm.DB {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
	if err != nil {
		t.Fatalf("Failed to connect to test database: %v", err)
	}

	// Run migrations
	if err := Migrate(db); err != nil {
		t.Fatalf("Failed to migrate test database: %v", err)
	}

	return db
}

// SetupTestDB creates a test database and returns cleanup function
func SetupTestDB() (*gorm.DB, func()) {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
	if err != nil {
		log.Fatalf("Failed to connect to test database: %v", err)
	}

	// Run migrations
	if err := Migrate(db); err != nil {
		log.Fatalf("Failed to migrate test database: %v", err)
	}

	cleanup := func() {
		sqlDB, _ := db.DB()
		sqlDB.Close()
	}

	return db, cleanup
}

// IntegrationTestDB creates a test PostgreSQL database for integration tests
func IntegrationTestDB(t *testing.T) *gorm.DB {
	// Skip integration tests if not in CI or if SKIP_INTEGRATION_TESTS is set
	if os.Getenv("CI") == "" && os.Getenv("SKIP_INTEGRATION_TESTS") != "" {
		t.Skip("Skipping integration tests")
	}

	// Use test database URL or create one
	dbURL := os.Getenv("TEST_DATABASE_URL")
	if dbURL == "" {
		dbURL = "postgres://bezbase_user:bezbase_password@localhost/bezbase_test?sslmode=disable"
	}

	db, err := Connect(dbURL)
	if err != nil {
		t.Skipf("Failed to connect to integration test database: %v", err)
	}

	// Run migrations
	if err := Migrate(db); err != nil {
		t.Fatalf("Failed to migrate integration test database: %v", err)
	}

	return db
}

// CleanupTestDB removes all data from test tables
func CleanupTestDB(db *gorm.DB) {
	// Delete all data from tables in reverse order to avoid foreign key constraints
	tables := []string{
		"rules",
		"auth_providers",
		"user_infos",
		"users",
		"roles",
	}

	for _, table := range tables {
		db.Exec(fmt.Sprintf("DELETE FROM %s", table))
	}
}

// CreateTestUser creates a test user for testing purposes
func CreateTestUser(db *gorm.DB) {
	// This function can be used to create standard test data
	// Implementation depends on your specific test needs
}

// TestConfig holds test configuration
type TestConfig struct {
	DatabaseURL string
	JWTSecret   string
}

// GetTestConfig returns test configuration
func GetTestConfig() *TestConfig {
	return &TestConfig{
		DatabaseURL: os.Getenv("TEST_DATABASE_URL"),
		JWTSecret:   "test-jwt-secret-key",
	}
}