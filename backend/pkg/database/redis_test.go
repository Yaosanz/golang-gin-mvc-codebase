package database_test

import (
	"context"
	"testing"
	"time"

	"go-starter-app/config"
	"go-starter-app/pkg/database"

	"github.com/joho/godotenv"
)

func TestNewRedis(t *testing.T) {
	// 1. Load .env (optional, for local dev)
	_ = godotenv.Load("../../.env")

	// 2. Override env for test (Redis Docker Desktop: localhost:6379)
	t.Setenv("REDIS_HOST", "127.0.0.1")
	t.Setenv("REDIS_PORT", "6379")
	t.Setenv("REDIS_PASSWORD", "")
	t.Setenv("REDIS_DB", "0")

	// 3. Load config
	cfg, err := config.NewConfig()
	if err != nil {
		t.Fatalf("failed to load config: %v", err)
	}

	// 4. Initialize Redis
	client := database.NewRedis(cfg)
	if client == nil {
		t.Fatal("Redis connection failed (nil). Pastikan Redis berjalan di Docker Desktop (localhost:6379)")
	}

	// 5. Ping (NewRedis already pings, but we double-check)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := client.Ping(ctx).Err(); err != nil {
		t.Fatalf("Redis ping failed: %v", err)
	}
}

func TestRedisSetGet(t *testing.T) {
	_ = godotenv.Load("../../.env")
	t.Setenv("REDIS_HOST", "127.0.0.1")
	t.Setenv("REDIS_PORT", "6379")
	t.Setenv("REDIS_PASSWORD", "")
	t.Setenv("REDIS_DB", "0")

	cfg, err := config.NewConfig()
	if err != nil {
		t.Fatalf("failed to load config: %v", err)
	}

	client := database.NewRedis(cfg)
	if client == nil {
		t.Skip("Redis tidak berjalan. Jalankan: docker run -d -p 6379:6379 redis:alpine")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	key := "test:gin:redis:hello"
	value := "world"

	if err := client.Set(ctx, key, value, 10*time.Second).Err(); err != nil {
		t.Fatalf("Redis Set failed: %v", err)
	}

	got, err := client.Get(ctx, key).Result()
	if err != nil {
		t.Fatalf("Redis Get failed: %v", err)
	}
	if got != value {
		t.Errorf("Redis Get: want %q, got %q", value, got)
	}

	// Cleanup
	_ = client.Del(ctx, key).Err()
}
