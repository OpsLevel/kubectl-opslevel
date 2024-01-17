package cmd

import (
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/opslevel/kubectl-opslevel/common"
	opslevel_jq_parser "github.com/opslevel/opslevel-jq-parser/v2024"
	_ "github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
)

// previewCmd represents the preview command
var previewCmd = &cobra.Command{
	Use:   "preview [SAMPLES_COUNT]",
	Short: "Preview the data found in your Kubernetes cluster",
	Long: `This command will print out all the data it can find in your Kubernetes cluster based on the settings in the configuration file.
By default a randomly selected sample will be shown.
If the optional argument SAMPLES_COUNT=0 this will print out everything.`,
	Args:       cobra.MaximumNArgs(1),
	ArgAliases: []string{"SAMPLES_COUNT"},
	Run: func(cmd *cobra.Command, args []string) {
		var (
			err         error
			sampleCount int = 5
		)
		if len(args) == 1 {
			sampleCount, _ = strconv.Atoi(args[0])
		}

		config, err := LoadConfig()
		cobra.CheckErr(err)

		client := createOpslevelClient()
		common.SyncCache(client)
		queue := make(chan opslevel_jq_parser.ServiceRegistration, 1)
		common.SetupControllersSync(config, queue)
		PrintServices(IsTextOutput(), sampleCount, queue)
	},
}

func init() {
	serviceCmd.AddCommand(previewCmd)
}

func PrintServices(isTextOutput bool, samples int, queue <-chan opslevel_jq_parser.ServiceRegistration) {
	services := common.AggregateServices(queue)
	// Deduplicate ServiceRegistrations

	// Sample the data
	sampled := common.GetSample[opslevel_jq_parser.ServiceRegistration](samples, *services)

	// Print
	if isTextOutput {
		fmt.Print("The following data was found in your Kubernetes cluster ...\n\n")
	}

	prettyJSON, err := json.MarshalIndent(sampled, "", "    ")
	cobra.CheckErr(err)
	fmt.Println(string(prettyJSON))
	fmt.Println()

	if isTextOutput {
		var servicesCount int = len(*services)
		if samples <= 0 || samples >= servicesCount {
			fmt.Println("This is the full list of services detected in your cluster.")
		} else {
			fmt.Printf("This is randomly selected list of %d / %d services detected in your cluster.\n", samples, servicesCount)
			fmt.Println("If you want to see the full list of services detected, pass SAMPLES_COUNT=0.")
		}
		fmt.Println("\nIf you're happy with the above data you can reconcile it with OpsLevel by running:\n\n OPSLEVEL_API_TOKEN=XXX kubectl opslevel service import\n\nOtherwise, please adjust the config file and rerun this command")
	}
}
