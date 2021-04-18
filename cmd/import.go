package cmd

import (
	"github.com/opslevel/kubectl-opslevel/config"
	"github.com/opslevel/kubectl-opslevel/common"
	"github.com/opslevel/kubectl-opslevel/opslevel"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/shurcooL/graphql"
	"github.com/rs/zerolog/log"
)

// importCmd represents the import command
var importCmd = &cobra.Command{
	Use:   "import",
	Short: "Create service entries from kubernetes data",
	Long: `Create service entries from kubernetes data`,
	Run: runImport,
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

	for _, service := range services {
		// fmt.Printf("Searching For: %s\n", service.Name)
		foundService, foundServiceErr := client.GetServiceWithAlias(service.Name)
		cobra.CheckErr(foundServiceErr)
		if (foundService.Id != nil) {
			// fmt.Printf("Found Existing Service: %s, %s\n", foundService.Name, foundService.Id)
			continue
		}
		// TODO: really sucks to have to transform these field the graphql types
		_, err := client.CreateService(opslevel.ServiceCreateInput{
			Name: graphql.String(service.Name),
			Product: graphql.String(service.Product),
			Description: graphql.String(service.Description),
			Languague: graphql.String(service.Language),
			Framework: graphql.String(service.Framework),
			// TODO: Tier
			// TODO: Owner
			// TODO: Lifecycle
		})
		cobra.CheckErr(err)
		// TODO: loop through service.Aliases and create them
		// TODO: loop through service.Tags and create them
	}
	log.Info().Msg("Import Complete")

	// Exploration

	// service, err := client.GetServiceWithAlias("coredns")
	// cobra.CheckErr(err)
	// fmt.Println(service.Id)
	// fmt.Println(service.Name)
	// fmt.Printf("%v", service.Owner)

	// service, err := client.CreateService(opslevel.ServiceCreateInput{
	// 	Name: "Bosun",
	// 	Owner: "DevOps",
	// })
	// cobra.CheckErr(err)
	// fmt.Println(service.Id)

	// service, err := client.UpdateService(opslevel.ServiceUpdateInput{
	// 	Id: "Z2lkOi8vb3BzbGV2ZWwvU2VydmljZS80MTQ1",
	// 	Name: "Vault Bosun",
	// 	Lifecycle: "pre-pre-alpha",
	// 	Tier: "tier_7",
	// })
	// cobra.CheckErr(err)
	// fmt.Println(service.Name)

	// team, teamErr := client.GetTeamWithId("Z2lkOi8vb3BzbGV2ZWwvVGVhbS83NzQ")
	// // team, teamErr := client.GetTeamWithAlias("infra")
	// cobra.CheckErr(teamErr)
	// fmt.Println(team.Id)
	// fmt.Println(team.Name)
	// fmt.Println(team.Manager.Email)

	// updateTeam, updateTeamErr := client.UpdateTeam(opslevel.TeamUpdateInput{
	// 	Alias: "infra",
	// 	Name: "DevOps",
	// })
	// cobra.CheckErr(updateTeamErr)
	// fmt.Println(updateTeam.Name)
	// fmt.Println(updateTeam.Manager.Email)

	// newTeam, newTeamErr := client.CreateTeam(opslevel.TeamCreateInput{
	// 	Name: "Devs",
	// })
	// cobra.CheckErr(newTeamErr)
	// fmt.Println("created")
	// fmt.Println(newTeam.Id)
	// fmt.Println(newTeam.Name)

	// teamId, teamAlias, err := client.DeleteTeam(opslevel.TeamDeleteInput{
	// 	Id: newTeam.Id,
	// })
	// cobra.CheckErr(err)
	// fmt.Println("deleted")
	// fmt.Println(teamId)
	// fmt.Println(teamAlias)


	// var mutation struct {
	// 	ServiceCreate ServiceCreatePayload `graphql:"serviceCreate(input: $input)"`
	// }

	// variables := map[string]interface{}{
	// 	"input": ServiceCreateInput{
	// 		Name: graphql.String("CoolService"),
	// 	},
	// }
	// //'{ "query": "mutation {  serviceCreate(input:{name:\"Brand New Service\"}) { service { id name } errors { message path } } }" }'
	// mutation := ServiceCreateMutation{}
	// err := client.Mutate(&mutation, variables)
	// if (err != nil) {
	// 	panic(err)
	// }
	// fmt.Println(mutation.ServiceCreate.Service.Id)
	// fmt.Println(mutation.ServiceCreate.Service.Name)

	// resp, respErr := client.Post(`{ "query": "mutation {  serviceCreate(input:{name:\"Brand New Service\",ownerAlias:\"infra\"}) { service { id name description owner { name alias } } errors { message path } } }" }`)
	// defer resp.Body.Close()
	// cobra.CheckErr(respErr)

    // fmt.Println("response Status:", resp.Status)
    // fmt.Println("response Headers:", resp.Header)
    // body, _ := ioutil.ReadAll(resp.Body)
    // fmt.Println("response Body:", string(body))
}
