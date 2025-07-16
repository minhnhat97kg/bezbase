# Database Migrations

This project uses [gormigrate](https://github.com/go-gormigrate/gormigrate) for database schema migrations. This provides versioned, reversible database migrations with better control than GORM's AutoMigrate.

## Migration Files

All migrations are defined in `/internal/database/migrations/migrations.go`. Each migration has:
- **ID**: Unique identifier (format: `YYYYMMDD_NNN_description`)
- **Migrate**: Function that applies the migration
- **Rollback**: Function that reverses the migration

## Current Migrations

1. **20250716_001_create_complete_schema** - Creates the complete database schema with all tables, indexes, and constraints:
   - **users table**: Core user entity with status, email verification, and login tracking
   - **user_info table**: Extended user profile information (name, bio, location, etc.)
   - **auth_providers table**: Multi-provider authentication support with username field
   - **Foreign key constraints**: Maintains referential integrity between tables
   - **Indexes**: Optimized for authentication and profile queries
   - **Composite unique constraint**: Ensures one provider per user per provider type

## Migration Commands

### Using the CLI Tool

Build the migration tool:
```bash
go build -o migrate-tool ./cmd/migrate/main.go
```

Check migration status:
```bash
DATABASE_URL="postgresql://bezbase_user:bezbase_password@localhost:5432/bezbase?sslmode=disable" ./migrate-tool -action=status
```

Run all pending migrations:
```bash
DATABASE_URL="postgresql://bezbase_user:bezbase_password@localhost:5432/bezbase?sslmode=disable" ./migrate-tool -action=migrate
```

Migrate to a specific migration:
```bash
DATABASE_URL="postgresql://bezbase_user:bezbase_password@localhost:5432/bezbase?sslmode=disable" ./migrate-tool -action=migrate -to=20250716_002_create_user_info_table
```

Rollback the last migration:
```bash
DATABASE_URL="postgresql://bezbase_user:bezbase_password@localhost:5432/bezbase?sslmode=disable" ./migrate-tool -action=rollback
```

### Programmatic Usage

```go
import (
    "bezbase/internal/database"
    "bezbase/internal/config"
)

// Load config and connect to database
cfg := config.Load()
db, err := database.Connect(cfg.DatabaseURL)
if err != nil {
    log.Fatal(err)
}

// Run migrations
if err := database.Migrate(db); err != nil {
    log.Fatal(err)
}

// Check migration status
if err := database.GetMigrationStatus(db); err != nil {
    log.Fatal(err)
}

// Rollback last migration
if err := database.RollbackMigration(db); err != nil {
    log.Fatal(err)
}
```

## Adding New Migrations

1. Add a new migration to the slice in `migrations.go`:

```go
{
    ID: "20250716_005_add_user_preferences",
    Migrate: func(tx *gorm.DB) error {
        // Add your migration logic here
        return tx.AutoMigrate(&UserPreferences{})
    },
    Rollback: func(tx *gorm.DB) error {
        // Add your rollback logic here
        return tx.Migrator().DropTable("user_preferences")
    },
}
```

2. Use timestamps in the format `YYYYMMDD_NNN` for ordering
3. Always provide both `Migrate` and `Rollback` functions
4. Test your migrations thoroughly before deploying

## Best Practices

- **Always backup your database** before running migrations in production
- **Test migrations locally** with the same database version as production
- **Use transactions** for complex migrations to ensure atomicity
- **Keep migrations small and focused** - one logical change per migration
- **Never modify existing migrations** - create new ones instead
- **Use descriptive names** for migration IDs

## Migration Table

Gormigrate creates a `migrations` table to track which migrations have been applied:

```sql
CREATE TABLE migrations (
    id VARCHAR(255) PRIMARY KEY
);
```

This table contains the IDs of all successfully applied migrations.

## Troubleshooting

### Migration Failed
If a migration fails, check:
1. Database connection and permissions
2. Migration syntax and logic
3. Database logs for specific errors
4. Whether the migration was partially applied

### Rollback Issues
If rollback fails:
1. Check the rollback function implementation
2. Ensure the database state allows rollback
3. May need to manually fix the database state

### Schema Conflicts
If you encounter schema conflicts:
1. Check if migrations are in the correct order
2. Ensure foreign key constraints are properly handled
3. Consider using database transactions in complex migrations