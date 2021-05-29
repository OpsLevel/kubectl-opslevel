package cmd

import (
	"fmt"

	yaml "gopkg.in/yaml.v2"

	"github.com/opslevel/kubectl-opslevel/config"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var configSimple = `#Simple Opslevel CLI Config
version: "1.0.0"
service:
  import:
  - selector:
      kind: deployment
    opslevel:
      name: .metadata.name
      owner: .metadata.namespace
      aliases:
      - '"k8s:\(.metadata.name)-\(.metadata.namespace)"'
      tags:
      - '{"imported": "kubectl-opslevel"}'
      - .spec.template.metadata.labels
`

var configSample = `#Sample Opslevel CLI Config
version: "1.0.0"
service:
  import:
  - selector:
      kind: deployment
      namespace: ""
      labels: {}
    opslevel:
      name: .metadata.name
      description: .metadata.annotations."opslevel.com/description"
      owner: .metadata.annotations."opslevel.com/owner"
      lifecycle: .metadata.annotations."opslevel.com/lifecycle"
      tier: .metadata.annotations."opslevel.com/tier"
      product: .metadata.annotations."opslevel.com/product"
      language: .metadata.annotations."opslevel.com/language"
      framework: .metadata.annotations."opslevel.com/framework"
      aliases:
      - '"k8s:\(.metadata.name)-\(.metadata.namespace)"'
      tags:
      - '{"imported": "kubectl-opslevel"}'
      - '.metadata.annotations | to_entries |  map(select(.key | startswith("opslevel.com/tags"))) | map({(.key | split(".")[2]): .value})'
      - .metadata.labels
      - .spec.template.metadata.labels
      tools:
      - '{"category": "other", "displayName": "my-cool-tool", "url": .metadata.annotations."example.com/my-cool-tool"} | if .url then . else empty end'
      - '.metadata.annotations | to_entries |  map(select(.key | startswith("opslevel.com/tools"))) | map({"category": .key | split(".")[2], "displayName": .key | split(".")[3], "url": .value})'
`

var configCmd = &cobra.Command{
	Use:   "config",
	Short: "Commands for working with the opslevel configuration",
	Long:  "Commands for working with the opslevel configuration",
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
		if viper.GetBool("simple") == true {
			fmt.Println(configSimple)
		} else {
			fmt.Println(configSample)
		}
	},
}

func init() {
	rootCmd.AddCommand(configCmd)
	configCmd.AddCommand(
		configViewCmd,
		configSampleCmd,
	)

	configSampleCmd.Flags().Bool("simple", false, "Adjust the sample config to a bit simpler")
	viper.BindPFlags(configSampleCmd.Flags())
}
