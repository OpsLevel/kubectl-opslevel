package cmd

import (
	"context"

	"github.com/opslevel/kubectl-opslevel/config"
	"github.com/opslevel/kubectl-opslevel/controller"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
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
	_, err := config.New()
	cobra.CheckErr(err)

	client := controller.CreateKubernetesClient()

    listOptions := metav1.ListOptions{
		LabelSelector: "k8s-app=kube-dns",
		FieldSelector: "metadata.name=coredns",
    }

	deployments, err2 := client.AppsV1().Deployments("").List(context.TODO(), listOptions)
	cobra.CheckErr(err2)

	for i, deployment := range deployments.Items {
		log.Info().Msgf("%d = %s\n", i, deployment.ObjectMeta.Name)
	}
}
