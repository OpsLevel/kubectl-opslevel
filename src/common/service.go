package common

import (
	"encoding/json"
	"fmt"

	"github.com/opslevel/kubectl-opslevel/config"
	"github.com/opslevel/kubectl-opslevel/k8sutils"
	"github.com/opslevel/opslevel-go/v2023"

	_ "github.com/rs/zerolog/log"
)

type SelectorParser struct {
	Excludes []JQParser
}

type ServiceRegistration struct {
	Name         string                                  `json:",omitempty"`
	Description  string                                  `json:",omitempty"`
	Owner        string                                  `json:",omitempty"`
	Lifecycle    string                                  `json:",omitempty"`
	Tier         string                                  `json:",omitempty"`
	Product      string                                  `json:",omitempty"`
	Language     string                                  `json:",omitempty"`
	Framework    string                                  `json:",omitempty"`
	System       string                                  `json:",omitempty"`
	Aliases      []string                                `json:",omitempty"`
	TagAssigns   []opslevel.TagInput                     `json:",omitempty"`
	TagCreates   []opslevel.TagInput                     `json:",omitempty"`
	Tools        []opslevel.ToolCreateInput              `json:",omitempty"` // This is a concrete class so fields are validated during `service preview`
	Repositories []opslevel.ServiceRepositoryCreateInput `json:",omitempty"` // This is a concrete class so fields are validated during `service preview`
}

func (s *ServiceRegistration) toPrettyJson() string {
	prettyJSON, _ := json.MarshalIndent(s, "", "    ")
	return string(prettyJSON)
}

func (s *ServiceRegistration) mergeData(o ServiceRegistration) {
	if s.Name == "" {
		s.Name = o.Name
	}
	if s.Description == "" {
		s.Description = o.Description
	}
	if s.Owner == "" {
		s.Owner = o.Owner
	}
	if s.Lifecycle == "" {
		s.Lifecycle = o.Lifecycle
	}
	if s.Tier == "" {
		s.Tier = o.Tier
	}
	if s.Product == "" {
		s.Product = o.Product
	}
	if s.Language == "" {
		s.Language = o.Language
	}
	if s.Framework == "" {
		s.Framework = o.Framework
	}
	if s.System == "" {
		s.System = o.System
	}
	s.Aliases = append(s.Aliases, o.Aliases...)

	s.Aliases = removeDuplicates(s.Aliases)
	s.TagAssigns = append(s.TagAssigns, removeOverlappedKeys(s.TagAssigns, o.TagAssigns)...)
	s.TagCreates = append(s.TagCreates, o.TagCreates...)
	s.TagAssigns = removeDuplicatesFromTagInputList(s.TagAssigns)
	s.TagAssigns = removeOverlappedKeys(s.TagAssigns, s.TagCreates)
	s.Tools = append(s.Tools, o.Tools...)
	s.Repositories = append(s.Repositories, o.Repositories...)
}

func parseField(field string, filter string, resources []byte) *JQResponseMulti {
	parser := NewJQParserMulti(filter)
	return parser.ParseMulti(field, resources)
}

func parseFieldArray(field string, filters []string, resources []byte) []*JQResponseMulti {
	var output []*JQResponseMulti
	for i, filter := range filters {
		output = append(output, parseField(fmt.Sprintf("%s[%d]", field, i+1), filter, resources))
	}
	return output
}

func contains(item opslevel.TagInput, data []opslevel.TagInput) bool {
	for _, v := range data {
		if item.Key == v.Key && item.Value == v.Value {
			return true
		}
	}
	return false
}

func removeDuplicatesFromTagInputList(data []opslevel.TagInput) []opslevel.TagInput {
	unique := []opslevel.TagInput{}
	for _, entry := range data {
		if !contains(entry, unique) {
			unique = append(unique, entry)
		}
	}
	return unique
}

// Also removes empty string values
func removeDuplicates(data []string) []string {
	keys := make(map[string]bool)
	list := []string{}

	for _, entry := range data {
		if entry == "" {
			continue
		}
		if _, value := keys[entry]; !value {
			keys[entry] = true
			list = append(list, entry)
		}
	}
	return list
}

func removeDuplicatesTags(data []opslevel.TagInput) (output []opslevel.TagInput) {
	keys := make(map[string]bool)

	for _, entry := range data {
		if entry.Key == "" {
			continue
		}
		if _, value := keys[entry.Key]; !value {
			keys[entry.Key] = true
			output = append(output, entry)
		}
	}
	return
}

