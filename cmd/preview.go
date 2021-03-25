package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// previewCmd represents the preview command
var previewCmd = &cobra.Command{
	Use:   "preview",
	Short: "Preview the service entries that will be created",
	Long: `Preview the service entries that will be created`,
	Run: runPreview,
}

func init() {
	serviceCmd.AddCommand(previewCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// previewCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// previewCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

func runPreview(cmd *cobra.Command, args []string) {
	fmt.Println("preview called")
}
