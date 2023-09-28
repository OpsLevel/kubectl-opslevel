package cmd

import (
	"sync"

	"github.com/opslevel/kubectl-opslevel/common"
	"github.com/opslevel/opslevel-go/v2023"

	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
)

var importCmd = &cobra.Command{
	Use:   "import",
	Short: "Create or Update service entries in OpsLevel",
	Long:  `This command will take the data found in your Kubernetes cluster and begin to reconcile it with OpsLevel`,
	Run: func(cmd *cobra.Command, args []string) {
		config, err := common.NewConfig()
		cobra.CheckErr(err)

		services, err := common.GetAllServices(config)
		cobra.CheckErr(err)

		client := createOpslevelClient()

		opslevel.Cache.CacheTiers(client)
		opslevel.Cache.CacheLifecycles(client)
		opslevel.Cache.CacheTeams(client)

		done := make(chan bool)
		queue := make(chan common.ServiceRegistration, concurrency)
		go createWorkerPool(concurrency, queue, done)
		go enqueue(services, queue)
		<-done
		log.Info().Msg("Import Complete")
	},
}

func init() {
	serviceCmd.AddCommand(importCmd)
}

func createWorkerPool(count int, queue chan common.ServiceRegistration, done chan<- bool) {
	log.Info().Msgf("Worker Concurrency == %v", count)
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
