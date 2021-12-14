package cmd

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/opslevel/kubectl-opslevel/common"
	"github.com/opslevel/kubectl-opslevel/config"
	"github.com/opslevel/kubectl-opslevel/jq"
	"github.com/opslevel/kubectl-opslevel/k8sutils"
	"github.com/opslevel/opslevel-go"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
)

var (
	resyncInterval int
	batchSize      int
)

var reconcileCmd = &cobra.Command{
	Use:   "reconcile",
	Short: "Run in the foreground as a kubernetes controller to reconcile data with service entries in OpsLevel",
	Long:  `Run in the foreground as a kubernetes controller to reconcile data with service entries in OpsLevel`,
	Run:   runReconcile,
}

func init() {
	serviceCmd.AddCommand(reconcileCmd)

	reconcileCmd.Flags().IntVar(&resyncInterval, "resync", 24, "The amount (in hours) before a full resync of the kubernetes cluster happens with OpsLevel. [default: 24]")
	reconcileCmd.Flags().IntVar(&batchSize, "batch", 500, "The max amount of k8s resources to batch process with jq. Helps to speedup initial startup. [default: 500]")
}

func runReconcile(cmd *cobra.Command, args []string) {
	config, configErr := config.New()
	cobra.CheckErr(configErr)

	jq.ValidateInstalled()

	k8sClient := k8sutils.CreateKubernetesClient()
	olClient := createOpslevelClient()

	opslevel.Cache.CacheTiers(olClient)
	opslevel.Cache.CacheLifecycles(olClient)
	opslevel.Cache.CacheTeams(olClient)

	resync := time.Hour * time.Duration(resyncInterval)
	reconcileQueue := make(chan common.ServiceRegistration, 1)

	for i, importConfig := range config.Service.Import {
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
		callback := createHandler(fmt.Sprintf("service.import[%d]", i), importConfig, reconcileQueue)
		controller := k8sutils.NewController(*gvr, resync, batchSize)
		controller.OnAdd = callback
		controller.OnUpdate = callback
		go controller.Start(1)
	}

	// Loop forever resyncing teams at resync interval
	ticker := time.NewTicker(resync)
	go func() {
		for {
			<-ticker.C
			olClient := createOpslevelClient()
			// has a mutex lock that will block TryGet in ReconcileService goroutine
			opslevel.Cache.CacheTiers(olClient)
			opslevel.Cache.CacheLifecycles(olClient)
			opslevel.Cache.CacheTeams(olClient)
		}
	}()

	// Loop forever waiting to reconcile 1 service at a time
	go func() {
		client := createOpslevelClient()
		for {
			for service := range reconcileQueue {
				common.ReconcileService(client, service)
			}
		}
	}()

	k8sutils.Start()
}

func createHandler(field string, config config.Import, queue chan common.ServiceRegistration) k8sutils.KubernetesControllerHandler {
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
