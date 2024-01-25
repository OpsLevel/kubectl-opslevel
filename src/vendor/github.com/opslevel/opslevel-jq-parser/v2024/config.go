package opslevel_jq_parser

import (
	"github.com/opslevel/opslevel-go/v2024"
	"gopkg.in/yaml.v3"
)

var SimpleConfig = `name: .metadata.name
owner: .metadata.namespace
aliases: # This are how we identify the services again during reconciliation - please make sure they are very unique
  - '"k8s:\(.metadata.name)-\(.metadata.namespace)"'
tags:
  assign: # tag with the same key name but with a different value will be updated on the service
    - '{"imported": "kubectl-opslevel"}'
    - .metadata.labels
  create: # tag with the same key name but with a different value with be added to the service
    - '{"environment": .spec.template.metadata.labels.environment}'
`

var SampleConfig = `name: .metadata.name
description: .metadata.annotations."opslevel.com/description"
owner: .metadata.annotations."opslevel.com/owner"
lifecycle: .metadata.annotations."opslevel.com/lifecycle"
tier: .metadata.annotations."opslevel.com/tier"
product: .metadata.annotations."opslevel.com/product"
language: .metadata.annotations."opslevel.com/language"
framework: .metadata.annotations."opslevel.com/framework"
system: .metadata.annotations."opslevel.com/system"
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
`

type TagRegistrationConfig struct {
	Assign []string `json:"assign" yaml:"assign"` // JQ expressions that return a single string or a map[string]string
	Create []string `json:"create" yaml:"create"` // JQ expressions that return a single string or a map[string]string
}

// ServiceRegistrationConfig represents the jq expressions configuration that can turn json data into a ServiceRegistration
type ServiceRegistrationConfig struct {
	Name         string                `json:"name" yaml:"name"`
	Description  string                `json:"description" yaml:"description"`
	Owner        string                `json:"owner" yaml:"owner"`
	Lifecycle    string                `json:"lifecycle" yaml:"lifecycle"`
	Tier         string                `json:"tier" yaml:"tier"`
	Product      string                `json:"product" yaml:"product"`
	Language     string                `json:"language" yaml:"language"`
	Framework    string                `json:"framework" yaml:"framework"`
	System       string                `json:"system" yaml:"system"`
	Aliases      []string              `json:"aliases" yaml:"aliases"` // JQ expressions that return a single string or a []string
	Tags         TagRegistrationConfig `json:"tags" yaml:"tags"`
	Tools        []string              `json:"tools" yaml:"tools"`               // JQ expressions that return a single map[string]string or a []map[string]string
	Repositories []string              `json:"repositories" yaml:"repositories"` // JQ expressions that return a single string or []string or map[string]string or a []map[string]string
}

// ServiceRegistration represents the parsed json data from a ServiceRegistrationConfig
type ServiceRegistration struct {
	Name         string                                  `json:",omitempty"`
	Description  string                                  `json:",omitempty"`
	Owner        string                                  `json:",omitempty"`
	Lifecycle    string                                  `json:",omitempty"`
	Tier         string                                  `json:",omitempty"`
	Product      string                                  `json:",omitempty"`
	Language     string                                  `json:",omitempty"`
	Framework    string                                  `json:",omitempty"`
	System       string                                  `json:",omitempty"`
	Aliases      []string                                `json:",omitempty"`
	TagAssigns   []opslevel.TagInput                     `json:",omitempty"`
	TagCreates   []opslevel.TagInput                     `json:",omitempty"`
	Tools        []opslevel.ToolCreateInput              `json:",omitempty"` // This is a concrete class so fields are validated during `service preview`
	Repositories []opslevel.ServiceRepositoryCreateInput `json:",omitempty"` // This is a concrete class so fields are validated during `service preview`
}

func NewServiceRegistrationConfig(data string) (*ServiceRegistrationConfig, error) {
	var output ServiceRegistrationConfig
	err := yaml.Unmarshal([]byte(data), &output)
	return &output, err
}
