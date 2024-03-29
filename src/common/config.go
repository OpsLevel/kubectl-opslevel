package common

import (
	_ "embed"

	"github.com/creasty/defaults"
	opslevel_jq_parser "github.com/opslevel/opslevel-jq-parser/v2024"
	opslevel_k8s_controller "github.com/opslevel/opslevel-k8s-controller/v2024"
	"gopkg.in/yaml.v3"
)

type Import struct {
	SelectorConfig opslevel_k8s_controller.K8SSelector          `yaml:"selector" json:"selector" mapstructure:"selector"`
	OpslevelConfig opslevel_jq_parser.ServiceRegistrationConfig `yaml:"opslevel" json:"opslevel" mapstructure:"opslevel"`
}

type Service struct {
	Import []Import `json:"import"`
}

type Config struct {
	Version string  `json:"version"`
	Service Service `json:"service"`
}

var ConfigCurrentVersion = "1.3.0"

//go:embed configs/config_sample.yaml
var ConfigSample string

//go:embed configs/config_simple.yaml
var ConfigSimple string

func ParseConfig(data string) (*Config, error) {
	var output Config
	if err := yaml.Unmarshal([]byte(data), &output); err != nil {
		return nil, err
	}
	if err := defaults.Set(&output); err != nil {
		return nil, err
	}
	return &output, nil
}
