package database

import (
	"fmt"
	"go-starter-app/config"
	"strconv"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	gormLogger "gorm.io/gorm/logger"
)

const (
	defaultMaxRetries = 3
	defaultRetryDelay = 2 * time.Second
)

// NewPostgres creates a new PostgresSQL database connection using the provided configuration.
// It includes connection retry logic and proper error handling.
func NewPostgres(cfg *config.Config) (*gorm.DB, error) {
	// Validate database configuration
	if err := validateDBConfig(cfg); err != nil {
		return nil, &Error{
			Component: "PGSQL",
			Type:      eValidError,
			Err:       fmt.Errorf("invalid database configuration: %w", err),
		}
	}

	dsn := PostgresDSN(cfg)
	gormConfig := createGormConfig(cfg)

	// Retry connection logic
	var db *gorm.DB
	var err error
	for retries := 0; retries < defaultMaxRetries; retries++ {
		db, err = gorm.Open(postgres.Open(dsn), gormConfig)
		if err == nil {
			break
		}
		time.Sleep(defaultRetryDelay)
	}
	if err != nil {
		return nil, &Error{
			Component: "PGSQL",
			Type:      eConFailed,
			Err:       fmt.Errorf("failed to connect to database after %d retries: %w", defaultMaxRetries, err),
		}
	}

	if err := configureConnectionPool(db, cfg); err != nil {
		return nil, &Error{
			Component: "PGSQL",
			Type:      eConError,
			Err:       fmt.Errorf("failed to configure connection pool: %w", err),
		}
	}

	return db, nil
}

// validateDBConfig validates the database configuration
func validateDBConfig(cfg *config.Config) error {
	required := map[string]string{
		"host":     cfg.Database().Host,
		"port":     cfg.Database().Port,
		"name":     cfg.Database().Name,
		"username": cfg.Database().Username,
	}

	for field, value := range required {
		if value == "" {
			return fmt.Errorf("database %s is required", field)
		}
	}

	// Validate numeric values
	if port, err := strconv.Atoi(cfg.Database().Port); err != nil || port <= 0 {
		return fmt.Errorf("invalid port number")
	}

	return nil
}

// PostgresDSN creates the database connection string
func PostgresDSN(cfg *config.Config) string {
	return fmt.Sprintf("host=%s port=%s dbname=%s sslmode=%s TimeZone=%s user=%s password=%s",
		cfg.Database().Host,
		cfg.Database().Port,
		cfg.Database().Name,
		cfg.Database().SSLMode,
		cfg.Database().Timezone,
		cfg.Database().Username,
		cfg.Database().Password,
	)
}

// createGormConfig creates the GORM configuration
func createGormConfig(cfg *config.Config) *gorm.Config {
	gormConfig := &gorm.Config{}
	if cfg.Database().Debug {
		gormConfig.Logger = gormLogger.Default.LogMode(gormLogger.Info)
	}
	return gormConfig
}
