package common

import (
	"encoding/json"
	"fmt"
	"strings"

	"golang.org/x/exp/maps"

	"github.com/google/go-cmp/cmp"
	"github.com/opslevel/opslevel-go/v2024"
	opslevel_jq_parser "github.com/opslevel/opslevel-jq-parser/v2024"
	"github.com/rs/zerolog/log"
)

type serviceAliasesResult string

const (
	serviceAliasesResult_NoAliasesMatched      serviceAliasesResult = "NoAliasesMatched"
	serviceAliasesResult_AliasMatched          serviceAliasesResult = "AliasMatched"
	serviceAliasesResult_MultipleServicesFound serviceAliasesResult = "MultipleServicesFound"
	serviceAliasesResult_APIErrorHappened      serviceAliasesResult = "APIErrorHappened"
	serviceAliasesResult_FoundServiceNoAlias   serviceAliasesResult = "FoundServiceNoAlias"
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

func (r *ServiceReconciler) Reconcile(registration opslevel_jq_parser.ServiceRegistration) error {
	if len(registration.Aliases) <= 0 {
		return fmt.Errorf("[%s] found 0 aliases from kubernetes data", registration.Name)
	}
	service, err := r.handleService(registration)
	if err != nil {
		return err
	}
	if service == nil {
		return nil
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

func (r *ServiceReconciler) ServiceNeedsUpdate(input opslevel.ServiceUpdateInput, service *opslevel.Service) bool {
	if input.Name != nil && *input.Name != service.Name {
		return true
	}
	if input.Product != nil && *input.Product != service.Product {
		return true
	}
	if input.Description != nil && *input.Description != service.Description {
		return true
	}
	if input.Language != nil && *input.Language != service.Language {
		return true
	}
	if input.Framework != nil && *input.Framework != service.Framework {
		return true
	}
	if input.TierAlias != nil && *input.TierAlias != service.Tier.Alias {
		return true
	}
	if input.LifecycleAlias != nil && *input.LifecycleAlias != service.Lifecycle.Alias {
		return true
	}
	if input.OwnerInput != nil && *input.OwnerInput.Alias != service.Owner.Alias {
		return true
	}
	return false
}

// This function has 4 outcomes that can happen while looping over the aliases list
// serviceAliasesResult_NoAliasesMatched - means that all API calls succeeded and none of the aliases matched an existing service
// serviceAliasesResult_AliasMatched - means that all the API calls succeeded and a single service was found matching 1 of N aliases
// serviceAliasesResult_MultipleServicesFound - means that all API calls succeeded but multiple services were returning means the list of aliases does not definitively describe a single service and might be a configuration problem
// serviceAliasesResult_APIErrorHappened - means that 1 of N aliases got an 4xx/5xx and thereforce we cannot say 100% that the services doesn't exist
// serviceAliasesResult_FoundServiceNoAlias - means that a service was found but that service has no alias (this should not be possible and can only happen from a bad code change.)
func (r *ServiceReconciler) lookupService(registration opslevel_jq_parser.ServiceRegistration) (*opslevel.Service, serviceAliasesResult) {
	var gotError error
	foundServices := map[string]*opslevel.Service{}
	for _, alias := range registration.Aliases {
		foundService, err := r.client.GetService(alias)
		if err != nil {
			gotError = err
			log.Warn().Err(err).Msgf("got an error when trying to get service with alias '%s'", alias)
		} else if foundService == nil {
			log.Warn().Msgf("unexpected happened: got service with alias '%s' but the result is nil", alias)
		} else if foundService.Id == "" {
			log.Warn().Msgf("unexpected happened: got service with alias '%s' but the result has no ID", alias)
		} else {
			// happy path
			foundServices[string(foundService.Id)] = foundService
		}
	}
	if gotError != nil {
		return nil, serviceAliasesResult_APIErrorHappened
	}
	foundServicesCount := len(foundServices)
	if foundServicesCount == 1 {
		keys := maps.Keys(foundServices)
		if len(keys) == 0 {
			return nil, serviceAliasesResult_FoundServiceNoAlias
		}
		key := keys[0]
		return foundServices[key], serviceAliasesResult_AliasMatched
	} else if foundServicesCount > 1 {
		return nil, serviceAliasesResult_MultipleServicesFound
	} else {
		return nil, serviceAliasesResult_NoAliasesMatched
	}
}

func (r *ServiceReconciler) handleService(registration opslevel_jq_parser.ServiceRegistration) (*opslevel.Service, error) {
	service, status := r.lookupService(registration)
	switch status {
	case serviceAliasesResult_NoAliasesMatched:
		if r.disableServiceCreation {
			log.Info().Msgf("[%s] Avoided creating a new service\n\tREASON: service creation is disabled", registration.Name)
			return nil, nil
		}

		newService, newServiceErr := r.createService(registration)
		if newServiceErr != nil {
			return nil, fmt.Errorf("[%s] api error during service creation ... skipping reconciliation.\n\tREASON: %v", registration.Name, newServiceErr)
		}
		service = newService
	case serviceAliasesResult_AliasMatched:
		r.updateService(service, registration)
	case serviceAliasesResult_MultipleServicesFound:
		aliases := ""
		if service != nil {
			aliases = fmt.Sprintf(`"%s"`, strings.Join(service.Aliases, `", "`))
		}
		return nil, fmt.Errorf("[%s] found multiple services with aliases = [%s].  cannot know which service to target for update ... skipping reconciliation", registration.Name, aliases)
	case serviceAliasesResult_APIErrorHappened:
		return nil, fmt.Errorf("[%s] api error during service lookup by alias.  unable to guarantee service was found or not ... skipping reconciliation", registration.Name)
	case serviceAliasesResult_FoundServiceNoAlias:
		return nil, fmt.Errorf("[%s] found matching service but it unexpectedly has no alias.  please submit a bug report. ... skipping reconciliation", registration.Name)
	}
	return service, nil
}

func (r *ServiceReconciler) createService(registration opslevel_jq_parser.ServiceRegistration) (*opslevel.Service, error) {
	serviceCreateInput := opslevel.ServiceCreateInput{
		Name:        registration.Name,
		Product:     opslevel.RefOf[string](registration.Product),
		Description: opslevel.RefOf[string](registration.Description),
		Language:    opslevel.RefOf[string](registration.Language),
		Framework:   opslevel.RefOf[string](registration.Framework),
	}
	if registration.System != "" {
		serviceCreateInput.Parent = opslevel.NewIdentifier(registration.System)
	}
	if v, ok := opslevel.Cache.TryGetTier(registration.Tier); ok {
		if v == nil {
			err := fmt.Errorf("the cache unexpectedly returned a tier that is nil - please submit a bug report")
			return nil, fmt.Errorf("[%s] Failed creating service\n\tREASON: %v", registration.Name, err.Error())
		}
		serviceCreateInput.TierAlias = opslevel.RefOf(v.Alias)
	} else if registration.Tier != "" {
		log.Warn().Msgf("[%s] Unable to find 'Tier' with alias '%s'", registration.Name, registration.Tier)
	}
	if v, ok := opslevel.Cache.TryGetLifecycle(registration.Lifecycle); ok {
		if v == nil {
			err := fmt.Errorf("the cache unexpectedly returned a lifecycle that is nil - please submit a bug report")
			return nil, fmt.Errorf("[%s] Failed creating service\n\tREASON: %v", registration.Name, err.Error())
		}
		serviceCreateInput.LifecycleAlias = opslevel.RefOf(v.Alias)
	} else if registration.Lifecycle != "" {
		log.Warn().Msgf("[%s] Unable to find 'Lifecycle' with alias '%s'", registration.Name, registration.Lifecycle)
	}
	if v, ok := opslevel.Cache.TryGetTeam(registration.Owner); ok {
		if v == nil {
			err := fmt.Errorf("the cache unexpectedly returned a team that is nil - please submit a bug report")
			return nil, fmt.Errorf("[%s] Failed creating service\n\tREASON: %v", registration.Name, err.Error())
		}
		serviceCreateInput.OwnerInput = opslevel.NewIdentifier(v.Alias)
	} else if registration.Owner != "" {
		log.Warn().Msgf("[%s] Unable to find 'Team' with alias '%s'", registration.Name, registration.Owner)
	}
	service, err := r.client.CreateService(serviceCreateInput)
	if err != nil {
		return service, fmt.Errorf("[%s] Failed creating service\n\tREASON: %v", registration.Name, err.Error())
	} else if service != nil {
		log.Info().Msgf("[%s] Created new service", service.Name)
		return service, nil
	} else {
		return nil, fmt.Errorf("[%s] unexpected happened: created service but the result is nil", registration.Name)
	}
}

func (r *ServiceReconciler) updateService(service *opslevel.Service, registration opslevel_jq_parser.ServiceRegistration) {
	if service == nil {
		log.Warn().Msgf("[%s] unexpected happened: service passed to be updated is nil", registration.Name)
		return
	}
	updateServiceInput := opslevel.ServiceUpdateInput{
		Id:          &service.Id,
		Product:     opslevel.RefOf[string](registration.Product),
		Description: opslevel.RefOf[string](registration.Description),
		Language:    opslevel.RefOf[string](registration.Language),
		Framework:   opslevel.RefOf[string](registration.Framework),
	}
	if registration.System != "" {
		updateServiceInput.Parent = opslevel.NewIdentifier(registration.System)
	}
	if v, ok := opslevel.Cache.TryGetTier(registration.Tier); ok {
		if v == nil {
			err := fmt.Errorf("the cache unexpectedly returned a tier that is nil - please submit a bug report")
			log.Warn().Msgf("[%s] unexpected happened: %v", service.Name, err)
		} else {
			updateServiceInput.TierAlias = opslevel.RefOf(v.Alias)
		}
	} else if registration.Tier != "" {
		log.Warn().Msgf("[%s] Unable to find 'Tier' with alias '%s'", service.Name, registration.Tier)
	}
	if v, ok := opslevel.Cache.TryGetLifecycle(registration.Lifecycle); ok {
		if v == nil {
			err := fmt.Errorf("the cache unexpectedly returned a lifecycle that is nil - please submit a bug report")
			log.Warn().Msgf("[%s] unexpected happened: %v", service.Name, err)
		} else {
			updateServiceInput.LifecycleAlias = opslevel.RefOf(v.Alias)
		}
	} else if registration.Lifecycle != "" {
		log.Warn().Msgf("[%s] Unable to find 'Lifecycle' with alias '%s'", service.Name, registration.Lifecycle)
	}
	if v, ok := opslevel.Cache.TryGetTeam(registration.Owner); ok {
		if v == nil {
			err := fmt.Errorf("the cache unexpectedly returned a team that is nil - please submit a bug report")
			log.Warn().Msgf("[%s] unexpected happened: %v", service.Name, err)
		} else {
			updateServiceInput.OwnerInput = opslevel.NewIdentifier(v.Alias)
		}
	} else if registration.Owner != "" {
		log.Warn().Msgf("[%s] Unable to find 'Team' with alias '%s'", service.Name, registration.Owner)
	}
	if r.ServiceNeedsUpdate(updateServiceInput, service) {
		updatedService, updateServiceErr := r.client.UpdateService(updateServiceInput)
		if updateServiceErr != nil {
			log.Error().Msgf("[%s] Failed updating service\n\tREASON: %v", service.Name, updateServiceErr.Error())
		} else if updatedService == nil {
			log.Warn().Msgf("[%s] unexpected happened: updated service but the result is nil - please submit a bug report", service.Name)
		} else if diff := cmp.Diff(service, updatedService); diff != "" {
			log.Info().Msgf("[%s] Updated Service - Diff:\n%s", service.Name, diff)
		}
	} else {
		log.Info().Msgf("[%s] No changes detected to fields - skipping update", service.Name)
	}
}

func (r *ServiceReconciler) handleAliases(service *opslevel.Service, registration opslevel_jq_parser.ServiceRegistration) {
	for _, alias := range registration.Aliases {
		if alias == "" || service.HasAlias(alias) {
			continue
		}
		err := r.client.CreateAlias(opslevel.AliasCreateInput{
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

func (r *ServiceReconciler) handleAssignTags(service *opslevel.Service, registration opslevel_jq_parser.ServiceRegistration) {
	if registration.TagAssigns == nil {
		return
	}
	if !r.ContainsAllTags(registration.TagAssigns, service.Tags.Nodes) {
		tags := map[string]string{}
		for _, tagAssign := range registration.TagAssigns {
			tags[tagAssign.Key] = tagAssign.Value
		}

		err := r.client.AssignTags(service, tags)
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

func (r *ServiceReconciler) handleCreateTags(service *opslevel.Service, registration opslevel_jq_parser.ServiceRegistration) {
	for _, tag := range registration.TagCreates {
		if service.HasTag(tag.Key, tag.Value) {
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
		} else {
			log.Info().Msgf("[%s] Created tag '%s = %s'", service.Name, tag.Key, tag.Value)
		}
	}
}

func (r *ServiceReconciler) handleTools(service *opslevel.Service, registration opslevel_jq_parser.ServiceRegistration) {
	for _, tool := range registration.Tools {
		toolEnv := ""
		if tool.Environment != nil {
			toolEnv = *tool.Environment
		}
		if service.HasTool(tool.Category, tool.DisplayName, toolEnv) {
			log.Debug().Msgf("[%s] Tool '{Category: %s, Environment: %s, Name: %s}' already exists on service ... skipping", service.Name, tool.Category, toolEnv, tool.DisplayName)
			continue
		}
		tool.ServiceId = &service.Id
		err := r.client.CreateTool(tool)
		if err != nil {
			log.Error().Msgf("[%s] Failed assigning tool '{Category: %s, Environment: %s, Name: %s}'\n\tREASON: %v", service.Name, tool.Category, toolEnv, tool.DisplayName, err.Error())
		} else {
			log.Info().Msgf("[%s] Ensured tool '{Category: %s, Environment: %s, Name: %s}'", service.Name, tool.Category, toolEnv, tool.DisplayName)
		}
	}
}

func (r *ServiceReconciler) handleRepositories(service *opslevel.Service, registration opslevel_jq_parser.ServiceRegistration) {
	for _, repositoryCreate := range registration.Repositories {
		if repositoryCreate.BaseDirectory == nil {
			log.Warn().Msgf("[%s] Repository with alias: '%+v' has no base directory so it cannot be looked up ... skipping", service.Name, repositoryCreate)
			continue
		}
		foundRepository, foundRepositoryErr := r.client.GetRepositoryWithAlias(*repositoryCreate.Repository.Alias)
		if foundRepositoryErr != nil {
			log.Warn().Msgf("[%s] Repository with alias: '%+v' not found so it cannot be attached to service ... skipping", service.Name, repositoryCreate)
			continue
		} else if foundRepository == nil {
			log.Warn().Msgf("[%s] Repository with alias: '%+v' call to GetRepositoryWithAlias unexpectedly returned nil - please submit a bug report ... skipping", service.Name, repositoryCreate)
			continue
		}
		serviceRepository := foundRepository.GetService(service.Id, *repositoryCreate.BaseDirectory)
		if serviceRepository != nil {
			if repositoryCreate.DisplayName != nil && serviceRepository.DisplayName != *repositoryCreate.DisplayName {
				repositoryUpdate := opslevel.ServiceRepositoryUpdateInput{
					Id:          serviceRepository.Id,
					DisplayName: repositoryCreate.DisplayName,
				}
				err := r.client.UpdateServiceRepository(repositoryUpdate)
				if err != nil {
					log.Error().Msgf("[%s] Failed updating repository '%+v'\n\tREASON: %v", service.Name, repositoryCreate, err.Error())
					continue
				} else {
					log.Info().Msgf("[%s] Updated repository '%+v'", service.Name, repositoryCreate)
					continue
				}
			}
			log.Debug().Msgf("[%s] Repository '%+v' already attached to service ... skipping", service.Name, repositoryCreate)
			continue
		}
		repositoryCreate.Service = opslevel.IdentifierInput{Id: &service.Id}
		err := r.client.CreateServiceRepository(repositoryCreate)
		if err != nil {
			log.Error().Msgf("[%s] Failed assigning repository '%+v'\n\tREASON: %v", service.Name, repositoryCreate, err.Error())
		} else {
			log.Info().Msgf("[%s] Attached repository '%+v'", service.Name, repositoryCreate)
		}
	}
}

func (r *ServiceReconciler) handleProperties(service *opslevel.Service, registration opslevel_jq_parser.ServiceRegistration) {
	for def, val := range registration.Properties {
		definition := opslevel.NewIdentifier(def)
		owner := opslevel.NewIdentifier(string(service.Id))
		value, err := opslevel.NewJSONInput(val)
		if err != nil {
			log.Error().Err(err).Msgf("[%s] NewJSONInput failed parsing property value: '%s' on definition: '%s'", service.Name, val, def)
			continue
		}
		input := opslevel.PropertyInput{
			Definition: *definition,
			Owner:      *owner,
			Value:      *value,
		}
		err = r.client.AssignPropertyHandler(input)
		if err != nil {
			log.Error().Err(err).Msgf("[%s] Failed assigning property with definition: '%s' and value: '%s'", service.Name, def, val)
			continue
		}
		log.Info().Msgf("[%s] Successfully assigned property with definition: '%s' and value: '%s'", service.Name, def, val)
	}
}
