package stash

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
)

func testCache(t *testing.T, c Cache) {
	t.Helper()
	ctx := context.Background()

	value, ok := c.Get(ctx, "missing")
	require.False(t, ok)
	require.Nil(t, value)
	c.Del(ctx, "missing")

	values := []struct {
		name  string
		value any
	}{
		{"bytes", []byte("foo")},
		{"string", "bar"},
		{"int", 42},
		{"float", 42.5},
		{"struct", struct{ Name string }{Name: "sunforge"}},
		{"nil", nil},
	}

	for _, entry := range values {
		c.Set(ctx, entry.name, entry.value)
	}
	for _, entry := range values {
		t.Run(entry.name, func(t *testing.T) {
			value, ok := c.Get(ctx, entry.name)
			require.True(t, ok)
			require.Equal(t, entry.value, value)

			c.Del(ctx, entry.name)
			value, ok = c.Get(ctx, entry.name)
			require.False(t, ok)
			require.Nil(t, value)
		})
	}
}
