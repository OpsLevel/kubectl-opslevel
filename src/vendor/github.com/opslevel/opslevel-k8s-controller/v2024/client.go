package opslevel_k8s_controller

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/go-logr/logr"
	"k8s.io/apimachinery/pkg/api/meta"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/dynamic/dynamicinformer"
	"k8s.io/klog/v2"

	"k8s.io/client-go/discovery"
	"k8s.io/client-go/discovery/cached/memory"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/restmapper"

	// This is here because of https://github.com/OpsLevel/kubectl-opslevel/issues/24
	_ "k8s.io/client-go/plugin/pkg/client/auth"
	"k8s.io/client-go/tools/clientcmd"
)

var (
	namespacesWereCached bool
	namespacesCache      []string
)

type K8SSelector struct {
	ApiVersion string   `json:"apiVersion" yaml:"apiVersion"`
	Kind       string   `json:"kind" yaml:"kind"`
	Namespaces []string `json:"namespaces,omitempty" yaml:"namespaces,omitempty"`
	Labels     []string `json:"labels,omitempty" yaml:"labels,omitempty"`
	Excludes   []string `json:"excludes,omitempty" yaml:"excludes,omitempty"`
}

func (selector *K8SSelector) GetListOptions() metav1.ListOptions {
	return metav1.ListOptions{
		LabelSelector: selector.LabelSelector(),
	}
}

func (selector *K8SSelector) LabelSelector() string {
	var labels []string
	for _, label := range selector.Labels {
		data := strings.Split(label, "=")
		labels = append(labels, fmt.Sprintf("%s=%s", data[0], data[1]))
	}
	return strings.Join(labels, ",")
}

type K8SClient struct {
	Client  kubernetes.Interface
	Dynamic dynamic.Interface
	Mapper  *restmapper.DeferredDiscoveryRESTMapper
}

// NewK8SClient
// This creates a wrapper which gives you an initialized and connected kubernetes client
// It then has a number of helper functions
func NewK8SClient() (*K8SClient, error) {
	loadingRules := clientcmd.NewDefaultClientConfigLoadingRules()
	configOverrides := &clientcmd.ConfigOverrides{}
	config, err := clientcmd.NewNonInteractiveDeferredLoadingClientConfig(loadingRules, configOverrides).ClientConfig()
	if err != nil {
		return nil, err
	}

	client1, err := kubernetes.NewForConfig(config)
	if err != nil {
		return nil, err
	}

	client2, err := dynamic.NewForConfig(config)
	if err != nil {
		return nil, err
	}

	dc, err := discovery.NewDiscoveryClientForConfig(config)
	if err != nil {
		return nil, err
	}
	mapper := restmapper.NewDeferredDiscoveryRESTMapper(memory.NewMemCacheClient(dc))

	// Suppress k8s client-go logs
	klog.SetLogger(logr.Discard())
	return &K8SClient{Client: client1, Dynamic: client2, Mapper: mapper}, nil
}

func (c *K8SClient) GetMapping(selector K8SSelector) (*meta.RESTMapping, error) {
	gv, gvErr := schema.ParseGroupVersion(selector.ApiVersion)
	if gvErr != nil {
		return nil, gvErr
	}
	gvk := gv.WithKind(selector.Kind)

	mapping, mappingErr := c.Mapper.RESTMapping(gvk.GroupKind(), gvk.Version)
	if mappingErr != nil {
		return nil, mappingErr
	}

	return mapping, nil
}

func (c *K8SClient) GetGVR(selector K8SSelector) (*schema.GroupVersionResource, error) {
	mapping, err := c.GetMapping(selector)
	if err != nil {
		return nil, err
	}
	return &mapping.Resource, nil
}

func (c *K8SClient) GetInformerFactory(resync time.Duration) dynamicinformer.DynamicSharedInformerFactory {
	return dynamicinformer.NewDynamicSharedInformerFactory(c.Dynamic, resync)
}

func (c *K8SClient) GetNamespaces(selector K8SSelector) ([]string, error) {
	if len(selector.Namespaces) > 0 {
		return selector.Namespaces, nil
	} else {
		if namespacesWereCached {
			return namespacesCache, nil
		}
		allNamespaces, err := c.GetAllNamespaces()
		if err != nil {
			return nil, err
		}
		namespacesWereCached = true
		namespacesCache = allNamespaces
		return namespacesCache, nil
	}
}

func (c *K8SClient) GetAllNamespaces() ([]string, error) {
	var output []string
	resources, queryErr := c.Client.CoreV1().Namespaces().List(context.Background(), metav1.ListOptions{})
	if queryErr != nil {
		return output, queryErr
	}
	for _, resource := range resources.Items {
		output = append(output, resource.Name)
	}
	return output, nil
}

func (c *K8SClient) Query(selector K8SSelector) (output []unstructured.Unstructured, err error) {
	aggregator := func(resource unstructured.Unstructured) {
		output = append(output, resource)
	}
	namespaces, err := c.GetNamespaces(selector)
	if err != nil {
		return
	}
	mapping, err := c.GetMapping(selector)
	if err != nil {
		err = fmt.Errorf("%s \n\t please ensure you are using a valid `ApiVersion` and `Kind` found in `kubectl api-resources --verbs=\"get,list\"`", err)
		return
	}
	options := selector.GetListOptions()
	dr := c.Dynamic.Resource(mapping.Resource)
	if mapping.Scope.Name() == meta.RESTScopeNameNamespace {
		for _, namespace := range namespaces {
			if err = c.list(dr.Namespace(namespace), options, aggregator); err != nil {
				return
			}
		}
	} else {
		if err = c.list(dr, options, aggregator); err != nil {
			return
		}
	}
	return
}

func (c *K8SClient) list(client dynamic.ResourceInterface, options metav1.ListOptions, aggregator func(resource unstructured.Unstructured)) (err error) {
	resources, err := client.List(context.Background(), options)
	if err != nil {
		return
	}
	for _, resource := range resources.Items {
		aggregator(resource)
	}
	return
}
