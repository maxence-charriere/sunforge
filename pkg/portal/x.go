package portal

import (
	"bytes"
	"context"
	"net/http"
	"time"

	"github.com/maxence-charriere/go-app/v11/pkg/errors"
	"github.com/maxence-charriere/sunforge/pkg/conduit"
	"golang.org/x/oauth2"
)

type XUser struct {
	Data struct {
		ID       string `json:"id"`
		Name     string `json:"name"`
		Username string `json:"username"`
	} `json:"data"`
}

type X struct {
	*conduit.OAuthClient
}

func MakeX(c *conduit.OAuthClient) X {
	c.Endpoint = oauth2.Endpoint{
		AuthURL:  "https://x.com/i/oauth2/authorize",
		TokenURL: "https://api.twitter.com/2/oauth2/token",
	}
	return X{OAuthClient: c}
}

func (c X) AuthAsUser(ctx context.Context, code string, pckeCodeVerifier string) (X, error) {
	userClient, err := c.OAuthClient.AuthAsUser(ctx, code, conduit.OAuthOption{
		Key:   "pcke_code_challenge",
		Value: pckeCodeVerifier,
	})
	if err != nil {
		return X{}, err
	}

	return X{OAuthClient: userClient}, nil
}

func (c X) AuthWithAccessToken(accessToken string, expiresAt time.Time) X {
	return X{OAuthClient: c.OAuthClient.AuthWithAccessToken(accessToken, expiresAt)}
}

func (c X) GetCurrentUser(ctx context.Context) (XUser, error) {
	var user XUser
	err := c.Get(ctx, "https://api.twitter.com/2/users/me", &user)
	return user, err
}

type XPostQuery struct {
	Text string `json:"text"`
}

func (c X) Post(ctx context.Context, q XPostQuery) error {
	body, err := c.Encode(q)
	if err != nil {
		return errors.New("encoding body failed").Wrap(err)
	}

	req, err := http.NewRequestWithContext(ctx,
		http.MethodPost,
		"https://api.twitter.com/2/tweets",
		bytes.NewBuffer(body),
	)
	req.Header.Set("Content-Type", "application/json")
	if err != nil {
		return errors.New("creating http request failed").Wrap(err)
	}
	return c.Do(ctx, req, nil)
}
