package config

import (
	"encoding/json"

	"github.com/opslevel/kubectl-opslevel/jq"

	_ "github.com/rs/zerolog/log"
)

type ServiceRegistrationConfig struct {
	Name string `default:".metadata.name"`
	Description string
	Owner string
	Lifecycle string
	Tier string
	Product string
	Language string
	Framework string
	Aliases []string // JQ expressions that return a single string or a string[]
	Tags []string // JQ expressions that return a single string or a map[string]string
}

type ServiceRegistrationParser struct {
	Name jq.JQ
	Description jq.JQ
	Owner jq.JQ
	Lifecycle jq.JQ
	Tier jq.JQ
	Product jq.JQ
	Language jq.JQ
	Framework jq.JQ
	Aliases []jq.JQ
	Tags []jq.JQ
}

type ServiceRegistration struct {
	Name string
	Description string `json:",omitempty"`
	Owner string `json:",omitempty"`
	Lifecycle string `json:",omitempty"`
	Tier string `json:",omitempty"`
	Product string `json:",omitempty"`
	Language string `json:",omitempty"`
	Framework string `json:",omitempty"`
	Aliases []string `json:",omitempty"`
	Tags map[string]string `json:",omitempty"`
}

func NewParser(config ServiceRegistrationConfig) (*ServiceRegistrationParser, error) {
	var err error
	createParser := func(filter string) *jq.JQ {
		if err != nil {
			return nil
		}
		var client jq.JQ
		client, err = jq.Create(filter)
		return &client
	}
	parser := ServiceRegistrationParser{}
	parser.Name = *createParser(config.Name)
	parser.Description = *createParser(config.Description)
	parser.Owner = *createParser(config.Owner)
	parser.Lifecycle = *createParser(config.Lifecycle)
	parser.Tier = *createParser(config.Tier)
	parser.Product = *createParser(config.Product)
	parser.Language = *createParser(config.Language)
	parser.Framework = *createParser(config.Framework)
	// TODO: Aliases & Tags
	if err != nil {
		return nil, err
	}
	return &parser, nil
}

func (parser *ServiceRegistrationParser) Parse(data []byte) (*ServiceRegistration, error) {
	var err error
	doParse := func(parser jq.JQ) string {
		if err != nil {
			return ""
		}
		var bytes []byte
		// TODO: if the parser == nil because the filter was bad - should just put the filter as the field?
		// TODO: if the parser's filter was "" then we shouldn't run this
		bytes, err = parser.Run(data)
		if err != nil {
			return ""
		}
		var output string
		jsonErr := json.Unmarshal(bytes, &output)
		// TODO: figure out how to properly handle failed JQ runs
		if (jsonErr != nil) {
			return ""
		}
		return output
	}
	service := ServiceRegistration{}
	service.Name = doParse(parser.Name)
	service.Description = doParse(parser.Description)
	service.Owner = doParse(parser.Owner)
	service.Lifecycle = doParse(parser.Lifecycle)
	service.Tier = doParse(parser.Tier)
	service.Product = doParse(parser.Product)
	service.Language = doParse(parser.Language)
	service.Framework = doParse(parser.Framework)
	// TODO: Aliases & Tags
	if err != nil {
		return nil, err
	}
	return &service, nil
}