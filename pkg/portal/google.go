package portal

import (
	"context"

	"github.com/maxence-charriere/sunforge/pkg/conduit"
	"golang.org/x/oauth2/google"
)

type Google struct {
	*conduit.OAuthClient
}

func MakeGoogle(c *conduit.OAuthClient) Google {
	c.Endpoint = google.Endpoint
	return Google{OAuthClient: c}
}

func (c Google) AuthAsUser(ctx context.Context, code string) (Google, error) {
	userClient, err := c.OAuthClient.AuthAsUser(ctx, code)
	if err != nil {
		return Google{}, err
	}
	return Google{OAuthClient: userClient}, nil
}

type GoogleUser struct {
	Sub           string `json:"sub"`
	Email         string `json:"email"`
	GivenName     string `json:"given_name"`
	Picture       string `json:"picture"`
	EmailVerified bool   `json:"email_verified"`
}

func (c Google) UserInfo(ctx context.Context) (GoogleUser, error) {
	var user GoogleUser
	err := c.Get(ctx, "https://openidconnect.googleapis.com/v1/userinfo", &user)
	return user, err
}
