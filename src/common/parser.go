package common

import (
	"encoding/json"
	"fmt"
	"github.com/opslevel/opslevel-jq-parser/v2023"
	"github.com/opslevel/opslevel-k8s-controller/v2023"
	"github.com/rs/zerolog/log"
	"time"
)

func NewParserHandler(config Import, queue chan opslevel_jq_parser.ServiceRegistration) func([]interface{}) {
	id := fmt.Sprintf("[%s/%s]", config.SelectorConfig.ApiVersion, config.SelectorConfig.Kind)

	parser := opslevel_jq_parser.NewJQServiceParser(config.OpslevelConfig)
	return func(items []interface{}) {
		for _, item := range items {
			data, err := json.Marshal(item)
			if err != nil {
				log.Error().Err(err).Msgf("%s - failed to marshal k8s resource", id)
				continue
			}
			registration, err := parser.Run(string(data))
			if err != nil {
				log.Error().Err(err).Msgf("%s - failed to parse k8s resource", id)
				continue
			}
			queue <- *registration
		}
	}
}

func AggregateServices(queue <-chan opslevel_jq_parser.ServiceRegistration) *[]opslevel_jq_parser.ServiceRegistration {
	services := make([]opslevel_jq_parser.ServiceRegistration, 0, 100)
	go func() {
		for service := range queue {
			services = append(services, service)
		}
	}()
	return &services
}

func SetupController(config *Config, queue chan opslevel_jq_parser.ServiceRegistration, resync time.Duration, batch int, runOnce bool) {
	for _, importConfig := range config.Service.Import {
		controller, err := opslevel_k8s_controller.NewK8SController(importConfig.SelectorConfig, resync, batch, runOnce)
		if err != nil {
			log.Error().Err(err).Msg("failed to create k8s controller")
			continue
		}
		callback := NewParserHandler(importConfig, queue)
		controller.OnAdd = callback
		controller.OnUpdate = callback
		controller.Start()
	}
}
