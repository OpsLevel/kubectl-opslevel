package config

import (
	"fmt"
	"bytes"
	"strings"
	"github.com/spf13/viper"
	"github.com/creasty/defaults"
)

var (
	ConfigFileName = "config.yaml"
)

type OpslevelKubernetesSelectorConfig struct {
	Kind string
	Namespace string
	Labels map[string]string
}

type Import struct {
    SelectorConfig OpslevelKubernetesSelectorConfig `mapstructure:"selector"`
    OpslevelConfig ServiceRegistrationConfig `mapstructure:"opslevel"`
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

func (selector *OpslevelKubernetesSelectorConfig) LabelSelector() string {
	var labels []string
    for key, value := range selector.Labels {
		labels = append(labels, fmt.Sprintf("%s=%s", key, value))
    }
    return strings.Join(labels, ",")
}