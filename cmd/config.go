package cmd

import (
	"fmt"
	yaml "gopkg.in/yaml.v2"

	"github.com/spf13/cobra"
	"github.com/opslevel/kubectl-opslevel/config"
)

var configCmd = &cobra.Command{
	Use:   "config",
	Short: "Commands for working with the opslevel configuration",
	Long: "Commands for working with the opslevel configuration",
}

var configViewCmd = &cobra.Command{
	Use:   "view",
	Short: "Print the final configuration result",
	Long: "Print the final configuration after loading all the overrides and defaults",
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
	Long: "Print a sample config file which could be used",
	Run: func(cmd *cobra.Command, args []string) {
		conf, err := config.Default()
		cobra.CheckErr(err)
		output, err2 := yaml.Marshal(conf)
		cobra.CheckErr(err2)
		fmt.Println(string(output))
	},
}

func init() {
	rootCmd.AddCommand(configCmd)
	configCmd.AddCommand(
		configViewCmd,
		configSampleCmd,
	)
}