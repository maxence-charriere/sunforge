package rest

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"sync"
	"time"

	"github.com/grokify/go-pkce"
	"github.com/maxence-charriere/go-app/v11/pkg/errors"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/clientcredentials"
)

type OAuthOption struct {
	Key   string
	Value string
}

type OAuthClient struct {
	ClientID     string
	ClientSecret string
	RedirectURL  string
	Endpoint     oauth2.Endpoint
	Transport    http.RoundTripper
	Decode       func([]byte, interface{}) error
	Encode       func(v any) ([]byte, error)
	Token        *oauth2.Token

	initOnce  sync.Once
	authMutex sync.Mutex
}

func (c *OAuthClient) Auth(ctx context.Context) error {
	c.initOnce.Do(c.init)

	oauth2Config := &clientcredentials.Config{
		ClientID:     c.ClientID,
		ClientSecret: c.ClientSecret,
		TokenURL:     c.Endpoint.TokenURL,
	}

	token, err := oauth2Config.Token(ctx)
	if err != nil {
		return errors.New("getting token failed").Wrap(err)
	}

	c.Token = token
	return nil
}

func (c *OAuthClient) AuthAsUser(ctx context.Context, code string, options ...OAuthOption) (*OAuthClient, error) {
	c.initOnce.Do(c.init)

	config := oauth2.Config{
		ClientID:     c.ClientID,
		ClientSecret: c.ClientSecret,
		Endpoint:     c.Endpoint,
		RedirectURL:  c.RedirectURL,
	}

	var exchangeOptions []oauth2.AuthCodeOption
	for _, o := range options {
		switch o.Key {
		case "pcke_code_challenge":
			exchangeOptions = append(exchangeOptions, oauth2.SetAuthURLParam(pkce.ParamCodeVerifier, o.Value))
		}
	}

	token, err := config.Exchange(ctx, code, exchangeOptions...)
	if err != nil {
		return nil, errors.New("exchanging user authorization code failed").Wrap(err)
	}

	return &OAuthClient{
		ClientID:     c.ClientID,
		ClientSecret: c.ClientSecret,
		RedirectURL:  c.RedirectURL,
		Endpoint:     c.Endpoint,
		Transport:    c.Transport,
		Encode:       c.Encode,
		Decode:       c.Decode,
		Token:        token,
	}, nil
}

func (c *OAuthClient) AuthWithAccessToken(accessToken string, expiresAt time.Time) *OAuthClient {
	c.initOnce.Do(c.init)

	return &OAuthClient{
		ClientID:     c.ClientID,
		ClientSecret: c.ClientSecret,
		RedirectURL:  c.RedirectURL,
		Endpoint:     c.Endpoint,
		Transport:    c.Transport,
		Encode:       c.Encode,
		Decode:       c.Decode,
		Token: &oauth2.Token{
			AccessToken: accessToken,
			Expiry:      expiresAt,
		},
	}
}

func (c *OAuthClient) Get(ctx context.Context, url string, out any) error {
	c.initOnce.Do(c.init)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return errors.New("creating http request failed").
			WithTag("url", url).
			Wrap(err)
	}
	return c.Do(ctx, req, out)
}

func (c *OAuthClient) Do(ctx context.Context, req *http.Request, out any) error {
	c.initOnce.Do(c.init)

	c.authMutex.Lock()
	if c.Token == nil || c.Token.Expiry.Before(time.Now()) {
		if err := c.Auth(ctx); err != nil {
			c.authMutex.Unlock()
			return errors.New("authenticating client failed").Wrap(err)
		}
	}
	c.authMutex.Unlock()

	req.Header.Set("Authorization", "Bearer "+c.Token.AccessToken)
	if c.ClientID != "" {
		req.Header.Set("Client-Id", c.ClientID)
	}

	res, err := c.Transport.RoundTrip(req)
	if err != nil {
		return errors.New("http request failed").
			WithTag("url", req.URL.String()).
			WithTag("method", req.Method).
			Wrap(err)
	}
	defer res.Body.Close()

	if res.StatusCode >= 400 {
		return errors.New("http request failed").
			WithTag("url", req.URL.String()).
			WithTag("method", req.Method).
			WithTag("code", res.StatusCode).
			WithTag("status", res.Status).
			Wrap(err)
	}

	b, err := io.ReadAll(res.Body)
	if err != nil {
		return errors.New("reading response failed").
			WithTag("url", req.URL.String()).
			WithTag("method", req.Method).
			Wrap(err)
	}

	if out == nil {
		return nil
	}

	if err := c.Decode(b, out); err != nil {
		return errors.New("decoding response failed").
			WithTag("url", req.URL.String()).
			WithTag("method", req.Method).
			Wrap(err)
	}
	return nil
}

func (c *OAuthClient) init() {
	if c.Transport == nil {
		c.Transport = http.DefaultTransport
	}
	if c.Decode == nil {
		c.Decode = json.Unmarshal
	}
	if c.Encode == nil {
		c.Encode = json.Marshal
	}
}
