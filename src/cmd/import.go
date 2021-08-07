package cmd

import (
	"encoding/json"
	"fmt"
	"sync"

	"github.com/opslevel/kubectl-opslevel/common"
	"github.com/opslevel/kubectl-opslevel/config"
	"github.com/opslevel/opslevel-go"

	"github.com/google/go-cmp/cmp"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
)

var importCmd = &cobra.Command{
	Use:   "import",
	Short: "Create or Update service entries in OpsLevel",
	Long:  `This command will take the data found in your Kubernetes cluster and begin to reconcile it with OpsLevel`,
	Run:   runImport,
}

func init() {
	serviceCmd.AddCommand(importCmd)
}

func runImport(cmd *cobra.Command, args []string) {
	config, configErr := config.New()
	cobra.CheckErr(configErr)

	services, servicesErr := common.QueryForServices(config)
	cobra.CheckErr(servicesErr)

	client := common.NewClient()
	CacheLookupTables(client)

	done := make(chan bool)
	queue := make(chan common.ServiceRegistration, concurrency)
	go createWorkerPool(concurrency, queue, done)
	go enqueue(services, queue)
	<-done
	log.Info().Msg("Import Complete")
}

// TODO: Helpers probably shouldn't be exported
// Helpers

func createWorkerPool(count int, queue chan common.ServiceRegistration, done chan<- bool) {
	var waitGroup sync.WaitGroup
	waitGroup.Add(count)
	for i := 0; i < count; i++ {
		go func(c *opslevel.Client, q chan common.ServiceRegistration, wg *sync.WaitGroup) {
			for data := range q {
				reconcileService(c, data)
			}
			wg.Done()
		}(common.NewClient(), queue, &waitGroup)
	}
	waitGroup.Wait()
	done <- true
}

func enqueue(services []common.ServiceRegistration, queue chan common.ServiceRegistration) {
	for _, service := range services {
		queue <- service
	}
	close(queue)
}

func reconcileService(client *opslevel.Client, service common.ServiceRegistration) {
	foundService, needsUpdate := FindService(client, service)
	if foundService == nil {
		newService, newServiceErr := CreateService(client, service)
		if newServiceErr != nil {
			log.Error().Msgf("[%s] Failed creating service\n\tREASON: %v", service.Name, newServiceErr.Error())
			return
		} else {
			log.Info().Msgf("[%s] Created new service", newService.Name)
		}
		foundService = newService
	}
	if needsUpdate {
		UpdateService(client, service, foundService)
	}
	go AssignAliases(client, service, foundService)
	go AssignTags(client, service, foundService)
	go CreateTags(client, service, foundService)
	go AssignTools(client, service, foundService)
	go AttachRepositories(client, service, foundService)
	log.Info().Msgf("[%s] Finished processing data", foundService.Name)
}

func FindService(client *opslevel.Client, registration common.ServiceRegistration) (*opslevel.Service, bool) {
	for _, alias := range registration.Aliases {
		foundService, err := client.GetServiceWithAlias(alias)
		if err == nil && foundService.Id != nil {
			log.Info().Msgf("[%s] Reconciling service found with alias '%s' ...", foundService.Name, alias)
			return foundService, true
		}
	}
	// TODO: last ditch effort - search for service with alias == registration.Name ?
	return nil, false
}

func GetTiers(client *opslevel.Client) (map[string]opslevel.Tier, error) {
	tiers := make(map[string]opslevel.Tier)
	tiersList, tiersErr := client.ListTiers()
	if tiersErr != nil {
		return tiers, tiersErr
	}
	for _, tier := range tiersList {
		tiers[string(tier.Alias)] = tier
	}
	return tiers, nil
}

func GetLifecycles(client *opslevel.Client) (map[string]opslevel.Lifecycle, error) {
	lifecycles := make(map[string]opslevel.Lifecycle)
	lifecyclesList, lifecyclesErr := client.ListLifecycles()
	if lifecyclesErr != nil {
		return lifecycles, lifecyclesErr
	}
	for _, lifecycle := range lifecyclesList {
		lifecycles[string(lifecycle.Alias)] = lifecycle
	}
	return lifecycles, nil
}

func GetTeams(client *opslevel.Client) (map[string]opslevel.Team, error) {
	teams := make(map[string]opslevel.Team)
	data, dataErr := client.ListTeams()
	if dataErr != nil {
		return teams, dataErr
	}
	for _, team := range data {
		teams[string(team.Alias)] = team
	}
	return teams, nil
}

// TODO: this makes this code hard to test
var (
	Tiers      map[string]opslevel.Tier
	Lifecycles map[string]opslevel.Lifecycle
	Teams      map[string]opslevel.Team
)

