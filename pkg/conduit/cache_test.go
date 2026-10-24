package conduit

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestHandlerWithCache(t *testing.T) {
	var h http.Handler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/bad" {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		w.Write([]byte("ok!"))
	})
	h = NewHandlerWithCache(h, time.Second*5, "/ignore")

	t.Run("get request has cache-control header set", func(t *testing.T) {
		w := httptest.NewRecorder()
		r := httptest.NewRequest(http.MethodGet, "/", nil)
		h.ServeHTTP(w, r)

		require.Equal(t, "max-age=5", w.Header().Get("Cache-Control"))
	})

	t.Run("get cached request has cache-control header set", func(t *testing.T) {
		for i := 0; i < 2; i++ {
			w := httptest.NewRecorder()
			r := httptest.NewRequest(http.MethodGet, "/get/cached", nil)
			h.ServeHTTP(w, r)

			require.Equal(t, "max-age=5", w.Header().Get("Cache-Control"))
		}
	})

	t.Run("non get request has no cache-control header", func(t *testing.T) {
		w := httptest.NewRecorder()
		r := httptest.NewRequest(http.MethodPost, "/", nil)
		h.ServeHTTP(w, r)

		require.Empty(t, w.Header().Get("Cache-Control"))
	})

	t.Run("ignored path has no cache-control header", func(t *testing.T) {
		w := httptest.NewRecorder()
		r := httptest.NewRequest(http.MethodGet, "/ignore", nil)
		h.ServeHTTP(w, r)

		require.Empty(t, w.Header().Get("Cache-Control"))
	})

	t.Run("error responses from get request are not cached", func(t *testing.T) {
		w := httptest.NewRecorder()
		r := httptest.NewRequest(http.MethodGet, "/bad", nil)
		h.ServeHTTP(w, r)

		require.Equal(t, "max-age=5", w.Header().Get("Cache-Control"))
	})
}

func TestCache(t *testing.T) {
	t.Run("non set value is not returned", func(t *testing.T) {
		var c Cache

		v, ok := c.Get("test")
		require.False(t, ok)
		require.Nil(t, v)
	})

	t.Run("expired value is not returned", func(t *testing.T) {
		var c Cache

		c.Set("test", -time.Second, 42)

		v, ok := c.Get("test")
		require.False(t, ok)
		require.Nil(t, v)
		require.Empty(t, c.store)
	})

	t.Run("value is returned", func(t *testing.T) {
		var c Cache

		c.Set("test", time.Minute, 42)
		v, ok := c.Get("test")
		require.True(t, ok)
		require.Equal(t, 42, v)
	})
}
