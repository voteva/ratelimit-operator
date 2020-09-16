package ratelimiter

import (
	istio_v1alpha3 "istio.io/client-go/pkg/apis/networking/v1alpha3"
	apiextensionsv1beta1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1beta1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	clientgoscheme "k8s.io/client-go/kubernetes/scheme"
	"github.com/voteva/ratelimit-operator/pkg/apis"
	v1 "github.com/voteva/ratelimit-operator/pkg/apis/operators/v1"
	"github.com/voteva/ratelimit-operator/pkg/utils"
	"sigs.k8s.io/controller-runtime/pkg/client/fake"
)

func buildRateLimiter() *v1.RateLimiter {
	logLevel := v1.INFO
	size := int32(1)

	return &v1.RateLimiter{
		ObjectMeta: metav1.ObjectMeta{
			Name:      utils.BuildRandomString(3),
			Namespace: utils.BuildRandomString(3),
		},
		Spec: v1.RateLimiterSpec{
			LogLevel: &logLevel,
			Size:     &size,
		},
		Status: v1.RateLimiterStatus{},
	}
}

func buildReconciler(rateLimiter *v1.RateLimiter) *ReconcileRateLimiter {
	scheme := buildScheme()
	objects := []runtime.Object{rateLimiter}
	client := fake.NewFakeClientWithScheme(scheme, objects...)
	return &ReconcileRateLimiter{Client: client, Scheme: scheme}
}

func buildEmptyReconciler() *ReconcileRateLimiter {
	scheme := buildScheme()
	objects := []runtime.Object{}
	client := fake.NewFakeClientWithScheme(scheme, objects...)
	return &ReconcileRateLimiter{Client: client, Scheme: scheme}
}

func buildScheme() *runtime.Scheme {
	scheme := runtime.NewScheme()
	_ = apis.AddToScheme(scheme)
	_ = clientgoscheme.AddToScheme(scheme)
	_ = apiextensionsv1beta1.AddToScheme(scheme)
	_ = istio_v1alpha3.AddToScheme(scheme)
	return scheme
}
