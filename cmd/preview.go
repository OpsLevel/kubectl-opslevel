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

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// previewCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// previewCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
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

	fmt.Printf("[\n    %v\n]\n", strings.Join(output, ", \n    "))
}
