package k8sutils

import (
	"context"
	"fmt"
	"os"
	"strings"
	"encoding/json"

	"github.com/rs/zerolog/log"

	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/go-logr/logr"
	"k8s.io/klog/v2"
)

type KubernetesSelector struct {
	Kind string
	Namespace string
	Labels map[string]string
}

type ClientWrapper struct {
	client kubernetes.Interface
}

func getKubernetesConfig() (*rest.Config, error) {
	config, err := rest.InClusterConfig()
	if err != nil {
		configPath := os.Getenv("KUBECONFIG")
		if configPath == "" {
			configPath = os.Getenv("HOME") + "/.kube/config"
		}
		config2, err2 := clientcmd.BuildConfigFromFlags("", configPath)
		if err2 != nil {
			return nil, err2
		}
		return config2, nil
	}
	return config, nil
}

func CreateKubernetesClient() ClientWrapper {
	config, err := getKubernetesConfig()
	if err != nil {
		log.Fatal().Msgf("Unable to create a kubernetes client: %v", err)
	}

	client, err2 := kubernetes.NewForConfig(config)
	if err2 != nil {
		log.Fatal().Msgf("Unable to create a kubernetes client: %v", err)
	}
	// Supress k8s client-go
	klog.SetLogger(logr.Discard())
	return ClientWrapper{client: client}
}

func (c *ClientWrapper) Query(selector KubernetesSelector, handler func(resource []byte) error) error {
	listOptions := metav1.ListOptions{
		LabelSelector: selector.LabelSelector(),
	}
	// TODO: use different client based on selector.Kind
	deployments, deploymentsErr := c.client.AppsV1().Deployments(selector.Namespace).List(context.TODO(), listOptions)
	if (deploymentsErr != nil) { return deploymentsErr }
	for _, resource := range deployments.Items {
		bytes, bytesErr := json.Marshal(resource)
		if (bytesErr != nil) { return bytesErr }
		handlerErr := handler(bytes)
		if (handlerErr != nil) { return handlerErr }
	}
	return nil
}

func (selector *KubernetesSelector) LabelSelector() string {
	var labels []string
    for key, value := range selector.Labels {
		labels = append(labels, fmt.Sprintf("%s=%s", key, value))
    }
    return strings.Join(labels, ",")
}
