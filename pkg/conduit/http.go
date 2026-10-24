package conduit

import (
	"context"
	"net/http"
	"strings"

	"github.com/maxence-charriere/go-app/v11/pkg/app"
	"github.com/maxence-charriere/go-app/v11/pkg/errors"
)

// PathAt returns the request path part at the given index.
func PathAt(r *http.Request, i int) string {
	path := r.URL.Path
	if path == "" {
		return ""
	}

	s := strings.Split(path[1:], "/")
	if i >= len(s) {
		return ""
	}
	return s[i]
}

func ListenAndServe(ctx context.Context, s *http.Server) error {
	go func() {
		<-ctx.Done()
		if err := s.Shutdown(context.Background()); err != nil {
			app.Log(errors.New("server failed to clean shutdown").Wrap(err))
		}
	}()

	switch err := s.ListenAndServe(); err {
	case nil,
		http.ErrServerClosed,
		context.Canceled:
		app.Logf("server stopped")
		return nil

	default:
		return errors.New("server crashed").Wrap(err)
	}
}
