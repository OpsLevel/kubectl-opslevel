package cmd

import (
	"fmt"

	"github.com/opslevel/kubectl-opslevel/opslevel"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
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
	Run: func(cmd *cobra.Command, args []string) {
		client := opslevel.NewClient(viper.GetString("apitoken"))
		list, err := client.ListLifecycles()
		if err == nil {
			for _, item := range list {
				fmt.Println(item.Alias)
			}
		}
	},
}

var tierCmd = &cobra.Command{
	Use:   "tiers",
	Short: "Lists the valid alias for tiers in your account",
	Long:  `Lists the valid alias for tiers in your account`,
	Run: func(cmd *cobra.Command, args []string) {
		client := opslevel.NewClient(viper.GetString("apitoken"))
		list, err := client.ListTiers()
		if err == nil {
			for _, item := range list {
				fmt.Println(item.Alias)
			}
		}
	},
}

var teamCmd = &cobra.Command{
	Use:   "teams",
	Short: "Lists the valid alias for teams in your account",
	Long:  `Lists the valid alias for teams in your account`,
	Run: func(cmd *cobra.Command, args []string) {
		client := opslevel.NewClient(viper.GetString("apitoken"))
		list, err := client.ListTeams()
		if err == nil {
			for _, item := range list {
				fmt.Println(item.Alias)
			}
		}
	},
}

func init() {
	accountCmd.AddCommand(lifecycleCmd)
	accountCmd.AddCommand(tierCmd)
	accountCmd.AddCommand(teamCmd)
	rootCmd.AddCommand(accountCmd)

	// TODO: should this be a global flag?
	lifecycleCmd.Flags().String("api-token", "", "The OpsLevel API Token. Overrides environment variable 'OL_APITOKEN'")
	tierCmd.Flags().String("api-token", "", "The OpsLevel API Token. Overrides environment variable 'OL_APITOKEN'")
	teamCmd.Flags().String("api-token", "", "The OpsLevel API Token. Overrides environment variable 'OL_APITOKEN'")
}
