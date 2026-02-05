package database

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

var Pool *pgxpool.Pool

// InitConnection initializes the PostgreSQL connection pool
func InitConnection() error {
	databaseURL := os.Getenv("DATABASE_URL")
	if databaseURL == "" {
		return fmt.Errorf("DATABASE_URL environment variable is not set")
	}

	config, err := pgxpool.ParseConfig(databaseURL)
	if err != nil {
		return fmt.Errorf("failed to parse database URL: %w", err)
	}

	// Configure connection pool settings
	config.MaxConns = 10
	config.MinConns = 2
	config.MaxConnLifetime = time.Hour
	config.MaxConnIdleTime = 25 * time.Minute

	// Create the pool
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	pool, err := pgxpool.NewWithConfig(ctx, config)
	if err != nil {
		return fmt.Errorf("failed to create connection pool: %w", err)
	}

	// Test the connection
	testCtx, testCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer testCancel()

	if err := pool.Ping(testCtx); err != nil {
		pool.Close()
		return fmt.Errorf("failed to ping database: %w", err)
	}

	Pool = pool
	log.Println("✅ Database connection pool initialized successfully")

	// Log pool stats
	stats := pool.Stat()
	log.Printf("📊 Pool Stats: Conns=%d, IdleConns=%d, MaxConns=%d, MinConns=%d",
		stats.AcquiredConns(), stats.IdleConns(), config.MaxConns, config.MinConns)

	return nil
}

// Close closes the connection pool
func Close() {
	if Pool != nil {
		Pool.Close()
		log.Println("✅ Database connection pool closed")
	}
}

// GetPool returns the connection pool
func GetPool() *pgxpool.Pool {
	return Pool
}

// HealthCheck performs a health check on the database
func HealthCheck(ctx context.Context) error {
	if Pool == nil {
		return fmt.Errorf("connection pool is not initialized")
	}

	conn, err := Pool.Acquire(ctx)
	if err != nil {
		return fmt.Errorf("failed to acquire connection: %w", err)
	}
	defer conn.Release()

	var version string
	err = conn.QueryRow(ctx, "SELECT version()").Scan(&version)
	if err != nil {
		return fmt.Errorf("health check query failed: %w", err)
	}

	return nil
}
