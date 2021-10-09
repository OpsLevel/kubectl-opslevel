package k8sutils

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/rs/zerolog/log"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/util/runtime"
	"k8s.io/apimachinery/pkg/util/wait"
	"k8s.io/client-go/dynamic/dynamicinformer"
	"k8s.io/client-go/tools/cache"
	"k8s.io/client-go/util/workqueue"
)

type KubernetesControllerHandler func([]interface{})

type KubernetesController struct {
	id       string
	factory  dynamicinformer.DynamicSharedInformerFactory
	queue    *workqueue.Type
	informer cache.SharedIndexInformer
	maxBatch int
	Channel  chan struct{}
	OnAdd    KubernetesControllerHandler
	OnUpdate KubernetesControllerHandler
	OnDelete KubernetesControllerHandler
}

type KubernetesEventType string

const (
	KubernetesEventTypeCreate KubernetesEventType = "create"
	KubernetesEventTypeUpdate KubernetesEventType = "update"
	KubernetesEventTypeDelete KubernetesEventType = "delete"
)

type KubernetesEvent struct {
	Key  string
	Type KubernetesEventType
}

func nullKubernetesControllerHandler(items []interface{}) {}

var onlyOneSignalHandler = make(chan struct{})

func setupSignalHandler() <-chan struct{} {
	close(onlyOneSignalHandler) // panics when called twice

	stop := make(chan struct{})
	c := make(chan os.Signal, 2)
	signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		<-c
		close(stop)
		<-c
		os.Exit(1) // second signal. Exit directly.
	}()

	return stop
}

func (c *KubernetesController) getLength() int {
	current := c.queue.Len()
	if current < c.maxBatch {
		return current
	} else {
		return c.maxBatch
	}
}

func (c *KubernetesController) processNextItem() bool {
	var matchType KubernetesEventType
	items := make([]interface{}, 0)
	length := c.getLength()
	for i := 0; i < length; i++ {
		item, quit := c.queue.Get()
		if quit {
			return false
		}
		event := item.(KubernetesEvent)
		if i == 0 {
			// First item determines what event type we process
			matchType = event.Type
		} else {
			if event.Type != matchType {
				c.queue.Add(item)
				c.queue.Done(item)
				continue
			}
		}
		obj, exists, getErr := c.informer.GetIndexer().GetByKey(event.Key)
		if getErr != nil {
			log.Warn().Msgf("error fetching object with key %s from informer cache: %v", event.Key, getErr)
			c.queue.Done(item)
			return true
		}
		if !exists {
			if matchType != KubernetesEventTypeDelete {
				log.Warn().Msgf("object with key %s doesn't exist in informer cache", event.Key)
			}
			c.queue.Done(item)
			return true
		}
		c.queue.Done(item)
		items = append(items, obj)
	}
	switch matchType {
	case KubernetesEventTypeCreate:
		c.OnAdd(items)
	case KubernetesEventTypeUpdate:
		c.OnUpdate(items)
	case KubernetesEventTypeDelete:
		c.OnDelete(items)
	}

	return true
}
func (c *KubernetesController) mainloop() {
	for c.processNextItem() {
	}
}

func (c *KubernetesController) Start(workers int) {
	defer runtime.HandleCrash()
	defer c.queue.ShutDown()
	c.factory.Start(c.Channel) // Starts all informers

	for _, ready := range c.factory.WaitForCacheSync(c.Channel) {
		if !ready {
			runtime.HandleError(fmt.Errorf("[%s] Timed out waiting for caches to sync", c.id))
			return
		}
		log.Info().Msgf("[%s] Informer is ready and synced", c.id)
	}
	if workers < 1 {
		workers = 1
	}
	for i := 0; i < workers; i++ {
		log.Info().Msgf("[%s] Creating worker #%d", c.id, i+1)
		go wait.Until(c.mainloop, time.Second, c.Channel)
	}

	<-c.Channel
}

func NewController(gvr schema.GroupVersionResource, resyncInterval time.Duration, maxBatch int) *KubernetesController {
	k8sClient := GetOrCreateKubernetesClient()
	queue := workqueue.New()
	factory := k8sClient.GetInformerFactory(resyncInterval)
	informer := factory.ForResource(gvr).Informer()
	informer.AddEventHandler(cache.ResourceEventHandlerFuncs{
		AddFunc: func(obj interface{}) {
			var err error
			var item KubernetesEvent
			item.Key, err = cache.MetaNamespaceKeyFunc(obj)
			item.Type = KubernetesEventTypeCreate
			if err == nil {
				log.Debug().Msgf("Queuing event: %+v", item)
				queue.Add(item)
			}
		},
		UpdateFunc: func(old, new interface{}) {
			var err error
			var item KubernetesEvent
			item.Key, err = cache.MetaNamespaceKeyFunc(old)
			item.Type = KubernetesEventTypeUpdate
			if err == nil {
				log.Debug().Msgf("Queuing event: %+v", item)
				queue.Add(item)
			}
		},
		DeleteFunc: func(obj interface{}) {
			var err error
			var item KubernetesEvent
			item.Key, err = cache.DeletionHandlingMetaNamespaceKeyFunc(obj)
			item.Type = KubernetesEventTypeDelete
			if err == nil {
				log.Debug().Msgf("Queuing event: %+v", item)
				queue.Add(item)
			}
		},
	})
	return &KubernetesController{
		id:       fmt.Sprintf("%s/%s/%s", gvr.Group, gvr.Version, gvr.Resource),
		queue:    queue,
		factory:  factory,
		informer: informer,
		maxBatch: maxBatch,
		Channel:  make(chan struct{}),
		OnAdd:    nullKubernetesControllerHandler,
		OnUpdate: nullKubernetesControllerHandler,
		OnDelete: nullKubernetesControllerHandler,
	}
}

func Start() {
	log.Info().Msg("Controller is Starting...")
	<-setupSignalHandler() // Block until signals
	log.Info().Msg("Controller is Stopping...")
}
