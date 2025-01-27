package db

import (
	"database/sql"
	"fmt"
	"os"

	_ "github.com/lib/pq"
)

// DB is the database connection global variable
var DB *sql.DB

// Connect connects to the database
func Connect() (*sql.DB, error) {
	// Data Source Name (DSN) is a string that contains the information needed to connect to the database
	dsn := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_NAME"),
	)

	// Open the database connection
	DB, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, fmt.Errorf("error connecting to dabase: %w", err)
	}

	// Ping the database to check if the connection is usable
	if err := DB.Ping(); err != nil {
		return nil, fmt.Errorf("error pinging database: %w", err)
	}

	// Return the database connection
	return DB, nil

}

// ExecuteMigrations executes the migrations
func ExecuteMigrations(conn *sql.DB) error {
	migrationsFiles := []string{
		"db/migrations/0001_user_channel.sql",
	}

	for _, file := range migrationsFiles {
		content, err := os.ReadFile(file)
		if err != nil {
			return fmt.Errorf("error reading migration file: %s,%w", file, err)
		}

		tx, err := conn.Begin()
		if err != nil {
			return fmt.Errorf("error starting transaction: %s,%w", file, err)
		}

		_, err = tx.Exec(string(content))
		if err != nil {
			tx.Rollback()
			return fmt.Errorf("error executing migration: %s,%w", file, err)
		}

		err = tx.Commit()
		if err != nil {
			return fmt.Errorf("error committing transaction: %s,%w", file, err)
		}

		fmt.Printf("Migration executed successfully: %s\n", file)
	}

	return nil
}
