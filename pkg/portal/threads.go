package portal

import (
	"context"
	"time"

	"github.com/maxence-charriere/sunforge/pkg/conduit"
	"golang.org/x/oauth2"
)

type ThreadsClient struct {
	*conduit.OAuthClient
}

func MakeThreadsClient(c *conduit.OAuthClient) ThreadsClient {
	c.Endpoint = oauth2.Endpoint{
		AuthURL:  "https://threads.net/oauth/authorize",
		TokenURL: "https://graph.threads.net/oauth/access_token",
	}
	return ThreadsClient{OAuthClient: c}
}

func (c ThreadsClient) AuthAsUser(ctx context.Context, code string) (ThreadsClient, error) {
	userClient, err := c.OAuthClient.AuthAsUser(ctx, code)
	if err != nil {
		return ThreadsClient{}, err
	}
	return ThreadsClient{OAuthClient: userClient}, nil
}

func (c ThreadsClient) AuthWithAccessToken(accessToken string, expiresAt time.Time) ThreadsClient {
	return ThreadsClient{OAuthClient: c.OAuthClient.AuthWithAccessToken(accessToken, expiresAt)}
}