// https://github.com/OpsLevel/kubectl-opslevel/issues/41
func removeOverlappedKeys(source []opslevel.TagInput, check []opslevel.TagInput) (output []opslevel.TagInput) {
	for _, tagAssign := range source {
		foundMatch := false
		for _, tagCreate := range check {
			if tagCreate.Key == tagAssign.Key {
				foundMatch = true
				break
			}
		}
		if !foundMatch {
			output = append(output, tagAssign)
		}
	}
	return
}

func convertToToolCreateInput(data map[string]string) (*opslevel.ToolCreateInput, error) {
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

func convertToServiceRepositoryCreateInput(data map[string]string) *opslevel.ServiceRepositoryCreateInput {
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
		Repository:    *opslevel.NewIdentifier(repoAlias),
		BaseDirectory: baseDirectory,
		DisplayName:   displayName,
	}
}

func getString(index int, data *JQResponseMulti) string {
	if index < len(data.Objects) {
		return data.Objects[index].StringObj
	}
	return ""
}

func getAliases(index int, data []*JQResponseMulti) []string {
	output := []string{}
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
	return removeDuplicates(output)
}

func getTags(index int, data []*JQResponseMulti) (output []opslevel.TagInput) {
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
				output = append(output, opslevel.TagInput{
					Key:   k,
					Value: v,
				})
			}
		case StringStringMapArray:
			for _, item := range parsedData.StringMapArray {
				for k, v := range item {
					if k == "" || v == "" {
						continue
					}
					output = append(output, opslevel.TagInput{
						Key:   k,
						Value: v,
					})
				}
			}
			// TODO: log warnings about a JQ filter that went unused because it returned an invalid type that we dont know how to handle
		}
	}
	return output
}

func getTools(index int, data []*JQResponseMulti) []opslevel.ToolCreateInput {
	output := []opslevel.ToolCreateInput{}
	count := len(data)
	for i := 0; i < count; i++ {
		if data[i].Objects == nil {
			continue
		}
		parsedData := data[i].Objects[index]
		switch parsedData.Type {
		case StringStringMap:
			if parsedData.StringMap == nil {
				continue
			}
			if input, err := convertToToolCreateInput(parsedData.StringMap); err == nil {
				output = append(output, *input)
			}
		case StringStringMapArray:
			for _, item := range parsedData.StringMapArray {
				if item == nil {
					continue
				}
				if input, err := convertToToolCreateInput(item); err == nil {
					output = append(output, *input)
				}
			}
		}
	}
	return output
}

func getRepositories(index int, data []*JQResponseMulti) []opslevel.ServiceRepositoryCreateInput {
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
			if input := convertToServiceRepositoryCreateInput(map[string]string{"repo": parsedData.StringObj}); input != nil {
				output = append(output, *input)
			}
		case StringArray:
			for _, item := range parsedData.StringArray {
				if item == "" {
					continue
				}
				if input := convertToServiceRepositoryCreateInput(map[string]string{"repo": item}); input != nil {
					output = append(output, *input)
				}
			}
		case StringStringMap:
			if parsedData.StringMap == nil {
				continue
			}
			if input := convertToServiceRepositoryCreateInput(parsedData.StringMap); input != nil {
				output = append(output, *input)
			}
		case StringStringMapArray:
			for _, item := range parsedData.StringMapArray {
				if item == nil {
					continue
				}
				if input := convertToServiceRepositoryCreateInput(item); input != nil {
					output = append(output, *input)
				}
			}
		}
	}
	return output
}

var (
	StartArray []byte = []byte(`[`)
	EndArray   []byte = []byte(`]`)
	JoinItem   []byte = []byte(`,`)
)

