package portal

import (
	"context"

	"github.com/maxence-charriere/sunforge/pkg/conduit"
	"golang.org/x/oauth2/google"
)

type GoogleClient struct {
	*conduit.OAuthClient
}

func MakeGoogleClient(c *conduit.OAuthClient) GoogleClient {
	c.Endpoint = google.Endpoint
	return GoogleClient{OAuthClient: c}
}

func (c GoogleClient) AuthAsUser(ctx context.Context, code string) (GoogleClient, error) {
	userClient, err := c.OAuthClient.AuthAsUser(ctx, code)
	if err != nil {
		return GoogleClient{}, err
	}
	return GoogleClient{OAuthClient: userClient}, nil
}

type GoogleUser struct {
	Sub           string `json:"sub"`
	Email         string `json:"email"`
	GivenName     string `json:"given_name"`
	Picture       string `json:"picture"`
	EmailVerified bool   `json:"email_verified"`
}

func (c GoogleClient) UserInfo(ctx context.Context) (GoogleUser, error) {
	var user GoogleUser
	err := c.Get(ctx, "https://openidconnect.googleapis.com/v1/userinfo", &user)
	return user, err
}
