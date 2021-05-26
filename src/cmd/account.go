package cmd

import (
	"fmt"

	"github.com/opslevel/kubectl-opslevel/common"
	"github.com/opslevel/opslevel-go"
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
	Run: func(cmd *cobra.Command, args []string) {
		client := common.NewClient()
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
		client := common.NewClient()
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
		client := common.NewClient()
		list, err := client.ListTeams()
		if err == nil {
			for _, item := range list {
				fmt.Println(item.Alias)
			}
		}
	},
}

var toolsCmd = &cobra.Command{
	Use:   "tools",
	Short: "Lists the valid alias for tools in your account",
	Long:  `Lists the valid alias for tools in your account`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println(opslevel.ToolCategoryAdmin)
		fmt.Println(opslevel.ToolCategoryCode)
		fmt.Println(opslevel.ToolCategoryContinuousIntegration)
		fmt.Println(opslevel.ToolCategoryDeployment)
		fmt.Println(opslevel.ToolCategoryErrors)
		fmt.Println(opslevel.ToolCategoryFeatureFlag)
		fmt.Println(opslevel.ToolCategoryHealthChecks)
		fmt.Println(opslevel.ToolCategoryIncidents)
		fmt.Println(opslevel.ToolCategoryLogs)
		fmt.Println(opslevel.ToolCategoryMetrics)
		fmt.Println(opslevel.ToolCategoryOrchestrator)
		fmt.Println(opslevel.ToolCategoryRunbooks)
		fmt.Println(opslevel.ToolCategoryStatusPage)
		fmt.Println(opslevel.ToolCategoryWiki)
		fmt.Println(opslevel.ToolCategoryOther)
	},
}

func init() {
	accountCmd.AddCommand(lifecycleCmd)
	accountCmd.AddCommand(tierCmd)
	accountCmd.AddCommand(teamCmd)
	accountCmd.AddCommand(toolsCmd)
	rootCmd.AddCommand(accountCmd)
}
