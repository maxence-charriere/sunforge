package stash

import (
	"context"
	"strconv"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func newTestExpiringCache(ttl time.Duration) (*ExpiringCache, func(time.Duration)) {
	now := time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC)
	return &ExpiringCache{
		ItemTTL: ttl,
		now:     func() time.Time { return now },
	}, func(d time.Duration) { now = now.Add(d) }
}

func TestExpiringCache(t *testing.T) {
	c, _ := newTestExpiringCache(time.Minute)
	testCache(t, c)
}

func TestExpiringCacheExpiration(t *testing.T) {
	ctx := context.Background()
	c, advance := newTestExpiringCache(time.Minute)
	c.Set(ctx, "key", "value")

	advance(time.Minute - time.Nanosecond)
	value, ok := c.Get(ctx, "key")
	require.True(t, ok)
	require.Equal(t, "value", value)

	// The value is still alive at its exact deadline, including during cleanup.
	advance(time.Nanosecond)
	c.Del(ctx, "missing")
	value, ok = c.Get(ctx, "key")
	require.True(t, ok)
	require.Equal(t, "value", value)

	// Reading does not refresh TTL.
	advance(time.Nanosecond)
	value, ok = c.Get(ctx, "key")
	require.False(t, ok)
	require.Nil(t, value)
	c.Del(ctx, "missing")
	require.Empty(t, c.items)
	require.Zero(t, c.queue.Len())
}

func TestExpiringCacheOverwriteRefreshesTTL(t *testing.T) {
	ctx := context.Background()
	c, advance := newTestExpiringCache(time.Minute)
	c.Set(ctx, "first", "original")
	advance(10 * time.Second)
	c.Set(ctx, "second", "other")
	advance(10 * time.Second)
	c.Set(ctx, "first", "replacement")

	advance(50*time.Second + time.Nanosecond)
	// The overwritten first key must move behind the second key for cleanup.
	c.Del(ctx, "missing")
	require.NotContains(t, c.items, "second")
	_, ok := c.Get(ctx, "second")
	require.False(t, ok)
	value, ok := c.Get(ctx, "first")
	require.True(t, ok)
	require.Equal(t, "replacement", value)

	advance(10 * time.Second)
	c.Del(ctx, "missing")
	require.Empty(t, c.items)
	_, ok = c.Get(ctx, "first")
	require.False(t, ok)
}

func TestExpiringCacheOverwriteKeepsStorageBounded(t *testing.T) {
	ctx := context.Background()
	c, _ := newTestExpiringCache(time.Hour)
	c.Set(ctx, "first", true)
	for i := 0; i < 1000; i++ {
		c.Set(ctx, "updated", i)
	}
	// Repeated writes must not accumulate obsolete queue entries.
	require.Len(t, c.items, 2)
	require.Equal(t, 2, c.queue.Len())
	value, ok := c.Get(ctx, "updated")
	require.True(t, ok)
	require.Equal(t, 999, value)
}

func TestExpiringCacheWriteCleanup(t *testing.T) {
	ctx := context.Background()
	for _, operation := range []string{"set", "delete", "delete missing"} {
		t.Run(operation, func(t *testing.T) {
			c, advance := newTestExpiringCache(time.Minute)
			c.Set(ctx, "expired", true)
			advance(30 * time.Second)
			c.Set(ctx, "live", true)
			advance(30*time.Second + time.Nanosecond)

			switch operation {
			case "set":
				c.Set(ctx, "new", true)
			case "delete":
				c.Del(ctx, "live")
			case "delete missing":
				c.Del(ctx, "missing")
			}
			require.NotContains(t, c.items, "expired")
			if operation != "delete" {
				value, ok := c.Get(ctx, "live")
				require.True(t, ok)
				require.Equal(t, true, value)
			}
			advance(time.Minute + time.Nanosecond)
			c.Del(ctx, "missing")
			require.Empty(t, c.items)
			require.Zero(t, c.queue.Len())
		})
	}
}

