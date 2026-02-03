package database

import (
	"context"
	"fmt"
	"go-starter-app/config"
	"log"
	"time"

	"github.com/redis/go-redis/v9"
)

var (
	RedisClient *redis.Client
)

func NewRedis(config *config.Config) *redis.Client {
	addr := fmt.Sprintf("%s:%s", config.Redis().Host, config.Redis().Port)
	RedisClient = redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: config.Redis().Password,
		DB:       config.Redis().DB,
	})

	// test connection
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := RedisClient.Ping(ctx).Err(); err != nil {
		log.Printf("failed to connect to redis: %v", err)
		return nil
	}

	fmt.Println("✅ Redis connected:", config.Redis().Host)
	return RedisClient
}
