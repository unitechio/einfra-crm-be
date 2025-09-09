
package cache

import (
	"context"
	"encoding/json"
	"time"
)

// Fetch retrieves data from the cache. If the data is not in the cache,
// it calls the `fetcher` function to get the data from the source and
// stores it in the cache before returning it.
func Fetch[T any](ctx context.Context, c Cache, key string, ttl time.Duration, fetcher func() (T, error)) (T, error) {
	var data T

	// 1. Try to get the data from the cache
	cached, err := c.Get(ctx, key)
	if err != nil {
		// If there is an error with the cache, it's better to fetch from the source
		// than to return an error.
		return fetcher()
	}

	// 2. If the data is in the cache, unmarshal it and return
	if cached != "" {
		if err := json.Unmarshal([]byte(cached), &data); err == nil {
			return data, nil // Cache hit
		}
	}

	// 3. If the data is not in the cache, fetch it from the source
	data, err = fetcher()
	if err != nil {
		return data, err
	}

	// 4. Marshal the data and store it in the cache
	serialized, err := json.Marshal(data)
	if err == nil {
		_ = c.Set(ctx, key, serialized, ttl) // Ignore error on set for now
	}

	return data, nil
}
