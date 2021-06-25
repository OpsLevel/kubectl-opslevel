package k8sutils

import (
	"context"
	"fmt"
	"strings"

	"github.com/go-logr/logr"
	"github.com/rs/zerolog/log"

	"k8s.io/apimachinery/pkg/api/meta"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/discovery"
	"k8s.io/client-go/discovery/cached/memory"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/restmapper"

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
	ApiVersion string
	Kind       string
	Namespace  NamespaceSelector
	Labels     map[string]string
}

type ClientWrapper struct {
	client  kubernetes.Interface
	dynamic dynamic.Interface
	mapper  restmapper.DeferredDiscoveryRESTMapper
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

	client1, client1Err := kubernetes.NewForConfig(config)
	if client1Err != nil {
		log.Fatal().Msgf("Unable to create a kubernetes client: %v", client1Err)
	}

	client2, client2Err2 := dynamic.NewForConfig(config)
	if client2Err2 != nil {
		log.Fatal().Msgf("Unable to create a dynamic kubernetes client: %v", client2Err2)
	}

	dc, dcErr := discovery.NewDiscoveryClientForConfig(config)
	if dcErr != nil {
		log.Fatal().Msgf("Unable to create a discovery kubernetes client: %v", dcErr)
	}
	mapper := restmapper.NewDeferredDiscoveryRESTMapper(memory.NewMemCacheClient(dc))

	// Supress k8s client-go
	klog.SetLogger(logr.Discard())
	return ClientWrapper{client: client1, dynamic: client2, mapper: *mapper}
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

func (c *ClientWrapper) Query(selector KubernetesSelector, namespaces []string, handler func(resource []byte) error) error {
	mapping, mappingErr := c.GetMapping(selector)
	if mappingErr != nil {
		return fmt.Errorf("%s \n\t Please ensure you are using a valid `ApiVersion` and `Kind` found in `kubectl api-resources --verbs=\"get,list\"`", mappingErr)
	}
	options := selector.GetListOptions()
	dr := c.dynamic.Resource(mapping.Resource)
	if mapping.Scope.Name() == meta.RESTScopeNameNamespace {
		for _, namespace := range selector.FilterNamespaces(namespaces) {
			listErr := List(dr.Namespace(namespace), options, handler)
			if listErr != nil {
				return listErr
			}
		}
	} else {
		listErr := List(dr, options, handler)
		if listErr != nil {
			return listErr
		}
	}
	return nil
}

func List(client dynamic.ResourceInterface, options metav1.ListOptions, handler func(resource []byte) error) error {
	resources, queryErr := client.List(context.TODO(), options)
	if queryErr != nil {
		return fmt.Errorf("%s `%s`", queryErr, "")
	}
	for _, resource := range resources.Items {
		bytes, bytesErr := resource.MarshalJSON()
		if bytesErr != nil {
			return bytesErr
		}
		if err := handler(bytes); err != nil {
			return err
		}
	}

	return nil
}

func (c *ClientWrapper) GetMapping(selector KubernetesSelector) (*meta.RESTMapping, error) {
	gv, gvErr := schema.ParseGroupVersion(selector.ApiVersion)
	if gvErr != nil {
		return nil, gvErr
	}
	gvk := gv.WithKind(selector.Kind)

	mapping, mappingErr := c.mapper.RESTMapping(gvk.GroupKind(), gvk.Version)
	if mappingErr != nil {
		return nil, mappingErr
	}

	return mapping, nil
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
