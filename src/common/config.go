package common

import (
	"fmt"
	"gopkg.in/yaml.v3"

	"github.com/opslevel/kubectl-opslevel/k8sutils"

	"github.com/creasty/defaults"
	"github.com/spf13/viper"
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

var ConfigCurrentVersion = "1.2.0"

// Make sure we only use spaces inside of these samples
var ConfigSimple = []byte(`#Simple Opslevel CLI Config
version: "1.2.0"
service:
  import:
    - selector: # This limits what data we look at in Kubernetes
        apiVersion: apps/v1 # only supports resources found in 'kubectl api-resources --verbs="get,list"'
        kind: Deployment
        excludes: # filters out resources if any expression returns truthy
          - .metadata.namespace == "kube-system"
          - .metadata.annotations."opslevel.com/ignore"
      opslevel: # This is how you map your kubernetes data to opslevel service
        name: .metadata.name
        owner: .metadata.namespace
        aliases: # This are how we identify the services again during reconciliation - please make sure they are very unique
          - '"k8s:\(.metadata.name)-\(.metadata.namespace)"'
        tags:
          assign: # tag with the same key name but with a different value will be updated on the service
            - '{"imported": "kubectl-opslevel"}'
            - .metadata.labels
          create: # tag with the same key name but with a different value with be added to the service
            - '{"environment": .spec.template.metadata.labels.environment}'
  collect:
    - selector: # This limits what data we look at in Kubernetes
        apiVersion: apps/v1 # only supports resources found in 'kubectl api-resources --verbs="get,list"'
        kind: Deployment
        excludes: # filters out resources if any expression returns truthy
          - .metadata.namespace == "kube-system"
          - .metadata.annotations."opslevel.com/ignore"
`)

var ConfigSample = []byte(`#Sample Opslevel CLI Config
version: "1.2.0"
service:
  import:
    - selector: # This limits what data we look at in Kubernetes
        apiVersion: apps/v1 # only supports resources found in 'kubectl api-resources --verbs="get,list"'
        kind: Deployment
        excludes: # filters out resources if any expression returns truthy
          - .metadata.namespace == "kube-system"
          - .metadata.annotations."opslevel.com/ignore"
      opslevel: # This is how you map your kubernetes data to opslevel service
        name: .metadata.name
        description: .metadata.annotations."opslevel.com/description"
        owner: .metadata.annotations."opslevel.com/owner"
        lifecycle: .metadata.annotations."opslevel.com/lifecycle"
        tier: .metadata.annotations."opslevel.com/tier"
        product: .metadata.annotations."opslevel.com/product"
        language: .metadata.annotations."opslevel.com/language"
        framework: .metadata.annotations."opslevel.com/framework"
        #system: .metadata.annotations."opslevel.com/system"
        aliases: # This are how we identify the services again during reconciliation - please make sure they are very unique
          - '"k8s:\(.metadata.name)-\(.metadata.namespace)"'
          - '"\(.metadata.namespace)-\(.metadata.name)"'
        tags:
          assign: # tag with the same key name but with a different value will be updated on the service
            - '{"imported": "kubectl-opslevel"}'
            # find annoations with format: opslevel.com/tags.<key name>: <value>
            - '.metadata.annotations | to_entries |  map(select(.key | startswith("opslevel.com/tags"))) | map({(.key | split(".")[2]): .value})'
            - .metadata.labels
          create: # tag with the same key name but with a different value with be added to the service
            - '{"environment": .spec.template.metadata.labels.environment}'
        tools:
          - '{"category": "other", "environment": "production", "displayName": "my-cool-tool", "url": .metadata.annotations."example.com/my-cool-tool"} | if .url then . else empty end'
          # find annotations with format: opslevel.com/tools.<category>.<displayname>: <url>
          - '.metadata.annotations | to_entries |  map(select(.key | startswith("opslevel.com/tools"))) | map({"category": .key | split(".")[2], "displayName": .key | split(".")[3], "url": .value})'
          # OR find annotations with format: opslevel.com/tools.<category>.<environment>.<displayname>: <url>
          # - '.metadata.annotations | to_entries |  map(select(.key | startswith("opslevel.com/tools"))) | map({"category": .key | split(".")[2], "environment": .key | split(".")[3], "displayName": .key | split(".")[4], "url": .value})'
        repositories: # attach repositories to the service using the opslevel repo alias - IE github.com:hashicorp/vault
          - '{"name": "My Cool Repo", "directory": "/", "repo": .metadata.annotations.repo} | if .repo then . else empty end'
          # if just the alias is returned as a single string we'll build the name for you and set the directory to "/"
          - .metadata.annotations.repo
          # find annotations with format: opslevel.com/repo.<displayname>.<repo.subpath.dots.turned.to.forwardslash>: <opslevel repo alias>
          - '.metadata.annotations | to_entries |  map(select(.key | startswith("opslevel.com/repo"))) | map({"name": .key | split(".")[2], "directory": .key | split(".")[3:] | join("/"), "repo": .value})'
  collect:
    - selector: # This limits what data we look at in Kubernetes
        apiVersion: apps/v1 # only supports resources found in 'kubectl api-resources --verbs="get,list"'
        kind: Deployment
        excludes: # filters out resources if any expression returns truthy
          - .metadata.namespace == "kube-system"
          - .metadata.annotations."opslevel.com/ignore"
`)

func NewConfig() (*Config, error) {
	v := &ConfigVersion{}
	viper.Unmarshal(&v)
	if v.Version != ConfigCurrentVersion {
		return nil, fmt.Errorf("Supported config version is '%s' but found '%s' | Please update config file or create a new sample with `kubectl opslevel config sample`", ConfigCurrentVersion, v.Version)
	}

	c := &Config{}
	viper.Unmarshal(&c)
	if err := defaults.Set(c); err != nil {
		return c, err
	}
	return c, nil
}

func GetConfig(data []byte) (*Config, error) {
	var output Config
	unmarshalErr := yaml.Unmarshal(data, &output)
	if unmarshalErr != nil {
		return nil, unmarshalErr
	}
	return &output, nil
}
