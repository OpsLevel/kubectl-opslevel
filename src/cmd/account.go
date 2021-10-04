package cmd

import (
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
)

var accountCmd = &cobra.Command{
	Use:   "account",
	Short: "Commands for interacting with the account API",
	Long:  `Commands for interacting with the account API`,
}

var lifecycleCmd = &cobra.Command{
	Use:   "lifecycles",
	Short: "Lists the valid alias for lifecycles in your account",
	Long:  `Lists the valid alias for lifecycles in your account`,
	Run:   movedToCLI,
}

var tierCmd = &cobra.Command{
	Use:   "tiers",
	Short: "Lists the valid alias for tiers in your account",
	Long:  `Lists the valid alias for tiers in your account`,
	Run:   movedToCLI,
}

var teamCmd = &cobra.Command{
	Use:   "teams",
	Short: "Lists the valid alias for teams in your account",
	Long:  `Lists the valid alias for teams in your account`,
	Run:   movedToCLI,
}

var toolsCmd = &cobra.Command{
	Use:   "tools",
	Short: "Lists the valid alias for tools in your account",
	Long:  `Lists the valid alias for tools in your account`,
	Run:   movedToCLI,
}

func movedToCLI(cmd *cobra.Command, args []string) {
	log.Error().Msg("This command has been moved to our CLI. https://www.opslevel.com/docs/api/cli/\nIt will be removed at a future date!")
}

func init() {
	accountCmd.AddCommand(lifecycleCmd)
	accountCmd.AddCommand(tierCmd)
	accountCmd.AddCommand(teamCmd)
	accountCmd.AddCommand(toolsCmd)
	rootCmd.AddCommand(accountCmd)
}
