package cmd

import (
	"context"
	"github.com/opslevel/kubectl-opslevel/common"
	opslevel_jq_parser "github.com/opslevel/opslevel-jq-parser/v2024"
	_ "github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
	"strconv"
	"time"
)

// previewCmd represents the preview command
var previewCmd = &cobra.Command{
	Use:        "preview",
	Short:      "[PREVIEW MODE] Create or Update service entries in OpsLevel",
	Long:       "This is a preview of the service commands",
	Args:       cobra.MaximumNArgs(1),
	ArgAliases: []string{"resync interval (seconds)"},
	Run: func(cmd *cobra.Command, args []string) {
		var err error
		previewResync := 45
		if len(args) > 0 {
			previewResync, err = strconv.Atoi(args[0])
			if err != nil {
				panic(err)
			}
		}
		resync := time.Second * time.Duration(previewResync)
		config, err := LoadConfig()
		cobra.CheckErr(err)

		queue := make(chan opslevel_jq_parser.ServiceRegistration, 1)
		ctx := common.InitSignalHandler(context.Background(), queue)
		// TODO: pick up from here
		//client := createOpslevelClient()
		common.SetupControllers(ctx, config, queue, resync)
		common.ReconcileServices(client, disableServiceCreation, queue)
	},
}

func init() {
	serviceCmd.AddCommand(previewCmd)
}
