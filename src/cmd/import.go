package cmd

import (
	"github.com/opslevel/kubectl-opslevel/common"
	"github.com/opslevel/kubectl-opslevel/config"
	"github.com/opslevel/opslevel-go"

	"github.com/google/go-cmp/cmp"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
)

// importCmd represents the import command
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

	client := common.NewClient()

	services, servicesErr := common.QueryForServices(config)
	cobra.CheckErr(servicesErr)

	CacheLookupTables(client)

	for _, service := range services {
		foundService, needsUpdate := FindService(client, service)
		if foundService == nil {
			newService, newServiceErr := CreateService(client, service)
			if newServiceErr != nil {
				log.Error().Msgf("Failed creating service: '%s' \n\tREASON: %v", service.Name, newServiceErr.Error())
				continue
			} else {
				log.Info().Msgf("Created new service: '%s'", newService.Name)
			}
			foundService = newService
		}
		if needsUpdate {
			UpdateService(client, service, foundService)
		}
		AssignAliases(client, service, foundService)
		AssignTags(client, service, foundService)
		CreateTags(client, service, foundService)
		AssignTools(client, service, foundService)
		log.Info().Msgf("===> Finished processing data for service: '%s'", foundService.Name)
	}
	log.Info().Msg("Import Complete")
}

// TODO: Helpers probably shouldn't be exported
// Helpers

func FindService(client *opslevel.Client, registration common.ServiceRegistration) (*opslevel.Service, bool) {
	for _, alias := range registration.Aliases {
		foundService, err := client.GetServiceWithAlias(alias)
		if err == nil && foundService.Id != nil {
			log.Info().Msgf("Reconciling service '%s' found with alias '%s' ...", foundService.Name, alias)
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
		log.Error().Msgf("===> Failed updating service: '%s' \n\tREASON: %v", service.Name, updateServiceErr.Error())
	} else {
		if diff := cmp.Diff(service, updatedService); diff != "" {
			log.Info().Msgf("===> Updated Service '%s' - Diff:\n%s", service.Name, diff)
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
			log.Error().Msgf("===> Failed assigning alias '%s' to service: '%s' \n\tREASON: %v", alias, service.Name, err.Error())
		} else {
			log.Info().Msgf("===> Assigned alias '%s' to service: '%s'", alias, service.Name)
		}
	}
}

func AssignTags(client *opslevel.Client, registration common.ServiceRegistration, service *opslevel.Service) {
	for tagKey, tagValue := range registration.TagAssigns {
		_, err := client.AssignTagForId(service.Id, tagKey, tagValue)
		if err != nil {
			log.Error().Msgf("===> Failed assigning tag '%s = %s' to service: '%s' \n\tREASON: %v", tagKey, tagValue, service.Name, err.Error())
		} else {
			log.Info().Msgf("===> Ensured tag '%s = %s' assigned to service: '%s'", tagKey, tagValue, service.Name)
		}
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
			log.Error().Msgf("===> Failed creating tag '%s = %s' on service: '%s' \n\tREASON: %v", tagKey, tagValue, service.Name, err.Error())
		} else {
			log.Info().Msgf("===> Created tag '%s = %s' on service: '%s'", tagKey, tagValue, service.Name)
		}
	}
}

func AssignTools(client *opslevel.Client, registration common.ServiceRegistration, service *opslevel.Service) {
	for _, tool := range registration.Tools {
		if service.HasTool(tool.Category, tool.DisplayName, tool.Environment) {
			log.Debug().Msgf("===> Tool '{Category: %s, Environment: %s, Name: %s}' already exists on service: '%s' ... skipping", tool.Category, tool.Environment, tool.DisplayName, service.Name)
			continue
		}
		tool.ServiceId = service.Id
		_, err := client.CreateTool(tool)
		if err != nil {
			log.Error().Msgf("===> Failed assigning tool '{Category: %s, Environment: %s, Name: %s}' to service: '%s' \n\tREASON: %v", tool.Category, tool.Environment, tool.DisplayName, service.Name, err.Error())
		} else {
			log.Info().Msgf("===> Ensured tool '{Category: %s, Environment: %s, Name: %s}' assigned to service: '%s'", tool.Category, tool.Environment, tool.DisplayName, service.Name)
		}
	}
}
