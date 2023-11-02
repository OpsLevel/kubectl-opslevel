package cmd

import (
	"encoding/json"
	"fmt"
	opslevel_common "github.com/opslevel/opslevel-common/v2023"
	opslevel_jq_parser "github.com/opslevel/opslevel-jq-parser/v2023"
	opslevel_k8s_controller "github.com/opslevel/opslevel-k8s-controller/v2023"
	"time"

	"github.com/opslevel/kubectl-opslevel/common"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
)

var (
	reconcileResyncInterval int
	reconcileBatchSize      int
)

var reconcileCmd = &cobra.Command{
	Use:   "reconcile",
	Short: "Run in the foreground as a kubernetes controller to reconcile data with service entries in OpsLevel",
	Long:  `Run in the foreground as a kubernetes controller to reconcile data with service entries in OpsLevel`,
	Run:   runReconcile,
}

func init() {
	serviceCmd.AddCommand(reconcileCmd)

	reconcileCmd.Flags().IntVar(&reconcileResyncInterval, "resync", 24, "The amount (in hours) before a full resync of the kubernetes cluster happens with OpsLevel. [default: 24]")
	reconcileCmd.Flags().IntVar(&reconcileBatchSize, "batch", 500, "The max amount of k8s resources to batch process with jq. Helps to speedup initial startup. [default: 500]")
}

func runReconcile(cmd *cobra.Command, args []string) {
	config, err := LoadConfig()
	cobra.CheckErr(err)

	common.SyncCache(createOpslevelClient())

	resync := time.Hour * time.Duration(reconcileResyncInterval)
	common.SyncCaches(createOpslevelClient(), resync)

	queue := make(chan opslevel_jq_parser.ServiceRegistration, 1)

	for i, importConfig := range config.Service.Import {
		selector := importConfig.SelectorConfig
		controller, err := opslevel_k8s_controller.NewK8SController(selector, resync, reconcileBatchSize, false)
		if err != nil {
			log.Error().Err(err).Msg("failed to create k8s controller")
			continue
		}
		callback := createHandler(fmt.Sprintf("service.import[%d]", i), importConfig, queue)
		controller.OnAdd = callback
		controller.OnUpdate = callback
		go controller.Start()
	}

	// Loop forever waiting to reconcile 1 service at a time
	go func() {
		reconciler := common.NewServiceReconciler(common.NewOpslevelClient(createOpslevelClient()))
		for registration := range queue {
			err := reconciler.Reconcile(registration)
			if err != nil {
				log.Error().Err(err).Msg("failed when reconciling service")
			}
		}
	}()

	opslevel_common.Run("Controller")
}

func createHandler(field string, config common.Import, queue chan opslevel_jq_parser.ServiceRegistration) func(items []interface{}) {
	id := fmt.Sprintf("%s/%s", config.SelectorConfig.ApiVersion, config.SelectorConfig.Kind)
	return func(items []interface{}) {
		var resources [][]byte
		for _, item := range items {
			data, _ := json.Marshal(item)
			resources = append(resources, data)
		}
		services, err := common.ProcessResources(field, config, resources)
		if err != nil {
			log.Error().Err(err)
			return
		}
		log.Info().Msgf("[%s] Processing '%d' service(s)", id, len(services))
		for _, service := range services {
			queue <- service
		}
	}
}
