package portal

import (
	"context"
	"encoding/json"
	"fmt"
	"net/url"
	"strings"
	"time"

	"github.com/maxence-charriere/go-app/v11/pkg/errors"
	"github.com/maxence-charriere/sunforge/pkg/conduit"
	"golang.org/x/oauth2"
)

var (
	patreonUserFields = formatPatreonFields("user",
		"about",
		"can_see_nsfw",
		"created",
		"email",
		"first_name",
		"full_name",
		"hide_pledges",
		"image_url",
		"is_email_verified",
		"last_name",
		"like_count",
		"social_connections",
		"thumb_url",
		"url",
		"vanity",
	)
	patreonMemberFields = formatPatreonFields("member",
		"campaign_lifetime_support_cents",
		"currently_entitled_amount_cents",
		"email",
		"full_name",
		"is_follower",
		"last_charge_date",
		"last_charge_status",
		"lifetime_support_cents",
		"next_charge_date",
		"note",
		"patron_status",
		"pledge_cadence",
		"pledge_relationship_start",
		"will_pay_amount_cents",
	)
	patreonCampaignFields = formatPatreonFields("campaign",
		"created_at",
		"creation_name",
		"google_analytics_id",
		"has_rss",
		"has_sent_rss_notify",
		"image_small_url",
		"image_url",
		"is_charged_immediately",
		"is_monthly",
		"is_nsfw",
		"main_video_embed",
		"main_video_url",
		"one_liner",
		"patron_count",
		"pay_per_name",
		"pledge_url",
		"published_at",
		"rss_artwork_url",
		"rss_feed_title",
		"show_earnings",
		"summary",
		"url",
		"vanity",
	)
)

type PatreonIdentity struct {
	User     PatreonUser
	Campaign PatreonCampaign
	Member   PatreonMember
}

type PatreonUser struct {
	ID                string    `json:"id"`
	About             string    `json:"about"`
	CanSeeNSFW        bool      `json:"can_see_nsfw"`
	Created           time.Time `json:"created"`
	Email             string    `json:"email"`
	FirstName         string    `json:"first_name"`
	FullName          string    `json:"full_name"`
	HidePledges       bool      `json:"hide_pledges"`
	ImageURL          string    `json:"image_url"`
	IsEmailVerified   bool      `json:"is_email_verified"`
	LastName          string    `json:"last_name"`
	LikeCount         int       `json:"like_count"`
	SocialConnections struct {
		DeviantArt PatreonSocialConnection `json:"deviantart"`
		Discord    PatreonSocialConnection `json:"discord"`
		Facebook   PatreonSocialConnection `json:"facebook"`
		Google     PatreonSocialConnection `json:"google"`
		Instagram  PatreonSocialConnection `json:"instagram"`
		Reddit     PatreonSocialConnection `json:"reddit"`
		Spotify    PatreonSocialConnection `json:"spotify"`
		Twitch     PatreonSocialConnection `json:"twitch"`
		Twitter    PatreonSocialConnection `json:"twitter"`
		Youtube    PatreonSocialConnection `json:"youtube"`
	} `json:"social_connections"`
	ThumbURL string `json:"thumb_url"`
	URL      string `json:"url"`
	Vanity   string `json:"vanity"`
}

type PatreonSocialConnection struct {
	Scopes []string `json:"scopes"`
	URL    string   `json:"url"`
	UserID string   `json:"user_id"`
}

type PatreonCampaign struct {
	ID                   string    `json:"id"`
	CreatedAt            time.Time `json:"created_at"`
	CreationName         string    `json:"creation_name"`
	DiscordServerID      string    `json:"discord_server_id"`
	GoogleAnalyticsID    string    `json:"google_analytics_id"`
	HasRSS               bool      `json:"has_rss"`
	HasSentRSSNotify     bool      `json:"has_sent_rss_notify"`
	ImageSmallURL        string    `json:"image_small_url"`
	ImageURL             string    `json:"image_url"`
	IsChargedImmediately bool      `json:"is_charged_immediately"`
	IsMonthly            bool      `json:"is_monthly"`
	IsNSFW               bool      `json:"is_nsfw"`
	MainVideoEmbed       string    `json:"main_video_embed"`
	MainVideoURL         string    `json:"main_video_url"`
	OneLiner             string    `json:"one_liner"`
	PatronCount          int       `json:"patron_count"`
	PayPerName           string    `json:"pay_per_name"`
	PledgeURL            string    `json:"pledge_url"`
	PublishedAt          time.Time `json:"published_at"`
	RSSArtworkURL        string    `json:"rss_artwork_url"`
	RSSFeedTitle         string    `json:"rss_feed_title"`
	ShowEarnings         bool      `json:"show_earnings"`
	Summary              string    `json:"summary"`
	Thanks_embed         string    `json:"thanks_embed"`
	ThanksMsg            string    `json:"thanks_msg"`
	ThanksVideoURL       string    `json:"thanks_video_url"`
	URL                  string    `json:"url"`
	Vanity               string    `json:"vanity"`
}

