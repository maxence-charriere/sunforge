package portal

import (
	"context"
	"io"
	"net/http"
	"strings"
	"testing"

	"github.com/maxence-charriere/sunforge/pkg/conduit"
	"github.com/stretchr/testify/require"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

type googleTransport func(*http.Request) (*http.Response, error)

func (f googleTransport) RoundTrip(r *http.Request) (*http.Response, error) {
	return f(r)
}

func TestGoogleSignIn(t *testing.T) {
	for _, scenario := range []struct {
		name      string
		tokenCode int
		userCode  int
		userBody  string
		wantError bool
	}{
		{
			name:      "success",
			tokenCode: http.StatusOK,
			userCode:  http.StatusOK,
			userBody:  `{"sub":"user-id","email":"user@example.com","given_name":"Alex","picture":"https://example.com/avatar.png","email_verified":true}`,
		},
		{
			name:      "rejected authorization code",
			tokenCode: http.StatusBadRequest,
			wantError: true,
		},
		{
			name:      "rejected user info request",
			tokenCode: http.StatusOK,
			userCode:  http.StatusUnauthorized,
			userBody:  `{}`,
			wantError: true,
		},
		{
			name:      "invalid user info response",
			tokenCode: http.StatusOK,
			userCode:  http.StatusOK,
			userBody:  `{`,
			wantError: true,
		},
	} {
		t.Run(scenario.name, func(t *testing.T) {
			userInfoRequested := false
			transport := googleTransport(func(r *http.Request) (*http.Response, error) {
				statusCode := scenario.tokenCode
				body := `{"access_token":"user-token","token_type":"Bearer","expires_in":3600}`
				switch r.URL.String() {
				case google.Endpoint.TokenURL:
					require.Equal(t, http.MethodPost, r.Method)
					require.NoError(t, r.ParseForm())
					require.Equal(t, "authorization_code", r.Form.Get("grant_type"))
					require.Equal(t, "auth-code", r.Form.Get("code"))
					require.Equal(t, "https://example.com/oauth/google", r.Form.Get("redirect_uri"))
					if statusCode != http.StatusOK {
						body = `{"error":"invalid_grant"}`
					}
				case "https://openidconnect.googleapis.com/v1/userinfo":
					userInfoRequested = true
					require.Equal(t, http.MethodGet, r.Method)
					require.Equal(t, "Bearer user-token", r.Header.Get("Authorization"))
					statusCode = scenario.userCode
					body = scenario.userBody
				default:
					t.Fatalf("unexpected request: %s", r.URL)
				}
				return &http.Response{
					StatusCode: statusCode,
					Header:     http.Header{"Content-Type": {"application/json"}},
					Body:       io.NopCloser(strings.NewReader(body)),
				}, nil
			})
			ctx := context.WithValue(context.Background(), oauth2.HTTPClient, &http.Client{Transport: transport})
			client := MakeGoogleClient(&conduit.OAuthClient{
				ClientID:     "client-id",
				ClientSecret: "client-secret",
				RedirectURL:  "https://example.com/oauth/google",
				Transport:    transport,
			})

			userClient, err := client.AuthAsUser(ctx, "auth-code")
			if scenario.tokenCode != http.StatusOK {
				require.Error(t, err)
				require.False(t, userInfoRequested)
				return
			}
			require.NoError(t, err)
			require.Nil(t, client.Token, "user authentication must not change the shared client token")

			user, err := userClient.UserInfo(ctx)
			require.True(t, userInfoRequested)
			if scenario.wantError {
				require.Error(t, err)
				return
			}
			require.NoError(t, err)
			require.Equal(t, GoogleUser{
				Sub:           "user-id",
				Email:         "user@example.com",
				GivenName:     "Alex",
				Picture:       "https://example.com/avatar.png",
				EmailVerified: true,
			}, user)
		})
	}
}
