package portal

import (
	"context"
	"path"

	"github.com/maxence-charriere/sunforge/pkg/conduit"
	"golang.org/x/oauth2/twitch"
)

type TwitchClient struct {
	*conduit.OAuthClient
}

func MakeTwitchClient(c *conduit.OAuthClient) TwitchClient {
	c.Endpoint = twitch.Endpoint
	return TwitchClient{OAuthClient: c}
}

func (c TwitchClient) AuthAsUser(ctx context.Context, code string) (TwitchClient, error) {
	userClient, err := c.OAuthClient.AuthAsUser(ctx, code)
	if err != nil {
		return TwitchClient{}, err
	}
	return TwitchClient{OAuthClient: userClient}, nil
}

type TwitchStream struct {
	ID           string   `json:"id"`
	UserID       string   `json:"user_id"`
	UserLogin    string   `json:"user_login"`
	UserName     string   `json:"user_name"`
	GameID       string   `json:"game_id"`
	GameName     string   `json:"game_name"`
	Type         string   `json:"type"`
	Title        string   `json:"title"`
	ViewerCount  int      `json:"viewer_count"`
	StartedAt    string   `json:"started_at"`
	Language     string   `json:"language"`
	ThumbnailUrl string   `json:"thumbnail_url"`
	TagIDs       []string `json:"tag_ids"`
	IsMature     bool     `json:"is_mature"`
}

func (c TwitchClient) Stream(ctx context.Context, url string) (TwitchStream, error) {
	login := path.Base(url)
	u := "https://api.twitch.tv/helix/streams?user_login=" + login

	var s struct {
		Data []TwitchStream `json:"data"`
	}

	err := c.Get(ctx, u, &s)
	if err != nil {
		return TwitchStream{}, err
	}

	for _, s := range s.Data {
		if s.UserLogin == login {
			return s, nil
		}
	}

	return TwitchStream{}, nil

}

type TwitchUser struct {
	ID              string `json:"id"`
	Login           string `json:"login"`
	DisplayName     string `json:"display_name"`
	ProfileImageURL string `json:"profile_image_url"`
	OfflineImageURL string `json:"offline_image_url"`
}

func (c TwitchClient) GetUser(ctx context.Context) (TwitchUser, error) {
	u := "https://api.twitch.tv/helix/users"
	var s struct {
		Data []TwitchUser `json:"data"`
	}

	err := c.Get(ctx, u, &s)
	if err != nil {
		return TwitchUser{}, err
	}

	for _, u := range s.Data {
		if u.ID != "" {
			return u, nil
		}
	}
	return TwitchUser{}, nil
}
