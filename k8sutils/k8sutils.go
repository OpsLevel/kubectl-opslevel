package k8sutils

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/go-logr/logr"
	"github.com/rs/zerolog/log"
	appsv1 "k8s.io/api/apps/v1"
	batchv1 "k8s.io/api/batch/v1"
	batchv1beta1 "k8s.io/api/batch/v1beta1"
	v1 "k8s.io/api/core/v1"
	networkingv1 "k8s.io/api/networking/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"

	// This is here because of https://github.com/OpsLevel/kubectl-opslevel/issues/24
	_ "k8s.io/client-go/plugin/pkg/client/auth"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/klog/v2"
)

type KubernetesSelector struct {
	Kind      string
	Namespace string
	Labels    map[string]string
}

type ClientWrapper struct {
	client kubernetes.Interface
}

func getKubernetesConfig() (*rest.Config, error) {
	loadingRules := clientcmd.NewDefaultClientConfigLoadingRules()
	configOverrides := &clientcmd.ConfigOverrides{}
	config, err := clientcmd.NewNonInteractiveDeferredLoadingClientConfig(loadingRules, configOverrides).ClientConfig()
	if err != nil {
		return nil, err
	}
	return config, nil
}

func CreateKubernetesClient() ClientWrapper {
	config, err := getKubernetesConfig()
	if err != nil {
		log.Fatal().Msgf("Unable to load kubernetes config: %v", err)
	}

	client, err2 := kubernetes.NewForConfig(config)
	if err2 != nil {
		log.Fatal().Msgf("Unable to create a kubernetes client: %v", err2)
	}
	// Supress k8s client-go
	klog.SetLogger(logr.Discard())
	return ClientWrapper{client: client}
}

// TODO: this feels horriable but i'm not sure of a better way to handle this

func (c *ClientWrapper) ForEachDeployment(namespace string, options metav1.ListOptions, handler func(resource appsv1.Deployment) error) error {
	resources, queryErr := c.client.AppsV1().Deployments(namespace).List(context.TODO(), options)
	if queryErr != nil {
		return queryErr
	}
	for _, resource := range resources.Items {
		if err := handler(resource); err != nil {
			return err
		}
	}
	return nil
}

func (c *ClientWrapper) ForEachStatefulSet(namespace string, options metav1.ListOptions, handler func(resource appsv1.StatefulSet) error) error {
	resources, queryErr := c.client.AppsV1().StatefulSets(namespace).List(context.TODO(), options)
	if queryErr != nil {
		return queryErr
	}
	for _, resource := range resources.Items {
		if err := handler(resource); err != nil {
			return err
		}
	}
	return nil
}

func (c *ClientWrapper) ForEachDaemonSet(namespace string, options metav1.ListOptions, handler func(resource appsv1.DaemonSet) error) error {
	resources, queryErr := c.client.AppsV1().DaemonSets(namespace).List(context.TODO(), options)
	if queryErr != nil {
		return queryErr
	}
	for _, resource := range resources.Items {
		if err := handler(resource); err != nil {
			return err
		}
	}
	return nil
}

func (c *ClientWrapper) ForEachJob(namespace string, options metav1.ListOptions, handler func(resource batchv1.Job) error) error {
	resources, queryErr := c.client.BatchV1().Jobs(namespace).List(context.TODO(), options)
	if queryErr != nil {
		return queryErr
	}
	for _, resource := range resources.Items {
		if err := handler(resource); err != nil {
			return err
		}
	}
	return nil
}

func (c *ClientWrapper) ForEachCronJob(namespace string, options metav1.ListOptions, handler func(resource batchv1beta1.CronJob) error) error {
	resources, queryErr := c.client.BatchV1beta1().CronJobs(namespace).List(context.TODO(), options)
	if queryErr != nil {
		return queryErr
	}
	for _, resource := range resources.Items {
		if err := handler(resource); err != nil {
			return err
		}
	}
	return nil
}

func (c *ClientWrapper) ForEachService(namespace string, options metav1.ListOptions, handler func(resource v1.Service) error) error {
	resources, queryErr := c.client.CoreV1().Services(namespace).List(context.TODO(), options)
	if queryErr != nil {
		return queryErr
	}
	for _, resource := range resources.Items {
		if err := handler(resource); err != nil {
			return err
		}
	}
	return nil
}

func (c *ClientWrapper) ForEachIngress(namespace string, options metav1.ListOptions, handler func(resource networkingv1.Ingress) error) error {
	resources, queryErr := c.client.NetworkingV1().Ingresses(namespace).List(context.TODO(), options)
	if queryErr != nil {
		return queryErr
	}
	for _, resource := range resources.Items {
		if err := handler(resource); err != nil {
			return err
		}
	}
	return nil
}

