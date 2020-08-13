package controller

import (
	"sigs.k8s.io/controller-runtime/pkg/handler"
	"sigs.k8s.io/controller-runtime/pkg/predicate"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
	"sigs.k8s.io/controller-runtime/pkg/source"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"github.com/golang-collections/go-datastructures/queue"
)

type StubController struct {
	Client client.Client
	Reconciler reconcile.Reconciler
	requestQueue *queue.Queue
	watches []watchTarget
}

func NewStubController(client client.Client, reconciler reconcile.Reconciler) *StubController {
	return &StubController{
		Client:       client,
		Reconciler:   reconciler,
		requestQueue: queue.New(100),
		watches: *new([]watchTarget),
	}
}

func (c *StubController) Watch(src source.Source, eventhandler handler.EventHandler, predicates ...predicate.Predicate) error {
	c.watches = append(c.watches, watchTarget{
		src:          src,
		eventhandler: eventhandler,
		predicates:   predicates,
	})
	return nil
}

// Start starts the controller.  Start blocks until stop is closed or a
// controller has an error starting.
func (c *StubController) Start(stop <-chan struct{}) error {
	return nil
}

func (c *StubController) Reconcile(request reconcile.Request) (reconcile.Result, error) {
	condition := false
	for ok := true; ok; ok = condition {
		result, err := c.Reconciler.Reconcile(request)
		switch {
		case err != nil || result.Requeue:
			condition = true
		default:
			condition = false
		}
	}

	return reconcile.Result{}, nil
}

type watchTarget struct {
	src source.Source
	eventhandler handler.EventHandler
	predicates []predicate.Predicate
}