func TestExpiringCacheDelete(t *testing.T) {
	ctx := context.Background()
	for _, key := range []string{"first", "middle", "last"} {
		t.Run(key, func(t *testing.T) {
			c, advance := newTestExpiringCache(time.Minute)
			for _, name := range []string{"first", "middle", "last"} {
				c.Set(ctx, name, name)
			}
			c.Del(ctx, key)
			c.Del(ctx, key)
			value, ok := c.Get(ctx, key)
			require.False(t, ok)
			require.Nil(t, value)
			require.Len(t, c.items, 2)
			require.Equal(t, 2, c.queue.Len())
			advance(time.Minute + time.Nanosecond)
			c.Del(ctx, "missing")
			require.Empty(t, c.items)
			require.Zero(t, c.queue.Len())
		})
	}
}

func TestExpiringCacheNilValue(t *testing.T) {
	ctx := context.Background()
	c, advance := newTestExpiringCache(time.Minute)
	c.Set(ctx, "key", nil)
	c.Set(ctx, "other", true)
	value, ok := c.Get(ctx, "key")
	require.True(t, ok)
	require.Nil(t, value)

	c.Set(ctx, "key", "replacement")
	value, ok = c.Get(ctx, "key")
	require.True(t, ok)
	require.Equal(t, "replacement", value)

	c.Set(ctx, "key", nil)
	advance(time.Minute + time.Nanosecond)
	c.Del(ctx, "missing")
	require.Empty(t, c.items)
	value, ok = c.Get(ctx, "key")
	require.False(t, ok)
	require.Nil(t, value)
}

func TestExpiringCacheDefaultTTL(t *testing.T) {
	ctx := context.Background()
	for _, firstOperation := range []string{"get", "set", "delete"} {
		t.Run(firstOperation, func(t *testing.T) {
			c, advance := newTestExpiringCache(0)
			switch firstOperation {
			case "get":
				_, ok := c.Get(ctx, "missing")
				require.False(t, ok)
			case "set":
				c.Set(ctx, "first", true)
			case "delete":
				c.Del(ctx, "missing")
			}

			c.Set(ctx, "key", "value")
			advance(time.Minute - time.Nanosecond)
			value, ok := c.Get(ctx, "key")
			require.True(t, ok)
			require.Equal(t, "value", value)
			advance(time.Nanosecond)
			value, ok = c.Get(ctx, "key")
			require.True(t, ok)
			require.Equal(t, "value", value)
			advance(time.Nanosecond)
			value, ok = c.Get(ctx, "key")
			require.False(t, ok)
			require.Nil(t, value)
		})
	}
}

func TestExpiringCacheZeroValue(t *testing.T) {
	var c ExpiringCache
	testCache(t, &c)
}

func TestExpiringCacheNegativeTTL(t *testing.T) {
	ctx := context.Background()
	c := ExpiringCache{ItemTTL: -time.Minute}
	c.Del(ctx, "missing")
	c.Set(ctx, "key", "value")
	value, ok := c.Get(ctx, "key")
	require.False(t, ok)
	require.Nil(t, value)
	require.Empty(t, c.items)
	require.Zero(t, c.queue.Len())
}

func TestExpiringCacheConcurrent(t *testing.T) {
	ctx := context.Background()
	var c ExpiringCache
	var workers sync.WaitGroup
	for worker := 0; worker < 8; worker++ {
		workers.Add(1)
		go func(worker int) {
			defer workers.Done()
			key := strconv.Itoa(worker)
			for i := 0; i < 100; i++ {
				c.Set(ctx, key, i)
				if value, ok := c.Get(ctx, key); !ok || value != i {
					t.Errorf("Get(%q) = (%v, %v), want (%d, true)", key, value, ok, i)
					return
				}
				c.Del(ctx, key)
				c.Set(ctx, "shared", i)
				c.Get(ctx, "shared")
				c.Del(ctx, "shared")
			}
		}(worker)
	}
	workers.Wait()
	require.Empty(t, c.items)
	require.Zero(t, c.queue.Len())
}
