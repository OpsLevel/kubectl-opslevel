package common

import (
	"encoding/json"

	"github.com/opslevel/kubectl-opslevel/config"
	"github.com/opslevel/kubectl-opslevel/jq"
	"github.com/opslevel/kubectl-opslevel/k8sutils"
	"github.com/opslevel/opslevel-go"

	_ "github.com/rs/zerolog/log"
)

type ServiceRegistrationParser struct {
	Name         JQParser
	Description  JQParser
	Owner        JQParser
	Lifecycle    JQParser
	Tier         JQParser
	Product      JQParser
	Language     JQParser
	Framework    JQParser
	Aliases      []JQParser
	TagAssigns   []JQParser
	TagCreates   []JQParser
	Tools        []JQParser
	Repositories []JQParser
}

type ServiceRegistration struct {
	Name         string
	Description  string                                  `json:",omitempty"`
	Owner        string                                  `json:",omitempty"`
	Lifecycle    string                                  `json:",omitempty"`
	Tier         string                                  `json:",omitempty"`
	Product      string                                  `json:",omitempty"`
	Language     string                                  `json:",omitempty"`
	Framework    string                                  `json:",omitempty"`
	Aliases      []string                                `json:",omitempty"`
	TagAssigns   map[string]string                       `json:",omitempty"`
	TagCreates   map[string]string                       `json:",omitempty"`
	Tools        []opslevel.ToolCreateInput              `json:",omitempty"` // This is a concrete class so fields are validated during `service preview`
	Repositories []opslevel.ServiceRepositoryCreateInput `json:",omitempty"` // This is a concrete class so fields are validated during `service preview`
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
	parser.Aliases = append(parser.Aliases, NewJQParser("\"k8s:\\(.metadata.name)-\\(.metadata.namespace)\""))
	for _, alias := range c.Aliases {
		parser.Aliases = append(parser.Aliases, NewJQParser(alias))
	}
	for _, tag := range c.Tags.Assign {
		parser.TagAssigns = append(parser.TagAssigns, NewJQParser(tag))
	}
	for _, tag := range c.Tags.Create {
		parser.TagCreates = append(parser.TagCreates, NewJQParser(tag))
	}
	for _, tool := range c.Tools {
		parser.Tools = append(parser.Tools, NewJQParser(tool))
	}
	for _, repository := range c.Repositories {
		parser.Repositories = append(parser.Repositories, NewJQParser(repository))
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
				if item == "" {
					continue
				}
				service.Aliases = append(service.Aliases, item)
			}
			break
			// TODO: log warnings about a JQ filter that went unused because it returned an invalid type that we dont know how to handle
		}
	}
	service.Aliases = removeDuplicates(service.Aliases)

	service.TagAssigns = map[string]string{}
	for _, tag := range parser.TagAssigns {
		output := tag.Parse(data)
		if output == nil {
			continue
		}
		switch output.Type {
		case StringStringMap:
			for k, v := range output.StringMap {
				if k == "" || v == "" {
					continue
				}
				service.TagAssigns[k] = v
			}
			break
		case StringStringMapArray:
			for _, item := range output.StringMapArray {
				for k, v := range item {
					if k == "" || v == "" {
						continue
					}
					service.TagAssigns[k] = v
				}
			}
			break
			// TODO: log warnings about a JQ filter that went unused because it returned an invalid type that we dont know how to handle
		}
	}

	service.TagCreates = map[string]string{}
	for _, tag := range parser.TagCreates {
		output := tag.Parse(data)
		if output == nil {
			continue
		}
		switch output.Type {
		case StringStringMap:
			for k, v := range output.StringMap {
				if k == "" || v == "" {
					continue
				}
				service.TagCreates[k] = v
			}
			break
		case StringStringMapArray:
			for _, item := range output.StringMapArray {
				for k, v := range item {
					if k == "" || v == "" {
						continue
					}
					service.TagCreates[k] = v
				}
			}
			break
			// TODO: log warnings about a JQ filter that went unused because it returned an invalid type that we dont know how to handle
		}
	}

	// https://github.com/OpsLevel/kubectl-opslevel/issues/41
	service.TagAssigns = removeOverlappedKeys(service.TagAssigns, service.TagCreates)

	service.Tools = []opslevel.ToolCreateInput{}
	for _, tool := range parser.Tools {
		output := tool.Parse(data)
		if output == nil {
			continue
		}
		switch output.Type {
		case StringStringMap:
			if input, err := ConvertToToolCreateInput(output.StringMap); err == nil {
				service.Tools = append(service.Tools, *input)
			}
			break
		case StringStringMapArray:
			for _, item := range output.StringMapArray {
				if input, err := ConvertToToolCreateInput(item); err == nil {
					service.Tools = append(service.Tools, *input)
				}
			}
			break
		}
	}

	service.Repositories = []opslevel.ServiceRepositoryCreateInput{}
	for _, repository := range parser.Repositories {
		output := repository.Parse(data)
		if output == nil {
			continue
		}
		switch output.Type {
		case String:
			if input := ConvertToServiceRepositoryCreateInput(map[string]string{"repo": output.StringObj}); input != nil {
				service.Repositories = append(service.Repositories, *input)
			}
			break
		case StringArray:
			for _, item := range output.StringArray {
				if input := ConvertToServiceRepositoryCreateInput(map[string]string{"repo": item}); input != nil {
					service.Repositories = append(service.Repositories, *input)
				}
			}
			break
		case StringStringMap:
			if input := ConvertToServiceRepositoryCreateInput(output.StringMap); input != nil {
				service.Repositories = append(service.Repositories, *input)
			}
			break
		case StringStringMapArray:
			for _, item := range output.StringMapArray {
				if input := ConvertToServiceRepositoryCreateInput(item); input != nil {
					service.Repositories = append(service.Repositories, *input)
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

func removeOverlappedKeys(source map[string]string, check map[string]string) map[string]string {
	output := make(map[string]string, len(source))
	for k := range source {
		if _, value := check[k]; !value {
			output[k] = source[k]
		}
	}
	return output
}

func ConvertToToolCreateInput(data map[string]string) (*opslevel.ToolCreateInput, error) {
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

func ConvertToServiceRepositoryCreateInput(data map[string]string) *opslevel.ServiceRepositoryCreateInput {
	var repoAlias string
	baseDirectory := ""
	displayName := ""
	if val, ok := data["repo"]; ok {
		repoAlias = val
	} else {
		return nil
	}
	if val, ok := data["directory"]; ok && val != "" {
		baseDirectory = val
	}
	if val, ok := data["name"]; ok {
		displayName = val
	}
	return &opslevel.ServiceRepositoryCreateInput{
		Repository:    *opslevel.NewIdFromAlias(repoAlias),
		BaseDirectory: baseDirectory,
		DisplayName:   displayName,
	}
}

func QueryForServices(c *config.Config) ([]ServiceRegistration, error) {
	var parser *ServiceRegistrationParser
	var services []ServiceRegistration
	k8sClient := k8sutils.CreateKubernetesClient()

	jq.ValidateInstalled()
	namespaces, namespacesErr := k8sClient.GetAllNamespaces()
	if namespacesErr != nil {
		return services, nil
	}

	for _, importConfig := range c.Service.Import {
		selector := importConfig.SelectorConfig
		parser = NewParser(importConfig.OpslevelConfig)
		processFoundResource := func(resource []byte) error {
			services = append(services, *parser.Parse(resource))
			return nil
		}
		k8sClient.Query(selector, namespaces, processFoundResource)
	}
	return services, nil
}
