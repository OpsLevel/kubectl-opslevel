package cmd

import (
	"fmt"
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
	apiToken     string
	apiTokenFile string
	cfgFile      string
	concurrency  int
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
	rootCmd.PersistentFlags().StringVar(&apiToken, "api-token", "", "The OpsLevel API Token. Overrides environment variable 'OL_APITOKEN' and the argument 'api-token-path'")
	rootCmd.PersistentFlags().StringVar(&apiTokenFile, "api-token-path", "", "Absolute path to a file containing the OpsLevel API Token. Overrides environment variable 'OL_APITOKEN'")
    rootCmd.PersistentFlags().String("apiurl", "https://api.opslevel.com/graphql", "The OpsLevel API Url. Overrides environment variable 'OPSLEVEL_API_URL'")
	rootCmd.PersistentFlags().IntP("workers", "w", -1, "Sets the number of workers for API call processing. The default is == # CPU cores (cgroup aware). Overrides environment variable 'OL_WORKERS'")

	viper.BindPFlags(rootCmd.PersistentFlags())
	viper.BindEnv("logFormat", "OPSLEVEL_LOG_FORMAT", "OL_LOG_FORMAT", "OL_LOGFORMAT")
	viper.BindEnv("logLevel", "OPSLEVEL_LOG_LEVEL", "OL_LOG_LEVEL", "OL_LOGLEVEL")
	viper.BindEnv("apiurl", "OPSLEVEL_API_URL", "OL_API_URL", "OL_APIURL")
	viper.BindEnv("apitoken", "OPSLEVEL_API_TOKEN", "OL_API_TOKEN", "OL_APITOKEN")
	viper.BindEnv("workers", "OPSLEVEL_WORKERS", "OL_WORKERS")
	cobra.OnInitialize(initConfig)
}

func initConfig() {
	readConfig()
	setupLogging()
	setupConcurrency()
	setupAPIToken()
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
	viper.SetEnvPrefix("OPSLEVEL")
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
}

// setupAPIToken evaluates several API token sources and sets the preferred token based on precedence.
//
// Precedence:
//   1. --api-token
//   2. --api-token-path
//   3. OL_APITOKEN
//
func setupAPIToken() {
	const key = "apitoken"

	if apiToken != "" {
		viper.Set(key, apiToken)
		return
	}

	if apiTokenFile == "" {
		return
	}

	b, err := os.ReadFile(apiTokenFile)
	cobra.CheckErr(fmt.Errorf("Failed to read provided api token file %s: %v", apiTokenFile, err))

	token := strings.TrimSpace(string(b))
	viper.Set(key, token)
}
