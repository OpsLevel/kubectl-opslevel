package cmd

import (
	// "fmt"
	"os"
	"github.com/spf13/cobra"

	homedir "github.com/mitchellh/go-homedir"
	"github.com/spf13/viper"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

var cfgFile string
var logFormat string
var logLevel string

var rootCmd = &cobra.Command{
	Use:   "kubectl-opslevel",
	Short: "Opslevel Commandline Tools",
	Long: `Opslevel Commandline Tools`,
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	cobra.CheckErr(rootCmd.Execute())
}

func init() {
	rootCmd.PersistentFlags().StringVarP(&cfgFile, "config", "c", "", "config file (default \"./opslevel.yaml\")")
	rootCmd.PersistentFlags().StringVar(&logFormat, "logFormat", "JSON", "The Log Format. (options [\"JSON\", \"TEXT\"])")
	rootCmd.PersistentFlags().StringVar(&logLevel, "logLevel", "INFO", "The Log Level. (options [\"ERROR\", \"WARN\", \"INFO\", \"DEBUG\"])")
	cobra.OnInitialize(initConfig)
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	setupLogging()
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Find home directory.
		home, err := homedir.Dir()
		cobra.CheckErr(err)

		viper.SetConfigName("opslevel")
		viper.SetConfigType("yaml")
		viper.AddConfigPath(".")
		viper.AddConfigPath(home)
	}

	viper.SetEnvPrefix("OL")
	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err != nil {
		// TODO: log error
		//fmt.Fprintln(os.Stderr, "Using config file:", viper.ConfigFileUsed())
	}
}

func setupLogging() {

	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix

	if logFormat == "TEXT" {
		log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
	}
	if logLevel == "ERROR" {
		zerolog.SetGlobalLevel(zerolog.ErrorLevel)
	} else if logLevel == "WARN" {
		zerolog.SetGlobalLevel(zerolog.WarnLevel)
	} else if logLevel == "DEBUG" {
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
	} else {
		zerolog.SetGlobalLevel(zerolog.InfoLevel)
	}
}