type PatreonMember struct {
	ID                           string    `json:"id"`
	CampaignLifetimeSupportCents int       `json:"campaign_lifetime_support_cents"`
	CurrentlyEntitledAmountCents int       `json:"currently_entitled_amount_cents"`
	FullName                     string    `json:"full_name"`
	IsFollower                   bool      `json:"is_follower"`
	LastChargeDate               time.Time `json:"last_charge_date"`
	LastChargeStatus             string    `json:"last_charge_status"`
	LifetimeSupportCents         int       `json:"lifetime_support_cents"`
	NextChargeDate               time.Time `json:"next_charge_date"`
	Note                         string    `json:"note"`
	PatronStatus                 string    `json:"patron_status"`
	PledgeCadence                int       `json:"pledge_cadence"`
	PledgeRelationshipStart      time.Time `json:"pledge_relationship_start"`
	WillPayAmountCents           int       `json:"will_pay_amount_cents"`
}

type PatreonWebhookMember struct {
	Member   PatreonMember
	Campaign PatreonCampaign
	User     PatreonUser
}

type patreonResponse struct {
	Data     patreonResource   `json:"data"`
	Included []patreonResource `json:"included"`
}

type patreonResource struct {
	ID         string      `json:"id"`
	Type       string      `json:"type"`
	Attributes interface{} `json:"attributes"`
}

func (r *patreonResource) UnmarshalJSON(d []byte) error {
	var res struct {
		ID         string          `json:"id"`
		Type       string          `json:"type"`
		Attributes json.RawMessage `json:"attributes"`
	}
	err := json.Unmarshal(d, &res)
	if err != nil {
		return err
	}
	r.ID = res.ID
	r.Type = res.Type

	switch res.Type {
	case "user":
		u := PatreonUser{ID: r.ID}
		err = json.Unmarshal([]byte(res.Attributes), &u)
		r.Attributes = u

	case "member":
		m := PatreonMember{ID: r.ID}
		err = json.Unmarshal([]byte(res.Attributes), &m)
		r.Attributes = m

	case "campaign":
		c := PatreonCampaign{ID: r.ID}
		err = json.Unmarshal([]byte(res.Attributes), &c)
		r.Attributes = c
	}
	return err
}

type PatreonClient struct {
	*conduit.OAuthClient
	WebhookSecret string
}

func MakePatreonClient(webhookSecret string, c *conduit.OAuthClient) PatreonClient {
	c.Endpoint = oauth2.Endpoint{
		AuthURL:  "https://www.patreon.com/api/oauth2/authorize",
		TokenURL: "https://www.patreon.com/api/oauth2/token",
	}

	return PatreonClient{
		OAuthClient:   c,
		WebhookSecret: webhookSecret,
	}
}

func (c *PatreonClient) AuthAsUser(ctx context.Context, code string) (PatreonClient, error) {
	userClient, err := c.OAuthClient.AuthAsUser(ctx, code)
	if err != nil {
		return PatreonClient{}, err
	}

	return PatreonClient{
		WebhookSecret: c.WebhookSecret,
		OAuthClient:   userClient,
	}, nil
}

func (c PatreonClient) Identity(ctx context.Context) (PatreonIdentity, error) {
	url := fmt.Sprintf("https://www.patreon.com/api/oauth2/v2/identity?include=%s&%s&%s&%s",
		url.QueryEscape(strings.Join([]string{
			"campaign",
			"memberships",
		}, ",")),
		patreonUserFields,
		patreonMemberFields,
		patreonCampaignFields,
	)

	var res patreonResponse
	if err := c.Get(ctx, url, &res); err != nil {
		return PatreonIdentity{}, err
	}

	user, ok := res.Data.Attributes.(PatreonUser)
	if !ok || user.ID == "" {
		return PatreonIdentity{}, errors.New("user is not set")
	}

	campaign, _ := patreonResourceByType(res.Included, "campaign").Attributes.(PatreonCampaign)
	member, _ := patreonResourceByType(res.Included, "member").Attributes.(PatreonMember)

	return PatreonIdentity{
		User:     user,
		Campaign: campaign,
		Member:   member,
	}, nil
}

func (c PatreonClient) ParseWebhookMember(v []byte) (PatreonWebhookMember, error) {
	var res patreonResponse
	if err := json.Unmarshal(v, &res); err != nil {
		return PatreonWebhookMember{}, errors.New("decoding failed").Wrap(err)
	}

	member, ok := res.Data.Attributes.(PatreonMember)
	if !ok {
		return PatreonWebhookMember{}, errors.New("member is not set")
	}

	campaign, ok := patreonResourceByType(res.Included, "campaign").Attributes.(PatreonCampaign)
	if !ok {
		return PatreonWebhookMember{}, errors.New("campaign is not set")
	}

	user, ok := patreonResourceByType(res.Included, "user").Attributes.(PatreonUser)
	if !ok {
		return PatreonWebhookMember{}, errors.New("user is not set")
	}

	return PatreonWebhookMember{
		Member:   member,
		Campaign: campaign,
		User:     user,
	}, nil
}

func formatPatreonFields(field string, attrs ...string) string {
	var b strings.Builder
	b.WriteString(url.QueryEscape(fmt.Sprintf("fields[%s]", field)))
	b.WriteByte('=')
	b.WriteString(url.QueryEscape(strings.Join(attrs, ",")))
	return b.String()
}

func patreonResourceByType(r []patreonResource, typ string) patreonResource {
	for _, rs := range r {
		if rs.Type == typ {
			return rs
		}
	}
	return patreonResource{}
}
