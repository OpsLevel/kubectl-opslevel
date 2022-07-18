package cmd

import (
	"bytes"
	"encoding/json"
	"fmt"

	yaml "gopkg.in/yaml.v3"

	"github.com/alecthomas/jsonschema"
	"github.com/opslevel/kubectl-opslevel/config"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// Make sure we only use spaces inside of these samples
var configSimple = []byte(`#Simple Opslevel CLI Config
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

var configSample = []byte(`#Sample Opslevel CLI Config
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
        aliases: # This are how we identify the services again during reconciliation - please make sure they are very unique
          - '"k8s:\(.metadata.name)-\(.metadata.namespace)"'
        tags:
          assign: # tag with the same key name but with a different value will be updated on the service
            - '{"imported": "kubectl-opslevel"}'
            # find annoations with format: opslevel.com/tags.<key name>: <value>
            - '.metadata.annotations | to_entries |  map(select(.key | startswith("opslevel.com/tags"))) | map({(.key | split(".")[2]): .value})'
            - .metadata.labels
          create: # tag with the same key name but with a different value with be added to the service
            - '{"environment": .spec.template.metadata.labels.environment}'
        tools:
          - '{"category": "other", "displayName": "my-cool-tool", "url": .metadata.annotations."example.com/my-cool-tool"} | if .url then . else empty end'
          # find annotations with format: opslevel.com/tools.<category>.<displayname>: <url> 
          - '.metadata.annotations | to_entries |  map(select(.key | startswith("opslevel.com/tools"))) | map({"category": .key | split(".")[2], "displayName": .key | split(".")[3], "url": .value})'
        repositories: # attach repositories to the service using the opslevel repo alias - IE github.com:hashicorp/vault
          - '{"name": "My Cool Repo", "directory": "/", "repo": .metadata.annotations.repo} | if .repo then . else empty end'
          # if just the alias is returned as a single string we'll build the name for you and set the directory to "/"
          - .metadata.annotations.repo
          # find annotations with format: opslevel.com/repo.<displayname>.<repo.subpath.dots.turned.to.forwardslash>: <opslevel repo alias> 
          - '.metadata.annotations | to_entries |  map(select(.key | startswith("opslevel.com/repos"))) | map({"name": .key | split(".")[2], "directory": .key | split(".")[3:] | join("/"), "repo": .value})'
  collect:
    - selector: # This limits what data we look at in Kubernetes
        apiVersion: apps/v1 # only supports resources found in 'kubectl api-resources --verbs="get,list"'
        kind: Deployment
        excludes: # filters out resources if any expression returns truthy
          - .metadata.namespace == "kube-system"
          - .metadata.annotations."opslevel.com/ignore"
`)

var configCmd = &cobra.Command{
	Use:   "config",
	Short: "Commands for working with the opslevel configuration",
	Long:  "Commands for working with the opslevel configuration",
}

var configSchemaCmd = &cobra.Command{
	Use:   "schema",
	Short: "Print the jsonschema for configuration file",
	Long:  "Print the jsonschema for configuration file",
	Run: func(cmd *cobra.Command, args []string) {
		schema := jsonschema.Reflect(&config.Config{})
		jsonBytes, jsonErr := json.MarshalIndent(schema, "", "  ")
		cobra.CheckErr(jsonErr)
		fmt.Println(string(jsonBytes))
	},
}

var configViewCmd = &cobra.Command{
	Use:   "view",
	Short: "Print the final configuration result",
	Long:  "Print the final configuration after loading all the overrides and defaults",
	Run: func(cmd *cobra.Command, args []string) {
		conf, err := config.New()
		cobra.CheckErr(err)
		output, err2 := yaml.Marshal(conf)
		cobra.CheckErr(err2)
		fmt.Println(string(output))
	},
}

var configSampleCmd = &cobra.Command{
	Use:   "sample",
	Short: "Print a sample config file",
	Long:  "Print a sample config file which could be used",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println(getSample(viper.GetBool("simple")))
	},
}

func init() {
	rootCmd.AddCommand(configCmd)
	configCmd.AddCommand(
		configSchemaCmd,
		configViewCmd,
		configSampleCmd,
	)

	configSampleCmd.Flags().Bool("simple", false, "Adjust the sample config to a bit simpler")
	viper.BindPFlags(configSampleCmd.Flags())
}

func getSample(simple bool) string {
	var sample []byte
	if simple == true {
		sample = configSimple
	} else {
		sample = configSample
	}
	// we use yaml unmarshal to prove that the samples are valid yaml
	var nodes yaml.Node
	unmarshalErr := yaml.Unmarshal(sample, &nodes)
	if unmarshalErr != nil {
		return unmarshalErr.Error()
	}
	var b bytes.Buffer
	yamlEncoder := yaml.NewEncoder(&b)
	yamlEncoder.SetIndent(2)
	encodeErr := yamlEncoder.Encode(&nodes)
	if encodeErr != nil {
		return encodeErr.Error()
	}
	return string(b.Bytes())
}
