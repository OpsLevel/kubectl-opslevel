package opslevel

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/hashicorp/go-retryablehttp"
	"github.com/hasura/go-graphql-client"
)

type Client struct {
	pageSize graphql.Int
	client   *graphql.Client
}

// Deprecated: Use NewGQLClient instead
func NewClient(apiToken string, options ...Option) *Client {
	options = append(options, SetAPIToken(apiToken))
	return NewGQLClient(options...)
}

func NewGQLClient(options ...Option) *Client {
	settings := newClientSettings(options...)

	retryClient := retryablehttp.NewClient()
	retryClient.RetryMax = settings.retries
	retryClient.Logger = nil

	standardClient := retryClient.StandardClient()
	var url string
	if strings.Contains(settings.url, "/LOCAL_TESTING/") {
		url = settings.url
	} else {
		url = fmt.Sprintf("%s/graphql", settings.url)
	}

	auth := fmt.Sprintf("Bearer %s", settings.token)
	modifier := graphql.RequestModifier(
		func(r *http.Request) {
			r.Header.Add("Authorization", auth)
			for key, value := range settings.headers {
				r.Header.Add(key, value)
			}
		})

	return &Client{
		pageSize: graphql.Int(settings.pageSize),
		client:   graphql.NewClient(url, standardClient).WithRequestModifier(modifier),
	}
}

func (client *Client) InitialPageVariables() PayloadVariables {
	return PayloadVariables{
		"after": "",
		"first": client.pageSize,
	}
}

func (client *Client) InitialPageVariablesPointer() *PayloadVariables {
	v := PayloadVariables{
		"after": "",
		"first": client.pageSize,
	}
	return &v
}

func (client *Client) Query(q interface{}, variables map[string]interface{}, options ...graphql.Option) error {
	return client.QueryCTX(context.Background(), q, variables, options...)
}

func (client *Client) QueryCTX(ctx context.Context, q interface{}, variables map[string]interface{}, options ...graphql.Option) error {
	return client.client.Query(ctx, q, variables, options...)
}

func (client *Client) Mutate(m interface{}, variables map[string]interface{}, options ...graphql.Option) error {
	return client.MutateCTX(context.Background(), m, variables, options...)
}

func (client *Client) MutateCTX(ctx context.Context, m interface{}, variables map[string]interface{}, options ...graphql.Option) error {
	return client.client.Mutate(ctx, m, variables, options...)
}

func (client *Client) ExecRaw(q string, variables map[string]interface{}, options ...graphql.Option) ([]byte, error) {
	return client.ExecRawCTX(context.Background(), q, variables, options...)
}

func (client *Client) ExecRawCTX(ctx context.Context, q string, variables map[string]interface{}, options ...graphql.Option) ([]byte, error) {
	return client.client.ExecRaw(ctx, q, variables, options...)
}

func (client *Client) Validate() error {
	var q struct {
		Account struct {
			Id ID
		}
	}
	err := client.Query(&q, nil)
	// TODO: we should probably use a custom OpsLevelClientError type - https://www.digitalocean.com/community/tutorials/creating-custom-errors-in-go
	if err != nil {
		return fmt.Errorf("client validation error: %s", err.Error())
	}
	return nil
}

func WithName(name string) graphql.Option {
	return graphql.OperationName(name)
}
