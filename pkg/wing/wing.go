package wing

import (
	"github.com/Sirupsen/logrus"
	"k8s.io/apimachinery/pkg/fields"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/cache"
	"k8s.io/client-go/util/workqueue"

	"github.com/jetstack/tarmak/pkg/apis/wing/v1alpha1"
	"github.com/jetstack/tarmak/pkg/wing/client"
)

type Wing struct {
	log       *logrus.Entry
	flags     *Flags
	clientset *client.Clientset
	stopCh    chan struct{}
}

type Flags struct {
	ManifestURL  string
	ServerURL    string
	ClusterName  string
	InstanceName string
}

func New(flags *Flags) *Wing {
	logger := logrus.New()
	logger.Level = logrus.DebugLevel

	t := &Wing{
		log:   logger.WithField("app", "wing"),
		flags: flags,
	}
	return t
}

func (w *Wing) Run(args []string) error {

	// create connection to wing server
	restConfig := &rest.Config{
		Host: w.flags.ServerURL,
		TLSClientConfig: rest.TLSClientConfig{
			Insecure: true,
		},
	}
	clientset, err := client.NewForConfig(restConfig)
	if err != nil {
		return err
	}
	w.clientset = clientset

	w.stopCh = make(chan struct{})

	// start watching for API server events that trigger applies
	w.watchForNotifications()

	// run converge loop
	w.convergeLoop()

	// Wait forever
	select {}
}

func (w *Wing) Must(err error) *Wing {
	if err != nil {
		w.log.Fatal(err)
	}
	return w
}

func (w *Wing) watchForNotifications() {

	// create the pod watcher
	instanceListWatcher := cache.NewListWatchFromClient(w.clientset.WingV1alpha1().RESTClient(), "instances", w.flags.ClusterName, fields.Everything())

	// create the workqueue
	queue := workqueue.NewRateLimitingQueue(workqueue.DefaultControllerRateLimiter())

	// Bind the workqueue to a cache with the help of an informer. This way we make sure that
	// whenever the cache is updated, the pod key is added to the workqueue.
	// Note that when we finally process the item from the workqueue, we might see a newer version
	// of the Pod than the version which was responsible for triggering the update.
	indexer, informer := cache.NewIndexerInformer(instanceListWatcher, &v1alpha1.Instance{}, 0, cache.ResourceEventHandlerFuncs{
		AddFunc: func(obj interface{}) {
			key, err := cache.MetaNamespaceKeyFunc(obj)
			if err == nil {
				queue.Add(key)
			}
		},
		UpdateFunc: func(old interface{}, new interface{}) {
			key, err := cache.MetaNamespaceKeyFunc(new)
			if err == nil {
				queue.Add(key)
			}
		},
		DeleteFunc: func(obj interface{}) {
			// IndexerInformer uses a delta queue, therefore for deletes we have to use this
			// key function.
			key, err := cache.DeletionHandlingMetaNamespaceKeyFunc(obj)
			if err == nil {
				queue.Add(key)
			}
		},
	}, cache.Indexers{})

	controller := NewController(queue, indexer, informer)

	// Now let's start the controller
	go controller.Run(1, w.stopCh)

}
