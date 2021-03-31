package controller

import (
	"os"

	"github.com/rs/zerolog/log"

	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/cache"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/workqueue"

	"github.com/go-logr/logr"
	"k8s.io/klog/v2"
)

type Controller struct {
	name      string
	clientset kubernetes.Interface
	queue     workqueue.RateLimitingInterface
	informer  cache.SharedIndexInformer
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

func CreateKubernetesClient() kubernetes.Interface {
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
	return client
}
