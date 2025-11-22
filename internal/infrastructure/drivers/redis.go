package drivers

import (
	"context"
	"log"

	"github.com/go-redis/redis/v8"
	"github.com/unitechio/einfra-be/internal/config"
)

// NewRedisClient creates a new Redis client.
func NewRedisClient(cfg config.RedisConfig) (*redis.Client, error) {
	client := redis.NewClient(&redis.Options{
		Addr:     cfg.Addr,
		Password: cfg.Password,
		DB:       cfg.DB,
	})

	// Ping the Redis server to check the connection.
	ctx := context.Background()
	if _, err := client.Ping(ctx).Result(); err != nil {
		return nil, err
	}

	log.Println("Successfully connected to Redis.")

	return client, nil
}
