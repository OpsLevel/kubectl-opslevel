package opslevel

import (
	"net/http"
	"context"

	"golang.org/x/oauth2"
	"github.com/shurcooL/graphql"
)

const defaultURL = "https://api.opslevel.com/graphql"

type ClientSettings struct {
	url string
	ctx context.Context
	httpClient *http.Client
}

type Client struct {
	url string
	ctx context.Context
	client *graphql.Client
}

type option func(*ClientSettings)

func SetURL(url string) option {
	return func(c *ClientSettings) {
		c.url = url
	}
}

func NewClient(apiToken string, options ...option) *Client {
	httpToken := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: apiToken, TokenType: "Bearer"},
	)
	settings := &ClientSettings{
		url: defaultURL,
		ctx: context.Background(),
		httpClient: oauth2.NewClient(context.Background(), httpToken),
	}
	for _, opt := range options {
		opt(settings)
	}
	return &Client{
		url: settings.url,
		ctx: settings.ctx,
		client: graphql.NewClient(settings.url, settings.httpClient),
	}
}

func (c *Client) Query(q interface{}, variables map[string]interface{}) error {
	return c.client.Query(c.ctx, q, variables)
}

func (c *Client) Mutate(m interface{}, variables map[string]interface{}) error {
	return c.client.Mutate(c.ctx, m, variables)
}