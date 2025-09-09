
package cache

import (
	"context"
	"time"

	lru "github.com/hashicorp/golang-lru/v2"
)

// LRUCache is an in-memory LRU implementation of the Cache interface.
// It uses the golang-lru library.
type LRUCache struct {
	lru *lru.Cache[string, string]
}

// NewLRUCache creates a new LRUCache with the given size.
func NewLRUCache(size int) (*LRUCache, error) {
	l, err := lru.New[string, string](size)
	if err != nil {
		return nil, err
	}
	return &LRUCache{lru: l}, nil
}

// Get retrieves a value from the cache.
func (c *LRUCache) Get(ctx context.Context, key string) (string, error) {
	// The context is not used in this implementation.
	val, ok := c.lru.Get(key)
	if !ok {
		return "", nil // Cache miss
	}
	return val, nil
}

// Set adds a value to the cache. The TTL is ignored in this implementation.
func (c *LRUCache) Set(ctx context.Context, key string, value interface{}, ttl time.Duration) error {
	// The context and TTL are not used in this implementation.
	c.lru.Add(key, value.(string)) // This implementation only supports string values
	return nil
}

// Delete removes a value from the cache.
func (c *LRUCache) Delete(ctx context.Context, key string) error {
	// The context is not used in this implementation.
	c.lru.Remove(key)
	return nil
}
