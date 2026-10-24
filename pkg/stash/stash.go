// Package stash provides caches for storing and retrieving values by key.
//
// [ExpiringCache] supports concurrent use and expires values after a configurable
// time-to-live. Its zero value is ready to use with a one-minute time-to-live.
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
