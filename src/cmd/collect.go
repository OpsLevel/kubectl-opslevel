package cmd

import (
	"encoding/json"
	"fmt"
	"sync"
	"time"

	opslevel_common "github.com/opslevel/opslevel-common/v2023"
	opslevel_k8s_controller "github.com/opslevel/opslevel-k8s-controller/v2023"

	"github.com/opslevel/kubectl-opslevel/common"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var collectResyncInterval int

var collectCmd = &cobra.Command{
	Use:   "collect",
	Short: "Acts as a kubernetes controller to collect resources for submission to OpsLevel as custom event check payloads",
	Long:  `Acts as a kubernetes controller to collect resources for submission to OpsLevel as custom event check payloads`,
	Run: func(cmd *cobra.Command, args []string) {
		config, err := LoadConfig()
		cobra.CheckErr(err)

		integrationUrl := viper.GetString("integration-url")
		if len(integrationUrl) <= 0 {
			cobra.CheckErr(fmt.Errorf("please specify --integration-url"))
		}

		resync := time.Hour * time.Duration(collectResyncInterval)
		queue := make(chan string, 1)
		var wg sync.WaitGroup
		for _, config := range config.Service.Collect {
			controller, err := opslevel_k8s_controller.NewK8SController(config.SelectorConfig, resync)
			if err != nil {
				log.Error().Err(err).Msg("failed to create k8s controller")
				continue
			}
			callback := createCollectHandler(config, queue)
			controller.OnAdd = callback
			controller.OnUpdate = callback
			go controller.Start(&wg)
		}
		PushPayloads(integrationUrl, queue)
		opslevel_common.Run("Controller")
	},
}

func init() {
	rootCmd.AddCommand(collectCmd)

	collectCmd.Flags().StringP("integration-url", "i", "", "OpsLevel integration url (OPSLEVEL_INTEGRATION_URL)")
	collectCmd.Flags().IntVar(&collectResyncInterval, "resync", 24, "The amount (in hours) before a full resync of the kubernetes cluster happens with OpsLevel. [default: 24]")

	err := viper.BindEnv("integration-url", "OPSLEVEL_INTEGRATION_URL")
	cobra.CheckErr(err)
}

func PushPayloads(integrationUrl string, queue <-chan string) {
	go func() {
		restClient := createRestClient()
		for {
			for payload := range queue {
				_, err := restClient.R().SetBody(payload).Post(integrationUrl)
				cobra.CheckErr(err)
			}
		}
	}()
}

func createCollectHandler(config common.Collect, queue chan string) func(item interface{}) {
	id := fmt.Sprintf("[%s/%s]", config.SelectorConfig.ApiVersion, config.SelectorConfig.Kind)
	return func(item interface{}) {
		data, err := json.Marshal(item)
		if err != nil {
			log.Error().Err(err).Msgf("%s - failed to marshal k8s resource", id)
			return
		}
		queue <- string(data)
	}
}
