package cmd

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"sort"
	"strconv"
	"time"

	"github.com/opslevel/kubectl-opslevel/common"
	"github.com/opslevel/kubectl-opslevel/config"

	_ "github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
)

// previewCmd represents the preview command
var previewCmd = &cobra.Command{
	Use:   "preview [SAMPLES_COUNT]",
	Short: "Preview the data found in your Kubernetes cluster returning SAMPLES_COUNT",
	Long: `This command will print out all the data it can find in your Kubernetes cluster based on the settings in the configuration file.
If SAMPLES_COUNT=0 this will print out everything.`,
	Run:        runPreview,
	Args:       cobra.MaximumNArgs(1),
	ArgAliases: []string{"samples"},
}

func init() {
	serviceCmd.AddCommand(previewCmd)
}

func runPreview(cmd *cobra.Command, args []string) {
	var samples = 5
	if len(args) > 0 {
		if parsedSamples, err := strconv.Atoi(args[0]); err == nil {
			samples = parsedSamples
		}
	}

	config, err := config.New()
	cobra.CheckErr(err)

	services, err2 := common.QueryForServices(config)
	cobra.CheckErr(err2)
	servicesCount := len(services)
	if samples < 1 {
		samples = servicesCount
	}

	fmt.Print("The following data was found in your Kubernetes cluster ...\n\n")
	if len(services) == 0 {
		fmt.Printf("[]\n")
	} else {
		prettyJSON, err := json.MarshalIndent(sample(services, samples), "", "    ")
		if err != nil {
			cobra.CheckErr(err)
		}
		fmt.Printf("%s\n", string(prettyJSON))
		if samples < servicesCount {
			fmt.Printf("\nShowing %v / %v resources\n", samples, servicesCount)
		}
	}

	fmt.Println("\nIf you're happy with the above data you can reconcile it with OpsLevel by running:\n\n OPSLEVEL_API_TOKEN=XXX kubectl opslevel service import\n\nOtherwise, please adjust the config file and rerun this command")
}

func sample(data []common.ServiceRegistration, samples int) []common.ServiceRegistration {
	max := len(data)
	if samples >= max {
		return data
	}
	output := make([]common.ServiceRegistration, samples)
	rand.Seed(time.Now().UTC().UnixNano())
	for i, index := range getSamples(0, max, samples) {
		output[i] = data[index]
	}
	return output
}

func getSamples(start int, end int, count int) []int {
	if end < start || (end-start) < count {
		return nil
	}
	nums := make([]int, 0)
	for len(nums) < count {
		num := rand.Intn((end - start)) + start
		exist := false
		for _, v := range nums {
			if v == num {
				exist = true
				break
			}
		}
		if !exist {
			nums = append(nums, num)
		}
	}
	sort.Ints(nums)
	return nums
}
