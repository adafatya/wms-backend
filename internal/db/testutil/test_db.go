package testutil

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"testing"

	"github.com/adafatya/wms-backend/internal/db/sqlc"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

func SetupTestDB(t *testing.T) (*sql.DB, sqlc.Querier, func()) {
	// Load .env.test from current directory or parent
	err := godotenv.Load(filepath.Join(getProjectRoot(), ".env.test"))
	if err != nil {
		t.Fatalf("failed to load .env.test: %v", err)
	}

	driver := os.Getenv("DB_DRIVER")
	source := os.Getenv("DB_SOURCE")

	// Ensure the test database exists
	err = ensureDBExists(driver, source)
	if err != nil {
		t.Fatalf("failed to ensure test db exists: %v", err)
	}

	// Connect to test database
	db, err := sql.Open(driver, source)
	if err != nil {
		t.Fatalf("failed to connect to test db: %v", err)
	}

	if err := db.Ping(); err != nil {
		t.Fatalf("failed to ping test db: %v", err)
	}

	// Drop all tables to ensure clean schema
	dropAllTables(db)

	// Run migrations
	err = runMigrations(db)
	if err != nil {
		t.Fatalf("failed to run migrations: %v", err)
	}

	querier := sqlc.New(db)

	cleanup := func() {
		TruncateTables(db)
		db.Close()
	}

	return db, querier, cleanup
}

func dropAllTables(db *sql.DB) {
	// Simple query to drop all tables in public schema
	rows, err := db.Query("SELECT table_name FROM information_schema.tables WHERE table_schema = 'public' AND table_type = 'BASE TABLE'")
	if err != nil {
		log.Fatalf("failed to list tables for drop: %v", err)
	}
	defer rows.Close()

	var tables []string
	for rows.Next() {
		var table string
		if err := rows.Scan(&table); err != nil {
			log.Fatalf("failed to scan table name: %v", err)
		}
		tables = append(tables, fmt.Sprintf("\"%s\"", table))
	}

	if len(tables) > 0 {
		query := fmt.Sprintf("DROP TABLE IF EXISTS %s CASCADE", strings.Join(tables, ", "))
		_, err = db.Exec(query)
		if err != nil {
			log.Fatalf("failed to drop tables: %v", err)
		}
	}
}

func TruncateTables(db *sql.DB) {
	// Simple query to truncate all tables except those we might want to keep
	// (e.g., schema_migrations if we used golang-migrate, but here we don't)
	rows, err := db.Query("SELECT table_name FROM information_schema.tables WHERE table_schema = 'public' AND table_type = 'BASE TABLE'")
	if err != nil {
		log.Fatalf("failed to list tables for truncation: %v", err)
	}
	defer rows.Close()

	var tables []string
	for rows.Next() {
		var table string
		if err := rows.Scan(&table); err != nil {
			log.Fatalf("failed to scan table name: %v", err)
		}
		tables = append(tables, fmt.Sprintf("\"%s\"", table))
	}

	if len(tables) > 0 {
		query := fmt.Sprintf("TRUNCATE TABLE %s RESTART IDENTITY CASCADE", strings.Join(tables, ", "))
		_, err = db.Exec(query)
		if err != nil {
			log.Fatalf("failed to truncate tables: %v", err)
		}
	}
}

func ensureDBExists(driver, source string) error {
	// Parse the connection string to extract the DB name and connect to "postgres" instead
	// Source format: postgresql://user:pass@host:port/dbname?sslmode=disable
	dbName := ""
	lastSlash := strings.LastIndex(source, "/")
	qMark := strings.Index(source, "?")
	if qMark == -1 {
		dbName = source[lastSlash+1:]
	} else {
		dbName = source[lastSlash+1 : qMark]
	}

	baseSource := source[:lastSlash+1] + "postgres"
	if qMark != -1 {
		baseSource += source[qMark:]
	}

	db, err := sql.Open(driver, baseSource)
	if err != nil {
		return err
	}
	defer db.Close()

	var exists bool
	query := fmt.Sprintf("SELECT EXISTS (SELECT 1 FROM pg_database WHERE datname = '%s')", dbName)
	err = db.QueryRow(query).Scan(&exists)
	if err != nil {
		return err
	}

	if !exists {
		_, err = db.Exec(fmt.Sprintf("CREATE DATABASE %s", dbName))
		if err != nil {
			return err
		}
	}

	return nil
}

func runMigrations(db *sql.DB) error {
	root := getProjectRoot()
	migrationsDir := filepath.Join(root, "migrations")

	files, err := filepath.Glob(filepath.Join(migrationsDir, "*.up.sql"))
	if err != nil {
		return err
	}
	sort.Strings(files)

	for _, file := range files {
		content, err := os.ReadFile(file)
		if err != nil {
			return err
		}

		_, err = db.Exec(string(content))
		if err != nil {
			return fmt.Errorf("error in file %s: %w", file, err)
		}
	}
	return nil
}

func getProjectRoot() string {
	// Simple root finder, adjusting for where tests are run
	dir, _ := os.Getwd()
	for {
		if _, err := os.Stat(filepath.Join(dir, "go.mod")); err == nil {
			return dir
		}
		parent := filepath.Dir(dir)
		if parent == dir {
			break
		}
		dir = parent
	}
	return ""
}
