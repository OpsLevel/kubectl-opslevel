package cmd

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/opslevel/kubectl-opslevel/common"

	"gopkg.in/yaml.v3"

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
		jsonBytes, err := json.MarshalIndent(schema, "", "    ")
		cobra.CheckErr(err)
		fmt.Println(string(jsonBytes))
	},
}

var configViewCmd = &cobra.Command{
	Use:   "view",
	Short: "Print the final configuration result",
	Long:  "Print the final configuration after loading all the overrides and defaults",
	Run: func(cmd *cobra.Command, args []string) {
		conf, err := LoadConfig()
		cobra.CheckErr(err)
		output, err := PrettyConfig(conf)
		cobra.CheckErr(err)
		fmt.Println(output)
	},
}

var configSampleCmd = &cobra.Command{
	Use:   "sample",
	Short: "Print a sample config file",
	Long:  "Print a sample config file which could be used",
	Run: func(cmd *cobra.Command, args []string) {
		var cfg *common.Config
		var err error
		if viper.GetBool("simple") {
			cfg, err = common.ParseConfig(common.ConfigSimple)
		} else {
			cfg, err = common.ParseConfig(common.ConfigSample)
		}
		pretty, err := PrettyConfig(cfg)
		cobra.CheckErr(err)
		fmt.Println(pretty)
	},
}

func PrettyConfig(cfg *common.Config) (string, error) {
	var buffer bytes.Buffer
	encoder := yaml.NewEncoder(&buffer)
	encoder.SetIndent(2)
	err := encoder.Encode(&cfg)
	if err != nil {
		return "", err
	}
	output := buffer.String()

	// wrap the version in quotes - assumes version is in format X.X.X (3 single digits)
	output = strings.Replace(output, "VER_STRING", common.ConfigCurrentVersion, 1)
	versionSubstring := output[9:14]
	output = strings.Replace(output, versionSubstring, fmt.Sprintf("\"%s\"", versionSubstring), 1)
	return output, err
}

func init() {
	rootCmd.AddCommand(configCmd)
	configCmd.AddCommand(configSchemaCmd, configViewCmd, configSampleCmd)

	configSampleCmd.Flags().Bool("simple", false, "Adjust the sample config to be less complex")
	err := viper.BindPFlags(configSampleCmd.Flags())
	cobra.CheckErr(err)
}

func readConfig() []byte {
	if cfgFile == "-" {
		buf := bytes.Buffer{}
		_, err := buf.ReadFrom(os.Stdin)
		if err != nil {
			panic(err)
		}
		return buf.Bytes()
	}
	res, err := os.ReadFile(cfgFile)
	if err != nil {
		// send a gentle error message instead of panicking if user forgot their config file
		if errors.As(os.ErrNotExist, &err) {
			fmt.Printf("config file not found: '%s' - try running:\n  kubectl opslevel config sample [--simple] > %s\n", cfgFile, cfgFile)
			fmt.Printf("make sure to edit and then preview it afterwards:\n  kubectl opslevel service preview\n")
			os.Exit(1)
		}
		panic(err)
	}
	return res
}

func LoadConfig() (*common.Config, error) {
	var (
		config      *common.Config
		configBytes []byte
		err         error
		help        = "Please update the config file or create a new one with a sample from `kubectl opslevel config sample`"
	)
	configBytes = readConfig()
	if len(configBytes) == 0 {
		return nil, fmt.Errorf("the config file is empty | %s", help)
	}
	config, err = common.ParseConfig(string(configBytes))
	if err != nil {
		return nil, err
	}
	if config.Version == "" {
		return nil, fmt.Errorf("could not parse version in the config file | %s", help)
	} else if config.Version != common.ConfigCurrentVersion {
		return nil, fmt.Errorf("supported config version is '%s' but found '%s' | %s",
			common.ConfigCurrentVersion, config.Version, help)
	}
	return config, nil
}
