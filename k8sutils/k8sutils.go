package k8sutils

import (
	"context"
	"os"
	"encoding/json"

	"github.com/opslevel/kubectl-opslevel/config"

	"github.com/rs/zerolog/log"

	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/cache"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/workqueue"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

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

func QueryForServices(c *config.Config) ([]config.ServiceRegistration, error) {
	var err error
	var parser *config.ServiceRegistrationParser
	var services []config.ServiceRegistration
	k8sClient := CreateKubernetesClient()
	
	for _, importConfig := range c.Service.Import {
		parser, err = config.NewParser(importConfig.OpslevelConfig)
		if (err != nil) { return nil, err }
		listOptions := metav1.ListOptions{
			LabelSelector: importConfig.SelectorConfig.LabelSelector(),
		}
		// TODO: use different client based on importConfig.SelectorConfig.Kind
		deployments, deploymentsErr := k8sClient.AppsV1().Deployments(importConfig.SelectorConfig.Namespace).List(context.TODO(), listOptions)
		if (deploymentsErr != nil) { return nil, deploymentsErr }
		for _, resource := range deployments.Items {
			bytes, bytesErr := json.Marshal(resource)
			if (bytesErr != nil) { return nil, err }
			service, serviceErr := parser.Parse(bytes)
			if (serviceErr != nil) { return nil, err }
			services = append(services, *service)
		}
	}
	return services, nil
}
