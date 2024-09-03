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
	client                   *OpslevelClient
	disableServiceCreation   bool
	disableServiceNameUpdate bool
}

func NewServiceReconciler(client *OpslevelClient, disableServiceCreation, disableServiceNameUpdate bool) *ServiceReconciler {
	return &ServiceReconciler{
		client:                   client,
		disableServiceCreation:   disableServiceCreation,
		disableServiceNameUpdate: disableServiceNameUpdate,
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

// This function has 4 outcomes that can happen while looping over the aliases list
// serviceAliasesResult_NoAliasesMatched - means that all API calls succeeded and none of the aliases matched an existing service
// serviceAliasesResult_AliasMatched - means that all the API calls succeeded and a single service was found matching 1 of N aliases
// serviceAliasesResult_MultipleServicesFound - means that all API calls succeeded but multiple services were returning means the list of aliases does not definitively describe a single service and might be a configuration problem
// serviceAliasesResult_APIErrorHappened - means that 1 of N aliases got a 4xx/5xx and thereforce we cannot say 100% that the services doesn't exist
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
			if len(foundService.Aliases) == 1 {
				// If this happens and there is only 1 alias to check we cannot assume the service doesn't exist
				// because it seems like our API has a race condition looking up the service
				return nil, serviceAliasesResult_APIErrorHappened
			}
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

// updateService uses compares each field (not foreign keys like Tools or Tags) value in the registration vs the value that is currently set on the service.
// if there are any updates needed, it will send a ServiceUpdateInput to the API.
func (r *ServiceReconciler) updateService(service *opslevel.Service, registration opslevel_jq_parser.ServiceRegistration) {
	if service == nil {
		log.Warn().Msgf("[%s] unexpected happened: service passed to be updated is nil", registration.Name)
		return
	}

	// updateServiceInput contains the changes needed to reconcile the service
	updateServiceInput := opslevel.ServiceUpdateInput{Id: &service.Id}
	// for each field - check if the value exists in the registration AND if the value has changed compared to what is currently set
	// cannot use cmp.Diff to compare field values, since that is used for comparing structs and not individual fields.
	// some fields like System need special comparisons, e.g. by using systemIdHasAlias
	// only purpose of cmp.Diff is to display an easy-to-read diff for the user to understand what happened and to check if there
	// is a need to submit an API update request, since the output of cmp.Diff is not really parseable.
	if registration.Description != "" && registration.Description != service.Description {
		updateServiceInput.Description = opslevel.RefOf(registration.Description)
	}
	if registration.Framework != "" && registration.Framework != service.Framework {
		updateServiceInput.Framework = opslevel.RefOf(registration.Framework)
	}
	if registration.Language != "" && registration.Language != service.Language {
		updateServiceInput.Language = opslevel.RefOf(registration.Language)
	}
	if registration.Lifecycle != "" && registration.Lifecycle != service.Lifecycle.Alias {
		if lifecycle, ok := opslevel.Cache.TryGetLifecycle(registration.Lifecycle); ok {
			if lifecycle == nil {
				err := fmt.Errorf("the cache unexpectedly returned a lifecycle that is nil - please submit a bug report")
				log.Warn().Msgf("[%s] unexpected happened: %v", service.Name, err)
			} else {
				updateServiceInput.LifecycleAlias = opslevel.RefOf(lifecycle.Alias)
			}
		} else if registration.Lifecycle != "" {
			log.Warn().Msgf("[%s] Unable to find 'Lifecycle' with alias '%s'", service.Name, registration.Lifecycle)
		}
	}
	if !r.disableServiceNameUpdate && registration.Name != "" && registration.Name != service.Name {
		updateServiceInput.Name = opslevel.RefOf(registration.Name)
	}
	if registration.Owner != "" && registration.Owner != service.Owner.Alias {
		if team, ok := opslevel.Cache.TryGetTeam(registration.Owner); ok {
			if team == nil {
				err := fmt.Errorf("the cache unexpectedly returned a team that is nil - please submit a bug report")
				log.Warn().Msgf("[%s] unexpected happened: %v", service.Name, err)
			} else {
				updateServiceInput.OwnerInput = opslevel.NewIdentifier(team.Alias)
			}
		} else if registration.Owner != "" {
			log.Warn().Msgf("[%s] Unable to find 'Team' with alias '%s'", service.Name, registration.Owner)
		}
	}
	// TODO: use the opslevel-go system cache here once it is added
	if registration.System != "" && !systemIdHasAlias(service.Parent, registration.System) {
		updateServiceInput.Parent = opslevel.NewIdentifier(registration.System)
	}
	if registration.Product != "" && registration.Product != service.Product {
		updateServiceInput.Product = opslevel.RefOf(registration.Product)
	}
	if registration.Tier != "" && registration.Tier != service.Tier.Alias {
		if tier, ok := opslevel.Cache.TryGetTier(registration.Tier); ok {
			if tier == nil {
				err := fmt.Errorf("the cache unexpectedly returned a tier that is nil - please submit a bug report")
				log.Warn().Msgf("[%s] unexpected happened: %v", service.Name, err)
			} else {
				updateServiceInput.TierAlias = opslevel.RefOf(tier.Alias)
			}
		} else if registration.Tier != "" {
			log.Warn().Msgf("[%s] Unable to find 'Tier' with alias '%s'", service.Name, registration.Tier)
		}
	}
	// if there is nothing in updateServiceInput aside from the service ID, do not send an update service API call
	if cmp.Equal(opslevel.ServiceUpdateInput{Id: &service.Id}, updateServiceInput) {
		log.Info().Msgf("[%s] No changes detected to fields - skipping update", service.Name)
		return
	}
	updateJSON, _ := json.Marshal(updateServiceInput)
	log.Info().Msgf("[%s] Detected Changes - Sending Update:\n%s", service.Name, string(updateJSON))
	updatedService, updateServiceErr := r.client.UpdateService(updateServiceInput)
	if updateServiceErr != nil {
		log.Error().Msgf("[%s] Failed updating service\n\tREASON: %v", service.Name, updateServiceErr.Error())
		return
	}
	if updatedService == nil {
		log.Warn().Msgf("[%s] unexpected happened: updated service but the result is nil - please submit a bug report", service.Name)
		return
	}
	serviceDiff := cmp.Diff(service, updatedService)
	log.Info().Msgf("[%s] Updated Service - Diff:\n%s", service.Name, serviceDiff)
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
	for _, inputRepository := range registration.Repositories {
		if inputRepository.Repository.Alias == nil || *inputRepository.Repository.Alias == "null" || *inputRepository.Repository.Alias == "" {
			continue
		}

		// setup logger to use for this repository
		repoLogger := log.With().Str("service", service.Name).Str("repo", *inputRepository.Repository.Alias).Logger()
		if inputRepository.BaseDirectory == nil {
			// case for when input is just 'github.com:OrgName/RepoName'
			repoLogger.Warn().Msgf("no base directory was associated with this repo - using root directory (/)")
			inputRepository.BaseDirectory = opslevel.RefOf("")
		}
		repoLogger = repoLogger.With().Str("base_dir", *inputRepository.BaseDirectory).Logger()

		// look up the repository in OpsLevel - exit if it does not exist
		foundRepository, foundRepositoryErr := r.client.GetRepositoryWithAlias(*inputRepository.Repository.Alias)
		if foundRepositoryErr != nil {
			repoLogger.Error().Err(foundRepositoryErr).Msgf("fetching repository in OpsLevel resulted in an error ... skipping")
			continue
		} else if foundRepository == nil {
			repoLogger.Warn().Msgf("repository not found in OpsLevel ... skipping")
			continue
		}

		// look up the ServiceRepository matching the base directory
		serviceRepository := foundRepository.GetService(service.Id, *inputRepository.BaseDirectory)

		// if the ServiceRepository is found, update the fields on the ServiceRepository (currently just the display name)
		if serviceRepository != nil {
			repoLogger.Debug().Msgf("found service repository (%s) has display name: '%s'", serviceRepository.Id, serviceRepository.DisplayName)

			// check to ensure that an update call is actually necessary
			repositoryUpdate := opslevel.ServiceRepositoryUpdateInput{Id: serviceRepository.Id, BaseDirectory: &serviceRepository.BaseDirectory}
			needsUpdate := false
			if inputRepository.DisplayName != nil && *inputRepository.DisplayName != "" && *inputRepository.DisplayName != serviceRepository.DisplayName {
				repositoryUpdate.DisplayName = inputRepository.DisplayName
				needsUpdate = true
			}
			if !needsUpdate {
				repoLogger.Info().Msgf("service repository (%s) does not require any updates", serviceRepository.Id)
				continue
			}

			// perform update
			serviceRepositoryUpdateErr := r.client.UpdateServiceRepository(repositoryUpdate)
			if serviceRepositoryUpdateErr != nil {
				repoLogger.Error().Err(serviceRepositoryUpdateErr).Msgf("failed updating service repository (%s)", serviceRepository.Id)
				continue
			}
			repoLogger.Info().Msgf("successfully updated service repository (%s)", serviceRepository.Id)
			continue
		}

		// if the ServiceRepository is not found, create the ServiceRepository (assign the repository to the current service)
		inputRepository.Service = opslevel.IdentifierInput{Id: &service.Id}
		err := r.client.CreateServiceRepository(inputRepository)
		if err != nil {
			repoLogger.Error().Err(err).Msgf("failed creating a new service repository")
			continue
		}
		repoLogger.Info().Msgf("successfully created a new service repository")
	}
}

func (r *ServiceReconciler) handleProperties(service *opslevel.Service, registration opslevel_jq_parser.ServiceRegistration) {
	for _, propertyInput := range registration.Properties {
		if propertyInput.Definition.Alias == nil {
			log.Warn().Msgf("[%s] Cannot assign property with no definition ... skipping", service.Name)
		}
		propertyInput.Owner = *opslevel.NewIdentifier(string(service.Id))
		err := r.client.AssignPropertyHandler(propertyInput)
		if err != nil {
			log.Error().Err(err).Msgf("[%s] Failed assigning property with definition: '%s' and value: '%s'", service.Name, *propertyInput.Definition.Alias, propertyInput.Value)
			continue
		}
		log.Info().Msgf("[%s] Successfully assigned property with definition: '%s' and value: '%s'", service.Name, *propertyInput.Definition.Alias, propertyInput.Value)
	}
}

func systemIdHasAlias(sys *opslevel.SystemId, alias string) bool {
	if sys == nil {
		return false
	}
	for _, existingAlias := range sys.Aliases {
		if alias == existingAlias {
			return true
		}
	}
	return false
}
