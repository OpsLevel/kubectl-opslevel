package cmd

import (
	"context"

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

		queue := make(chan opslevel_jq_parser.ServiceRegistration, 1)
		ctx := common.InitSignalHandler(context.Background(), queue)
		client := createOpslevelClient()
		common.SyncCache(client)
		common.SetupControllers(ctx, config, queue, 0)
		common.ReconcileServices(client, disableServiceCreation, enableServiceNameUpdate, queue, 0)
		log.Info().Msg("Import Complete")
	},
}

func init() {
	serviceCmd.AddCommand(importCmd)
}
