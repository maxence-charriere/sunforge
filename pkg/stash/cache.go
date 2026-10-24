package stash

import "context"

// Cache stores values by key.
type Cache interface {
	// Get returns the value with the given key and whether it is cached.
	Get(ctx context.Context, key string) (any, bool)

	// Set stores the value at the given key.
	Set(ctx context.Context, key string, value any)

	// Del deletes the value at the given key.
	Del(ctx context.Context, key string)
}
