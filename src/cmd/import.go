package cmd

import (
	"sync"

	"github.com/opslevel/kubectl-opslevel/common"
	"github.com/opslevel/kubectl-opslevel/config"
	"github.com/opslevel/kubectl-opslevel/jq"
	"github.com/opslevel/opslevel-go"

	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
)

var importCmd = &cobra.Command{
	Use:   "import",
	Short: "Create or Update service entries in OpsLevel",
	Long:  `This command will take the data found in your Kubernetes cluster and begin to reconcile it with OpsLevel`,
	Run:   runImport,
}

func init() {
	serviceCmd.AddCommand(importCmd)
}

func runImport(cmd *cobra.Command, args []string) {
	config, configErr := config.New()
	cobra.CheckErr(configErr)

	jq.ValidateInstalled()

	services, servicesErr := common.GetAllServices(config)
	cobra.CheckErr(servicesErr)

	common.GetOrCreateAliasCache().CacheAll(createOpslevelClient())

	log.Info().Msgf("Worker Concurrency == %v", concurrency)
	done := make(chan bool)
	queue := make(chan common.ServiceRegistration, concurrency)
	go createWorkerPool(concurrency, queue, done)
	go enqueue(services, queue)
	<-done
	log.Info().Msg("Import Complete")
}

// TODO: Helpers probably shouldn't be exported
// Helpers

func createWorkerPool(count int, queue chan common.ServiceRegistration, done chan<- bool) {
	var waitGroup sync.WaitGroup
	waitGroup.Add(count)
	for i := 0; i < count; i++ {
		go func(c *opslevel.Client, q chan common.ServiceRegistration, wg *sync.WaitGroup) {
			for data := range q {
				common.ReconcileService(c, data)
			}
			wg.Done()
		}(createOpslevelClient(), queue, &waitGroup)
	}
	waitGroup.Wait()
	done <- true
}

func enqueue(services []common.ServiceRegistration, queue chan common.ServiceRegistration) {
	for _, service := range services {
		queue <- service
	}
	close(queue)
}
