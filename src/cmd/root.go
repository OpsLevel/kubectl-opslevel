package cmd

import (
	"os"
	"runtime"
	"strings"

	"github.com/spf13/cobra"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/spf13/viper"

	// https://github.com/golang/go/issues/33803
	"go.uber.org/automaxprocs/maxprocs"
)

var (
	cfgFile     string
	concurrency int
)

var rootCmd = &cobra.Command{
	Use:     "kubectl-opslevel",
	Aliases: []string{"kubectl opslevel"},
	Short:   "Opslevel Commandline Tools",
	Long:    `Opslevel Commandline Tools`,
}

func Execute() {
	cobra.CheckErr(rootCmd.Execute())
}

func init() {
	rootCmd.PersistentFlags().StringVarP(&cfgFile, "config", "c", "./opslevel-k8s.yaml", "")
	rootCmd.PersistentFlags().String("logFormat", "TEXT", "overrides environment variable 'OL_LOGFORMAT' (options [\"JSON\", \"TEXT\"])")
	rootCmd.PersistentFlags().String("logLevel", "INFO", "overrides environment variable 'OL_LOGLEVEL' (options [\"ERROR\", \"WARN\", \"INFO\", \"DEBUG\"])")
	rootCmd.PersistentFlags().String("api-token", "", "The OpsLevel API Token. Overrides environment variable 'OL_APITOKEN'")
	rootCmd.PersistentFlags().IntP("workers", "w", -1, "Sets the number of workers for API call processing. The default is == # CPU cores (cgroup aware). Overrides environment variable 'OL_WORKERS'")

	viper.BindPFlags(rootCmd.PersistentFlags())
	viper.BindPFlag("apitoken", rootCmd.PersistentFlags().Lookup("api-token"))
	cobra.OnInitialize(initConfig)
}

func initConfig() {
	readConfig()
	setupLogging()
	setupConcurrency()
}

func readConfig() {
	if cfgFile != "" {
		if cfgFile == "." {
			viper.SetConfigType("yaml")
			viper.ReadConfig(os.Stdin)
			return
		} else {
			viper.SetConfigFile(cfgFile)
		}
	} else {
		home, err := os.UserHomeDir()
		cobra.CheckErr(err)

		viper.SetConfigName("opslevel")
		viper.SetConfigType("yaml")
		viper.AddConfigPath(".")
		viper.AddConfigPath(home)
	}
	viper.SetEnvPrefix("OL")
	viper.AutomaticEnv()
	viper.ReadInConfig()
}

func setupLogging() {
	logFormat := strings.ToLower(viper.GetString("logFormat"))
	logLevel := strings.ToLower(viper.GetString("logLevel"))

	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix

	if logFormat == "text" {
		output := zerolog.ConsoleWriter{Out: os.Stderr}
		log.Logger = log.Output(output)
	}

	switch {
	case logLevel == "error":
		zerolog.SetGlobalLevel(zerolog.ErrorLevel)
	case logLevel == "warn":
		zerolog.SetGlobalLevel(zerolog.WarnLevel)
	case logLevel == "debug":
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
	default:
		zerolog.SetGlobalLevel(zerolog.InfoLevel)
	}
}

func setupConcurrency() {
	maxprocs.Set(maxprocs.Logger(log.Debug().Msgf))

	concurrency = viper.GetInt("workers")
	if concurrency <= 0 {
		concurrency = runtime.GOMAXPROCS(0)
	}

	log.Info().Msgf("Worker Concurrency == %v", concurrency)
}
