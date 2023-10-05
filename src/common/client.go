package common

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/google/go-cmp/cmp"
	"github.com/opslevel/opslevel-go/v2023"
	"github.com/rs/zerolog/log"
)

func ReconcileService(client *opslevel.Client, service ServiceRegistration) {
	if len(service.Aliases) <= 0 {
		log.Warn().Msgf("[%s] found 0 aliases from kubernetes data", service.Name)
		return
	}
	log.Trace().Msgf("[%s] Parsed Data: \n%s", service.Name, service.toPrettyJson())
	foundService, foundServiceStatus := validateServiceAliases(client, service)
	switch foundServiceStatus {
	case serviceAliasesResult_NoAliasesMatched:
		newService, newServiceErr := createService(client, service)
		if newServiceErr != nil {
			log.Warn().Msgf("[%s] api error during service creation ... skipping reconciliation.\n\tREASON: %v", service.Name, newServiceErr)
			return
		}
		foundService = newService
	case serviceAliasesResult_AliasMatched:
		updateService(client, service, foundService)

	case serviceAliasesResult_MultipleServicesFound:
		log.Warn().Msgf("[%s] found multiple services with aliases = [\"%s\"].  cannot know which service to target for update ... skipping reconciliation", service.Name, strings.Join(service.Aliases, "\", \""))
		return
	case serviceAliasesResult_APIErrorHappened:
		log.Warn().Msgf("[%s] api error during service lookup by alias.  unable to guarentee service was found or not ... skipping reconciliation", service.Name)
		return
	}

	handleAliases(client, service, foundService)
	handleTags(client, service, foundService)
	handleTools(client, service, foundService)
	handleSystem(client, service, foundService)
	handleRepositories(client, service, foundService)
	log.Info().Msgf("[%s] Finished processing data", foundService.Name)
}

type serviceAliasesResult string

const (
	serviceAliasesResult_NoAliasesMatched      serviceAliasesResult = "NoAliasesMatched"
	serviceAliasesResult_AliasMatched          serviceAliasesResult = "AliasMatched"
	serviceAliasesResult_MultipleServicesFound serviceAliasesResult = "MultipleServicesFound"
	serviceAliasesResult_APIErrorHappened      serviceAliasesResult = "APIErrorHappened"
)

// This function has 4 outcomes that can happen while looping over the aliases list
// serviceAliasesResult_NoAliasesMatched - means that all API calls succeeded and none of the aliases matched an existing service
// serviceAliasesResult_AliasMatched - means that all the API calls succeeded and a single service was found matching 1 of N aliases
// serviceAliasesResult_MultipleServicesFound - means that all API calls succeeded but multiple services were returning means the list of aliases does not definitively describe a single service and might be a configuration problem
// serviceAliasesResult_APIErrorHappened - means that 1 of N aliases got an 4xx/5xx and thereforce we cannot say 100% that the services doesn't exist
func validateServiceAliases(client *opslevel.Client, registration ServiceRegistration) (*opslevel.Service, serviceAliasesResult) {
	var gotError error
	foundServices := map[string]*opslevel.Service{}
	for _, alias := range registration.Aliases {
		foundService, err := client.GetServiceWithAlias(alias)
		if err != nil {
			gotError = err
		} else {
			if foundService.Id != "" {
				foundServices[string(foundService.Id)] = foundService
			}
		}
	}
	if gotError != nil {
		return nil, serviceAliasesResult_APIErrorHappened
	}
	foundServicesCount := len(foundServices)
	if foundServicesCount > 1 {
		return nil, serviceAliasesResult_MultipleServicesFound
	}
	if foundServicesCount < 1 {
		return nil, serviceAliasesResult_NoAliasesMatched
	}
	output := []*opslevel.Service{}
	for _, value := range foundServices {
		output = append(output, value)
	}
	return output[0], serviceAliasesResult_AliasMatched
}

func serviceNeedsUpdate(input opslevel.ServiceUpdateInput, service *opslevel.Service) bool {
	if input.Name != "" && input.Name != service.Name {
		return true
	}
	if input.Product != "" && input.Product != service.Product {
		return true
	}
	if input.Description != "" && input.Description != service.Description {
		return true
	}
	if input.Language != "" && input.Language != service.Language {
		return true
	}
	if input.Framework != "" && input.Framework != service.Framework {
		return true
	}
	if input.Tier != "" && input.Tier != service.Tier.Alias {
		return true
	}
	if input.Lifecycle != "" && input.Lifecycle != service.Lifecycle.Alias {
		return true
	}
	if input.Owner != "" && input.Owner != service.Owner.Alias {
		return true
	}
	return false
}

