package cmd

import (
	"os"
	"strings"
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

func Execute() {
	cobra.CheckErr(rootCmd.Execute())
}

func init() {
	rootCmd.PersistentFlags().StringVarP(&cfgFile, "config", "c", "", "config file (default \"./opslevel.yaml\")")
	rootCmd.PersistentFlags().StringVar(&logFormat, "logFormat", "JSON", "The Log Format. (options [\"JSON\", \"TEXT\"])")
	rootCmd.PersistentFlags().StringVar(&logLevel, "logLevel", "INFO", "The Log Level. (options [\"ERROR\", \"WARN\", \"INFO\", \"DEBUG\"])")
	cobra.OnInitialize(initConfig)
}

func initConfig() {
	if cfgFile != "" {
		viper.SetConfigFile(cfgFile)
	} else {
		home, err := homedir.Dir()
		cobra.CheckErr(err)

		viper.SetConfigName("opslevel")
		viper.SetConfigType("yaml")
		viper.AddConfigPath(".")
		viper.AddConfigPath(home)
	}

	viper.SetEnvPrefix("OL")
	viper.AutomaticEnv()
	viper.BindPFlags(rootCmd.Flags())

	if err := viper.ReadInConfig(); err != nil {
		panic(err)
	}

	logFormat = strings.ToLower(viper.GetString("logFormat"))
	logLevel = strings.ToLower(viper.GetString("logLevel"))

	setupLogging()
}

func setupLogging() {

	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix

	if logFormat == "text" {
		output := zerolog.ConsoleWriter{Out: os.Stderr}
		log.Logger = log.Output(output)

	}
	if logLevel == "error" {
		zerolog.SetGlobalLevel(zerolog.ErrorLevel)
	} else if logLevel == "warn" {
		zerolog.SetGlobalLevel(zerolog.WarnLevel)
	} else if logLevel == "debug" {
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
	} else {
		zerolog.SetGlobalLevel(zerolog.InfoLevel)
	}
}
