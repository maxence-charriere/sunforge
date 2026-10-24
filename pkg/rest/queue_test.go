package rest

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestQueue(t *testing.T) {
	q := newQueue[string](1)
	require.Equal(t, 0, q.Len())
	require.Equal(t, 1, q.Cap())

	_, ok := q.Dequeue()
	require.False(t, ok)

	j := "test"

	q.Enqueue(j)
	require.Equal(t, 1, q.Len())

	j2, ok := q.Dequeue()
	require.True(t, ok)
	require.Equal(t, j, j2)
	require.Equal(t, 0, q.Len())
}
