package conduit

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestMux(t *testing.T) {
	utests := []struct {
		scenario     string
		method       string
		path         string
		target       string
		expectedCode int
	}{
		{
			scenario:     "method is handled",
			method:       http.MethodGet,
			path:         "/test",
			target:       "/test",
			expectedCode: http.StatusOK,
		},
		{
			scenario:     "method is not handled",
			method:       http.MethodPost,
			path:         "/test",
			target:       "/test",
			expectedCode: http.StatusNotFound,
		},
		{
			scenario:     "path is not handled",
			method:       http.MethodPost,
			path:         "/test",
			target:       "/hello",
			expectedCode: http.StatusNotFound,
		},
	}

	for _, u := range utests {
		t.Run(u.scenario, func(t *testing.T) {
			mux := NewMux()
			server := httptest.NewServer(mux)
			defer server.Close()

			mux.HandleFunc(u.method, u.path, func(w http.ResponseWriter, r *http.Request) {
			})

			res, err := server.Client().Get(server.URL + u.target)
			require.NoError(t, err)
			require.Equal(t, u.expectedCode, res.StatusCode)
		})
	}
}
