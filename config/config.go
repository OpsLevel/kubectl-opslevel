package config

import (
	"bytes"
	"github.com/spf13/viper"
	"github.com/creasty/defaults"
)

var (
	ConfigFileName = "config.yaml"
)

type OpslevelServiceRegistration struct {
	Name string `default:".metadata.name"`
	Description string
	Owner string
	Lifecycle string
	Tier string
	Product string
	Language string
	Framework string
	Aliases []string
	Tags []string
}

type OpslevelKubernetesSelector struct {
	Kind string
	Namespace string
	Labels map[string]string
}

type Import struct {
    Selector OpslevelKubernetesSelector
    Opslevel OpslevelServiceRegistration
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