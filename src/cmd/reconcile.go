package cmd

import (
	"time"

	opslevel_common "github.com/opslevel/opslevel-common/v2023"
	opslevel_jq_parser "github.com/opslevel/opslevel-jq-parser/v2023"

	"github.com/opslevel/kubectl-opslevel/common"
	"github.com/spf13/cobra"
)

var reconcileResyncInterval int

var reconcileCmd = &cobra.Command{
	Use:   "reconcile",
	Short: "Run in the foreground as a kubernetes controller to reconcile data with service entries in OpsLevel",
	Long:  `Run in the foreground as a kubernetes controller to reconcile data with service entries in OpsLevel`,
	Run: func(cmd *cobra.Command, args []string) {
		config, err := LoadConfig()
		cobra.CheckErr(err)

		client := createOpslevelClient()
		common.SyncCache(client)
		resync := time.Hour * time.Duration(reconcileResyncInterval)
		common.SyncCaches(createOpslevelClient(), resync)
		queue := make(chan opslevel_jq_parser.ServiceRegistration, 1)
		common.SetupControllers(config, queue, resync)
		common.ReconcileServices(client, queue)
		opslevel_common.Run("Controller")
	},
}

func init() {
	serviceCmd.AddCommand(reconcileCmd)
	reconcileCmd.Flags().IntVar(&reconcileResyncInterval, "resync", 24, "The amount (in hours) before a full resync of the kubernetes cluster happens with OpsLevel. [default: 24]")
}
