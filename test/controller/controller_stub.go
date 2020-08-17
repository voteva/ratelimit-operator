package controller

import (
	"k8s.io/client-go/util/workqueue"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	"sigs.k8s.io/controller-runtime/pkg/predicate"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
	"sigs.k8s.io/controller-runtime/pkg/source"
	"time"
)

type StubController struct {
	Client     client.Client
	Reconciler reconcile.Reconciler
	Queue      workqueue.RateLimitingInterface
	watches    []watchTarget
}

func NewStubController(client client.Client, reconciler reconcile.Reconciler) *StubController {
	return &StubController{
		Client:     client,
		Reconciler: reconciler,
		Queue:      NewStubQueue(),
		watches:    *new([]watchTarget),
	}
}

func (c *StubController) Watch(src source.Source, eventhandler handler.EventHandler, predicates ...predicate.Predicate) error {
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

	for &nextRequest != nil {
		c.reconcileRequest(*nextRequest)
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
	var nextRequest *reconcile.Request
	if !c.Queue.ShuttingDown() && c.Queue.Len() > 0 {
		get, _ := c.Queue.Get()
		nextRequest = get.(*reconcile.Request)
	}
	return nextRequest
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