func createService(client *opslevel.Client, registration ServiceRegistration) (*opslevel.Service, error) {
	serviceCreateInput := opslevel.ServiceCreateInput{
		Name:        registration.Name,
		Product:     registration.Product,
		Description: registration.Description,
		Language:    registration.Language,
		Framework:   registration.Framework,
	}
	if v, ok := opslevel.Cache.TryGetTier(registration.Tier); ok {
		serviceCreateInput.Tier = string(v.Alias)
	} else if registration.Tier != "" {
		log.Warn().Msgf("[%s] Unable to find 'Tier' with alias '%s'", registration.Name, registration.Tier)
	}
	if v, ok := opslevel.Cache.TryGetLifecycle(registration.Lifecycle); ok {
		serviceCreateInput.Lifecycle = string(v.Alias)
	} else if registration.Lifecycle != "" {
		log.Warn().Msgf("[%s] Unable to find 'Lifecycle' with alias '%s'", registration.Name, registration.Lifecycle)
	}
	if v, ok := opslevel.Cache.TryGetTeam(registration.Owner); ok {
		serviceCreateInput.Owner = string(v.Alias)
	} else if registration.Owner != "" {
		log.Warn().Msgf("[%s] Unable to find 'Team' with alias '%s'", registration.Name, registration.Owner)
	}
	service, err := client.CreateService(serviceCreateInput)
	if err != nil {
		log.Error().Msgf("[%s] Failed creating service\n\tREASON: %v", registration.Name, err.Error())
	} else {
		log.Info().Msgf("[%s] Created new service", service.Name)
	}
	return service, err
}

func updateService(client *opslevel.Client, registration ServiceRegistration, service *opslevel.Service) {
	updateServiceInput := opslevel.ServiceUpdateInput{
		Id:          service.Id,
		Product:     registration.Product,
		Description: registration.Description,
		Language:    registration.Language,
		Framework:   registration.Framework,
	}
	if v, ok := opslevel.Cache.TryGetTier(registration.Tier); ok {
		updateServiceInput.Tier = string(v.Alias)
	} else if registration.Tier != "" {
		log.Warn().Msgf("[%s] Unable to find 'Tier' with alias '%s'", service.Name, registration.Tier)
	}
	if v, ok := opslevel.Cache.TryGetLifecycle(registration.Lifecycle); ok {
		updateServiceInput.Lifecycle = string(v.Alias)
	} else if registration.Lifecycle != "" {
		log.Warn().Msgf("[%s] Unable to find 'Lifecycle' with alias '%s'", service.Name, registration.Lifecycle)
	}
	if v, ok := opslevel.Cache.TryGetTeam(registration.Owner); ok {
		updateServiceInput.Owner = string(v.Alias)
	} else if registration.Owner != "" {
		log.Warn().Msgf("[%s] Unable to find 'Team' with alias '%s'", service.Name, registration.Owner)
	}
	if serviceNeedsUpdate(updateServiceInput, service) {
		updatedService, updateServiceErr := client.UpdateService(updateServiceInput)
		if updateServiceErr != nil {
			log.Error().Msgf("[%s] Failed updating service\n\tREASON: %v", service.Name, updateServiceErr.Error())
		} else {
			if diff := cmp.Diff(service, updatedService); diff != "" {
				log.Info().Msgf("[%s] Updated Service - Diff:\n%s", service.Name, diff)
			}
		}
	} else {
		log.Info().Msgf("[%s] No changes detected to fields - skipping update", service.Name)
	}
}

func handleAliases(client *opslevel.Client, registration ServiceRegistration, service *opslevel.Service) {
	for _, alias := range registration.Aliases {
		if alias == "" || service.HasAlias(alias) {
			continue
		}
		_, err := client.CreateAlias(opslevel.AliasCreateInput{
			Alias:   alias,
			OwnerId: service.Id,
		})
		if err != nil {
			log.Error().Msgf("[%s] Failed assigning alias '%s'\n\tREASON: %v", service.Name, alias, err.Error())
		} else {
			log.Info().Msgf("[%s] Assigned alias '%s'", service.Name, alias)
		}
	}
}

func handleTags(client *opslevel.Client, registration ServiceRegistration, service *opslevel.Service) {
	assignTags(client, registration, service)
	createTags(client, registration, service)
}

func containsAllTags(tagAssigns []opslevel.TagInput, serviceTags []opslevel.Tag) bool {
	found := map[int]bool{}
	for i, expected := range tagAssigns {
		found[i] = false
		for _, match := range serviceTags {
			if expected.Key == match.Key && expected.Value == match.Value {
				found[i] = true
				break
			}
		}
	}
	for _, value := range found {
		if !value {
			return false
		}
	}
	return true
}

func assignTags(client *opslevel.Client, registration ServiceRegistration, service *opslevel.Service) {
	if registration.TagAssigns == nil {
		return
	}
	if !containsAllTags(registration.TagAssigns, service.Tags.Nodes) {
		tags := map[string]string{}
		for _, tagAssign := range registration.TagAssigns {
			tags[tagAssign.Key] = tagAssign.Value
		}

		_, err := client.AssignTags(string(service.Id), tags)
		jsonBytes, _ := json.Marshal(registration.TagAssigns)
		if err != nil {
			log.Error().Msgf("[%s] Failed assigning tags: %s\n\tREASON: %v", service.Name, string(jsonBytes), err.Error())
		} else {
			log.Info().Msgf("[%s] Assigned tags: %s", service.Name, string(jsonBytes))
		}
	} else {
		log.Info().Msgf("[%s] All tags already assigned to service.", service.Name)
	}
}

