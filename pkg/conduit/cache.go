package conduit

import (
	"fmt"
	"net/http"
	"sync"
	"time"
)

// Returns an http.Handler that set responses to GET requests with a
// Control-Cache header to max-age with the given duration.
//
// Paths specified as not cached will be ignored.
func NewHandlerWithCache(h http.Handler, expireDuration time.Duration, notCachedPaths ...string) http.Handler {
	ignoredPaths := make(map[string]struct{}, len(notCachedPaths))
	for _, path := range notCachedPaths {
		ignoredPaths[path] = struct{}{}
	}

	var cache Cache
	cacheControl := fmt.Sprintf("max-age=%.f", expireDuration.Seconds())

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			h.ServeHTTP(w, r)
			return
		}

		if _, ok := ignoredPaths[r.URL.Path]; ok {
			h.ServeHTTP(w, r)
			return
		}

		if res, ok := cache.Get(r.URL.Path); ok {
			res := res.(cachedResponse)
			for k := range res.header {
				w.Header().Set(k, res.header.Get(k))
			}
			w.Write(res.body)
			return
		}

		writer := responseWriter{ResponseWriter: w}
		writer.Header().Set("Cache-Control", cacheControl)
		h.ServeHTTP(&writer, r)

		if writer.statusCode >= 300 {
			return
		}

		cache.Set(r.URL.Path, expireDuration/2, cachedResponse{
			header: writer.Header(),
			body:   writer.body.Bytes(),
		})
	})
}

type cachedResponse struct {
	header http.Header
	body   []byte
}

// A cache that stores values with an expiration time.
type Cache struct {
	mutex sync.Mutex
	store map[string]cacheItem
}

// Set a value to the given key.
func (c *Cache) Set(key string, duration time.Duration, v any) {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	if c.store == nil {
		c.store = make(map[string]cacheItem)
	}

	c.store[key] = cacheItem{
		ExpireAt: time.Now().Add(duration),
		Value:    v,
	}
}

// Get a value with the given key.
func (c *Cache) Get(key string) (any, bool) {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	item, ok := c.store[key]
	if !ok {
		return nil, false
	}

	if time.Now().After(item.ExpireAt) {
		delete(c.store, key)
		return nil, false
	}

	return item.Value, true
}

type cacheItem struct {
	ExpireAt time.Time
	Value    any
}
