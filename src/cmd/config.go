package cmd

import (
	"encoding/json"
	"fmt"
	"github.com/opslevel/kubectl-opslevel/common"

	yaml "gopkg.in/yaml.v3"

	"github.com/alecthomas/jsonschema"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

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
		schema := jsonschema.Reflect(&common.Config{})
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
		conf, err := common.NewConfig()
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
		var cfg *common.Config
		if viper.GetBool("simple") {
			cfg, _ = common.GetConfig(common.ConfigSimple)
		} else {
			cfg, _ = common.GetConfig(common.ConfigSample)
		}
		output, err := yaml.Marshal(cfg)
		cobra.CheckErr(err)
		fmt.Println(string(output))
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
