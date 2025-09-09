
package cache

import (
	"context"
	"time"
)

// Cache is the interface for a cache store.
// It allows for different implementations like Redis or in-memory LRU.
type Cache interface {
	Get(ctx context.Context, key string) (string, error)
	Set(ctx context.Context, key string, value interface{}, ttl time.Duration) error
	Delete(ctx context.Context, key string) error
}
