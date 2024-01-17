package common

import (
	"context"
	"encoding/json"
	"fmt"
	"sync"
	"time"

	"github.com/opslevel/opslevel-go/v2024"
	opslevel_jq_parser "github.com/opslevel/opslevel-jq-parser/v2024"
	opslevel_k8s_controller "github.com/opslevel/opslevel-k8s-controller/v2024"
	"github.com/rs/zerolog/log"
)

func AggregateServices(queue <-chan opslevel_jq_parser.ServiceRegistration) *[]opslevel_jq_parser.ServiceRegistration {
	services := make([]opslevel_jq_parser.ServiceRegistration, 0, 100)
	for registration := range queue {
		services = append(services, registration)
	}
	return &services
}

func ReconcileServices(client *opslevel.Client, disableServiceCreation bool, queue <-chan opslevel_jq_parser.ServiceRegistration) {
	reconciler := NewServiceReconciler(NewOpslevelClient(client), disableServiceCreation)
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

func newController(importConfig Import, queue chan<- opslevel_jq_parser.ServiceRegistration, sync time.Duration) (*opslevel_k8s_controller.K8SController, error) {
	controller, err := opslevel_k8s_controller.NewK8SController(importConfig.SelectorConfig, sync)
	if err != nil {
		return nil, err
	}
	callback := NewParserHandler(importConfig, queue)
	controller.OnAdd = callback
	controller.OnUpdate = callback
	return controller, err
}

// SetupControllersSync returns the pointer to a K8sController that will run only once.
func SetupControllersSync(config *Config, queue chan<- opslevel_jq_parser.ServiceRegistration) {
	go func() {
		var wg sync.WaitGroup
		for _, importConfig := range config.Service.Import {
			controller, err := newController(importConfig, queue, 0)
			if err != nil {
				log.Error().Err(err).Msg("failed to create k8s controller")
				continue
			}
			wg.Add(1)
			controller.RunOnce(&wg)
		}
		wg.Wait()
		close(queue)
	}()
}

// SetupControllers returns the pointer to a K8sController that will run until the context is cancelled.
func SetupControllers(config *Config, queue chan<- opslevel_jq_parser.ServiceRegistration, resync time.Duration, ctx context.Context) {
	go func() {
		for _, importConfig := range config.Service.Import {
			controller, err := newController(importConfig, queue, resync)
			if err != nil {
				log.Error().Err(err).Msg("failed to create k8s controller")
				continue
			}
			controller.Run(ctx)
		}
	}()
}
