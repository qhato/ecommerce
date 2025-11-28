package database

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/qhato/ecommerce/pkg/logger"
)

// DB wraps the pgxpool.Pool
type DB struct {
	pool *pgxpool.Pool
}

// Config holds database configuration
type Config struct {
	Host           string
	Port           int
	User           string
	Password       string
	Database       string
	SSLMode        string
	MaxConnections int
	MaxIdleConns   int
	MaxLifetime    time.Duration
	MaxIdleTime    time.Duration
}

// New creates a new database connection pool
func New(ctx context.Context, cfg Config) (*DB, error) {
	// Build connection string
	dsn := fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		cfg.Host, cfg.Port, cfg.User, cfg.Password, cfg.Database, cfg.SSLMode,
	)

	// Parse config
	poolConfig, err := pgxpool.ParseConfig(dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to parse database config: %w", err)
	}

	// Set pool configuration
	poolConfig.MaxConns = int32(cfg.MaxConnections)
	poolConfig.MinConns = int32(cfg.MaxIdleConns)
	poolConfig.MaxConnLifetime = cfg.MaxLifetime
	poolConfig.MaxConnIdleTime = cfg.MaxIdleTime

	// Set connection timeout
	poolConfig.ConnConfig.ConnectTimeout = 10 * time.Second

	// Create pool
	pool, err := pgxpool.NewWithConfig(ctx, poolConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to create connection pool: %w", err)
	}

	// Test connection
	if err := pool.Ping(ctx); err != nil {
		pool.Close()
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	logger.Info("Database connection pool created successfully")

	return &DB{pool: pool}, nil
}

// Close closes the database connection pool
func (db *DB) Close() {
	if db.pool != nil {
		db.pool.Close()
		logger.Info("Database connection pool closed")
	}
}

// Pool returns the underlying connection pool
func (db *DB) Pool() *pgxpool.Pool {
	return db.pool
}

// Ping tests the database connection
func (db *DB) Ping(ctx context.Context) error {
	return db.pool.Ping(ctx)
}

// Begin starts a new transaction
func (db *DB) Begin(ctx context.Context) (pgx.Tx, error) {
	return db.pool.Begin(ctx)
}

// BeginTx starts a new transaction with options
func (db *DB) BeginTx(ctx context.Context, txOptions pgx.TxOptions) (pgx.Tx, error) {
	return db.pool.BeginTx(ctx, txOptions)
}

// Exec executes a query without returning any rows
func (db *DB) Exec(ctx context.Context, query string, args ...interface{}) error {
	_, err := db.pool.Exec(ctx, query, args...)
	return err
}

// Query executes a query that returns rows
func (db *DB) Query(ctx context.Context, query string, args ...interface{}) (pgx.Rows, error) {
	return db.pool.Query(ctx, query, args...)
}

// QueryRow executes a query that returns at most one row
func (db *DB) QueryRow(ctx context.Context, query string, args ...interface{}) pgx.Row {
	return db.pool.QueryRow(ctx, query, args...)
}

// Stats returns database pool statistics
func (db *DB) Stats() *pgxpool.Stat {
	return db.pool.Stat()
}

// Health checks the database health and returns status
func (db *DB) Health(ctx context.Context) error {
	ctx, cancel := context.WithTimeout(ctx, 2*time.Second)
	defer cancel()

	if err := db.Ping(ctx); err != nil {
		return fmt.Errorf("database unhealthy: %w", err)
	}

	stats := db.Stats()
	logger.WithFields(logger.Fields{
		"total_conns":    stats.TotalConns(),
		"acquired_conns": stats.AcquiredConns(),
		"idle_conns":     stats.IdleConns(),
	}).Debug("Database health check passed")

	return nil
}
