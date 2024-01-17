package cmd

import (
	"context"
	"time"

	opslevel_jq_parser "github.com/opslevel/opslevel-jq-parser/v2024"

	"github.com/opslevel/kubectl-opslevel/common"
	"github.com/spf13/cobra"
)

var reconcileResyncInterval int

var reconcileCmd = &cobra.Command{
	Use:   "reconcile",
	Short: "Run as a kubernetes controller to reconcile data in OpsLevel",
	Long:  `Run in the foreground as a kubernetes controller to reconcile data with service entries in OpsLevel`,
	Run: func(cmd *cobra.Command, args []string) {
		var (
			err    error
			resync time.Duration = time.Hour * time.Duration(reconcileResyncInterval)
		)
		config, err := LoadConfig()
		cobra.CheckErr(err)

		client := createOpslevelClient()
		common.SyncCache(client)
		common.SyncCaches(createOpslevelClient(), resync)
		queue := make(chan opslevel_jq_parser.ServiceRegistration, 1)
		reconcilerCtx, _ := context.WithCancel(context.Background()) // TODO: cancel on SIGINT/SIGTERM
		common.SetupControllers(config, queue, resync, reconcilerCtx)
		common.ReconcileServices(client, disableServiceCreation, queue)
	},
}

func init() {
	serviceCmd.AddCommand(reconcileCmd)
	reconcileCmd.Flags().IntVar(&reconcileResyncInterval, "resync", 24, "The amount (in hours) before a full resync of the kubernetes cluster happens with OpsLevel.")
}
