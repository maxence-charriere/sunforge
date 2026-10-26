package portal

import (
	"context"
	"time"

	"github.com/maxence-charriere/sunforge/pkg/conduit"
	"golang.org/x/oauth2"
)

type Threads struct {
	*conduit.OAuthClient
}

func MakeThreads(c *conduit.OAuthClient) Threads {
	c.Endpoint = oauth2.Endpoint{
		AuthURL:  "https://threads.net/oauth/authorize",
		TokenURL: "https://graph.threads.net/oauth/access_token",
	}
	return Threads{OAuthClient: c}
}

func (c Threads) AuthAsUser(ctx context.Context, code string) (Threads, error) {
	userClient, err := c.OAuthClient.AuthAsUser(ctx, code)
	if err != nil {
		return Threads{}, err
	}
	return Threads{OAuthClient: userClient}, nil
}

func (c Threads) AuthWithAccessToken(accessToken string, expiresAt time.Time) Threads {
	return Threads{OAuthClient: c.OAuthClient.AuthWithAccessToken(accessToken, expiresAt)}
}
