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

	tiers, tiersErr := GetTiers(client)
	cobra.CheckErr(tiersErr)
	lifecycles, lifecyclesErr := GetLifecycles(client)
	cobra.CheckErr(lifecyclesErr)

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
			Languague:   service.Language,
			Framework:   service.Framework,
			// TODO: Owner
		}
		if v, ok := tiers[service.Tier]; ok {
			serviceCreateInput.Tier = string(v.Alias)
		}
		if v, ok := lifecycles[service.Lifecycle]; ok {
			serviceCreateInput.Lifecycle = string(v.Alias)
		}

		newService, err := client.CreateService(serviceCreateInput)
		cobra.CheckErr(err)
		client.CreateAliases(newService.Id, service.Aliases)
		client.AssignTagsForId(newService.Id, service.Tags)
	}
	log.Info().Msg("Import Complete")
}

func GetTiers(client *opslevel.Client) (map[string]opslevel.Tier, error) {
	tiersList, tiersErr := client.ListTiers()
	if tiersErr != nil {
		return nil, tiersErr
	}
	tiers := make(map[string]opslevel.Tier)
	for _, tier := range tiersList {
		tiers[string(tier.Alias)] = tier
	}
	return tiers, nil
}

func GetLifecycles(client *opslevel.Client) (map[string]opslevel.Lifecycle, error) {
	lifecyclesList, lifecyclesErr := client.ListLifecycles()
	if lifecyclesErr != nil {
		return nil, lifecyclesErr
	}
	lifecycles := make(map[string]opslevel.Lifecycle)
	for _, lifecycle := range lifecyclesList {
		lifecycles[string(lifecycle.Alias)] = lifecycle
	}
	return lifecycles, nil
}
