package cmd

import (
	"github.com/opslevel/kubectl-opslevel/common"
	opslevel_jq_parser "github.com/opslevel/opslevel-jq-parser/v2023"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
)

var importCmd = &cobra.Command{
	Use:   "import",
	Short: "Create or Update service entries in OpsLevel",
	Long:  `This command will take the data found in your Kubernetes cluster and begin to reconcile it with OpsLevel`,
	Run: func(cmd *cobra.Command, args []string) {
		_, err := LoadConfig()
		cobra.CheckErr(err)

		common.SyncCache(createOpslevelClient())
		queue := make(chan opslevel_jq_parser.ServiceRegistration, 1)

		// Start K8S Poller / Controller sending data to channel in goroutine - closeing after sync is done

		// Reconcile - TODO: before this used to process this in parallel is that needed?
		reconciler := common.NewServiceReconciler(common.NewOpslevelClient(createOpslevelClient()))
		for registration := range queue {
			err := reconciler.Reconcile(registration)
			if err != nil {
				log.Error().Err(err).Msg("failed when reconciling service")
			}
		}

		log.Info().Msg("Import Complete")
	},
}

func init() {
	serviceCmd.AddCommand(importCmd)
}
