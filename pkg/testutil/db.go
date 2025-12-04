package testutil

import (
	"database/sql"
	"fmt"
	"testing"

	_ "github.com/lib/pq"
)

// TestDB represents a test database connection
type TestDB struct {
	DB     *sql.DB
	DBName string
}

// SetupTestDB creates a test database and returns a connection
func SetupTestDB(t *testing.T) *TestDB {
	t.Helper()

	// Connect to postgres to create test database
	adminDB, err := sql.Open("postgres", "postgres://postgres:postgres@localhost:5432/postgres?sslmode=disable")
	if err != nil {
		t.Fatalf("Failed to connect to postgres: %v", err)
	}
	defer adminDB.Close()

	// Generate unique test database name
	dbName := fmt.Sprintf("test_ecommerce_%d", t.Name())

	// Drop if exists and create fresh database
	_, err = adminDB.Exec(fmt.Sprintf("DROP DATABASE IF EXISTS %s", dbName))
	if err != nil {
		t.Fatalf("Failed to drop test database: %v", err)
	}

	_, err = adminDB.Exec(fmt.Sprintf("CREATE DATABASE %s", dbName))
	if err != nil {
		t.Fatalf("Failed to create test database: %v", err)
	}

	// Connect to test database
	connStr := fmt.Sprintf("postgres://postgres:postgres@localhost:5432/%s?sslmode=disable", dbName)
	testDB, err := sql.Open("postgres", connStr)
	if err != nil {
		t.Fatalf("Failed to connect to test database: %v", err)
	}

	return &TestDB{
		DB:     testDB,
		DBName: dbName,
	}
}

// TeardownTestDB cleans up the test database
func (tdb *TestDB) Teardown(t *testing.T) {
	t.Helper()

	tdb.DB.Close()

	// Connect to postgres to drop test database
	adminDB, err := sql.Open("postgres", "postgres://postgres:postgres@localhost:5432/postgres?sslmode=disable")
	if err != nil {
		t.Logf("Warning: Failed to connect to postgres for cleanup: %v", err)
		return
	}
	defer adminDB.Close()

	_, err = adminDB.Exec(fmt.Sprintf("DROP DATABASE IF EXISTS %s", tdb.DBName))
	if err != nil {
		t.Logf("Warning: Failed to drop test database: %v", err)
	}
}

// RunInTransaction runs a function inside a transaction and rolls back
func (tdb *TestDB) RunInTransaction(t *testing.T, fn func(*sql.Tx) error) {
	t.Helper()

	tx, err := tdb.DB.Begin()
	if err != nil {
		t.Fatalf("Failed to begin transaction: %v", err)
	}

	defer func() {
		_ = tx.Rollback()
	}()

	if err := fn(tx); err != nil {
		t.Fatalf("Transaction failed: %v", err)
	}
}

// CreateTestSchema runs migration SQL to create tables
func (tdb *TestDB) CreateTestSchema(t *testing.T, schema string) {
	t.Helper()

	_, err := tdb.DB.Exec(schema)
	if err != nil {
		t.Fatalf("Failed to create test schema: %v", err)
	}
}
