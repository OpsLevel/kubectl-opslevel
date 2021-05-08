package cmd

import (
	"github.com/spf13/cobra"
)

var serviceCmd = &cobra.Command{
	Use:   "service",
	Short: "Commands for interacting with the service API",
	Long:  `Commands for interacting with the service API`,
}

func init() {
	rootCmd.AddCommand(serviceCmd)
}
