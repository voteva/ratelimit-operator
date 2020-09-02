package controller

import (
	"k8s.io/apimachinery/pkg/api/meta"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/util/workqueue"
	"sigs.k8s.io/controller-runtime/pkg/cache"
	"sigs.k8s.io/controller-runtime/pkg/cache/informertest"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	"sigs.k8s.io/controller-runtime/pkg/predicate"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
	"sigs.k8s.io/controller-runtime/pkg/runtime/inject"
	"sigs.k8s.io/controller-runtime/pkg/source"
	"time"
)

type StubController struct {
	Client     client.Client
	Reconciler reconcile.Reconciler
	Queue      workqueue.RateLimitingInterface
	watches    []watchTarget
	scheme *runtime.Scheme
	setFields func(i interface{}) error
}

func NewStubController(client client.Client, reconciler reconcile.Reconciler, cache cache.Cache, scheme runtime.Scheme) controller.Controller {

	return &StubController{
		Client:     client,
		Reconciler: reconciler,
		Queue:      NewStubQueue(),
		watches:    *new([]watchTarget),
		setFields: func(i interface{}) error {
			if _, err := inject.SchemeInto(&scheme, i); err != nil {
				return err
			}
			if _, err := inject.CacheInto(cache, i); err != nil {
				return err
			}
			if _, err := inject.MapperInto(createMapping(scheme), i); err != nil {
				return err
			}
			return nil
		},
	}
}

func createMapping(scheme runtime.Scheme) meta.RESTMapper{
	restMapper := meta.NewDefaultRESTMapper(nil)
	for gvk, _ := range scheme.AllKnownTypes() {
		restMapper.Add(gvk, meta.RESTScopeNamespace)
	}
	return restMapper
}

func (c *StubController) Watch(src source.Source, eventhandler handler.EventHandler, predicates ...predicate.Predicate) error {

	// Inject Cache into arguments
	if err := c.setFields(src); err != nil {
		return err
	}
	if err := c.setFields(eventhandler); err != nil {
		return err
	}
	for _, pr := range predicates {
		if err := c.setFields(pr); err != nil {
			return err
		}
	}

	c.watches = append(c.watches, watchTarget{
		src:          src,
		eventhandler: eventhandler,
		predicates:   predicates,
	})
	return src.Start(eventhandler, c.Queue, predicates...)
}

// Start starts the controller.  Start blocks until stop is closed or a
// controller has an error starting.
func (c *StubController) Start(stop <-chan struct{}) error {
	return nil
}

func (c *StubController) Reconcile(request reconcile.Request) (reconcile.Result, error) {
	var nextRequest *reconcile.Request
	if &request != nil {
		nextRequest = &request
	} else {
		nextRequest = c.nextRequest()
	}

	for nextRequest != nil {
		c.reconcileRequest(*nextRequest)
		//after ok reconcile forget(delete) request in queue
		c.Queue.Forget(nextRequest)
		nextRequest = c.nextRequest()
	}

	return reconcile.Result{}, nil
}

func (c *StubController) reconcileRequest(request reconcile.Request) {
	needRequeue := false
	for ok := true; ok; ok = needRequeue {
		result, err := c.Reconciler.Reconcile(request)
		switch {
		case err != nil || result.Requeue:
			needRequeue = true
		default:
			needRequeue = false
		}
	}
}

func (c *StubController) nextRequest() *reconcile.Request {
	var nextRequest reconcile.Request
	if !c.Queue.ShuttingDown() && c.Queue.Len() > 0 {
		get, _ := c.Queue.Get()
		nextRequest = get.(reconcile.Request)
		return &nextRequest
	} else {
		return nil
	}
}

type watchTarget struct {
	src          source.Source
	eventhandler handler.EventHandler
	predicates   []predicate.Predicate
}

type StubRateLimitingQueue struct {
	wrappedQueue workqueue.Interface
}

func (s StubRateLimitingQueue) Add(item interface{}) {
	s.wrappedQueue.Add(item)
}

func (s StubRateLimitingQueue) Len() int {
	return s.wrappedQueue.Len()
}

func (s StubRateLimitingQueue) Get() (item interface{}, shutdown bool) {
	return s.wrappedQueue.Get()
}

func (s StubRateLimitingQueue) Done(item interface{}) {
	s.wrappedQueue.Done(item)
}

func (s StubRateLimitingQueue) ShutDown() {
	s.wrappedQueue.ShutDown()
}

func (s StubRateLimitingQueue) ShuttingDown() bool {
	return s.wrappedQueue.ShuttingDown()
}

func (s StubRateLimitingQueue) AddAfter(item interface{}, duration time.Duration) {
	s.wrappedQueue.Add(item)
}

func (s StubRateLimitingQueue) AddRateLimited(item interface{}) {
	s.wrappedQueue.Add(item)
}

func (s StubRateLimitingQueue) Forget(item interface{}) {

}

func (s StubRateLimitingQueue) NumRequeues(item interface{}) int {
	return 0
}

func NewStubQueue() StubRateLimitingQueue {
	return StubRateLimitingQueue{wrappedQueue: workqueue.New()}
}

//isn't finished
func GetUnitTestEnv(scheme runtime.Scheme, reconciler reconcile.Reconciler, initObjs...runtime.Object) (*controller.Controller, *client.Client) {
	var cacheInformers cache.Cache = &informertest.FakeInformers{
		Scheme: &scheme,
	}
	stubClient := NewStubClient(&scheme, &cacheInformers, initObjs...)
	stubController := NewStubController(stubClient, reconciler, cacheInformers, scheme)

	return &stubController, &stubClient
}
