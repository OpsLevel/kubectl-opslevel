package cmd

import (
	"github.com/spf13/cobra"
)

var (
	resyncInterval int
	batchSize      int
)

var runCmd = &cobra.Command{
	Use:   "run",
	Short: "Commands for running a k8s controller",
}

func init() {
	rootCmd.AddCommand(runCmd)

	runCmd.PersistentFlags().IntVar(&resyncInterval, "resync", 24, "The amount (in hours) before a full resync of the kubernetes cluster happens with OpsLevel. [default: 24]")
	runCmd.PersistentFlags().IntVar(&batchSize, "batch", 500, "The max amount of k8s resources to batch process with jq. Helps to speedup initial startup. [default: 500]")
}