func CacheLookupTables(client *opslevel.Client) {
	log.Info().Msg("Caching 'Tiers' lookup table from OpsLevel API ...")
	tiers, tiersErr := GetTiers(client)
	if tiersErr != nil {
		log.Warn().Msgf("===> Failed to retrive tiers from OpsLevel API - Unable to assign field 'Tier' to services. REASON: %s", tiersErr.Error())
	}
	Tiers = tiers

	log.Info().Msg("Caching 'Lifecycles' lookup table from OpsLevel API ...")
	lifecycles, lifecyclesErr := GetLifecycles(client)
	if lifecyclesErr != nil {
		log.Warn().Msgf("===> Failed to retrive lifecycles from OpsLevel API - Unable to assign field 'Lifecycle' to services. REASON: %s", lifecyclesErr.Error())
	}
	Lifecycles = lifecycles

	log.Info().Msg("Caching 'Teams' lookup table from OpsLevel API ...")
	teams, teamsErr := GetTeams(client)
	if teamsErr != nil {
		log.Warn().Msgf("===> Failed to retrive teams from OpsLevel API - Unable to assign field 'Owner' to services. REASON: %s", teamsErr.Error())
	}
	Teams = teams
}

func CreateService(client *opslevel.Client, registration common.ServiceRegistration) (*opslevel.Service, error) {
	serviceCreateInput := opslevel.ServiceCreateInput{
		Name:        registration.Name,
		Product:     registration.Product,
		Description: registration.Description,
		Language:    registration.Language,
		Framework:   registration.Framework,
	}
	if v, ok := Tiers[registration.Tier]; ok {
		serviceCreateInput.Tier = string(v.Alias)
	}
	if v, ok := Lifecycles[registration.Lifecycle]; ok {
		serviceCreateInput.Lifecycle = string(v.Alias)
	}
	if v, ok := Teams[registration.Owner]; ok {
		serviceCreateInput.Owner = string(v.Alias)
	}
	return client.CreateService(serviceCreateInput)
}

func UpdateService(client *opslevel.Client, registration common.ServiceRegistration, service *opslevel.Service) {
	updateServiceInput := opslevel.ServiceUpdateInput{
		Id:           service.Id,
		Product:      registration.Product,
		Descripition: registration.Description,
		Language:     registration.Language,
		Framework:    registration.Framework,
	}
	if v, ok := Tiers[registration.Tier]; ok {
		updateServiceInput.Tier = string(v.Alias)
	}
	if v, ok := Lifecycles[registration.Lifecycle]; ok {
		updateServiceInput.Lifecycle = string(v.Alias)
	}
	if v, ok := Teams[registration.Owner]; ok {
		updateServiceInput.Owner = string(v.Alias)
	}
	updatedService, updateServiceErr := client.UpdateService(updateServiceInput)
	if updateServiceErr != nil {
		log.Error().Msgf("[%s] Failed updating service\n\tREASON: %v", service.Name, updateServiceErr.Error())
	} else {
		if diff := cmp.Diff(service, updatedService); diff != "" {
			log.Info().Msgf("[%s] Updated Service - Diff:\n%s", service.Name, diff)
		}
	}
}

func AssignAliases(client *opslevel.Client, registration common.ServiceRegistration, service *opslevel.Service) {
	for _, alias := range registration.Aliases {
		if service.HasAlias(alias) {
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

func AssignTags(client *opslevel.Client, registration common.ServiceRegistration, service *opslevel.Service) {
	_, err := client.AssignTagsForId(service.Id, registration.TagAssigns)
	jsonBytes, _ := json.Marshal(registration.TagAssigns)
	if err != nil {
		log.Error().Msgf("[%s] Failed assigning tags: %s\n\tREASON: %v", service.Name, string(jsonBytes), err.Error())
	} else {
		log.Info().Msgf("[%s] Assigned tags: %s", service.Name, string(jsonBytes))
	}
}

func CreateTags(client *opslevel.Client, registration common.ServiceRegistration, service *opslevel.Service) {
	for tagKey, tagValue := range registration.TagCreates {
		if service.HasTag(tagKey, tagValue) {
			continue
		}
		input := opslevel.TagCreateInput{
			Id:    service.Id,
			Key:   tagKey,
			Value: tagValue,
		}
		_, err := client.CreateTag(input)
		if err != nil {
			log.Error().Msgf("[%s] Failed creating tag '%s = %s'\n\tREASON: %v", service.Name, tagKey, tagValue, err.Error())
		} else {
			log.Info().Msgf("[%s] Created tag '%s = %s'", service.Name, tagKey, tagValue)
		}
	}
}

func AssignTools(client *opslevel.Client, registration common.ServiceRegistration, service *opslevel.Service) {
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

func AttachRepositories(client *opslevel.Client, registration common.ServiceRegistration, service *opslevel.Service) {
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
			log.Error().Msgf("[%s] Failed assigning repository '$s'\n\tREASON: %v", service.Name, repositoryAsString, err.Error())
		} else {
			log.Info().Msgf("[%s] Attached repository '%s'", service.Name, repositoryAsString)
		}
	}
}
