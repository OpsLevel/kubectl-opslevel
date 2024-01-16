package cmd

import (
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/opslevel/kubectl-opslevel/common"
	opslevel_common "github.com/opslevel/opslevel-common/v2023"
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
			sampleCount, err = strconv.Atoi(args[0])
			if sampleCount < 0 {
				err = fmt.Errorf("SAMPLES_COUNT must be >= 0, got %d", sampleCount)
			}
			cobra.CheckErr(err)
		}

		config, err := LoadConfig()
		cobra.CheckErr(err)

		client := createOpslevelClient()
		common.SyncCache(client)
		queue := make(chan opslevel_jq_parser.ServiceRegistration, 1)
		common.SetupControllers(config, queue, 0)
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
	sampled := opslevel_common.GetSample[opslevel_jq_parser.ServiceRegistration](samples, *services)

	// Print
	if isTextOutput {
		fmt.Print("The following data was found in your Kubernetes cluster ...\n\n")
	}

	prettyJSON, err := json.MarshalIndent(sampled, "", "    ")
	cobra.CheckErr(err)
	fmt.Println(string(prettyJSON))

	if isTextOutput {
		var (
			servicesCount int    = len(*services)
			midText       string = "a random sample of "
		)
		if samples <= 0 || samples >= servicesCount {
			samples = servicesCount
			midText = ""
		}
		fmt.Printf("\nShowing %s%d / %d resources\n", midText, samples, servicesCount)
		fmt.Println("\nIf you're happy with the above data you can reconcile it with OpsLevel by running:\n\n OPSLEVEL_API_TOKEN=XXX kubectl opslevel service import\n\nOtherwise, please adjust the config file and rerun this command")
	}
}
