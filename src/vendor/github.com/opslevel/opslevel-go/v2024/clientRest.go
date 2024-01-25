package opslevel

import (
	"github.com/go-resty/resty/v2"
)

type RestResponse struct {
	Result  string `json:"result"`
	Message string `json:"message"`
}

func NewRestClient(options ...Option) *resty.Client {
	client := resty.New()
	settings := newClientSettings(options...)
	client.SetBaseURL(settings.url)
	client.SetHeader("Accept", "application/json")
	for key, value := range settings.headers {
		client.SetHeader(key, value)
	}
	client.SetTimeout(settings.timeout)
	return client
}
