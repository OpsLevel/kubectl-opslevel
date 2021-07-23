package common

import (
	"encoding/json"

	"github.com/opslevel/kubectl-opslevel/config"
	"github.com/opslevel/kubectl-opslevel/jq"
	"github.com/opslevel/kubectl-opslevel/k8sutils"
	"github.com/opslevel/opslevel-go"

	_ "github.com/rs/zerolog/log"
)

type SelectorParser struct {
	Excludes []JQParser
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

func parseField(filter string, resources []byte) *JQResponseMulti {
	parser := NewJQParserMulti(filter)
	return parser.ParseMulti(resources)
}

func parseFieldArray(filters []string, resources []byte) []*JQResponseMulti {
	var output []*JQResponseMulti
	for _, filter := range filters {
		output = append(output, parseField(filter, resources))
	}
	return output
}

func aggregateAliases(index int, data []*JQResponseMulti) []string {
	output := []string{}
	count := len(data)
	for i := 0; i < count; i++ {
		if data[i].Objects == nil {
			continue
		}
		parsedData := data[i].Objects[index]
		switch parsedData.Type {
		case String:
			output = append(output, parsedData.StringObj)
		case StringArray:
			for _, item := range parsedData.StringArray {
				if item == "" {
					continue
				}
				output = append(output, item)
			}
			// TODO: log warnings about a JQ filter that went unused because it returned an invalid type that we dont know how to handle
		}
	}
	return output
}

func aggregateMap(index int, data []*JQResponseMulti) map[string]string {
	output := map[string]string{}
	count := len(data)
	for i := 0; i < count; i++ {
		if data[i].Objects == nil {
			continue
		}
		parsedData := data[i].Objects[index]
		switch parsedData.Type {
		case StringStringMap:
			for k, v := range parsedData.StringMap {
				if k == "" || v == "" {
					continue
				}
				output[k] = v
			}
		case StringStringMapArray:
			for _, item := range parsedData.StringMapArray {
				for k, v := range item {
					if k == "" || v == "" {
						continue
					}
					output[k] = v
				}
			}
			// TODO: log warnings about a JQ filter that went unused because it returned an invalid type that we dont know how to handle
		}
	}
	return output
}

func aggregateTools(index int, data []*JQResponseMulti) []opslevel.ToolCreateInput {
	output := []opslevel.ToolCreateInput{}
	count := len(data)
	for i := 0; i < count; i++ {
		if data[i].Objects == nil {
			continue
		}
		parsedData := data[i].Objects[index]
		switch parsedData.Type {
		case StringStringMap:
			if input, err := ConvertToToolCreateInput(parsedData.StringMap); err == nil {
				output = append(output, *input)
			}
		case StringStringMapArray:
			for _, item := range parsedData.StringMapArray {
				if input, err := ConvertToToolCreateInput(item); err == nil {
					output = append(output, *input)
				}
			}
		}
	}
	return output
}

func aggregateRepositories(index int, data []*JQResponseMulti) []opslevel.ServiceRepositoryCreateInput {
	output := []opslevel.ServiceRepositoryCreateInput{}
	count := len(data)
	for i := 0; i < count; i++ {
		if data[i].Objects == nil {
			continue
		}
		parsedData := data[i].Objects[index]
		switch parsedData.Type {
		case String:
			if parsedData.StringObj == "" {
				continue
			}
			if input := ConvertToServiceRepositoryCreateInput(map[string]string{"repo": parsedData.StringObj}); input != nil {
				output = append(output, *input)
			}
		case StringArray:
			for _, item := range parsedData.StringArray {
				if item == "" {
					continue
				}
				if input := ConvertToServiceRepositoryCreateInput(map[string]string{"repo": item}); input != nil {
					output = append(output, *input)
				}
			}
		case StringStringMap:
			if input := ConvertToServiceRepositoryCreateInput(parsedData.StringMap); input != nil {
				output = append(output, *input)
			}
		case StringStringMapArray:
			for _, item := range parsedData.StringMapArray {
				if input := ConvertToServiceRepositoryCreateInput(item); input != nil {
					output = append(output, *input)
				}
			}
		}
	}
	return output
}

func GetString(index int, data *JQResponseMulti) string {
	if index < len(data.Objects) {
		return data.Objects[index].StringObj
	}
	return ""
}

// TODO: bubble up errors
func Parse(c config.ServiceRegistrationConfig, count int, resources []byte) ([]ServiceRegistration, error) {
	services := make([]ServiceRegistration, count)

	// Parse
	Names := parseField(c.Name, resources)
	Descriptions := parseField(c.Description, resources)
	Owners := parseField(c.Owner, resources)
	Lifecycles := parseField(c.Lifecycle, resources)
	Tiers := parseField(c.Tier, resources)
	Products := parseField(c.Product, resources)
	Languages := parseField(c.Language, resources)
	Frameworks := parseField(c.Framework, resources)
	Aliases := parseFieldArray(c.Aliases, resources)
	Aliases = append(Aliases, parseField("\"k8s:\\(.metadata.name)-\\(.metadata.namespace)\"", resources))
	TagAssigns := parseFieldArray(c.Tags.Assign, resources)
	TagCreates := parseFieldArray(c.Tags.Create, resources)
	Tools := parseFieldArray(c.Tools, resources)
	Repositories := parseFieldArray(c.Repositories, resources)

	// Aggregate
	for i := 0; i < count; i++ {
		service := &services[i]

		service.Name = GetString(i, Names)
		service.Description = GetString(i, Descriptions)
		service.Owner = GetString(i, Owners)
		service.Lifecycle = GetString(i, Lifecycles)
		service.Tier = GetString(i, Tiers)
		service.Product = GetString(i, Products)
		service.Language = GetString(i, Languages)
		service.Framework = GetString(i, Frameworks)
		service.Aliases = aggregateAliases(i, Aliases)
		service.Aliases = removeDuplicates(service.Aliases)
		service.TagAssigns = aggregateMap(i, TagAssigns)
		service.TagCreates = aggregateMap(i, TagCreates)
		// https://github.com/OpsLevel/kubectl-opslevel/issues/41
		service.TagAssigns = removeOverlappedKeys(service.TagAssigns, service.TagCreates)
		service.Tools = aggregateTools(i, Tools)
		service.Repositories = aggregateRepositories(i, Repositories)
	}

	return services, nil
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
	var services []ServiceRegistration
	k8sClient := k8sutils.CreateKubernetesClient()

	jq.ValidateInstalled()
	namespaces, namespacesErr := k8sClient.GetAllNamespaces()
	if namespacesErr != nil {
		return services, nil
	}
	for _, importConfig := range c.Service.Import {
		selector := importConfig.SelectorConfig
		if selectorErr := selector.Validate(); selectorErr != nil {
			return services, selectorErr
		}

		resources, queryErr := k8sClient.Query(selector, namespaces)
		if queryErr != nil {
			return services, queryErr
		}

		resources = filterResources(selector, resources)
		parsedServices, parseError := Parse(importConfig.OpslevelConfig, len(resources), joinResources(resources))
		if parseError != nil {
			return services, parseError
		}
		services = append(services, parsedServices...)
	}
	return services, nil
}

func anyIsTrue(resourceIndex int, filters []*JQResponseMulti) bool {
	filtersCount := len(filters)
	for filterIndex := 0; filterIndex < filtersCount; filterIndex++ {
		results := filters[filterIndex].Objects
		if results == nil {
			return false
		}
		parsedData := results[resourceIndex]
		switch parsedData.Type {
		case Bool:
			if parsedData.BoolObj {
				return true
			}
		case BoolArray:
			for _, value := range parsedData.BoolArray {
				if value {
					return true
				}
			}
		}
	}
	return false

}

func filterResources(selector k8sutils.KubernetesSelector, resources [][]byte) [][]byte {
	var output [][]byte
	resourceCount := len(resources)
	// Parse
	filterResults := parseFieldArray(selector.Excludes, joinResources(resources))

	// Aggregate
	for resourceIndex := 0; resourceIndex < resourceCount; resourceIndex++ {
		if anyIsTrue(resourceIndex, filterResults) {
			continue
		}
		output = append(output, resources[resourceIndex])
	}
	return output
}

var StartArray []byte = []byte(`[`)
var EndArray []byte = []byte(`]`)
var JoinItem []byte = []byte(`,`)

func joinResources(resources [][]byte) []byte {
	var output []byte
	output = append(output, StartArray...)
	count := len(resources) - 1
	for i, item := range resources {
		output = append(output, item...)
		if i < count {
			output = append(output, JoinItem...)
		}
	}
	output = append(output, EndArray...)
	return output
}
