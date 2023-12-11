package common

import (
	"encoding/json"
	"fmt"
	"github.com/opslevel/opslevel-go/v2023"
	"github.com/opslevel/opslevel-jq-parser/v2023"
	"github.com/opslevel/opslevel-k8s-controller/v2023"
	"github.com/rs/zerolog/log"
	"sync"
	"time"
)

func AggregateServices(queue <-chan opslevel_jq_parser.ServiceRegistration) *[]opslevel_jq_parser.ServiceRegistration {
	services := make([]opslevel_jq_parser.ServiceRegistration, 0, 100)
	for registration := range queue {
		services = append(services, registration)
	}
	return &services
}

func ReconcileServices(client *opslevel.Client, queue <-chan opslevel_jq_parser.ServiceRegistration) {
	reconciler := NewServiceReconciler(NewOpslevelClient(client))
	for registration := range queue {
		err := reconciler.Reconcile(registration)
		if err != nil {
			log.Error().Err(err).Msg("failed when reconciling service")
		}
	}
}

func NewParserHandler(config Import, queue chan<- opslevel_jq_parser.ServiceRegistration) func(interface{}) {
	id := fmt.Sprintf("[%s/%s]", config.SelectorConfig.ApiVersion, config.SelectorConfig.Kind)

	parser := opslevel_jq_parser.NewJQServiceParser(config.OpslevelConfig)
	return func(item interface{}) {
		data, err := json.Marshal(item)
		if err != nil {
			log.Error().Err(err).Msgf("%s - failed to marshal k8s resource", id)
			return
		}
		registration, err := parser.Run(string(data))
		if err != nil {
			log.Error().Err(err).Msgf("%s - failed to parse k8s resource", id)
			return
		}
		queue <- *registration
	}
}

func SetupControllers(config *Config, queue chan<- opslevel_jq_parser.ServiceRegistration, resync time.Duration) {
	go func() {
		var wg sync.WaitGroup
		for _, importConfig := range config.Service.Import {
			controller, err := opslevel_k8s_controller.NewK8SController(importConfig.SelectorConfig, resync)
			if err != nil {
				log.Error().Err(err).Msg("failed to create k8s controller")
				continue
			}
			callback := NewParserHandler(importConfig, queue)
			controller.OnAdd = callback
			controller.OnUpdate = callback
			controller.Start(&wg)
		}
		if resync <= 0 {
			wg.Wait()
			close(queue)
		}
	}()
}
