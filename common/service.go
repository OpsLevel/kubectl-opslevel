package common

import (
	"github.com/opslevel/kubectl-opslevel/config"
	"github.com/opslevel/kubectl-opslevel/k8sutils"

	_ "github.com/rs/zerolog/log"
)

type ServiceRegistrationParser struct {
	Name        JQParser
	Description JQParser
	Owner       JQParser
	Lifecycle   JQParser
	Tier        JQParser
	Product     JQParser
	Language    JQParser
	Framework   JQParser
	Aliases     []JQParser
	Tags        []JQParser
}

type ServiceRegistration struct {
	Name        string
	Description string            `json:",omitempty"`
	Owner       string            `json:",omitempty"`
	Lifecycle   string            `json:",omitempty"`
	Tier        string            `json:",omitempty"`
	Product     string            `json:",omitempty"`
	Language    string            `json:",omitempty"`
	Framework   string            `json:",omitempty"`
	Aliases     []string          `json:",omitempty"`
	Tags        map[string]string `json:",omitempty"`
}

func NewParser(c config.ServiceRegistrationConfig) *ServiceRegistrationParser {
	parser := ServiceRegistrationParser{}
	parser.Name = NewJQParser(c.Name)
	parser.Description = NewJQParser(c.Description)
	parser.Owner = NewJQParser(c.Owner)
	parser.Lifecycle = NewJQParser(c.Lifecycle)
	parser.Tier = NewJQParser(c.Tier)
	parser.Product = NewJQParser(c.Product)
	parser.Language = NewJQParser(c.Language)
	parser.Framework = NewJQParser(c.Framework)
	for _, alias := range c.Aliases {
		parser.Aliases = append(parser.Aliases, NewJQParser(alias))
	}
	for _, tag := range c.Tags {
		parser.Tags = append(parser.Tags, NewJQParser(tag))
	}
	return &parser
}

func GetString(parser JQParser, data []byte) string {
	output := parser.Parse(data)
	if output == nil {
		return ""
	}
	if output.Type == String {
		return output.StringObj
	}
	return ""
}

func (parser *ServiceRegistrationParser) Parse(data []byte) *ServiceRegistration {
	service := ServiceRegistration{}
	service.Name = GetString(parser.Name, data)
	service.Description = GetString(parser.Description, data)
	service.Owner = GetString(parser.Owner, data)
	service.Lifecycle = GetString(parser.Lifecycle, data)
	service.Tier = GetString(parser.Tier, data)
	service.Product = GetString(parser.Product, data)
	service.Language = GetString(parser.Language, data)
	service.Framework = GetString(parser.Framework, data)
	// TODO: need to treat service.Aliases as a "hash set" to not append duplicates
	for _, alias := range parser.Aliases {
		output := alias.Parse(data)
		if output == nil {
			continue
		}
		switch output.Type {
		case String:
			service.Aliases = append(service.Aliases, output.StringObj)
			break
		case StringArray:
			for _, item := range output.StringArray {
				service.Aliases = append(service.Aliases, item)
			}
			break
		}
	}
	service.Tags = map[string]string{}
	for _, tag := range parser.Tags {
		output := tag.Parse(data)
		if output == nil {
			continue
		}
		switch output.Type {
		case StringStringMap:
			for k, v := range output.StringMap {
				service.Tags[k] = v
			}
			break
		case StringStringMapArray:
			for _, item := range output.StringMapArray {
				for k, v := range item {
					service.Tags[k] = v
				}
			}
			break
		}
	}
	return &service
}

func QueryForServices(c *config.Config) ([]ServiceRegistration, error) {
	var parser *ServiceRegistrationParser
	var services []ServiceRegistration
	k8sClient := k8sutils.CreateKubernetesClient()

	for _, importConfig := range c.Service.Import {
		parser = NewParser(importConfig.OpslevelConfig)
		process := func(resource []byte) error {
			services = append(services, *parser.Parse(resource))
			return nil
		}
		if err := k8sClient.Query(importConfig.SelectorConfig, process); err != nil {
			return services, err
		}
	}
	return services, nil
}
