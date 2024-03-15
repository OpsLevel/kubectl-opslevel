package common

import (
	"encoding/json"
	"fmt"

	"github.com/google/go-cmp/cmp"

	"github.com/opslevel/opslevel-go/v2024"
	opslevel_jq_parser "github.com/opslevel/opslevel-jq-parser/v2024"
	"github.com/rs/zerolog/log"
)

type serviceAliasesResult string

// TODO: we could be better served using error types and wrapping them
const (
	serviceAliasesResult_NoAliasesMatched      serviceAliasesResult = "NoAliasesMatched"
	serviceAliasesResult_AliasMatched          serviceAliasesResult = "AliasMatched"
	serviceAliasesResult_MultipleServicesFound serviceAliasesResult = "MultipleServicesFound"
	serviceAliasesResult_APIErrorHappened      serviceAliasesResult = "APIErrorHappened"
)

type ServiceReconciler struct {
	client                 *OpslevelClient
	disableServiceCreation bool
}

func NewServiceReconciler(client *OpslevelClient, disableServiceCreation bool) *ServiceReconciler {
	return &ServiceReconciler{
		client:                 client,
		disableServiceCreation: disableServiceCreation,
	}
}

// Reconcile looks up services matching the aliases provided in the registration. If it does not find one, it will create one.
// Reconcile is push-only, meaning it will never remove data contained in the service but not defined in the registration.
func (r *ServiceReconciler) Reconcile(registration opslevel_jq_parser.ServiceRegistration) error {
	if registration.Name == "" {
		return fmt.Errorf("cannot reconcile service with no name")
	}
	if len(registration.Aliases) == 0 {
		return fmt.Errorf("[%s] found 0 aliases from kubernetes data", registration.Name)
	}
	service, status := r.lookupService(registration)
	switch status {
	case serviceAliasesResult_APIErrorHappened:
		return fmt.Errorf("[%s] api error during service lookup by alias.  unable to guarantee service was found or not ... skipping reconciliation", registration.Name)
	case serviceAliasesResult_MultipleServicesFound:
		return fmt.Errorf("[%s] found multiple services with aliases = [%s]. cannot know which service to target for update ... skipping reconciliation", registration.Name, registration.Aliases)
	case serviceAliasesResult_AliasMatched:
		if service == nil {
			return fmt.Errorf("[%s] unexpected nil before update - submit a bug report ... skipping reconciliation", registration.Name)
		}
		err := r.updateService(service, registration)
		if err != nil {
			return err
		}
	default:
		// happy path
		if r.disableServiceCreation {
			log.Info().Msgf("avoiding creating a new service.  service creation is disabled")
			return nil
		}
		var err error
		service, err = r.createService(registration)
		if err != nil {
			return err
		}
		if service == nil {
			return fmt.Errorf("[%s] unexpected nil after create - submit a bug report ... skipping reconciliation", registration.Name)
		}
	}

	// We don't care about errors at this point because they will just be logged
	r.handleAliases(service, registration)
	r.handleAssignTags(service, registration)
	r.handleCreateTags(service, registration)
	r.handleTools(service, registration)
	r.handleRepositories(service, registration)
	r.handleProperties(service, registration)
	return nil
}

func (r *ServiceReconciler) ContainsAllTags(tagAssigns []opslevel.TagInput, serviceTags []opslevel.Tag) bool {
	if len(tagAssigns) > len(serviceTags) {
		return false
	}
	serviceTagsMap := make(map[string]bool)
	for _, tag := range serviceTags {
		serviceTagsMap[tag.Key+tag.Value] = true
	}
	for _, tag := range tagAssigns {
		if _, ok := serviceTagsMap[tag.Key+tag.Value]; !ok {
			return false
		}
	}
	return true
}

