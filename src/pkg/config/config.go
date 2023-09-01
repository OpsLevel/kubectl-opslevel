package config

import (
	"errors"
	"fmt"
	"gopkg.in/yaml.v3"

	"github.com/opslevel/kubectl-opslevel/k8sutils"

	"github.com/creasty/defaults"
)

var (
	ConfigCurrentVersion = "1.3.0"
)

type TagRegistrationConfig struct {
	Assign []string `json:"assign"` // JQ expressions that return a single string or a map[string]string
	Create []string `json:"create"` // JQ expressions that return a single string or a map[string]string
}

type ServiceRegistrationConfig struct {
	Name         string                `json:"name"`
	Description  string                `json:"description"`
	Owner        string                `json:"owner"`
	Lifecycle    string                `json:"lifecycle"`
	Tier         string                `json:"tier"`
	Product      string                `json:"product"`
	Language     string                `json:"language"`
	Framework    string                `json:"framework"`
	System       string                `json:"system"`
	Aliases      []string              `json:"aliases"` // JQ expressions that return a single string or a []string
	Tags         TagRegistrationConfig `json:"tags"`
	Tools        []string              `json:"tools"`        // JQ expressions that return a single map[string]string or a []map[string]string
	Repositories []string              `json:"repositories"` // JQ expressions that return a single string or []string or map[string]string or a []map[string]string
}

type Import struct {
	SelectorConfig k8sutils.KubernetesSelector `yaml:"selector" json:"selector" mapstructure:"selector"`
	OpslevelConfig ServiceRegistrationConfig   `yaml:"opslevel" json:"opslevel" mapstructure:"opslevel"`
}

type Collect struct {
	SelectorConfig k8sutils.KubernetesSelector `yaml:"selector" json:"selector" mapstructure:"selector"`
}

type Service struct {
	Import  []Import  `json:"import"`
	Collect []Collect `json:"collect"`
}

type Config struct {
	Version string  `json:"version"`
	Service Service `json:"service"`
}

type ConfigVersion struct {
	Version string
}

func NewConfig(data []byte) (*Config, error) {
	v := &ConfigVersion{}
	if err := yaml.Unmarshal(data, &v); err != nil {
		return nil, err
	}

	if v.Version != ConfigCurrentVersion {
		return nil, errors.New(fmt.Sprintf("Supported config version is '%s' but found '%s' | Please update config file or create a new sample with `kubectl opslevel config sample`", ConfigCurrentVersion, v.Version))
	}

	c := &Config{}
	if err := yaml.Unmarshal(data, &c); err != nil {
		return nil, err
	}
	if err := defaults.Set(c); err != nil {
		return c, err
	}
	return c, nil
}
