package cmd

import (
	"encoding/json"
	"fmt"
	"github.com/opslevel/kubectl-opslevel/common"
	opslevel_common "github.com/opslevel/opslevel-common/v2023"
	"github.com/opslevel/opslevel-jq-parser/v2023"
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

		config, err := LoadConfig()
		cobra.CheckErr(err)

		client := createOpslevelClient()
		common.SyncCache(client)
		queue := make(chan opslevel_jq_parser.ServiceRegistration, 1)
		common.SetupControllers(config, queue, 0)
		common.PrintServices(IsTextOutput(), sampleCount, queue)
	},
}

func init() {
	serviceCmd.AddCommand(previewCmd)
}

func PrintServices(isTextOutput bool, samples int, queue <-chan opslevel_jq_parser.ServiceRegistration) {
	services := AggregateServices(queue)
	// Deduplicate ServiceRegistrations

	// Sample the data
	sampled := opslevel_common.GetSample(samples, *services)

	// Print
	if isTextOutput {
		fmt.Print("The following data was found in your Kubernetes cluster ...\n\n")
	}

	prettyJSON, err := json.MarshalIndent(sampled, "", "    ")
	cobra.CheckErr(err)
	fmt.Printf("%s\n", string(prettyJSON))

	if isTextOutput {
		servicesCount := len(*services)
		if samples < servicesCount {
			if samples == 0 {
				samples = servicesCount
			}
			fmt.Printf("\nShowing %v / %v resources\n", samples, servicesCount)
		}
		fmt.Println("\nIf you're happy with the above data you can reconcile it with OpsLevel by running:\n\n OPSLEVEL_API_TOKEN=XXX kubectl opslevel service import\n\nOtherwise, please adjust the config file and rerun this command")
	}
}
