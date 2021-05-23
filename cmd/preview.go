package cmd

import (
	"encoding/json"
	"fmt"

	"github.com/opslevel/kubectl-opslevel/common"
	"github.com/opslevel/kubectl-opslevel/config"

	_ "github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
)

// previewCmd represents the preview command
var previewCmd = &cobra.Command{
	Use:   "preview",
	Short: "Preview the data found in your Kubernetes cluster",
	Long:  `This command will print out all the data it can find in your Kubernetes cluster based on the settings in the configuration file`,
	Run:   runPreview,
}

func init() {
	serviceCmd.AddCommand(previewCmd)
}

func runPreview(cmd *cobra.Command, args []string) {
	config, err := config.New()
	cobra.CheckErr(err)

	fmt.Println("The following data was found in your Kubernetes cluster ...")

	services, err2 := common.QueryForServices(config)
	cobra.CheckErr(err2)

	prettyJSON, err := json.MarshalIndent(services, "", "    ")
	if err != nil {
		fmt.Printf("[]\n")
	}
	fmt.Printf("%s\n", string(prettyJSON))

	fmt.Println("\nIf you're happy with the above data you can reconcile it with OpsLevel by running:\n\n OL_APITOKEN=XXX kubectl opslevel service import\n\nOtherwise, please adjust the config file and rerun this command")
}
