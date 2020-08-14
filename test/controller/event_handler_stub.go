package controller

import (
	"k8s.io/apimachinery/pkg/api/meta"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/util/workqueue"
	"sigs.k8s.io/controller-runtime/pkg/event"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	"sigs.k8s.io/controller-runtime/pkg/predicate"
)

// StubEventHandler adapts a handler.EventHandler interface to a cache.ResourceEventHandler interface
type StubEventHandler struct {
	EventHandler handler.EventHandler
	Queue        workqueue.RateLimitingInterface
	Predicates   []predicate.Predicate
}

// OnAdd creates CreateEvent and calls Create on StubEventHandler
func (e StubEventHandler) OnAdd(obj interface{}) {
	c := event.CreateEvent{}

	// Pull metav1.Object out of the object
	if o, err := meta.Accessor(obj); err == nil {
		c.Meta = o
	} else {

		return
	}

	// Pull the runtime.Object out of the object
	if o, ok := obj.(runtime.Object); ok {
		c.Object = o
	} else {
		return
	}

	for _, p := range e.Predicates {
		if !p.Create(c) {
			return
		}
	}

	// Invoke create handler
	e.EventHandler.Create(c, e.Queue)
}

// OnUpdate creates UpdateEvent and calls Update on StubEventHandler
func (e StubEventHandler) OnUpdate(oldObj, newObj interface{}) {
	u := event.UpdateEvent{}

	// Pull metav1.Object out of the object
	if o, err := meta.Accessor(oldObj); err == nil {
		u.MetaOld = o
	} else {
		return
	}

	// Pull the runtime.Object out of the object
	if o, ok := oldObj.(runtime.Object); ok {
		u.ObjectOld = o
	} else {
		return
	}

	// Pull metav1.Object out of the object
	if o, err := meta.Accessor(newObj); err == nil {
		u.MetaNew = o
	} else {
		return
	}

	// Pull the runtime.Object out of the object
	if o, ok := newObj.(runtime.Object); ok {
		u.ObjectNew = o
	} else {
		return
	}

	for _, p := range e.Predicates {
		if !p.Update(u) {
			return
		}
	}

	// Invoke update handler
	e.EventHandler.Update(u, e.Queue)
}

// OnDelete creates DeleteEvent and calls Delete on StubEventHandler
func (e StubEventHandler) OnDelete(obj interface{}) {
	d := event.DeleteEvent{}


	// Pull metav1.Object out of the object
	if o, err := meta.Accessor(obj); err == nil {
		d.Meta = o
	} else {
		return
	}

	// Pull the runtime.Object out of the object
	if o, ok := obj.(runtime.Object); ok {
		d.Object = o
	} else {
		return
	}

	for _, p := range e.Predicates {
		if !p.Delete(d) {
			return
		}
	}

	// Invoke delete handler
	e.EventHandler.Delete(d, e.Queue)
}
