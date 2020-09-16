package ratelimiterconfig

import (
	istio_v1alpha3 "istio.io/client-go/pkg/apis/networking/v1alpha3"
	corev1 "k8s.io/api/core/v1"
	apiextensionsv1beta1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1beta1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	clientgoscheme "k8s.io/client-go/kubernetes/scheme"
	"github.com/voteva/ratelimit-operator/pkg/apis"
	v1 "github.com/voteva/ratelimit-operator/pkg/apis/operators/v1"
	"github.com/voteva/ratelimit-operator/pkg/utils"
	"sigs.k8s.io/controller-runtime/pkg/client/fake"
)

func buildRateLimiterConfig(rl *v1.RateLimiter) *v1.RateLimiterConfig {
	host := utils.BuildRandomString(3)
	failureModeDeny := true
	rateLimitRequestTimeout := "0.25s"
	return &v1.RateLimiterConfig{
		ObjectMeta: metav1.ObjectMeta{
			Name:      utils.BuildRandomString(3),
			Namespace: rl.Namespace,
		},
		Spec: v1.RateLimiterConfigSpec{
			ApplyTo:     v1.GATEWAY,
			Host:        &host,
			Port:        int32(utils.BuildRandomInt(2)),
			RateLimiter: rl.Name,
			Descriptors: []v1.Descriptor{{
				Key: utils.BuildRandomString(3),
			}},
			RateLimitRequestTimeout: &rateLimitRequestTimeout,
			FailureModeDeny:         &failureModeDeny,
		},
	}
}

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

func buildReconciler(rateLimiter *v1.RateLimiter) *ReconcileRateLimiterConfig {
	scheme := buildScheme()
	objects := []runtime.Object{rateLimiter}
	client := fake.NewFakeClientWithScheme(scheme, objects...)
	return &ReconcileRateLimiterConfig{
		Client:      client,
		Scheme:      scheme,
		rateLimiter: rateLimiter,
		configMap:   buildConfigMap(rateLimiter),
	}
}

func buildEmptyReconciler() *ReconcileRateLimiterConfig {
	scheme := buildScheme()
	objects := []runtime.Object{}
	client := fake.NewFakeClientWithScheme(scheme, objects...)
	return &ReconcileRateLimiterConfig{
		Client: client,
		Scheme: scheme,
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

func buildNamespacedName(rateLimiterConfig *v1.RateLimiterConfig) types.NamespacedName {
	return types.NamespacedName{
		Name:      rateLimiterConfig.Name,
		Namespace: rateLimiterConfig.Namespace,
	}
}