func joinResourceBytes(resources [][]byte) []byte {
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

func FilterResources(selector k8sutils.KubernetesSelector, resources [][]byte) [][]byte {
	var output [][]byte
	resourceCount := len(resources)
	// Parse
	filterResults := parseFieldArray("selector.excludes", selector.Excludes, joinResourceBytes(resources))

	// Aggregate
	for resourceIndex := 0; resourceIndex < resourceCount; resourceIndex++ {
		if anyIsTrue(resourceIndex, filterResults) {
			continue
		}
		output = append(output, resources[resourceIndex])
	}
	return output
}

func aliasOverlaps(a []string, b []string) bool {
	for _, i := range a {
		for _, j := range b {
			if i == j {
				return true
			}
		}
	}
	return false
}

// TODO: bubble up errors better
func parseResources(field string, c config.ServiceRegistrationConfig, count int, resources []byte) ([]ServiceRegistration, error) {
	services := make([]ServiceRegistration, count)

	// Parse
	Names := parseField(fmt.Sprintf("%s.name", field), c.Name, resources)
	Descriptions := parseField(fmt.Sprintf("%s.description", field), c.Description, resources)
	Owners := parseField(fmt.Sprintf("%s.owner", field), c.Owner, resources)
	Lifecycles := parseField(fmt.Sprintf("%s.lifecycle", field), c.Lifecycle, resources)
	Tiers := parseField(fmt.Sprintf("%s.tier", field), c.Tier, resources)
	Products := parseField(fmt.Sprintf("%s.product", field), c.Product, resources)
	Languages := parseField(fmt.Sprintf("%s.language", field), c.Language, resources)
	Frameworks := parseField(fmt.Sprintf("%s.framework", field), c.Framework, resources)
	Systems := parseField(fmt.Sprintf("%s.system", field), c.System, resources)
	Aliases := parseFieldArray(fmt.Sprintf("%s.aliases", field), c.Aliases, resources)
	if len(Aliases) < 1 {
		Aliases = append(Aliases, parseField("Auto Added Alias", "\"k8s:\\(.metadata.name)-\\(.metadata.namespace)\"", resources))
	}
	TagAssigns := parseFieldArray(fmt.Sprintf("%s.tags.assign", field), c.Tags.Assign, resources)
	TagCreates := parseFieldArray(fmt.Sprintf("%s.tags.create", field), c.Tags.Create, resources)
	Tools := parseFieldArray(fmt.Sprintf("%s.tools", field), c.Tools, resources)
	Repositories := parseFieldArray(fmt.Sprintf("%s.repository", field), c.Repositories, resources)

	// Aggregate
	for i := 0; i < count; i++ {
		service := &services[i]

		service.Name = getString(i, Names)
		service.Description = getString(i, Descriptions)
		service.Owner = getString(i, Owners)
		service.Lifecycle = getString(i, Lifecycles)
		service.Tier = getString(i, Tiers)
		service.Product = getString(i, Products)
		service.Language = getString(i, Languages)
		service.Framework = getString(i, Frameworks)
		service.System = getString(i, Systems)
		service.Aliases = getAliases(i, Aliases)
		service.TagAssigns = getTags(i, TagAssigns)
		service.TagCreates = getTags(i, TagCreates)
		service.TagCreates = removeDuplicatesTags(service.TagCreates)
		service.TagAssigns = removeOverlappedKeys(service.TagAssigns, service.TagCreates)
		service.Tools = getTools(i, Tools)
		service.Repositories = getRepositories(i, Repositories)
	}

	return services, nil
}

func dedupServices(input []ServiceRegistration) ([]ServiceRegistration, error) {
	var output []ServiceRegistration
	for _, source := range input {
		wasMerged := false
		for i, dest := range output {
			if aliasOverlaps(source.Aliases, dest.Aliases) {
				dest.mergeData(source)
				output[i] = dest
				wasMerged = true
				break
			}
		}
		if !wasMerged {
			output = append(output, source)
		}
	}
	return output, nil
}

func getServices(c *config.Config) ([]ServiceRegistration, error) {
	var services []ServiceRegistration
	k8sClient := k8sutils.CreateKubernetesClient()
	for i, importConfig := range c.Service.Import {
		selector := importConfig.SelectorConfig
		if selectorErr := selector.Validate(); selectorErr != nil {
			return services, selectorErr
		}

		resources, queryErr := k8sClient.Query(selector)
		if queryErr != nil {
			return services, queryErr
		}

		parsedServices, parsedServicesErr := ProcessResources(fmt.Sprintf("service.import[%d]", i+1), importConfig, resources)
		if parsedServicesErr != nil {
			return services, parsedServicesErr
		}

		services = append(services, parsedServices...)
	}
	return services, nil
}

func GetAllServices(c *config.Config) ([]ServiceRegistration, error) {
	services, err := getServices(c)
	if err != nil {
		return nil, err
	}
	return dedupServices(services)
}

func ProcessResources(field string, config config.Import, resources [][]byte) ([]ServiceRegistration, error) {
	filtered := FilterResources(config.SelectorConfig, resources)
	if len(filtered) < 1 {
		return []ServiceRegistration{}, nil
	}
	parsed, parseError := parseResources(field, config.OpslevelConfig, len(filtered), joinResourceBytes(filtered))
	if parseError != nil {
		return nil, parseError
	}
	deduped, dedupErr := dedupServices(parsed)
	if dedupErr != nil {
		return nil, dedupErr
	}
	return deduped, nil
}
