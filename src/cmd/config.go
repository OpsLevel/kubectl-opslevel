package cmd

import (
	"encoding/json"
	"fmt"

	yaml "gopkg.in/yaml.v3"

	"github.com/alecthomas/jsonschema"
	"github.com/opslevel/kubectl-opslevel/pkg/config"
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
		conf := getCfgFile()
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
		fmt.Println(config.GetSample(viper.GetBool("simple")))
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