// lookupService has 4 outcomes that can happen while looping over the aliases list
// serviceAliasesResult_NoAliasesMatched - means that all API calls succeeded and none of the aliases matched an existing service
// serviceAliasesResult_AliasMatched - means that all the API calls succeeded and a single service was found matching 1 of N aliases
// serviceAliasesResult_MultipleServicesFound - means that all API calls succeeded but multiple services were returning means the list of aliases does not definitively describe a single service and might be a configuration problem
// serviceAliasesResult_APIErrorHappened - means that 1 of N aliases got a 4xx/5xx and thereforce we cannot say 100% that the services doesn't exist
func (r *ServiceReconciler) lookupService(registration opslevel_jq_parser.ServiceRegistration) (*opslevel.Service, serviceAliasesResult) {
	var foundService *opslevel.Service
	for _, alias := range registration.Aliases {
		gotService, err := r.client.GetService(alias)
		if err != nil {
			log.Warn().Err(err).Msgf("got an error when trying to get service with alias '%s'", alias)
			return nil, serviceAliasesResult_APIErrorHappened
		} else if gotService == nil {
			log.Debug().Msgf("did not find a service with alias '%s'", alias)
			continue
		} else if foundService != nil {
			log.Debug().Msgf("found another service with the same alias '%s' (%s)", alias, gotService.Id)
			return nil, serviceAliasesResult_MultipleServicesFound
		}
		// happy path
		foundService = gotService
	}
	if foundService == nil {
		return nil, serviceAliasesResult_NoAliasesMatched
	}
	return foundService, serviceAliasesResult_AliasMatched // happy path
}

func (r *ServiceReconciler) createService(registration opslevel_jq_parser.ServiceRegistration) (*opslevel.Service, error) {
	serviceInput := opslevel.ServiceCreateInput{
		Name: registration.Name,
	}
	if registration.Description != "" {
		serviceInput.Description = opslevel.RefOf(registration.Description)
	}
	if registration.Framework != "" {
		serviceInput.Framework = opslevel.RefOf(registration.Framework)
	}
	if registration.Language != "" {
		serviceInput.Language = opslevel.RefOf(registration.Language)
	}
	if registration.Lifecycle != "" {
		serviceInput.LifecycleAlias = opslevel.RefOf(registration.Lifecycle)
	}
	if registration.Owner != "" {
		serviceInput.OwnerInput = opslevel.NewIdentifier(registration.Owner)
	}
	if registration.System != "" {
		serviceInput.Parent = opslevel.NewIdentifier(registration.System)
	}
	if registration.Product != "" {
		serviceInput.Product = opslevel.RefOf(registration.Product)
	}
	if registration.Tier != "" {
		serviceInput.TierAlias = opslevel.RefOf(registration.Tier)
	}
	created, err := r.client.CreateService(serviceInput)
	if err != nil {
		log.Error().Msgf("[%s] API error: '%s'", registration.Name, err.Error())
		return nil, err
	}
	toJSON, _ := json.Marshal(&created)
	log.Info().Msgf("[%s] Created Service:\n%s", registration.Name, string(toJSON))
	return created, nil
}

func (r *ServiceReconciler) updateService(service *opslevel.Service, registration opslevel_jq_parser.ServiceRegistration) error {
	serviceInput := opslevel.ServiceUpdateInput{
		Id: &service.Id,
	}
	if registration.Description != "" {
		serviceInput.Description = opslevel.RefOf(registration.Description)
	}
	if registration.Framework != "" {
		serviceInput.Framework = opslevel.RefOf(registration.Framework)
	}
	if registration.Language != "" {
		serviceInput.Language = opslevel.RefOf(registration.Language)
	}
	if registration.Lifecycle != "" {
		serviceInput.LifecycleAlias = opslevel.RefOf(registration.Lifecycle)
	}
	if registration.Name != "" {
		serviceInput.Name = opslevel.RefOf(registration.Name)
	}
	if registration.Owner != "" {
		serviceInput.OwnerInput = opslevel.NewIdentifier(registration.Owner)
	}
	if registration.System != "" {
		serviceInput.Parent = opslevel.NewIdentifier(registration.System)
	}
	if registration.Product != "" {
		serviceInput.Product = opslevel.RefOf(registration.Product)
	}
	if registration.Tier != "" {
		serviceInput.TierAlias = opslevel.RefOf(registration.Tier)
	}
	inputDiff := cmp.Diff(serviceInput, opslevel.ServiceUpdateInput{Id: &service.Id})
	if inputDiff == "" {
		log.Info().Msgf("[%s] No changes - skipping", service.Name)
		return nil
	}
	log.Info().Msgf("[%s] Needs changes - Diff:\n%s", service.Name, inputDiff)
	updated, err := r.client.UpdateService(serviceInput)
	if err != nil {
		log.Error().Msgf("[%s] API error: '%s'", service.Name, err.Error())
		return err
	}
	diff := cmp.Diff(service, updated)
	log.Info().Msgf("[%s] Updated Service - Diff:\n%s", service.Name, diff)
	return nil
}

