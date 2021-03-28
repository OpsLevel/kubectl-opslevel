package cmd

import (
	"github.com/spf13/cobra"
)

var serviceCmd = &cobra.Command{
	Use:   "service",
	Short: "Tools for interacting with the service API",
	Long: `Tools for interacting with the service API`,
}

func init() {
	rootCmd.AddCommand(serviceCmd)
}
