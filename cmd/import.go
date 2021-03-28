package cmd

import (
	"github.com/spf13/cobra"
	"github.com/rs/zerolog/log"
)

// importCmd represents the import command
var importCmd = &cobra.Command{
	Use:   "import",
	Short: "Create service entries from kubernetes data",
	Long: `Create service entries from kubernetes data`,
	Run: runImport,
}

func init() {
	serviceCmd.AddCommand(importCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// importCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// importCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

func runImport(cmd *cobra.Command, args []string) {
	log.Info().Msgf("import called")
}