func (r *ServiceReconciler) handleAliases(service *opslevel.Service, registration opslevel_jq_parser.ServiceRegistration) {
	for _, alias := range registration.Aliases {
		if service.HasAlias(alias) {
			continue
		}
		err := r.client.CreateAlias(opslevel.AliasCreateInput{
			Alias:   alias,
			OwnerId: service.Id,
		})
		if err != nil {
			log.Error().Msgf("[%s] Failed assigning alias '%s'\n\tREASON: %v", service.Name, alias, err.Error())
			continue
		}
		log.Info().Msgf("[%s] Assigned alias '%s'", service.Name, alias)
	}
}

func (r *ServiceReconciler) handleAssignTags(service *opslevel.Service, registration opslevel_jq_parser.ServiceRegistration) {
	if registration.TagAssigns == nil || (service.Tags != nil && r.ContainsAllTags(registration.TagAssigns, service.Tags.Nodes)) {
		log.Info().Msgf("[%s] 0/%d tags need to be assigned to service.", service.Name, len(registration.TagAssigns))
		return
	}
	tags := make(map[string]string)
	for _, tagAssign := range registration.TagAssigns {
		tags[tagAssign.Key] = tagAssign.Value
	}

	err := r.client.AssignTags(service, tags)
	jsonBytes, _ := json.Marshal(registration.TagAssigns)
	if err != nil {
		log.Error().Msgf("[%s] Failed assigning tags: %s\n\tREASON: %v", service.Name, string(jsonBytes), err.Error())
		return
	}
	log.Info().Msgf("[%s] Assigned tags: %s", service.Name, string(jsonBytes))
}

func (r *ServiceReconciler) handleCreateTags(service *opslevel.Service, registration opslevel_jq_parser.ServiceRegistration) {
	for _, tag := range registration.TagCreates {
		if service.Tags != nil && service.HasTag(tag.Key, tag.Value) {
			continue
		}
		input := opslevel.TagCreateInput{
			Id:    &service.Id,
			Key:   tag.Key,
			Value: tag.Value,
		}
		err := r.client.CreateTag(input)
		if err != nil {
			log.Error().Msgf("[%s] Failed creating tag '%s = %s'\n\tREASON: %v", service.Name, tag.Key, tag.Value, err.Error())
			continue
		}
		log.Info().Msgf("[%s] Created tag '%s = %s'", service.Name, tag.Key, tag.Value)
	}
}

func (r *ServiceReconciler) handleTools(service *opslevel.Service, registration opslevel_jq_parser.ServiceRegistration) {
	for _, tool := range registration.Tools {
		toolEnv := ""
		if tool.Environment != nil {
			toolEnv = *tool.Environment
		}
		if service.Tools != nil && service.HasTool(tool.Category, tool.DisplayName, toolEnv) {
			log.Debug().Msgf("[%s] Tool '{Category: %s, Environment: %s, Name: %s}' already exists on service ... skipping", service.Name, tool.Category, toolEnv, tool.DisplayName)
			continue
		}
		tool.ServiceId = &service.Id
		err := r.client.CreateTool(tool)
		if err != nil {
			log.Error().Msgf("[%s] Failed assigning tool '{Category: %s, Environment: %s, Name: %s}'\n\tREASON: %v", service.Name, tool.Category, toolEnv, tool.DisplayName, err.Error())
			continue
		}
		log.Info().Msgf("[%s] Ensured tool '{Category: %s, Environment: %s, Name: %s}'", service.Name, tool.Category, toolEnv, tool.DisplayName)
	}
}

