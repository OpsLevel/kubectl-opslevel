package cmd

import (
	"github.com/spf13/cobra"
)

var accountCmd = &cobra.Command{
	Use:   "account",
	Short: "Commands for interacting with the account API",
	Long:  `Commands for interacting with the account API`,
}

func init() {
	rootCmd.AddCommand(accountCmd)
}
