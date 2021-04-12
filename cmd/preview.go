package cmd

import (
	"fmt"
	"encoding/json"
	"strings"

	"github.com/opslevel/kubectl-opslevel/config"
	"github.com/opslevel/kubectl-opslevel/k8sutils"

	_ "github.com/rs/zerolog/log"
	"github.com/spf13/cobra"

	
)

// previewCmd represents the preview command
var previewCmd = &cobra.Command{
	Use:   "preview",
	Short: "Preview the service entries that will be created",
	Long: `Preview the service entries that will be created`,
	Run: runPreview,
}

func init() {
	serviceCmd.AddCommand(previewCmd)
}

func runPreview(cmd *cobra.Command, args []string) {
	config, err := config.New()
	cobra.CheckErr(err)

	services, err2 := k8sutils.QueryForServices(config)
	cobra.CheckErr(err2)

	var output []string
	for _, service := range services {
		serviceBytes, serviceBytesErr := json.Marshal(service)
		cobra.CheckErr(serviceBytesErr)
		output = append(output, string(serviceBytes))
	}
	// TODO: is there a better way to format this?
	fmt.Printf("[\n    %v\n]\n", strings.Join(output, ", \n    "))
}
