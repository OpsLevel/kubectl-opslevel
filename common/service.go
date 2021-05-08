package common

import (
	"encoding/json"

	"github.com/opslevel/kubectl-opslevel/config"
	"github.com/opslevel/kubectl-opslevel/k8sutils"
	"github.com/opslevel/kubectl-opslevel/opslevel"

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
	Tools       []JQParser
}

type ServiceRegistration struct {
	Name        string
	Description string                     `json:",omitempty"`
	Owner       string                     `json:",omitempty"`
	Lifecycle   string                     `json:",omitempty"`
	Tier        string                     `json:",omitempty"`
	Product     string                     `json:",omitempty"`
	Language    string                     `json:",omitempty"`
	Framework   string                     `json:",omitempty"`
	Aliases     []string                   `json:",omitempty"`
	Tags        map[string]string          `json:",omitempty"`
	Tools       []opslevel.ToolCreateInput `json:",omitempty"` // This is a concrete class so fields are validated during `service preview`
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
	for _, tool := range c.Tools {
		parser.Tools = append(parser.Tools, NewJQParser(tool))
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
	// TODO: the following chunks should probably be extracted into named functions for clarity
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
			// TODO: log warnings about a JQ filter that went unused because it returned an invalid type that we dont know how to handle
		}
	}
	service.Aliases = removeDuplicates(service.Aliases)
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
			// TODO: log warnings about a JQ filter that went unused because it returned an invalid type that we dont know how to handle
		}
	}
	service.Tools = []opslevel.ToolCreateInput{}
	for _, tool := range parser.Tools {
		output := tool.Parse(data)
		if output == nil {
			continue
		}
		switch output.Type {
		case StringStringMap:
			if tool, err := ConvertToTool(output.StringMap); err == nil {
				service.Tools = append(service.Tools, *tool)
			}
			break
		case StringStringMapArray:
			for _, item := range output.StringMapArray {
				if tool, err := ConvertToTool(item); err == nil {
					service.Tools = append(service.Tools, *tool)
				}
			}
			break
		}
	}
	return &service
}

func removeDuplicates(data []string) []string {
	keys := make(map[string]bool)
	list := []string{}

	for _, entry := range data {
		if _, value := keys[entry]; !value {
			keys[entry] = true
			list = append(list, entry)
		}
	}
	return list
}

func ConvertToTool(data map[string]string) (*opslevel.ToolCreateInput, error) {
	bytes, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}
	tool := &opslevel.ToolCreateInput{}
	if unmarshalErr := json.Unmarshal(bytes, tool); unmarshalErr != nil {
		return nil, unmarshalErr
	}
	return tool, nil
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
