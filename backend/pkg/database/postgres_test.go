package database_test

import (
	"testing"

	"go-starter-app/config"
	"go-starter-app/pkg/database"

	"github.com/joho/godotenv"
)

func TestNewPostgres(t *testing.T) {
	// 1. Load .env (opsional, untuk local dev)
	_ = godotenv.Load("../../.env")

	// 2. FORCE env untuk test (WAJIB)
	t.Setenv("DB_HOST", "localhost")
	t.Setenv("DB_PORT", "5432")
	t.Setenv("DB_USERNAME", "postgres")
	t.Setenv("DB_PASSWORD", "")
	t.Setenv("DB_NAME", "corpu")
	t.Setenv("DB_SSL_MODE", "disable")
	t.Setenv("DB_DEBUG", "false")

	// 3. Load config
	cfg, err := config.NewConfig()
	if err != nil {
		t.Fatalf("failed to load config: %v", err)
	}

	// 4. Initialize postgres
	db, err := database.NewPostgres(cfg)
	if err != nil {
		t.Fatalf("failed to initialize postgres: %v", err)
	}

	// 5. Validate connection
	sqlDB, err := db.DB()
	if err != nil {
		t.Fatalf("failed to get sql.DB: %v", err)
	}

	if err := sqlDB.Ping(); err != nil {
		t.Fatalf("failed to ping database: %v", err)
	}
}
