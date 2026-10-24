package vi

import (
	"crypto/rand"
	"encoding/base64"
	"io"
	"net/url"
	"strings"

	"github.com/maxence-charriere/go-app/v11/pkg/errors"
)

// OAuthProvider defines the configuration for an OAuth provider.
type OAuthProvider struct {
	Slug              string   `json:",omitempty"`
	Name              string   `json:",omitempty"`
	AuthorizeEndpoint string   `json:",omitempty"`
	Scopes            []string `json:",omitempty"`
}

// AuthorizeURL builds the authorization URL for the provider using the given
// query parameters.
func (p OAuthProvider) AuthorizeURL(query map[string]string) string {
	var b strings.Builder
	b.WriteString(p.AuthorizeEndpoint)
	b.WriteString("?response_type=code&scope=")
	b.WriteString(strings.Join(p.Scopes, " "))

	for k, v := range query {
		b.WriteString("&")
		b.WriteString(k)
		b.WriteString("=")
		b.WriteString(v)
	}

	u, err := url.Parse(b.String())
	if err != nil {
		return ""
	}
	return u.String()
}

// OAuthCode holds data related to an OAuth authorization code exchange.
type OAuthCode struct {
	Code             string `json:""`
	PCKECodeVerifier string `json:",omitempty"`
	Provider         string `json:",omitempty"`
	State            string `json:"-"`
}

// GenerateToken returns a random token of length n.
func GenerateToken(n int) (string, error) {
	v := make([]byte, n)
	if _, err := io.ReadFull(rand.Reader, v); err != nil {
		return "", errors.New("generating token failed").
			WithTag("len", n).
			Wrap(err)
	}
	return base64.RawURLEncoding.EncodeToString(v)[:n], nil
}
