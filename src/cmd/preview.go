package cmd

import (
	"encoding/json"
	"fmt"
	"github.com/opslevel/kubectl-opslevel/common"
	_ "github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
	"strconv"
)

// previewCmd represents the preview command
var previewCmd = &cobra.Command{
	Use:   "preview [SAMPLES_COUNT]",
	Short: "Preview the data found in your Kubernetes cluster returning SAMPLES_COUNT",
	Long: `This command will print out all the data it can find in your Kubernetes cluster based on the settings in the configuration file.
If the optional argument SAMPLES_COUNT=0 this will print out everything.`,
	Args:       cobra.MaximumNArgs(1),
	ArgAliases: []string{"samples"},
	Run: func(cmd *cobra.Command, args []string) {
		sampleCount := 5
		if len(args) > 0 {
			if parsedSamples, err := strconv.Atoi(args[0]); err == nil {
				sampleCount = parsedSamples
			}
		}

		config, err := common.NewConfig()
		cobra.CheckErr(err)

		services, err := common.GetAllServices(config)
		cobra.CheckErr(err)

		if IsTextOutput() {
			fmt.Print("The following data was found in your Kubernetes cluster ...\n\n")
		}
		sampled := common.GetSample(sampleCount, services)
		prettyJSON, err := json.MarshalIndent(sampled, "", "    ")
		cobra.CheckErr(err)
		fmt.Printf("%s\n", string(prettyJSON))

		if IsTextOutput() {
			servicesCount := len(services)
			if sampleCount < servicesCount {
				fmt.Printf("\nShowing %v / %v resources\n", sampleCount, servicesCount)
			}
			fmt.Println("\nIf you're happy with the above data you can reconcile it with OpsLevel by running:\n\n OPSLEVEL_API_TOKEN=XXX kubectl opslevel service import\n\nOtherwise, please adjust the config file and rerun this command")
		}
	},
}

func init() {
	serviceCmd.AddCommand(previewCmd)
}