func toJSON[T any](object T) string {
	s, _ := json.Marshal(object)
	return string(s)
}

func (r *ServiceReconciler) handleRepositories(service *opslevel.Service, registration opslevel_jq_parser.ServiceRegistration) {
	repoLog := log.With().Str("where", "handleRepositories").Str("service", service.Name).Logger()
	for _, newRepo := range registration.Repositories {
		if newRepo.Repository.Alias == nil || *newRepo.Repository.Alias == "" {
			repoLog.Warn().Msgf("repository (%s) has no alias ... skipping", toJSON(newRepo))
			continue
		}
		newRepoLog := repoLog.With().Str("repository", *newRepo.Repository.Alias).Logger()
		foundRepository, foundRepositoryErr := r.client.GetRepositoryWithAlias(*newRepo.Repository.Alias)
		if foundRepositoryErr != nil {
			newRepoLog.Error().Err(foundRepositoryErr).Msgf("repository get with alias error ... skipping")
			continue
		}
		// update repository case
		if foundRepository != nil {
			if newRepo.BaseDirectory == nil && newRepo.DisplayName == nil {
				// TODO: both fields cannot be unset in opslevel-go because of pointer + omitempty
				newRepoLog.Debug().Msgf("repository update has no base directory or display name ... skipping")
				continue
			}
			serviceRepository := foundRepository.GetService(service.Id, *newRepo.BaseDirectory)
			if serviceRepository == nil {
				newRepoLog.Warn().Msgf("tried to get service repository, got nil ... skipping")
				continue
			}
			if newRepo.BaseDirectory != nil && *newRepo.BaseDirectory == serviceRepository.BaseDirectory && newRepo.DisplayName != nil && *newRepo.DisplayName == serviceRepository.DisplayName {
				newRepoLog.Info().Msgf("repository (%s) already attached to service ... skipping", toJSON(newRepo))
				continue
			}
			updateInput := opslevel.ServiceRepositoryUpdateInput{
				Id:            serviceRepository.Id,
				BaseDirectory: newRepo.BaseDirectory,
				DisplayName:   newRepo.DisplayName,
			}
			err := r.client.UpdateServiceRepository(updateInput)
			if err != nil {
				newRepoLog.Error().Err(err).Msgf("failed updating repository (%s) ... skipping", toJSON(updateInput))
				continue
			}
			newRepoLog.Info().Msgf("successfully updated repository on service (%s)", toJSON(updateInput))
			continue
		}
		// create repository case
		err := r.client.CreateServiceRepository(newRepo)
		if err != nil {
			newRepoLog.Error().Err(err).Msgf("failed assigning repository (%s)", toJSON(newRepo))
			continue
		}
		newRepoLog.Info().Msgf("successfully attached repository (%s)", toJSON(newRepo))
	}
}

func (r *ServiceReconciler) handleProperties(service *opslevel.Service, registration opslevel_jq_parser.ServiceRegistration) {
	for _, propertyInput := range registration.Properties {
		propertyInput.Owner = *opslevel.NewIdentifier(string(service.Id))
		err := r.client.AssignPropertyHandler(propertyInput)
		if err != nil {
			// TODO: nil here
			log.Error().Err(err).Msgf("[%s] Failed assigning property with definition: '%s' and value: '%s'", service.Name, *propertyInput.Definition.Alias, propertyInput.Value)
			continue
		}
		// TODO: nil here
		log.Info().Msgf("[%s] Successfully assigned property with definition: '%s' and value: '%s'", service.Name, *propertyInput.Definition.Alias, propertyInput.Value)
	}
}
