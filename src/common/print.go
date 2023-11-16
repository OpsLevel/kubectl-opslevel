package common

import (
	"encoding/json"
	"fmt"
	"github.com/opslevel/opslevel-common/v2023"
	"github.com/opslevel/opslevel-jq-parser/v2023"
	"github.com/spf13/cobra"
)

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
