package opslevel

import (
	"fmt"
	"os"
	"runtime"
	"strings"
	"time"
)

type ClientSettings struct {
	url      string
	token    string
	timeout  time.Duration
	retries  int
	headers  map[string]string
	pageSize int // Only Used by GQL
}

type Option func(*ClientSettings)

func newClientSettings(options ...Option) *ClientSettings {
	settings := &ClientSettings{
		url:     "https://app.opslevel.com",
		token:   os.Getenv("OPSLEVEL_API_TOKEN"),
		timeout: time.Second * 10,
		retries: 10,

		pageSize: 100,
		headers: map[string]string{
			"User-Agent":         buildUserAgent(""),
			"GraphQL-Visibility": "public",
		},
	}
	for _, opt := range options {
		opt(settings)
	}
	return settings
}

func SetAPIToken(apiToken string) Option {
	return func(c *ClientSettings) {
		c.token = apiToken
	}
}

func SetURL(url string) Option {
	return func(c *ClientSettings) {
		c.url = strings.TrimRight(url, "/")
	}
}

func SetHeader(key string, value string) Option {
	return func(c *ClientSettings) {
		c.headers[key] = value
	}
}

func SetHeaders(headers map[string]string) Option {
	return func(c *ClientSettings) {
		for key, value := range headers {
			c.headers[key] = value
		}
	}
}

func SetUserAgentExtra(extra string) Option {
	return SetHeader("User-Agent", buildUserAgent(extra))
}

func SetTimeout(amount time.Duration) Option {
	return func(c *ClientSettings) {
		c.timeout = amount
	}
}

func SetMaxRetries(amount int) Option {
	return func(c *ClientSettings) {
		c.retries = amount
	}
}

func SetAPIVisibility(visibility string) Option {
	return SetHeader("GraphQL-Visibility", visibility)
}

func SetPageSize(size int) Option {
	return func(c *ClientSettings) {
		c.pageSize = size
	}
}

/*
Return a string suitable for use as a User-Agent header.
The string will be of the form:

<agent_name>/<agent_version> go/<go_ver> <plat_name>/<plat_ver> client/<code_extras> user/<user_extras>
*/
func buildUserAgent(extra string) string {
	base := fmt.Sprintf("opslevel-go/%s go/%s %s/%s", clientVersion, runtime.Version(), runtime.GOOS, runtime.GOARCH)
	if extra != "" {
		base = fmt.Sprintf("%s client/%s", base, extra)
	}
	if value, present := os.LookupEnv("OPSLEVEL_USER_AGENT_EXTRAS"); present {
		base = fmt.Sprintf("%s user/%s", base, value)
	}
	return base
}
