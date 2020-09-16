package ratelimiter_test

import (
	"context"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	istio_v1alpha3 "istio.io/client-go/pkg/apis/networking/v1alpha3"
	appsv1 "k8s.io/api/apps/v1"
	v1 "k8s.io/api/core/v1"
	apiextensionsv1beta1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1beta1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	clientgoscheme "k8s.io/client-go/kubernetes/scheme"
	"github.com/voteva/ratelimit-operator/pkg/apis"
	v12 "github.com/voteva/ratelimit-operator/pkg/apis/operators/v1"
	"github.com/voteva/ratelimit-operator/pkg/controller/ratelimiter"
	"ratelimit-operator/test/controller"
	"sigs.k8s.io/controller-runtime/pkg/cache"
	"sigs.k8s.io/controller-runtime/pkg/cache/informertest"
	"sigs.k8s.io/controller-runtime/pkg/client/fake"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
	"sigs.k8s.io/controller-runtime/pkg/runtime/inject"
	"sigs.k8s.io/controller-runtime/pkg/source"
)

var (
	name                    = "test-ratelimit-operator"
	namespace               = "test-ratelimit"
	preparedRateLimiter     = getPreparedRateLimiterResource()
	preparedDeployment      = getPreparedDeploymentResource()
	preparedRedisDeployment = getPreparedRedisDeploymentResource(name + "-redis")
	preparedService         = getPreparedServiceResource()
	preparedConfigMap       = getPreparedConfigMapResource()
)

var scheme *runtime.Scheme

func init() {
	scheme = runtime.NewScheme()
	apis.AddToScheme(scheme)
	clientgoscheme.AddToScheme(scheme)
	apiextensionsv1beta1.AddToScheme(scheme)
	istio_v1alpha3.AddToScheme(scheme)
}

