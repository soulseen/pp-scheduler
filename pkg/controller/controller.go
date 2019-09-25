package controller

import (
	"fmt"
	log "github.com/golang/glog"
	"github.com/soulseen/pp-scheduler/pkg/prioritize"
	"github.com/soulseen/pp-scheduler/pkg/sqlite"
	apiv1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/fields"
	"k8s.io/apimachinery/pkg/util/runtime"
	utilruntime "k8s.io/apimachinery/pkg/util/runtime"
	"k8s.io/apimachinery/pkg/util/wait"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/cache"
	"k8s.io/kubernetes/staging/src/k8s.io/client-go/util/workqueue"
	"k8s.io/sample-controller/pkg/signals"
	"os"
	"sigs.k8s.io/controller-runtime/pkg/client/config"
	"time"
)

const (
	maxRetries = 5

	defaultResync = 600 * time.Second

	DefaultLabelKey = "ks-pipeline"
)

type Controller struct {
	queue    workqueue.RateLimitingInterface
	informer cache.Controller
	indexer  cache.Indexer
}

func (c *Controller) Run(threadiness int, stopCh <-chan struct{}) {
	// catch error
	defer utilruntime.HandleCrash()
	// close queue and shutdown work.
	defer c.queue.ShutDown()

	log.Info("start controller...")

	// run Informer
	go c.informer.Run(stopCh)

	// waiting for cache successful
	if !cache.WaitForCacheSync(stopCh, c.informer.HasSynced) {
		utilruntime.HandleError(fmt.Errorf("cache time out."))
		return
	}

	for i := 0; i < threadiness; i++ {
		go wait.Until(c.runWorker, time.Second, stopCh)
	}

	<-stopCh
	log.Info("Stopping Pod controller")
}

func (c *Controller) runWorker() {
	for c.processNextItem() {
	}
}

func (c *Controller) processNextItem() bool {

	// get next.
	key, quit := c.queue.Get()
	if quit {
		return false
	}

	// remove Key
	defer c.queue.Done(key)

	// process Key
	err := c.processItem(key.(string))
	c.handleErr(err, key)

	return true
}

func (c *Controller) processItem(key string) error {
	log.Infof("Start procrss event %s", key)
	obj, exists, err := c.indexer.GetByKey(key)
	pod := obj.(*apiv1.Pod)
	if err != nil {
		log.Errorf("Fetching object with key %s from store failed with %v", key, err)
		return err
	}

	if exists {
		name := pod.GetName()
		keys := prioritize.ParseMark(pod.Labels)
		bindingNode := pod.Spec.NodeName
		for _, key := range keys {
			err := addUpdateKeyNode(key, bindingNode)
			if err != nil {
				log.Error("update key: ", key, " error: ", err)
			}
		}

		log.Info("Process success for Pod: ", name)

	} else {
		log.Info("Pod: ", key, " does not exist anymore")
	}
	return nil
}

func addUpdateKeyNode(key, node string) error {
	res, err := sqlite.KeyNodeCilent.KeyNodeSearch(key, node)
	if err != nil {
		log.Errorf("Search object with key %s from sqlite error: ", err)
		return err
	}
	if len(res) != 0 {
		log.Info("update record: ", key, ":", node, ":", res[0].Count+1)
		sqlite.KeyNodeCilent.KeyNodeUpdate(res[0].Id, res[0].Count+1)
		return nil
	} else {
		log.Info("insert record: ", key, ":", node)
		sqlite.KeyNodeCilent.KeyNodeInsert(key, node, 1)
		return nil
	}
}

func RunController() {
	// Get a config to talk to the apiserver
	log.Info("setting up client for manager")

	cfg, err := config.GetConfig()
	if err != nil {
		log.Error(err, "unable to set up client config")
		os.Exit(1)
	}

	// creates the clientset
	clientset, err := kubernetes.NewForConfig(cfg)
	if err != nil {
		log.Fatalf("build clientset error: %s", err.Error())
	}

	setLabel := map[string]string{"spec.schedulerName": "ks-scheduler"}

	selectPod := fields.SelectorFromSet(setLabel)

	podListWatcher := cache.NewListWatchFromClient(clientset.CoreV1().RESTClient(), "pods", apiv1.NamespaceAll, selectPod)

	stopCh := signals.SetupSignalHandler()

	queue := workqueue.NewRateLimitingQueue(workqueue.DefaultControllerRateLimiter())

	indexer, informer := cache.NewIndexerInformer(podListWatcher, &apiv1.Pod{}, defaultResync, cache.FilteringResourceEventHandler{
		FilterFunc: func(obj interface{}) bool {
			switch t := obj.(type) {
			case *apiv1.Pod:
				return assignedPod(t)
			default:
				utilruntime.HandleError(fmt.Errorf("unable to handle object : %T", obj))
				return false
			}
		},
		Handler: cache.ResourceEventHandlerFuncs{
			AddFunc: func(obj interface{}) {
				switch t := obj.(type) {
				case *apiv1.Pod:
					res := assignedPod(t)
					if res == true {
						log.Info("add queue with pod name: ", t.GetName())
						key, err := cache.MetaNamespaceKeyFunc(obj)
						if err == nil {
							queue.Add(key)
						}
					}
				default:
					utilruntime.HandleError(fmt.Errorf("unable to handle object : %T", obj))
					return
				}
			},
		},
	}, cache.Indexers{})

	ctrl := Controller{
		queue,
		informer,
		indexer,
	}

	//ctrl.Run(stopCh)
	go ctrl.Run(3, stopCh)

	// Wait forever
	select {}

}

// assignedPod selects pods that are assigned (scheduled and running).
func assignedPod(pod *apiv1.Pod) bool {
	if len(pod.Spec.NodeName) == 0 {
		return false
	}
	if pod.Status.Phase != apiv1.PodRunning {
		return false
	}
	if _, ok := pod.Labels[DefaultLabelKey]; !ok {
		return false
	}
	return true
}

// handleErr checks if an error happened and makes sure we will retry later.
func (c *Controller) handleErr(err error, key interface{}) {
	if err == nil {
		c.queue.Forget(key)
		return
	}

	// This controller retries 5 times if something goes wrong. After that, it stops trying.
	if c.queue.NumRequeues(key) < maxRetries {
		log.Info("Error syncing pod ", key, ". error:", err)

		c.queue.AddRateLimited(key)
		return
	}

	c.queue.Forget(key)
	runtime.HandleError(err)
	log.Infof("Dropping pod %q out of the queue: %v", key, err)
}
