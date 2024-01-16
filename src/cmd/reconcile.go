package cmd

import (
	"fmt"
	"time"

	opslevel_jq_parser "github.com/opslevel/opslevel-jq-parser/v2024"
	"github.com/rs/zerolog/log"

	"github.com/opslevel/kubectl-opslevel/common"
	"github.com/spf13/cobra"
)

var (
	reconcileResyncInterval        int
	reconcileResyncIntervalMinutes int
)

var reconcileCmd = &cobra.Command{
	Use:   "reconcile",
	Short: "Run as a  kubernetes controller to reconcile data in OpsLevel",
	Long:  `Run in the foreground as a kubernetes controller to reconcile data with service entries in OpsLevel`,
	Run: func(cmd *cobra.Command, args []string) {
		var (
			err    error
			resync time.Duration
		)
		if reconcileResyncIntervalMinutes == 0 {
			resync = time.Hour * time.Duration(reconcileResyncInterval)
		} else {
			log.Info().Msgf("Overriding resync interval to be %d minutes", reconcileResyncIntervalMinutes)
			resync = time.Minute * time.Duration(reconcileResyncIntervalMinutes)
		}
		if resync <= 0 {
			err = fmt.Errorf("resync interval must be >= 0, got %d", resync)
		}
		cobra.CheckErr(err)
		config, err := LoadConfig()
		cobra.CheckErr(err)

		client := createOpslevelClient()
		common.SyncCache(client)
		common.SyncCaches(createOpslevelClient(), resync)
		queue := make(chan opslevel_jq_parser.ServiceRegistration, 1)
		common.SetupControllers(config, queue, resync)
		common.ReconcileServices(client, disableServiceCreation, queue)
	},
}

func init() {
	serviceCmd.AddCommand(reconcileCmd)
	reconcileCmd.Flags().IntVar(&reconcileResyncInterval, "resync", 24, "The amount (in hours) before a full resync of the kubernetes cluster happens with OpsLevel.")
	reconcileCmd.Flags().IntVar(&reconcileResyncIntervalMinutes, "resync-minutes", 0, "Overrides --resync if set. Not recommended for use in production.")
}
