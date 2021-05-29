package k8sutils

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/go-logr/logr"
	"github.com/rs/zerolog/log"

	// commented out because we are just dealing with []byte but we might need them in the future
	// appsv1 "k8s.io/api/apps/v1"
	// batchv1 "k8s.io/api/batch/v1"
	// batchv1beta1 "k8s.io/api/batch/v1beta1"
	// v1 "k8s.io/api/core/v1"
	// networkingv1 "k8s.io/api/networking/v1"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"

	// This is here because of https://github.com/OpsLevel/kubectl-opslevel/issues/24
	_ "k8s.io/client-go/plugin/pkg/client/auth"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/klog/v2"
)

type NamespaceSelector struct {
	Include []string
	Exclude []string
}

type KubernetesSelector struct {
	Kind      string
	Namespace NamespaceSelector
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

func (c *ClientWrapper) ForEachDeployment(namespace string, options metav1.ListOptions, handler func(resource []byte) error) error {
	resources, queryErr := c.client.AppsV1().Deployments(namespace).List(context.TODO(), options)
	if queryErr != nil {
		return queryErr
	}
	for _, resource := range resources.Items {
		bytes, bytesErr := json.Marshal(resource)
		if bytesErr != nil {
			return bytesErr
		}
		if err := handler(bytes); err != nil {
			return err
		}
	}
	return nil
}

func (c *ClientWrapper) ForEachStatefulSet(namespace string, options metav1.ListOptions, handler func(resource []byte) error) error {
	resources, queryErr := c.client.AppsV1().StatefulSets(namespace).List(context.TODO(), options)
	if queryErr != nil {
		return queryErr
	}
	for _, resource := range resources.Items {
		bytes, bytesErr := json.Marshal(resource)
		if bytesErr != nil {
			return bytesErr
		}
		if err := handler(bytes); err != nil {
			return err
		}
	}
	return nil
}

func (c *ClientWrapper) ForEachDaemonSet(namespace string, options metav1.ListOptions, handler func(resource []byte) error) error {
	resources, queryErr := c.client.AppsV1().DaemonSets(namespace).List(context.TODO(), options)
	if queryErr != nil {
		return queryErr
	}
	for _, resource := range resources.Items {
		bytes, bytesErr := json.Marshal(resource)
		if bytesErr != nil {
			return bytesErr
		}
		if err := handler(bytes); err != nil {
			return err
		}
	}
	return nil
}

func (c *ClientWrapper) ForEachJob(namespace string, options metav1.ListOptions, handler func(resource []byte) error) error {
	resources, queryErr := c.client.BatchV1().Jobs(namespace).List(context.TODO(), options)
	if queryErr != nil {
		return queryErr
	}
	for _, resource := range resources.Items {
		bytes, bytesErr := json.Marshal(resource)
		if bytesErr != nil {
			return bytesErr
		}
		if err := handler(bytes); err != nil {
			return err
		}
	}
	return nil
}

func (c *ClientWrapper) ForEachCronJob(namespace string, options metav1.ListOptions, handler func(resource []byte) error) error {
	resources, queryErr := c.client.BatchV1beta1().CronJobs(namespace).List(context.TODO(), options)
	if queryErr != nil {
		return queryErr
	}
	for _, resource := range resources.Items {
		bytes, bytesErr := json.Marshal(resource)
		if bytesErr != nil {
			return bytesErr
		}
		if err := handler(bytes); err != nil {
			return err
		}
	}
	return nil
}

func (c *ClientWrapper) ForEachService(namespace string, options metav1.ListOptions, handler func(resource []byte) error) error {
	resources, queryErr := c.client.CoreV1().Services(namespace).List(context.TODO(), options)
	if queryErr != nil {
		return queryErr
	}
	for _, resource := range resources.Items {
		bytes, bytesErr := json.Marshal(resource)
		if bytesErr != nil {
			return bytesErr
		}
		if err := handler(bytes); err != nil {
			return err
		}
	}
	return nil
}

func (c *ClientWrapper) ForEachIngress(namespace string, options metav1.ListOptions, handler func(resource []byte) error) error {
	resources, queryErr := c.client.NetworkingV1().Ingresses(namespace).List(context.TODO(), options)
	if queryErr != nil {
		return queryErr
	}
	for _, resource := range resources.Items {
		bytes, bytesErr := json.Marshal(resource)
		if bytesErr != nil {
			return bytesErr
		}
		if err := handler(bytes); err != nil {
			return err
		}
	}
	return nil
}

func (c *ClientWrapper) ForEachConfigMap(namespace string, options metav1.ListOptions, handler func(resource []byte) error) error {
	resources, queryErr := c.client.CoreV1().ConfigMaps(namespace).List(context.TODO(), options)
	if queryErr != nil {
		return queryErr
	}
	for _, resource := range resources.Items {
		bytes, bytesErr := json.Marshal(resource)
		if bytesErr != nil {
			return bytesErr
		}
		if err := handler(bytes); err != nil {
			return err
		}
	}
	return nil
}

func (c *ClientWrapper) ForEachSecret(namespace string, options metav1.ListOptions, handler func(resource []byte) error) error {
	resources, queryErr := c.client.CoreV1().Secrets(namespace).List(context.TODO(), options)
	if queryErr != nil {
		return queryErr
	}
	for _, resource := range resources.Items {
		bytes, bytesErr := json.Marshal(resource)
		if bytesErr != nil {
			return bytesErr
		}
		if err := handler(bytes); err != nil {
			return err
		}
	}
	return nil
}

func (c *ClientWrapper) GetAllNamespaces() ([]string, error) {
	var output []string
	resources, queryErr := c.client.CoreV1().Namespaces().List(context.TODO(), metav1.ListOptions{})
	if queryErr != nil {
		return output, queryErr
	}
	for _, resource := range resources.Items {
		output = append(output, resource.Name)
	}
	return output, nil
}

func (c *ClientWrapper) Query(kind string, namespace string, options metav1.ListOptions, handler func(resource []byte) error) error {
	switch strings.ToLower(kind) {
	case "deployment":
		c.ForEachDeployment(namespace, options, handler)
		break
	case "statefulset":
		c.ForEachStatefulSet(namespace, options, handler)
		break
	case "daemonset":
		c.ForEachDaemonSet(namespace, options, handler)
		break
	case "job":
		c.ForEachJob(namespace, options, handler)
		break
	case "cronjob":
		c.ForEachCronJob(namespace, options, handler)
		break
	case "service":
		c.ForEachService(namespace, options, handler)
		break
	case "ingress":
		c.ForEachIngress(namespace, options, handler)
		break
	case "configmap":
		c.ForEachConfigMap(namespace, options, handler)
		break
	case "secret":
		c.ForEachSecret(namespace, options, handler)
		break
	}
	return nil
}

func (selector *KubernetesSelector) GetListOptions() metav1.ListOptions {
	return metav1.ListOptions{
		LabelSelector: selector.LabelSelector(),
	}
}

func (selector *KubernetesSelector) LabelSelector() string {
	var labels []string
	for key, value := range selector.Labels {
		labels = append(labels, fmt.Sprintf("%s=%s", key, value))
	}
	return strings.Join(labels, ",")
}

// TODO: write a test for this to ensure this continues to work properly
func (selector *KubernetesSelector) FilterNamespaces(namespaces []string) []string {
	var output []string
	include := purge(selector.Namespace.Include)
	exclude := purge(selector.Namespace.Exclude)

	useInclude := len(include) > 0
	useExclude := len(exclude) > 0

	for _, namespace := range namespaces {
		if useInclude && !contains(include, namespace) {
			continue
		}
		if useExclude && contains(exclude, namespace) {
			continue
		}
		output = append(output, namespace)
	}
	return output
}

// removes empty strings from the []string
func purge(s []string) []string {
	var r []string
	for _, str := range s {
		if str != "" {
			r = append(r, str)
		}
	}
	return r
}

func contains(s []string, str string) bool {
	for _, v := range s {
		if v == str {
			return true
		}
	}
	return false
}
