package common

import (
	opslevel_jq_parser "github.com/opslevel/opslevel-jq-parser/v2023"
	opslevel_k8s_controller "github.com/opslevel/opslevel-k8s-controller/v2023"
	_ "github.com/rs/zerolog/log"
)

/*
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
	s.Aliases = append(s.Aliases, o.Aliases...)
	s.Aliases = removeDuplicates(s.Aliases)
	s.TagAssigns = append(s.TagAssigns, removeOverlappedKeys(s.TagAssigns, o.TagAssigns)...)
	s.TagCreates = append(s.TagCreates, o.TagCreates...)
	s.TagAssigns = removeDuplicatesFromTagInputList(s.TagAssigns)
	s.TagAssigns = removeOverlappedKeys(s.TagAssigns, s.TagCreates)
	s.Tools = append(s.Tools, o.Tools...)
	s.Repositories = append(s.Repositories, o.Repositories...)
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
	var unique []opslevel.TagInput
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
	var list []string
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

func FilterResources(selector opslevel_common.K8SSelector, resources [][]byte) [][]byte {
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

func getServices(c *Config) ([]ServiceRegistration, error) {
	var services []ServiceRegistration
	k8sClient, err := opslevel_common.NewK8SClient()
	if err != nil {
		return services, err
	}
	for i, importConfig := range c.Service.Import {
		selector := importConfig.SelectorConfig

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

func GetAllServices(c *Config) ([]ServiceRegistration, error) {
	services, err := getServices(c)
	if err != nil {
		return nil, err
	}
	return dedupServices(services)
}

func ProcessResources(field string, config Import, resources [][]byte) ([]ServiceRegistration, error) {
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
*/

func FilterResources(selector opslevel_k8s_controller.K8SSelector, resources [][]byte) [][]byte {
	return resources
}

func ProcessResources(field string, config Import, resources [][]byte) ([]opslevel_jq_parser.ServiceRegistration, error) {
	return nil, nil
}
