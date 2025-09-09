
package cache

import (
	"context"
	"time"

	"github.com/go-redis/redis/v8"
)

// RedisCache is a Redis implementation of the Cache interface.
// It uses the go-redis client to interact with a Redis server.
type RedisCache struct {
	client *redis.Client
}

// NewRedisCache creates a new RedisCache.
func NewRedisCache(client *redis.Client) *RedisCache {
	return &RedisCache{client: client}
}

// Get retrieves a value from the cache.
func (r *RedisCache) Get(ctx context.Context, key string) (string, error) {
	val, err := r.client.Get(ctx, key).Result()
	if err == redis.Nil {
		return "", nil // Cache miss
	}
	return val, err
}

// Set adds a value to the cache with a specified TTL.
func (r *RedisCache) Set(ctx context.Context, key string, value interface{}, ttl time.Duration) error {
	return r.client.Set(ctx, key, value, ttl).Err()
}

// Delete removes a value from the cache.
func (r *RedisCache) Delete(ctx context.Context, key string) error {
	return r.client.Del(ctx, key).Err()
}
