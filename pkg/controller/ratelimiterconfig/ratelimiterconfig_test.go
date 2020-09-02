package ratelimiterconfig_test

import (
	"context"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"istio.io/client-go/pkg/apis/networking/v1alpha3"
	istio_v1alpha3 "istio.io/client-go/pkg/apis/networking/v1alpha3"
	v13 "k8s.io/api/core/v1"
	apiextensionsv1beta1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1beta1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	clientgoscheme "k8s.io/client-go/kubernetes/scheme"
	"ratelimit-operator/pkg/apis"
	v1 "ratelimit-operator/pkg/apis/operators/v1"
	v12 "ratelimit-operator/pkg/apis/operators/v1"
	"ratelimit-operator/pkg/controller/ratelimiterconfig"
	"sigs.k8s.io/controller-runtime/pkg/client/fake"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
	"time"
)

var scheme *runtime.Scheme

func init() {
	scheme = runtime.NewScheme()
	apis.AddToScheme(scheme)
	clientgoscheme.AddToScheme(scheme)
	apiextensionsv1beta1.AddToScheme(scheme)
	istio_v1alpha3.AddToScheme(scheme)
}

var _ = Describe("Ratelimiterconfig", func() {

	Context("when RateLimiterConfig (primary resource) added", func() {
		It("reconcile RateLimitConfig. create secondary resources", func() {
			namespace := "test-namespace"
			ratelimiterconfigName := "test-ratelimiterconfig-name"
			ratelimiterName := "test-ratelimiter-name"

			rateLimiterConfigResource := getPreparedRateLimiterConfigResource(ratelimiterconfigName, namespace, ratelimiterName)
			rateLimiterResource := getPreparedRateLimiterResource(ratelimiterName, namespace)
			configMapResource := getPreparedConfigMapResource(ratelimiterName, namespace)

			objects := []runtime.Object{rateLimiterConfigResource, rateLimiterResource, configMapResource}
			client := fake.NewFakeClientWithScheme(scheme, objects...)
			reconciler := &ratelimiterconfig.ReconcileRateLimiterConfig{Client: client, Scheme: scheme}
			request := &reconcile.Request{NamespacedName: types.NamespacedName{Namespace: namespace, Name: ratelimiterconfigName}}

			//reconcile loop (creation of primary resource should be a reason of reconcile call)
			doReconcile(reconciler, request)

			//check for created objects
			envFltr := &v1alpha3.EnvoyFilter{}

			err := reconciler.Client.Get(context.TODO(), request.NamespacedName, envFltr)
			Expect(err).To(BeNil(), "EnvoyFilter is nil")
		})
	})

	Context("when RateLimiterConfig (primary resource) added. Create RateLimit resource after delay", func() {
		It("reconcile RateLimitConfig. create secondary resources", func() {
			namespace := "test-namespace"
			ratelimiterconfigName := "test-ratelimiterconfig-name"
			ratelimiterName := "test-ratelimiter-name"

			rateLimiterConfigResource := getPreparedRateLimiterConfigResource(ratelimiterconfigName, namespace, ratelimiterName)
			rateLimiterResource := getPreparedRateLimiterResource(ratelimiterName, namespace)
			configMapResource := getPreparedConfigMapResource(ratelimiterName, namespace)

			objects := []runtime.Object{rateLimiterConfigResource}
			client := fake.NewFakeClientWithScheme(scheme, objects...)
			reconciler := &ratelimiterconfig.ReconcileRateLimiterConfig{Client: client, Scheme: scheme}
			request := &reconcile.Request{NamespacedName: types.NamespacedName{Namespace: namespace, Name: ratelimiterconfigName}}

			//initation of RateLimiter resource creation after delay
			time.AfterFunc(7*time.Second, func() {
				client.Create(context.TODO(), rateLimiterResource)
				client.Create(context.TODO(), configMapResource)
			})
			//reconcile loop (creation of primary resource should be a reason of reconcile call)
			doReconcile(reconciler, request)

			//check for created objects
			envFltr := &v1alpha3.EnvoyFilter{}

			err := reconciler.Client.Get(context.TODO(), request.NamespacedName, envFltr)
			Expect(err).To(BeNil(), "EnvoyFilter is nil")
		})
	})

})

//Decides whether do Reconcile() again. Decision based on Result and error according to reconcile loop rules
func doReconcile(reconciler *ratelimiterconfig.ReconcileRateLimiterConfig, request *reconcile.Request) {
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

func getPreparedRateLimiterConfigResource(Name string, Namespace string, RateLimiter string) *v1.RateLimiterConfig {
	return &v1.RateLimiterConfig{
		ObjectMeta: metav1.ObjectMeta{
			Name:      Name,
			Namespace: Namespace,
		},
		Spec: v1.RateLimiterConfigSpec{
			ApplyTo:     "GATEWAY",
			RateLimiter: RateLimiter,
		},
	}
}

func getPreparedEnvoyFilterResource(instance *v1.RateLimiterConfig) *v1alpha3.EnvoyFilter {
	return &v1alpha3.EnvoyFilter{
		ObjectMeta: metav1.ObjectMeta{
			Name:      instance.Name,
			Namespace: instance.Namespace,
		},
	}
}

func getPreparedRateLimiterResource(Name string, Namespace string) *v12.RateLimiter {
	return &v12.RateLimiter{
		ObjectMeta: metav1.ObjectMeta{
			Name:      Name,
			Namespace: Namespace,
		},
	}
}

func getPreparedConfigMapResource(Name string, Namespace string) *v13.ConfigMap {
	return &v13.ConfigMap{
		ObjectMeta: metav1.ObjectMeta{
			Name:      Name,
			Namespace: Namespace,
		},
	}
}
