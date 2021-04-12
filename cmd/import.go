package cmd

import (
	"fmt"
	"io/ioutil"

	_ "github.com/opslevel/kubectl-opslevel/config"
	_ "github.com/opslevel/kubectl-opslevel/k8sutils"
	"github.com/opslevel/kubectl-opslevel/opslevel"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	_ "github.com/rs/zerolog/log"
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
	// config, err := config.New()
	// cobra.CheckErr(err)

	apiToken := viper.GetString("apitoken")
	// exampleService := &opslevel.ServiceRegistration {
	// 	Name: "my-cool-service",
	// }
	client := opslevel.NewClient(apiToken)

	resp, respErr := client.Post(`{ "query": "mutation {  serviceCreate(input:{name:\"Brand New Service\",ownerAlias:\"infra\"}) { service { id name description owner { name alias } } errors { message path } } }" }`)
	defer resp.Body.Close()
	cobra.CheckErr(respErr)

    fmt.Println("response Status:", resp.Status)
    fmt.Println("response Headers:", resp.Header)
    body, _ := ioutil.ReadAll(resp.Body)
    fmt.Println("response Body:", string(body))
}