func (c *ClientWrapper) ForEachConfigMap(namespace string, options metav1.ListOptions, handler func(resource v1.ConfigMap) error) error {
	resources, queryErr := c.client.CoreV1().ConfigMaps(namespace).List(context.TODO(), options)
	if queryErr != nil {
		return queryErr
	}
	for _, resource := range resources.Items {
		if err := handler(resource); err != nil {
			return err
		}
	}
	return nil
}

func (c *ClientWrapper) ForEachSecret(namespace string, options metav1.ListOptions, handler func(resource v1.Secret) error) error {
	resources, queryErr := c.client.CoreV1().Secrets(namespace).List(context.TODO(), options)
	if queryErr != nil {
		return queryErr
	}
	for _, resource := range resources.Items {
		if err := handler(resource); err != nil {
			return err
		}
	}
	return nil
}

func (c *ClientWrapper) Query(selector KubernetesSelector, handler func(resource []byte) error) error {
	listOptions := metav1.ListOptions{
		LabelSelector: selector.LabelSelector(),
	}
	// TODO: Change handler to use []byte instead of a k8s types
	switch strings.ToLower(selector.Kind) {
	case "deployment":
		handler := func(resource appsv1.Deployment) error {
			bytes, bytesErr := json.Marshal(resource)
			if bytesErr != nil {
				return bytesErr
			}
			handlerErr := handler(bytes)
			if handlerErr != nil {
				return handlerErr
			}
			return nil
		}
		c.ForEachDeployment(selector.Namespace, listOptions, handler)
		break
	case "statefulset":
		handler := func(resource appsv1.StatefulSet) error {
			bytes, bytesErr := json.Marshal(resource)
			if bytesErr != nil {
				return bytesErr
			}
			handlerErr := handler(bytes)
			if handlerErr != nil {
				return handlerErr
			}
			return nil
		}
		c.ForEachStatefulSet(selector.Namespace, listOptions, handler)
		break
	case "daemonset":
		handler := func(resource appsv1.DaemonSet) error {
			bytes, bytesErr := json.Marshal(resource)
			if bytesErr != nil {
				return bytesErr
			}
			handlerErr := handler(bytes)
			if handlerErr != nil {
				return handlerErr
			}
			return nil
		}
		c.ForEachDaemonSet(selector.Namespace, listOptions, handler)
		break
	case "job":
		handler := func(resource batchv1.Job) error {
			bytes, bytesErr := json.Marshal(resource)
			if bytesErr != nil {
				return bytesErr
			}
			handlerErr := handler(bytes)
			if handlerErr != nil {
				return handlerErr
			}
			return nil
		}
		c.ForEachJob(selector.Namespace, listOptions, handler)
		break
	case "cronjob":
		handler := func(resource batchv1beta1.CronJob) error {
			bytes, bytesErr := json.Marshal(resource)
			if bytesErr != nil {
				return bytesErr
			}
			handlerErr := handler(bytes)
			if handlerErr != nil {
				return handlerErr
			}
			return nil
		}
		c.ForEachCronJob(selector.Namespace, listOptions, handler)
		break
	case "service":
		handler := func(resource v1.Service) error {
			bytes, bytesErr := json.Marshal(resource)
			if bytesErr != nil {
				return bytesErr
			}
			handlerErr := handler(bytes)
			if handlerErr != nil {
				return handlerErr
			}
			return nil
		}
		c.ForEachService(selector.Namespace, listOptions, handler)
		break
	case "ingress":
		handler := func(resource networkingv1.Ingress) error {
			bytes, bytesErr := json.Marshal(resource)
			if bytesErr != nil {
				return bytesErr
			}
			handlerErr := handler(bytes)
			if handlerErr != nil {
				return handlerErr
			}
			return nil
		}
		c.ForEachIngress(selector.Namespace, listOptions, handler)
		break
	case "configmap":
		handler := func(resource v1.ConfigMap) error {
			bytes, bytesErr := json.Marshal(resource)
			if bytesErr != nil {
				return bytesErr
			}
			handlerErr := handler(bytes)
			if handlerErr != nil {
				return handlerErr
			}
			return nil
		}
		c.ForEachConfigMap(selector.Namespace, listOptions, handler)
		break
	case "secret":
		handler := func(resource v1.Secret) error {
			bytes, bytesErr := json.Marshal(resource)
			if bytesErr != nil {
				return bytesErr
			}
			handlerErr := handler(bytes)
			if handlerErr != nil {
				return handlerErr
			}
			return nil
		}
		c.ForEachSecret(selector.Namespace, listOptions, handler)
		break
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
