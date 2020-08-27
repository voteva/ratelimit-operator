package ratelimiterconfig

import (
	istio_v1alpha3 "istio.io/client-go/pkg/apis/networking/v1alpha3"
	corev1 "k8s.io/api/core/v1"
	apiextensionsv1beta1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1beta1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	clientgoscheme "k8s.io/client-go/kubernetes/scheme"
	"ratelimit-operator/pkg/apis"
	v1 "ratelimit-operator/pkg/apis/operators/v1"
	"ratelimit-operator/pkg/utils"
	"sigs.k8s.io/controller-runtime/pkg/client/fake"
)

func buildRateLimiterConfig(rl *v1.RateLimiter) *v1.RateLimiterConfig {
	host := utils.BuildRandomString(3)
	return &v1.RateLimiterConfig{
		ObjectMeta: metav1.ObjectMeta{
			Name:      utils.BuildRandomString(3),
			Namespace: rl.Namespace,
		},
		Spec: v1.RateLimiterConfigSpec{
			RateLimiter: rl.Name,
			ApplyTo:     v1.GATEWAY,
			Host:        &host,
			Port:        int32(utils.BuildRandomInt(2)),
			RateLimitProperty: v1.RateLimitProperty{
				Domain: utils.BuildRandomString(3),
				Descriptors: []v1.Descriptor{{
					Key: utils.BuildRandomString(3),
				}},
			},
			FailureModeDeny: true,
		},
	}
}

func buildRateLimiter() *v1.RateLimiter {
	logLevel := v1.INFO
	port := int32(utils.BuildRandomInt(2))
	size := int32(1)

	return &v1.RateLimiter{
		ObjectMeta: metav1.ObjectMeta{
			Name:      utils.BuildRandomString(3),
			Namespace: utils.BuildRandomString(3),
		},
		Spec: v1.RateLimiterSpec{
			LogLevel: &logLevel,
			Port:     &port,
			Size:     &size,
		},
		Status: v1.RateLimiterStatus{},
	}
}

func buildReconciler(rateLimiter *v1.RateLimiter) *ReconcileRateLimiterConfig {
	scheme := buildScheme()
	objects := []runtime.Object{rateLimiter}
	client := fake.NewFakeClientWithScheme(scheme, objects...)
	return &ReconcileRateLimiterConfig{
		client:      client,
		scheme:      scheme,
		rateLimiter: rateLimiter,
		configMap:   buildConfigMap(rateLimiter),
	}
}

func buildEmptyReconciler() *ReconcileRateLimiterConfig {
	scheme := buildScheme()
	objects := []runtime.Object{}
	client := fake.NewFakeClientWithScheme(scheme, objects...)
	return &ReconcileRateLimiterConfig{
		client: client,
		scheme: scheme,
	}
}

func buildScheme() *runtime.Scheme {
	scheme := runtime.NewScheme()
	_ = apis.AddToScheme(scheme)
	_ = clientgoscheme.AddToScheme(scheme)
	_ = apiextensionsv1beta1.AddToScheme(scheme)
	_ = istio_v1alpha3.AddToScheme(scheme)
	return scheme
}

func buildConfigMap(instance *v1.RateLimiter) *corev1.ConfigMap {
	configMap := &corev1.ConfigMap{
		ObjectMeta: metav1.ObjectMeta{
			Name:      instance.Name,
			Namespace: instance.Namespace,
		},
		Data: map[string]string{},
	}
	return configMap
}
