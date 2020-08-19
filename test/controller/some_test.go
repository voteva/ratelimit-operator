package controller

import (
	istio_v1alpha3 "istio.io/client-go/pkg/apis/networking/v1alpha3"
	apiextensionsv1beta1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1beta1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/kubernetes/scheme"
	"ratelimit-operator/pkg/apis"
	v12 "ratelimit-operator/pkg/apis/operators/v1"
	"ratelimit-operator/pkg/controller/ratelimiter"
	"sigs.k8s.io/controller-runtime/pkg/cache"
	"sigs.k8s.io/controller-runtime/pkg/cache/informertest"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
	"sigs.k8s.io/controller-runtime/pkg/source"
	"testing"
)

func TestWithStubsExample(t *testing.T) {

	//define scheme with types
	runtimeScheme := runtime.NewScheme()
	apis.AddToScheme(runtimeScheme)
	scheme.AddToScheme(runtimeScheme)
	apiextensionsv1beta1.AddToScheme(runtimeScheme)
	istio_v1alpha3.AddToScheme(runtimeScheme)

	//init informers cache
	var cacheInformers cache.Cache = &informertest.FakeInformers{
		Scheme: runtimeScheme,
	}

	//create stub client
	client := NewStubClient(runtimeScheme, &cacheInformers)
	reconcileRateLimiter := ratelimiter.ReconcileRateLimiter{
		Client: client,
		Scheme: runtimeScheme,
	}
	//compose stub controller with previous resources
	controller := NewStubController(client, reconcileRateLimiter, cacheInformers, *runtimeScheme)


	//set watches for resources mutation (according to production controller)
	controller.Watch(&source.Kind{Type: &v12.RateLimiter{}}, &handler.EnqueueRequestForObject{})

	request := reconcile.Request{NamespacedName: types.NamespacedName{
		Namespace: "test-namespace",
		Name:      "test-name",
	}}

	//run reconcile and then check the state/objects through stub client (GET)
	controller.Reconcile(request)

}
