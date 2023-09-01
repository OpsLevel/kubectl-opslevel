package cmd

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/opslevel/kubectl-opslevel/common"
	"github.com/opslevel/kubectl-opslevel/k8sutils"
	"github.com/opslevel/kubectl-opslevel/pkg/config"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var collectCmd = &cobra.Command{
	Use:   "collect",
	Short: "Acts as a kubernetes controller to collect resources for submission to OpsLevel as custom event check payloads",
	Long:  `Acts as a kubernetes controller to collect resources for submission to OpsLevel as custom event check payloads`,
	Run:   runCollect,
}

func init() {
	runCmd.AddCommand(collectCmd)

	collectCmd.Flags().StringP("integration-url", "i", "", "OpsLevel integration url (OPSLEVEL_INTEGRATION_URL)")

	viper.BindEnv("integration-url", "OPSLEVEL_INTEGRATION_URL")
}

func runCollect(cmd *cobra.Command, args []string) {
	config := getCfgFile()

	integrationUrl := viper.GetString("integration-url")
	if len(integrationUrl) <= 0 {
		cobra.CheckErr(fmt.Errorf("please specify --integration-url"))
	}

	k8sClient := k8sutils.CreateKubernetesClient()

	resync := time.Hour * time.Duration(resyncInterval)
	collectQueue := make(chan string, 1)

	for i, importConfig := range config.Service.Collect {
		selector := importConfig.SelectorConfig
		if selectorErr := selector.Validate(); selectorErr != nil {
			log.Fatal().Err(selectorErr)
			return
		}
		gvr, err := k8sClient.GetGVR(selector)
		if err != nil {
			log.Error().Err(err)
			continue
		}
		callback := createCollectHandler(fmt.Sprintf("service.import[%d]", i), importConfig, collectQueue)
		controller := k8sutils.NewController(*gvr, resync, batchSize)
		controller.OnAdd = callback
		controller.OnUpdate = callback
		go controller.Start(1)
	}

	// Loop forever waiting to collect 1 payload at a time
	go func() {
		restClient := createRestClient()
		for {
			for payload := range collectQueue {
				restClient.R().SetBody(payload).Post(integrationUrl)
			}
		}
	}()

	k8sutils.Start()
}

func createCollectHandler(field string, config config.Collect, queue chan string) k8sutils.KubernetesControllerHandler {
	id := fmt.Sprintf("%s/%s", config.SelectorConfig.ApiVersion, config.SelectorConfig.Kind)
	return func(items []interface{}) {
		var resources [][]byte
		for _, item := range items {
			data, _ := json.Marshal(item)
			resources = append(resources, data)
		}
		filtered := common.FilterResources(config.SelectorConfig, resources)
		log.Info().Msgf("[%s] Processing '%d' payload(s)", id, len(filtered))
		for _, payload := range filtered {
			queue <- string(payload)
		}
	}
}
