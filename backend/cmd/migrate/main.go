package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"bezbase/internal/config"
	"bezbase/internal/database"
)

func main() {
	var action = flag.String("action", "migrate", "Action to perform: migrate, rollback, status")
	var migrationID = flag.String("to", "", "Migration ID to migrate to (optional)")
	flag.Parse()

	// Load configuration
	cfg := config.Load()

	// Initialize database
	db, err := database.Connect(cfg.Database.URL)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	sqlDB, err := db.DB()
	if err != nil {
		log.Fatalf("Failed to get underlying sql.DB: %v", err)
	}
	defer sqlDB.Close()

	switch *action {
	case "migrate":
		if *migrationID != "" {
			if err := database.MigrateTo(db, *migrationID); err != nil {
				log.Fatalf("Migration failed: %v", err)
			}
		} else {
			if err := database.Migrate(db); err != nil {
				log.Fatalf("Migration failed: %v", err)
			}
		}
		fmt.Println("Migration completed successfully")

	case "rollback":
		if err := database.RollbackMigration(db); err != nil {
			log.Fatalf("Rollback failed: %v", err)
		}
		fmt.Println("Rollback completed successfully")

	case "status":
		if err := database.GetMigrationStatus(db); err != nil {
			log.Fatalf("Failed to get migration status: %v", err)
		}

	default:
		fmt.Printf("Unknown action: %s\n", *action)
		fmt.Println("Available actions: migrate, rollback, status")
		os.Exit(1)
	}
}
