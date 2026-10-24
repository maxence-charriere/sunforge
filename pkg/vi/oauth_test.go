package vi

import (
	"strconv"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestOAuthProviderAuthorizeURL(t *testing.T) {
	p := OAuthProvider{
		Slug:              "test",
		AuthorizeEndpoint: "https://test.io/auth",
		Scopes:            []string{"a", "b"},
	}

	authorizeURL := p.AuthorizeURL(map[string]string{
		"client_id": "4242",
	})

	require.Equal(t, "https://test.io/auth?response_type=code&scope=a b&client_id=4242", authorizeURL)
}

func TestGenerateToken(t *testing.T) {
	utests := []int{
		8,
		16,
		32,
		64,
		128,
		256,
	}

	for _, u := range utests {
		t.Run(strconv.Itoa(u), func(t *testing.T) {
			token, err := GenerateToken(u)
			require.NoError(t, err)
			require.Len(t, token, u)
			t.Log(token)
		})
	}
}
