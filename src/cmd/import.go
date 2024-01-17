package cmd

import (
	"github.com/opslevel/kubectl-opslevel/common"
	opslevel_jq_parser "github.com/opslevel/opslevel-jq-parser/v2024"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
)

var importCmd = &cobra.Command{
	Use:   "import",
	Short: "Create or Update service entries in OpsLevel",
	Long:  `This command will take the data found in your Kubernetes cluster and begin to reconcile it with OpsLevel`,
	Run: func(cmd *cobra.Command, args []string) {
		config, err := LoadConfig()
		cobra.CheckErr(err)

		client := createOpslevelClient()
		common.SyncCache(client)
		queue := make(chan opslevel_jq_parser.ServiceRegistration, 1)
		common.SetupControllersSync(config, queue)
		common.ReconcileServices(client, disableServiceCreation, queue)
		log.Info().Msg("Import Complete")
	},
}

func init() {
	serviceCmd.AddCommand(importCmd)
}
