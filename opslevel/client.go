package opslevel

import (
	"fmt"
	"bytes"
	"net/http"
)

type Client struct {
	httpClient http.Client
	apiToken string
}

func NewClient(apiToken string) *Client {
	return &Client {
		httpClient: http.Client{},
		apiToken: fmt.Sprintf("Bearer %v", apiToken),
	}
}

func (c *Client) Post(query string) (*http.Response, error) {
	req, err := http.NewRequest("POST", "https://api.opslevel.com/graphql", bytes.NewBuffer([]byte(query)))
	if err != nil {
		return nil, err
	}
	return c.Do(req)
}

func (c *Client) Do(req *http.Request) (*http.Response, error) {
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Accept", "application/json")
	req.Header.Add("Authorization", c.apiToken)
	return c.httpClient.Do(req)
}