package common

import (
	"github.com/opslevel/kubectl-opslevel/config"
	"github.com/opslevel/kubectl-opslevel/k8sutils"

	_ "github.com/rs/zerolog/log"
)



type ServiceRegistrationParser struct {
	Name JQParser
	Description JQParser
	Owner JQParser
	Lifecycle JQParser
	Tier JQParser
	Product JQParser
	Language JQParser
	Framework JQParser
	Aliases []JQParser
	Tags []JQParser
}

type ServiceRegistration struct {
	Name string
	Description string `json:",omitempty"`
	Owner string `json:",omitempty"`
	Lifecycle string `json:",omitempty"`
	Tier string `json:",omitempty"`
	Product string `json:",omitempty"`
	Language string `json:",omitempty"`
	Framework string `json:",omitempty"`
	Aliases []string `json:",omitempty"`
	Tags map[string]string `json:",omitempty"`
}

func NewParser(c config.ServiceRegistrationConfig) *ServiceRegistrationParser {
	parser := ServiceRegistrationParser{}
	parser.Name = NewJQParser(c.Name)
	parser.Description = NewJQParser(c.Description)
	parser.Owner = NewJQParser(c.Owner)
	parser.Lifecycle = NewJQParser(c.Lifecycle)
	parser.Tier = NewJQParser(c.Tier)
	parser.Product = NewJQParser(c.Product)
	parser.Language = NewJQParser(c.Language)
	parser.Framework = NewJQParser(c.Framework)
	// TODO: Aliases & Tags
	return &parser
}

func (parser *ServiceRegistrationParser) Parse(data []byte) *ServiceRegistration {
	service := ServiceRegistration{}
	service.Name = parser.Name.Parse(data)
	service.Description = parser.Description.Parse(data)
	service.Owner = parser.Owner.Parse(data)
	service.Lifecycle = parser.Lifecycle.Parse(data)
	service.Tier = parser.Tier.Parse(data)
	service.Product = parser.Product.Parse(data)
	service.Language = parser.Language.Parse(data)
	service.Framework = parser.Framework.Parse(data)
	// TODO: Aliases & Tags
	return &service
}

func QueryForServices(c *config.Config) ([]ServiceRegistration, error) {
	var parser *ServiceRegistrationParser
	var services []ServiceRegistration
	k8sClient := k8sutils.CreateKubernetesClient()
	
	for _, importConfig := range c.Service.Import {
		parser = NewParser(importConfig.OpslevelConfig)
		process := func(resource []byte) error {
			services = append(services, *parser.Parse(resource))
			return nil
		}
		if err := k8sClient.Query(importConfig.SelectorConfig, process); err != nil {
			return services, err
		}
	}
	return services, nil
}