package portal

import (
	"context"
	"io"
	"net/http"
	"strings"
	"testing"
	"time"

	"github.com/maxence-charriere/sunforge/pkg/conduit"
	"github.com/stretchr/testify/require"
	"golang.org/x/oauth2"
)

type twitchTransport func(*http.Request) (*http.Response, error)

func (f twitchTransport) RoundTrip(r *http.Request) (*http.Response, error) {
	return f(r)
}

func TestTwitchStream(t *testing.T) {
	for _, scenario := range []struct {
		name      string
		status    int
		body      string
		want      TwitchStream
		wantError bool
	}{
		{
			name:   "live stream",
			status: http.StatusOK,
			body:   `{"data":[{"id":"stream-id","user_login":"streamer","title":"Live","viewer_count":42}]}`,
			want:   TwitchStream{ID: "stream-id", UserLogin: "streamer", Title: "Live", ViewerCount: 42},
		},
		{
			name:   "offline",
			status: http.StatusOK,
			body:   `{"data":[]}`,
		},
		{
			name:   "different user",
			status: http.StatusOK,
			body:   `{"data":[{"id":"other-stream","user_login":"other"}]}`,
		},
		{
			name:      "request rejected",
			status:    http.StatusUnauthorized,
			body:      `{}`,
			wantError: true,
		},
		{
			name:      "invalid response",
			status:    http.StatusOK,
			body:      `{`,
			wantError: true,
		},
	} {
		t.Run(scenario.name, func(t *testing.T) {
			client := MakeTwitchClient(&conduit.OAuthClient{
				ClientID: "client-id",
				Token: &oauth2.Token{
					AccessToken: "access-token",
					Expiry:      time.Now().Add(time.Hour),
				},
				Transport: twitchTransport(func(r *http.Request) (*http.Response, error) {
					require.Equal(t, http.MethodGet, r.Method)
					require.Equal(t, "https://api.twitch.tv/helix/streams?user_login=streamer", r.URL.String())
					require.Equal(t, "Bearer access-token", r.Header.Get("Authorization"))
					require.Equal(t, "client-id", r.Header.Get("Client-Id"))
					return &http.Response{
						StatusCode: scenario.status,
						Body:       io.NopCloser(strings.NewReader(scenario.body)),
					}, nil
				}),
			})

			stream, err := client.Stream(context.Background(), "https://www.twitch.tv/streamer")
			if scenario.wantError {
				require.Error(t, err)
				return
			}
			require.NoError(t, err)
			require.Equal(t, scenario.want, stream)
		})
	}
}