var _ = Describe("RateLimit controller", func() {

	Context("when RateLimit resource added", func() {

		It("reconcile RateLimit. create secondary resources", func() {

			objects := []runtime.Object{preparedRateLimiter}
			client := fake.NewFakeClientWithScheme(scheme, objects...)
			reconciler := &ratelimiter.ReconcileRateLimiter{Client: client, Scheme: scheme}
			request := &reconcile.Request{NamespacedName: types.NamespacedName{Namespace: namespace, Name: name}}

			//reconcile loop (creation of primary resource should be a reason of reconcile call)
			doReconcile(reconciler, request)

			//check for created objects
			dep := &appsv1.Deployment{}
			confMap := &v1.ConfigMap{}
			srv := &v1.Service{}
			redisDep := &appsv1.Deployment{}
			var err error

			err = reconciler.Client.Get(context.TODO(), request.NamespacedName, dep)
			Expect(err).To(BeNil(), "Deployment is nil")

			err = reconciler.Client.Get(context.TODO(), request.NamespacedName, redisDep)
			Expect(err).To(BeNil(), "Redis Deployment is nil")

			err = reconciler.Client.Get(context.TODO(), request.NamespacedName, confMap)
			Expect(err).To(BeNil(), "ConfigMap is nil")

			err = reconciler.Client.Get(context.TODO(), request.NamespacedName, srv)
			Expect(err).To(BeNil(), "Service is nil")
		})
	})

	Context("when Deployment (secondary resource) has been removed", func() {
		It("reconcile RateLimit. restore deleted Deployment for rate limit service", func() {
			objects := []runtime.Object{
				preparedRateLimiter,
				preparedDeployment,
				preparedRedisDeployment,
				preparedService,
				preparedConfigMap,
			}
			client := fake.NewFakeClientWithScheme(scheme, objects...)
			reconciler := &ratelimiter.ReconcileRateLimiter{Client: client, Scheme: scheme}
			request := &reconcile.Request{NamespacedName: types.NamespacedName{Namespace: namespace, Name: name}}

			//imitation of resource deletion
			client.Delete(context.TODO(), getPreparedDeploymentResource())

			//check resource really deleted
			dep := &appsv1.Deployment{}
			err := client.Get(context.TODO(), request.NamespacedName, dep)
			Expect(err).NotTo(BeNil(), "Deployment is not nil")

			//reconcile loop (deletion of secondary resource fires an event with request -> should call reconcile)
			doReconcile(reconciler, request)

			err = reconciler.Client.Get(context.TODO(), request.NamespacedName, dep)
			Expect(err).To(BeNil(), "Deployment is nil")
		})
	})

	Context("when Redis Deployment (secondary resource) has been removed", func() {
		It("reconcile RateLimit. restore deleted Redis Deployment for rate limit service", func() {
			objects := []runtime.Object{
				preparedRateLimiter,
				preparedDeployment,
				preparedRedisDeployment,
				preparedService,
				preparedConfigMap,
			}
			client := fake.NewFakeClientWithScheme(scheme, objects...)
			reconciler := &ratelimiter.ReconcileRateLimiter{Client: client, Scheme: scheme}
			request := &reconcile.Request{NamespacedName: types.NamespacedName{Namespace: namespace, Name: name}}

			//imitation of resource deletion
			client.Delete(context.TODO(), preparedRedisDeployment)

			//check resource really deleted
			dep := &appsv1.Deployment{}
			err := client.Get(context.TODO(), types.NamespacedName{Namespace: preparedRedisDeployment.Namespace, Name: preparedRedisDeployment.Name}, dep)
			Expect(err).NotTo(BeNil(), "Deployment is not nil")

			//reconcile loop (deletion of secondary resource fires an event with request -> should call reconcile)
			doReconcile(reconciler, request)

			dep = &appsv1.Deployment{}
			err = reconciler.Client.Get(context.TODO(), types.NamespacedName{Namespace: preparedRedisDeployment.Namespace, Name: preparedRedisDeployment.Name}, dep)
			Expect(err).To(BeNil(), "Deployment is nil")
		})
	})

	Context("when Service (secondary resource) has been removed", func() {
		It("reconcile RateLimit. restore deleted Service for rate limit service", func() {
			objects := []runtime.Object{
				preparedRateLimiter,
				preparedDeployment,
				preparedRedisDeployment,
				preparedService,
				preparedConfigMap,
			}
			client := fake.NewFakeClientWithScheme(scheme, objects...)
			reconciler := &ratelimiter.ReconcileRateLimiter{Client: client, Scheme: scheme}
			request := &reconcile.Request{NamespacedName: types.NamespacedName{Namespace: namespace, Name: name}}

			//imitation of resource deletion
			client.Delete(context.TODO(), preparedService)

			//check resource really deleted
			svc := &v1.Service{}
			err := client.Get(context.TODO(), types.NamespacedName{Namespace: preparedService.Namespace, Name: preparedService.Name}, svc)
			Expect(err).NotTo(BeNil(), "Service is not nil")

			//reconcile loop (deletion of secondary resource fires an event with request -> should call reconcile)
			doReconcile(reconciler, request)

			svc = &v1.Service{}
			err = reconciler.Client.Get(context.TODO(), types.NamespacedName{Namespace: preparedService.Namespace, Name: preparedService.Name}, svc)
			Expect(err).To(BeNil(), "Service is nil")
		})
	})

	Context("when ConfigMap (secondary resource) has been removed", func() {
		It("reconcile RateLimit. restore deleted ConfigMap for rate limit service", func() {
			objects := []runtime.Object{
				preparedRateLimiter,
				preparedDeployment,
				preparedRedisDeployment,
				preparedService,
				preparedConfigMap,
			}
			client := fake.NewFakeClientWithScheme(scheme, objects...)
			reconciler := &ratelimiter.ReconcileRateLimiter{Client: client, Scheme: scheme}
			request := &reconcile.Request{NamespacedName: types.NamespacedName{Namespace: namespace, Name: name}}

			//imitation of resource deletion
			client.Delete(context.TODO(), preparedConfigMap)

			//check resource really deleted
			cfgMap := &v1.ConfigMap{}
			err := client.Get(context.TODO(), types.NamespacedName{Namespace: preparedConfigMap.Namespace, Name: preparedConfigMap.Name}, cfgMap)
			Expect(err).NotTo(BeNil(), "ConfigMap is not nil")

			//reconcile loop (deletion of secondary resource fires an event with request -> should call reconcile)
			doReconcile(reconciler, request)

			cfgMap = &v1.ConfigMap{}
			err = reconciler.Client.Get(context.TODO(), types.NamespacedName{Namespace: preparedConfigMap.Namespace, Name: preparedConfigMap.Name}, cfgMap)
			Expect(err).To(BeNil(), "Service is nil")
		})
	})

	Context("when RateLimit resource added. Controller stub", func() {

		It("reconcile RateLimit. create secondary resources", func() {

			var cacheInformers cache.Cache = &informertest.FakeInformers{
				Scheme: scheme,
			}

			objects := []runtime.Object{preparedRateLimiter}
			client := controller.NewStubClient(scheme, &cacheInformers, objects...)
			reconcileRateLimiter := ratelimiter.ReconcileRateLimiter{
				Client: client,
				Scheme: scheme,
			}
			controller := controller.NewStubController(client, reconcileRateLimiter, cacheInformers, *scheme)

			rateLimiterSource := &source.Kind{Type: &v12.RateLimiter{}}
			inject.CacheInto(cacheInformers, rateLimiterSource)
			
			configMapSource := &source.Kind{Type: &v1.ConfigMap{}}
			inject.CacheInto(cacheInformers, configMapSource)


			controller.Watch(rateLimiterSource, &handler.EnqueueRequestForObject{})
			controller.Watch(configMapSource, &handler.EnqueueRequestForOwner{
				OwnerType:    &v12.RateLimiter{},
				IsController: true,
			})

			request := &reconcile.Request{NamespacedName: types.NamespacedName{Namespace: namespace, Name: name}}

			//reconcile loop (creation of primary resource should be a reason of reconcile call)
			controller.Reconcile(*request)

			//check for created objects
			dep := &appsv1.Deployment{}
			confMap := &v1.ConfigMap{}
			srv := &v1.Service{}
			redisDep := &appsv1.Deployment{}
			var err error

			err = client.Get(context.TODO(), request.NamespacedName, dep)
			Expect(err).To(BeNil(), "Deployment is nil")

			err = client.Get(context.TODO(), request.NamespacedName, redisDep)
			Expect(err).To(BeNil(), "Redis Deployment is nil")

			err = client.Get(context.TODO(), request.NamespacedName, confMap)
			Expect(err).To(BeNil(), "ConfigMap is nil")

			err = client.Get(context.TODO(), request.NamespacedName, srv)
			Expect(err).To(BeNil(), "Service is nil")
		})
	})
})

//Decides whether do Reconcile() again. Decision based on Result and error according to reconcile loop rules
func doReconcile(reconciler *ratelimiter.ReconcileRateLimiter, request *reconcile.Request) {
	condition := false
	for ok := true; ok; ok = condition {
		result, err := reconciler.Reconcile(*request)
		switch {
		case err != nil || result.Requeue:
			condition = true
		default:
			condition = false
		}
	}
}

func getPreparedRateLimiterResource() *v12.RateLimiter {
	return &v12.RateLimiter{
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: namespace,
		},
	}
}

func getPreparedDeploymentResource() *appsv1.Deployment {
	return &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: namespace,
		},
	}
}

func getPreparedRedisDeploymentResource(RedisDeploymentName string) *appsv1.Deployment {
	return &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name:      RedisDeploymentName,
			Namespace: namespace,
		},
	}
}

func getPreparedServiceResource() *v1.Service {
	return &v1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: namespace,
		},
	}
}

func getPreparedConfigMapResource() *v1.ConfigMap {
	return &v1.ConfigMap{
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: namespace,
		},
	}
}
