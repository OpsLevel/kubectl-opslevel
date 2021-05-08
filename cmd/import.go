package cmd

import (
	"github.com/opslevel/kubectl-opslevel/common"
	"github.com/opslevel/kubectl-opslevel/config"
	"github.com/opslevel/kubectl-opslevel/opslevel"

	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// importCmd represents the import command
var importCmd = &cobra.Command{
	Use:   "import",
	Short: "Create service entries from kubernetes data",
	Long:  `Create service entries from kubernetes data`,
	Run:   runImport,
}

func init() {
	serviceCmd.AddCommand(importCmd)

	// TODO: should this be a global flag?
	importCmd.Flags().String("api-token", "", "The OpsLevel API Token. Overrides environment variable 'OL_APITOKEN'")
}

func runImport(cmd *cobra.Command, args []string) {
	config, configErr := config.New()
	cobra.CheckErr(configErr)

	client := opslevel.NewClient(viper.GetString("apitoken"))

	services, servicesErr := common.QueryForServices(config)
	cobra.CheckErr(servicesErr)

	tiers, _ := GetTiers(client)
	lifecycles, _ := GetLifecycles(client)
	teams, _ := GetTeams(client)

	for _, service := range services {
		// fmt.Printf("Searching For: %s\n", service.Name)
		foundService, foundServiceErr := client.GetServiceWithAlias(service.Name)
		cobra.CheckErr(foundServiceErr)
		if foundService.Id != nil {
			// fmt.Printf("Found Existing Service: %s, %s\n", foundService.Name, foundService.Id)
			continue
		}
		serviceCreateInput := opslevel.ServiceCreateInput{
			Name:        service.Name,
			Product:     service.Product,
			Description: service.Description,
			Language:    service.Language,
			Framework:   service.Framework,
		}
		if v, ok := tiers[service.Tier]; ok {
			serviceCreateInput.Tier = string(v.Alias)
		}
		if v, ok := lifecycles[service.Lifecycle]; ok {
			serviceCreateInput.Lifecycle = string(v.Alias)
		}
		if v, ok := teams[service.Owner]; ok {
			serviceCreateInput.Owner = string(v.Alias)
		}

		newService, err := client.CreateService(serviceCreateInput)
		cobra.CheckErr(err)
		client.CreateAliases(newService.Id, service.Aliases)
		client.AssignTagsForId(newService.Id, service.Tags)
		for _, tool := range service.Tools {
			tool.ServiceId = newService.Id
			if _, createToolErr := client.CreateTool(tool); createToolErr != nil {
				cobra.CheckErr(createToolErr)
				break
			}
		}
	}
	log.Info().Msg("Import Complete")
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
