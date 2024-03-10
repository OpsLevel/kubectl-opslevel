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
		logLocal := log.With().Str("where", "parser_handler_func").Str("id", id).Logger()
		data, err := json.Marshal(item)
		if err != nil {
			logLocal.Error().Err(err).Msg("failed to marshal k8s resource")
			return
		}
		registration, err := parser.Run(string(data))
		if err != nil {
			logLocal.Error().Err(err).Msg("failed to parse k8s resource")
			return
		}
		queue <- *registration
	}
}

func SetupControllers(ctx context.Context, config *Config, queue chan<- opslevel_jq_parser.ServiceRegistration, resync time.Duration) {
	go func() {
		var wg *sync.WaitGroup
		if resync <= 0 {
			wg = &sync.WaitGroup{}
		}
		for _, importConfig := range config.Service.Import {
			logger := log.Logger.With().Str("from", "opslevel-k8s-controller").Logger()
			controller, err := opslevel_k8s_controller.NewK8SController(&logger, importConfig.SelectorConfig, resync)
			if err != nil {
				log.Error().Err(err).Msg("failed to create k8s controller")
				continue
			}
			callback := NewParserHandler(importConfig, queue)
			controller.OnAdd = callback
			controller.OnUpdate = callback
			if wg != nil {
				wg.Add(1)
			}
			controller.Start(ctx, wg)
		}
		if resync <= 0 {
			wg.Wait()
			close(queue)
		}
	}()
}
