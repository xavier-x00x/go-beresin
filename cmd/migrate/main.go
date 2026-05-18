package main

import (
	"context"
	"log"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/joho/godotenv"
	"go-beresin/pkg/database"
)

func main() {
	// Load .env file
	_ = godotenv.Load()

	if len(os.Args) < 2 {
		log.Fatalf("Usage: go run cmd/migrate/main.go [up|down]")
	}

	command := os.Args[1]
	if command != "up" && command != "down" {
		log.Fatalf("Invalid command: %s. Use 'up' or 'down'.", command)
	}

	ctx := context.Background()

	// Initialize database connection pool
	pool, err := database.InitPool(ctx)
	if err != nil {
		log.Fatalf("Database initialization failed: %v", err)
	}
	defer pool.Close()

	// Ensure the schema_migrations table exists
	_, err = pool.Exec(ctx, `
		CREATE TABLE IF NOT EXISTS schema_migrations (
			version VARCHAR(255) PRIMARY KEY,
			migrated_at TIMESTAMPTZ DEFAULT NOW()
		);
	`)
	if err != nil {
		log.Fatalf("Failed to initialize schema_migrations table: %v", err)
	}

	// Read migration files from the migrations directory
	migrationsDir := "./migrations"
	files, err := os.ReadDir(migrationsDir)
	if err != nil {
		log.Fatalf("Failed to read migrations directory: %v", err)
	}

	var sqlFiles []string
	suffix := "." + command + ".sql"
	for _, file := range files {
		if !file.IsDir() && strings.HasSuffix(file.Name(), suffix) {
			sqlFiles = append(sqlFiles, file.Name())
		}
	}

	// Sort files to ensure correct order
	if command == "down" {
		sort.Sort(sort.Reverse(sort.StringSlice(sqlFiles)))
	} else {
		sort.Strings(sqlFiles)
	}

	if len(sqlFiles) == 0 {
		log.Printf("[INFO] No migrations found for command: %s", command)
		return
	}

	log.Printf("[INFO] Found %d migration file(s) for command '%s'", len(sqlFiles), command)

	// Execute migrations
	for _, filename := range sqlFiles {
		filePath := filepath.Join(migrationsDir, filename)
		
		// Extract version name, e.g., "000001_init_schema.up.sql" -> "000001_init_schema"
		parts := strings.Split(filename, ".")
		version := parts[0]

		// Check if migration has already been applied/rolled back
		var exists bool
		err = pool.QueryRow(ctx, "SELECT EXISTS(SELECT 1 FROM schema_migrations WHERE version = $1)", version).Scan(&exists)
		if err != nil {
			log.Fatalf("Failed to query migration status for %s: %v", version, err)
		}

		if command == "up" && exists {
			log.Printf("[INFO] Migration %s is already applied. Skipping.", filename)
			continue
		}
		if command == "down" && !exists {
			log.Printf("[INFO] Migration %s has not been applied or is already rolled back. Skipping.", filename)
			continue
		}

		log.Printf("[INFO] Executing migration: %s", filename)

		content, err := os.ReadFile(filePath)
		if err != nil {
			log.Fatalf("Failed to read migration file %s: %v", filename, err)
		}

		// Execute the SQL queries in a transaction
		tx, err := pool.Begin(ctx)
		if err != nil {
			log.Fatalf("Failed to begin transaction: %v", err)
		}

		_, err = tx.Exec(ctx, string(content))
		if err != nil {
			tx.Rollback(ctx)
			log.Fatalf("Migration failed on %s: %v", filename, err)
		}

		// Update schema_migrations table
		if command == "up" {
			_, err = tx.Exec(ctx, "INSERT INTO schema_migrations (version) VALUES ($1)", version)
		} else {
			_, err = tx.Exec(ctx, "DELETE FROM schema_migrations WHERE version = $1", version)
		}
		if err != nil {
			tx.Rollback(ctx)
			log.Fatalf("Failed to update schema_migrations for %s: %v", version, err)
		}

		err = tx.Commit(ctx)
		if err != nil {
			log.Fatalf("Failed to commit transaction: %v", err)
		}

		log.Printf("[SUCCESS] Finished execution of %s", filename)
	}

	log.Printf("[SUCCESS] All migrations successfully completed.")
}
