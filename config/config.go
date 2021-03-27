package config

import (
	"bytes"
	"github.com/spf13/viper"
	"github.com/creasty/defaults"
)

var (
	ConfigFileName = "config.yaml"
)

type Import struct {
    Kind string
	Namespace string 
	OLServiceName string `default:"{$.metadata.name}" mapstructure:"ol_service_name" yaml:"ol_service_name"`
	OLAlias string `mapstructure:"ol_alias" yaml:"ol_alias"`
	OLProduct string `mapstructure:"ol_product" yaml:"ol_product"`
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