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
	Short: "Preview the service entries that will be created",
	Long:  `Preview the service entries that will be created`,
	Run:   runPreview,
}

func init() {
	serviceCmd.AddCommand(previewCmd)
}

func runPreview(cmd *cobra.Command, args []string) {
	config, err := config.New()
	cobra.CheckErr(err)

	services, err2 := common.QueryForServices(config)
	cobra.CheckErr(err2)

	prettyJSON, err := json.MarshalIndent(services, "", "    ")
	if err != nil {
		fmt.Printf("[]\n")
	}
	fmt.Printf("%s\n", string(prettyJSON))
}
