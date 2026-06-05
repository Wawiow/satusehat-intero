package database

import (
	"database/sql"
	"fmt"
	"log"

	_ "modernc.org/sqlite"
)

// InitDB initializes the SQLite database
func InitDB(filepath string) (*sql.DB, error) {
	if filepath == "" {
		filepath = "app.db"
	}

	db, err := sql.Open("sqlite", filepath)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	// Ping database to ensure connectivity
	if err = db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	if err := initTables(db); err != nil {
		return nil, fmt.Errorf("failed to initialize tables: %w", err)
	}

	log.Printf("Successfully connected to SQLite database at: %s", filepath)
	return db, nil
}

func initTables(db *sql.DB) error {
	query := `
	CREATE TABLE IF NOT EXISTS patients (
		nik TEXT PRIMARY KEY,
		id TEXT,
		ihs_number TEXT,
		name TEXT,
		gender TEXT,
		birth_date TEXT,
		phone TEXT,
		address TEXT,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP
	);

	CREATE TABLE IF NOT EXISTS locations (
		id TEXT PRIMARY KEY,
		identifier_value TEXT,
		name TEXT,
		description TEXT,
		phone TEXT,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP
	);

	CREATE TABLE IF NOT EXISTS practitioners (
		nik TEXT PRIMARY KEY,
		id TEXT,
		ihs_number TEXT,
		name TEXT,
		gender TEXT,
		birth_date TEXT,
		phone TEXT,
		address TEXT,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP
	);

	CREATE TABLE IF NOT EXISTS encounters (
		id TEXT PRIMARY KEY,
		identifier_value TEXT,
		status TEXT,
		subject_id TEXT,
		location_id TEXT,
		start_time TEXT,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP
	);
	`
	if _, err := db.Exec(query); err != nil {
		return err
	}

	// Programmatic migrations to add new columns if they were missing in existing DBs
	_, _ = db.Exec("ALTER TABLE patients ADD COLUMN ihs_number TEXT;")
	_, _ = db.Exec("ALTER TABLE practitioners ADD COLUMN ihs_number TEXT;")

	return nil
}
