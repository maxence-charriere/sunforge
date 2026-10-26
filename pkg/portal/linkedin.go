package portal

import (
	"bytes"
	"context"
	"net/http"
	"time"

	"github.com/maxence-charriere/go-app/v11/pkg/errors"
	"github.com/maxence-charriere/sunforge/pkg/conduit"
	"golang.org/x/oauth2/linkedin"
)

type LinkedInClient struct {
	*conduit.OAuthClient
}

func MakeLinkedInClient(c *conduit.OAuthClient) LinkedInClient {
	c.Endpoint = linkedin.Endpoint
	return LinkedInClient{OAuthClient: c}
}

func (c LinkedInClient) AuthAsUser(ctx context.Context, code string) (LinkedInClient, error) {
	userClient, err := c.OAuthClient.AuthAsUser(ctx, code)
	if err != nil {
		return LinkedInClient{}, err
	}
	return LinkedInClient{OAuthClient: userClient}, nil
}

func (c LinkedInClient) AuthWithAccessToken(accessToken string, expiresAt time.Time) LinkedInClient {
	return LinkedInClient{OAuthClient: c.OAuthClient.AuthWithAccessToken(accessToken, expiresAt)}
}

type LinkedInMemberDetails struct {
	Sub           string `json:"sub"`
	Name          string `json:"name"`
	GivenName     string `json:"given_name"`
	FamilyName    string `json:"family_name"`
	Picture       string `json:"picture"`
	Email         string `json:"email"`
	EmailVerified bool   `json:"email_verified"`
}

func (c LinkedInClient) MemberDetails(ctx context.Context) (LinkedInMemberDetails, error) {
	var member LinkedInMemberDetails
	if err := c.Get(ctx, "https://api.linkedin.com/v2/userinfo", &member); err != nil {
		return LinkedInMemberDetails{}, err
	}
	return member, nil
}

type LinkedInUgcPostsQuery struct {
	Author          string                             `json:"author"`
	LifecycleState  string                             `json:"lifecycleState"`
	SpecificContent LinkedInUgcSpecificContent         `json:"specificContent"`
	Visibility      LinkedInUgcMemberNetworkVisibility `json:"visibility"`
}

type LinkedInUgcSpecificContent struct {
	ShareContent LinkedInUgcShareContent `json:"com.linkedin.ugc.ShareContent"`
}

type LinkedInUgcShareContent struct {
	ShareCommentary    LinkedInUgcText    `json:"shareCommentary"`
	ShareMediaCategory string             `json:"shareMediaCategory"`
	Media              []LinkedInUgcMedia `json:"media"`
}

type LinkedInUgcText struct {
	Text string `json:"text"`
}

type LinkedInUgcMemberNetworkVisibility struct {
	MemberNetworkVisibility string `json:"com.linkedin.ugc.MemberNetworkVisibility"`
}

type LinkedInUgcMedia struct {
	Status      string          `json:"status"`
	OriginalUrl string          `json:"originalUrl"`
	Title       LinkedInUgcText `json:"title,omitempty"`
}

func (c LinkedInClient) UgcPosts(ctx context.Context, q LinkedInUgcPostsQuery) error {
	body, err := c.Encode(q)
	if err != nil {
		return errors.New("encoding body failed").Wrap(err)
	}

	req, err := http.NewRequestWithContext(ctx,
		http.MethodPost,
		"https://api.linkedin.com/v2/ugcPosts",
		bytes.NewBuffer(body),
	)
	if err != nil {
		return errors.New("making http request failed").Wrap(err)
	}
	req.Header.Set("X-Restli-Protocol-Version", "2.0.0")
	return c.Do(ctx, req, nil)
}
