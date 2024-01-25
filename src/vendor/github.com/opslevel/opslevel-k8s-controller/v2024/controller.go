package opslevel_k8s_controller

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/rs/zerolog/log"
	"k8s.io/apimachinery/pkg/util/runtime"
	"k8s.io/client-go/dynamic/dynamicinformer"
	"k8s.io/client-go/tools/cache"
	"k8s.io/client-go/util/workqueue"
)

type K8SControllerHandler func(interface{})

func nullKubernetesControllerHandler(item interface{}) {}

type K8SController struct {
	id       string
	factory  dynamicinformer.DynamicSharedInformerFactory
	queue    *workqueue.Type
	informer cache.SharedIndexInformer
	filter   *K8SFilter
	OnAdd    K8SControllerHandler
	OnUpdate K8SControllerHandler
	OnDelete K8SControllerHandler
}

type K8SControllerEventType string

const (
	ControllerEventTypeCreate K8SControllerEventType = "create"
	ControllerEventTypeUpdate K8SControllerEventType = "update"
	ControllerEventTypeDelete K8SControllerEventType = "delete"
)

type K8SEvent struct {
	Key  string
	Type K8SControllerEventType
}

func (c *K8SController) mainloop(item interface{}) {
	log.Debug().Str("queue_addr", fmt.Sprintf("%p", &c.queue)).Int("queue_len", c.queue.Len()).Msg("mainloop: running from top")
	var (
		indexer cache.Indexer = c.informer.GetIndexer()
		event   K8SEvent
	)

	if _, ok := item.(K8SEvent); !ok {
		log.Warn().Msgf("mainloop: cannot create K8SEvent from unknown interface '%T'", item)
		return
	}
	event = item.(K8SEvent)
	obj, exists, err := indexer.GetByKey(event.Key)
	if err != nil {
		log.Warn().Msgf("error fetching object with key '%s' from informer cache: '%v'", event.Key, err)
		return
	}
	if !exists {
		log.Debug().Msgf("object with key '%s' skipped because it was not found", event.Key)
		return
	}
	if c.filter.Matches(obj) {
		log.Debug().Msgf("object with key '%s' skipped because it matches filter", event.Key)
		return
	}
	switch event.Type {
	case ControllerEventTypeCreate:
		c.OnAdd(obj)
	case ControllerEventTypeUpdate:
		c.OnUpdate(obj)
	case ControllerEventTypeDelete:
		c.OnDelete(obj)
	default:
		log.Warn().Msgf("no event handler for '%s', event type '%s'", event.Key, event.Type)
	}
}

// Start starts the informer factory and sends events in the queue to the main loop.
// If a wait group is passed, Start will decrement it once it processes all events
// in the queue after one loop.
// If a wait group is not passed, Start will run continuously until the passed context
// is interrupted.
func (c *K8SController) Start(ctx context.Context, wg *sync.WaitGroup) {
	c.factory.Start(nil) // Starts all informers
	for _, ready := range c.factory.WaitForCacheSync(nil) {
		if !ready {
			runtime.HandleError(fmt.Errorf("[%s] Timed out waiting for caches to sync", c.id))
			return
		}
		log.Info().Msgf("[%s] Informer is ready and synced", c.id)
	}
	go func() {
		defer runtime.HandleCrash()
		if wg != nil {
			defer wg.Done()
		}
		var item interface{}
		var quit bool
	body:
		for {
			// This is a blocking operation performed in a goroutine to get the next event
			itemCh := make(chan interface{}, 1)
			quitCh := make(chan bool, 1)
			go func() {
				item, quit = c.queue.Get()
				itemCh <- item
				quitCh <- quit
			}()

			// Stop execution if the context is cancelled before get event finishes
			select {
			case <-ctx.Done():
				log.Debug().Msg("Breaking: on signal")
				break body
			case item = <-itemCh:
				quit = <-quitCh
			}
			if quit {
				log.Debug().Msg("Breaking: on quit")
				break
			}
			c.mainloop(item)
			c.queue.Done(item)
		}
	}()
	if wg != nil {
		c.queue.ShutDownWithDrain()
	}
}

func NewK8SController(selector K8SSelector, resyncInterval time.Duration) (*K8SController, error) {
	k8sClient, err := NewK8SClient()
	if err != nil {
		return nil, err
	}
	gvr, err := k8sClient.GetGVR(selector)
	if err != nil {
		return nil, err
	}

	queue := workqueue.New()
	filter := NewK8SFilter(selector)
	factory := k8sClient.GetInformerFactory(resyncInterval)
	informer := factory.ForResource(*gvr).Informer()
	_, err = informer.AddEventHandler(cache.ResourceEventHandlerFuncs{
		AddFunc: func(obj interface{}) {
			key, err := cache.MetaNamespaceKeyFunc(obj)
			if err != nil {
				return
			}
			queue.Add(K8SEvent{
				Key:  key,
				Type: ControllerEventTypeCreate,
			})
		},
		UpdateFunc: func(old, new interface{}) {
			key, err := cache.MetaNamespaceKeyFunc(old)
			if err != nil {
				return
			}
			queue.Add(K8SEvent{
				Key:  key,
				Type: ControllerEventTypeUpdate,
			})
		},
		DeleteFunc: func(obj interface{}) {
			key, err := cache.DeletionHandlingMetaNamespaceKeyFunc(obj)
			if err != nil {
				return
			}
			queue.Add(K8SEvent{
				Key:  key,
				Type: ControllerEventTypeDelete,
			})
		},
	})
	return &K8SController{
		id:       fmt.Sprintf("%s/%s/%s", gvr.Group, gvr.Version, gvr.Resource),
		queue:    queue,
		factory:  factory,
		informer: informer,
		filter:   filter,
		OnAdd:    nullKubernetesControllerHandler,
		OnUpdate: nullKubernetesControllerHandler,
		OnDelete: nullKubernetesControllerHandler,
	}, err
}