func createTags(client *opslevel.Client, registration ServiceRegistration, service *opslevel.Service) {
	for _, tag := range registration.TagCreates {
		if service.HasTag(tag.Key, tag.Value) {
			continue
		}
		input := opslevel.TagCreateInput{
			Id:    service.Id,
			Key:   tag.Key,
			Value: tag.Value,
		}
		_, err := client.CreateTag(input)
		if err != nil {
			log.Error().Msgf("[%s] Failed creating tag '%s = %s'\n\tREASON: %v", service.Name, tag.Key, tag.Value, err.Error())
		} else {
			log.Info().Msgf("[%s] Created tag '%s = %s'", service.Name, tag.Key, tag.Value)
		}
	}
}

func handleTools(client *opslevel.Client, registration ServiceRegistration, service *opslevel.Service) {
	for _, tool := range registration.Tools {
		if service.HasTool(tool.Category, tool.DisplayName, tool.Environment) {
			log.Debug().Msgf("[%s] Tool '{Category: %s, Environment: %s, Name: %s}' already exists on service ... skipping", service.Name, tool.Category, tool.Environment, tool.DisplayName)
			continue
		}
		tool.ServiceId = service.Id
		_, err := client.CreateTool(tool)
		if err != nil {
			log.Error().Msgf("[%s] Failed assigning tool '{Category: %s, Environment: %s, Name: %s}'\n\tREASON: %v", service.Name, tool.Category, tool.Environment, tool.DisplayName, err.Error())
		} else {
			log.Info().Msgf("[%s] Ensured tool '{Category: %s, Environment: %s, Name: %s}'", service.Name, tool.Category, tool.Environment, tool.DisplayName)
		}
	}
}

func handleRepositories(client *opslevel.Client, registration ServiceRegistration, service *opslevel.Service) {
	for _, repositoryCreate := range registration.Repositories {
		repositoryAsString := fmt.Sprintf("{Alias: %s, Directory: %s, Name: %s}", repositoryCreate.Repository.Alias, repositoryCreate.BaseDirectory, repositoryCreate.DisplayName)
		foundRepository, foundRepositoryErr := client.GetRepositoryWithAlias(string(repositoryCreate.Repository.Alias))
		if foundRepositoryErr != nil {
			log.Warn().Msgf("[%s] Repository with alias: '%s' not found so it cannot be attached to service ... skipping", service.Name, repositoryAsString)
			continue
		}
		serviceRepository := foundRepository.GetService(service.Id, repositoryCreate.BaseDirectory)
		if serviceRepository != nil {
			if repositoryCreate.DisplayName != "" && serviceRepository.DisplayName != repositoryCreate.DisplayName {
				repositoryUpdate := opslevel.ServiceRepositoryUpdateInput{
					Id:          serviceRepository.Id,
					DisplayName: repositoryCreate.DisplayName,
				}
				_, err := client.UpdateServiceRepository(repositoryUpdate)
				if err != nil {
					log.Error().Msgf("[%s] Failed updating repository '%s'\n\tREASON: %v", service.Name, repositoryAsString, err.Error())
					continue
				} else {
					log.Info().Msgf("[%s] Updated repository '%s'", service.Name, repositoryAsString)
					continue
				}
			}
			log.Debug().Msgf("[%s] Repository '%s' already attached to service ... skipping", service.Name, repositoryAsString)
			continue
		}
		repositoryCreate.Service = opslevel.IdentifierInput{Id: service.Id}
		_, err := client.CreateServiceRepository(repositoryCreate)
		if err != nil {
			log.Error().Msgf("[%s] Failed assigning repository '%s'\n\tREASON: %v", service.Name, repositoryAsString, err.Error())
		} else {
			log.Info().Msgf("[%s] Attached repository '%s'", service.Name, repositoryAsString)
		}
	}
}

func handleSystem(client *opslevel.Client, registration ServiceRegistration, service *opslevel.Service) {
	if registration.System != "" {
		system, foundSystemErr := client.GetSystem(registration.System)
		if foundSystemErr != nil {
			log.Warn().Msgf("[%s] System with alias: '%s' not found so service cannot be attached to it ... skipping", service.Name, system)
			return
		}

		assignServiceErr := system.AssignService(client, service.Aliases[0])
		if assignServiceErr != nil {
			log.Error().Msgf("[%s] Failed assigning service to system '%s'\n\tREASON: %v", service.Name, system, assignServiceErr.Error())
			return
		}

		log.Info().Msgf("[%s] Attached service '%s'", service.Name, system)
	}
}
