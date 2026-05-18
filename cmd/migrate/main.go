package main

import (
	"context"
	"log"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"go-beresin/pkg/database"
)

func main() {
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
	// Up migrations: 000001_..., 000002_...
	// Down migrations: usually reverse order or same, but for init_schema, there's only one.
	// For robust down migrations we sort in reverse order.
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
		log.Printf("[INFO] Executing migration: %s", filename)

		content, err := os.ReadFile(filePath)
		if err != nil {
			log.Fatalf("Failed to read migration file %s: %v", filename, err)
		}

		// Execute the SQL queries
		_, err = pool.Exec(ctx, string(content))
		if err != nil {
			log.Fatalf("Migration failed on %s: %v", filename, err)
		}
		log.Printf("[SUCCESS] Finished execution of %s", filename)
	}

	log.Printf("[SUCCESS] All migrations successfully completed.")
}
