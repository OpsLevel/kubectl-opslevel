package config

import (
	"bytes"

	"github.com/opslevel/kubectl-opslevel/k8sutils"

	"github.com/creasty/defaults"
	"github.com/spf13/viper"
)

var (
	ConfigFileName = "config.yaml"
)

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
	Tags        []string // JQ expressions that return a single string or a map[string]string
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
	Service Service `json:"service"`
}

func New() (*Config, error) {
	c := &Config{}
	viper.Unmarshal(&c)
	if err := defaults.Set(c); err != nil {
		return c, err
	}
	return c, nil
}

func Default() (*Config, error) {
	c := &Config{}
	v := viper.New()
	v.SetConfigType("yaml")
	v.ReadConfig(bytes.NewBuffer([]byte(ConfigSample)))
	v.Unmarshal(&c)
	if err := defaults.Set(c); err != nil {
		return c, err
	}
	return c, nil
}
