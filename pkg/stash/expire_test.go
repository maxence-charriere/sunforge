package stash

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestExpire(t *testing.T) {
	testCache(t, &Expire{
		ItemTTL: time.Minute,
	})
}

func TestExpireExpire(t *testing.T) {
	ctx := context.TODO()

	c := Expire{
		ItemTTL: time.Minute,
	}

	c.Set(ctx, "/hello", "hello")
	require.Len(t, c.items, 1)

	c.Set(ctx, "/world", "world")
	require.Len(t, c.items, 2)

	c.Set(ctx, "/goodbye", "goodbye")
	require.Len(t, c.items, 3)

	c.items["/hello"].expiresAt = time.Now().Add(-time.Second)
	c.items["/world"].expiresAt = time.Now().Add(-time.Second)

	c.Set(ctx, "/goodmorning", "goodmorning")
	require.Len(t, c.items, 2)
	require.Len(t, c.queue, 2)
}

func TestExpireSetSameKey(t *testing.T) {
	ctx := context.TODO()

	c := Expire{
		ItemTTL: time.Minute,
	}

	c.Set(ctx, "/foo", "foo")
	require.Len(t, c.queue, 1)
	require.Len(t, c.items, 1)

	c.Set(ctx, "/bar", "bar")
	require.Len(t, c.queue, 2)
	require.Len(t, c.items, 2)

	c.Set(ctx, "/bar", "barre")
	require.Len(t, c.queue, 3)
	require.Len(t, c.items, 2)

	c.Set(ctx, "/foo", "fooo")
	require.Len(t, c.queue, 2)
	require.Len(t, c.items, 2)
}

func TestExpireNilValue(t *testing.T) {
	ctx := context.Background()
	c := Expire{ItemTTL: time.Minute}
	c.Set(ctx, "nil", nil)
	c.Set(ctx, "other", "value")

	value, ok := c.Get(ctx, "nil")
	require.True(t, ok)
	require.Nil(t, value)

	c.items["nil"].expiresAt = time.Now().Add(-time.Second)
	value, ok = c.Get(ctx, "nil")
	require.False(t, ok)
	require.Nil(t, value)

	c.Set(ctx, "trigger", true)
	require.NotContains(t, c.items, "nil")
	value, ok = c.Get(ctx, "other")
	require.True(t, ok)
	require.Equal(t, "value", value)
}

func TestExpireOverwriteNilValue(t *testing.T) {
	ctx := context.Background()
	c := Expire{ItemTTL: time.Minute}
	c.Set(ctx, "first", true)
	c.Set(ctx, "key", nil)
	c.Set(ctx, "key", "replacement")

	// Cleaning up the old entry must preserve its replacement.
	c.Del(ctx, "first")
	value, ok := c.Get(ctx, "key")
	require.True(t, ok)
	require.Equal(t, "replacement", value)

	c.Set(ctx, "key", nil)
	c.Del(ctx, "key")
	value, ok = c.Get(ctx, "key")
	require.False(t, ok)
	require.Nil(t, value)
}
