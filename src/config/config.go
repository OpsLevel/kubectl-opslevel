package config

import (
	"errors"
	"fmt"

	"github.com/opslevel/kubectl-opslevel/k8sutils"

	"github.com/creasty/defaults"
	"github.com/spf13/viper"
)

var (
	ConfigCurrentVersion = "1.0.0"
)

type TagRegistrationConfig struct {
	Assign []string // JQ expressions that return a single string or a map[string]string
	Create []string // JQ expressions that return a single string or a map[string]string
}

type ServiceRegistrationConfig struct {
	Name        string `default:".metadata.name"`
	Description string
	Owner       string
	Lifecycle   string
	Tier        string
	Product     string
	Language    string
	Framework   string
	Aliases     []string // JQ expressions that return a single string or a string[]
	Tags        TagRegistrationConfig
	Tools       []string // JQ expressions that return a single map[string]string or a []map[string]string
}

type Import struct {
	SelectorConfig k8sutils.KubernetesSelector `yaml:"selector" json:"selector" mapstructure:"selector"`
	OpslevelConfig ServiceRegistrationConfig   `yaml:"opslevel" json:"opslevel" mapstructure:"opslevel"`
}

type Service struct {
	Import []Import `json:"import"`
}

type Config struct {
	Version string
	Service Service `json:"service"`
}

type ConfigVersion struct {
	Version string
}

func New() (*Config, error) {
	v := &ConfigVersion{}
	viper.Unmarshal(&v)
	if v.Version != ConfigCurrentVersion {
		return nil, errors.New(fmt.Sprintf("Supported config version is '%s' but found '%s' | Please update config file or create a new sample with `kubectl opslevel config sample`", ConfigCurrentVersion, v.Version))
	}

	c := &Config{}
	viper.Unmarshal(&c)
	if err := defaults.Set(c); err != nil {
		return c, err
	}
	return c, nil
}
