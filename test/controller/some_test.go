package controller

import (
	istio_v1alpha3 "istio.io/client-go/pkg/apis/networking/v1alpha3"
	v1 "k8s.io/api/core/v1"
	apiextensionsv1beta1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1beta1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/kubernetes/scheme"
	"ratelimit-operator/pkg/apis"
	v12 "ratelimit-operator/pkg/apis/operators/v1"
	"ratelimit-operator/pkg/controller/ratelimiter"
	"sigs.k8s.io/controller-runtime/pkg/cache"
	"sigs.k8s.io/controller-runtime/pkg/cache/informertest"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	"sigs.k8s.io/controller-runtime/pkg/runtime/inject"
	"sigs.k8s.io/controller-runtime/pkg/source"
	"testing"
)

func TestNewStubClient(t *testing.T) {

	c := v1.ConfigMap{}

	eventHandler := StubEventHandler{
		EventHandler: &handler.EnqueueRequestForOwner{
			OwnerType:    &c,
			IsController: true,
		},
		Queue: NewStubQueue(),
	}

	isController := true
	rateLimiter_1 := v12.RateLimiter{
		TypeMeta: metav1.TypeMeta{
			Kind: "RateLimiter",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      "test-ratelimiter",
			Namespace: "test-namespace",
			OwnerReferences: []metav1.OwnerReference{{
				APIVersion: "v1",
				Kind:       "ConfigMap",
				Controller: &isController,
			}},
		},
	}

	rateLimiter_2 := v12.RateLimiter{
		TypeMeta: metav1.TypeMeta{
			Kind: "RateLimiter",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      "test-ratelimiter",
			Namespace: "test-namespace",
			OwnerReferences: []metav1.OwnerReference{metav1.OwnerReference{
				APIVersion: "v1",
				Kind:       "ConfigMap",
				Controller: &isController,
			}},
		},
	}

	eventHandler.OnUpdate(&rateLimiter_1, &rateLimiter_2)
}

func TestNewStubClient_2(t *testing.T) {

	//runtimeScheme := runtime.NewScheme()
	//apis.AddToScheme(runtimeScheme)
	//scheme.AddToScheme(runtimeScheme)
	//apiextensionsv1beta1.AddToScheme(runtimeScheme)
	//istio_v1alpha3.AddToScheme(runtimeScheme)
	//
	//var cv2 cache.Cache = &informertest.FakeInformers{
	//	Scheme: runtimeScheme,
	//}
	//
	//client := NewStubClient(runtimeScheme, &cv2)
	//reconcileRateLimiter := ratelimiter.ReconcileRateLimiter{}
	//fakeInformer, _ := cache_.FakeInformerFor(&v12.RateLimiter{})
	//
	//controller := NewStubController(client, reconcileRateLimiter)
	//rateLimiterSource := &source.Kind{Type: &v12.RateLimiter{}}
	//inject.CacheInto(cv2, rateLimiterSource)
	//
	//rateLimiterSource2 := &source.Kind{Type: &v12.RateLimiter{}}
	//inject.CacheInto(cv2, rateLimiterSource2)
	//
	//controller.Watch(rateLimiterSource, &handler.EnqueueRequestForObject{})
	//controller.Watch(rateLimiterSource, &handler.EnqueueRequestForObject{})
	//
	//fakeInformer.Add(&v12.RateLimiter{})
}

func TestNewStubClient_3(t *testing.T) {
	runtimeScheme := runtime.NewScheme()
	apis.AddToScheme(runtimeScheme)
	scheme.AddToScheme(runtimeScheme)
	apiextensionsv1beta1.AddToScheme(runtimeScheme)
	istio_v1alpha3.AddToScheme(runtimeScheme)

	var cacheInformers cache.Cache = &informertest.FakeInformers{
		Scheme: runtimeScheme,
	}

	client := NewStubClient(runtimeScheme, &cacheInformers)
	reconcileRateLimiter := ratelimiter.ReconcileRateLimiter{}
	controller := NewStubController(client, reconcileRateLimiter)

	rateLimiterSource := &source.Kind{Type: &v12.RateLimiter{}}
	inject.CacheInto(cacheInformers, rateLimiterSource)


	controller.Watch(rateLimiterSource, &handler.EnqueueRequestForObject{})


}
