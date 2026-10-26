package portal

import (
	"context"

	"github.com/maxence-charriere/sunforge/pkg/conduit"
	"golang.org/x/oauth2/google"
)

type YoutubeChannel struct {
	Items []struct {
		Kind    string `json:"kind"`
		ID      string `json:"id"`
		Snippet struct {
			Title       string `json:"title"`
			Description string `json:"description"`
			CustomUrl   string `json:"customUrl"`
			Thumbnails  struct {
				Default struct {
					URL string `json:"url"`
				} `json:"default"`
				Medium struct {
					URL string `json:"url"`
				} `json:"medium"`
				High struct {
					URL string `json:"url"`
				} `json:"high"`
			} `json:"thumbnails"`
		} `json:"snippet"`
	} `json:"items"`
}

type YouTube struct {
	*conduit.OAuthClient
}

func MakeYouTube(c *conduit.OAuthClient) YouTube {
	c.Endpoint = google.Endpoint
	return YouTube{OAuthClient: c}
}

func (c YouTube) AuthAsUser(ctx context.Context, code string) (YouTube, error) {
	userClient, err := c.OAuthClient.AuthAsUser(ctx, code, conduit.OAuthOption{
		Key: "pcke_code_challenge",
	})
	if err != nil {
		return YouTube{}, err
	}

	return YouTube{OAuthClient: userClient}, nil
}

func (c YouTube) GetCurrentUserChannel(ctx context.Context) (YoutubeChannel, error) {
	var channel YoutubeChannel
	err := c.Get(ctx, "https://www.googleapis.com/youtube/v3/channels?part=snippet&mine=true", &channel)
	return channel, err
}
