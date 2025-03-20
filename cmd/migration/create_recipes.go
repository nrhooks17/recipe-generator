package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/jackc/pgx/v5/pgxpool"
)

// Database connection parameters - adjust these as needed
const (
	host     = "localhost"
	port     = 5432
	user     = "postgres"
	password = "postgres"
	dbname   = "recipe_generator"
)

func main() {
	// Connect to the database using pgxpool
	connStr := fmt.Sprintf("postgres://%s:%s@%s:%d/%s", 
		user, password, host, port, dbname)
	
	pool, err := pgxpool.New(context.Background(), connStr)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer pool.Close()

	log.Println("Successfully connected to the database")

    // need to add all sql files here for when I setup a new database

	// Execute the SQL files
	err = executeSQLFile(pool, "../../internal/database/migration/truncate_recipe_tables.sql")
	if err != nil {
		log.Fatalf("Failed to execute truncate script: %v", err)
	}
	
	
	err = executeSQLFile(pool, "../../internal/database/migration/create_recipes.sql")
	if err != nil {
		log.Fatalf("Failed to execute create recipes script: %v", err)
	}
	
	log.Println("Migration completed successfully")
}

// executeSQLFile reads and executes an SQL file
func executeSQLFile(pool *pgxpool.Pool, filePath string) error {
	// Get the absolute path
	absPath, err := filepath.Abs(filePath)
	if err != nil {
		return fmt.Errorf("failed to get absolute path: %w", err)
	}
	
	// Read the SQL file
	content, err := os.ReadFile(absPath)
	if err != nil {
		return fmt.Errorf("failed to read SQL file %s: %w", filePath, err)
	}
	
	// Execute the SQL
	log.Printf("Executing SQL file: %s", filePath)
	_, err = pool.Exec(context.Background(), string(content))
	if err != nil {
		return fmt.Errorf("failed to execute SQL from %s: %w", filePath, err)
	}
	
	return nil
}



