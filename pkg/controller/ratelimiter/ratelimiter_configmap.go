package ratelimiter

import (
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	operatorsv1alpha1 "ratelimit-operator/pkg/apis/operators/v1alpha1"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
)

func (r *ReconcileRateLimiter) ConfigMapForRateLimiter(m *operatorsv1alpha1.RateLimiter) *corev1.ConfigMap {
	configMap := &corev1.ConfigMap{
		ObjectMeta: metav1.ObjectMeta{
			Name:      m.Name,
			Namespace: m.Namespace,
		},
		Data: map[string]string{
			"LOG_LEVEL":            "DEBUG",
			"REDIS_SOCKET_TYPE":    "tcp",
			"REDIS_URL":            "localhost:6379",
			"RUNTIME_ROOT":         "/data/ratelimit",
			"RUNTIME_SUBDIRECTORY": "config",
			"USE_STATSD":           "false",
			"rate_limit.property": `
				domain: test
				descriptors:
				  - key: custom-rl-header
					value: setting1
					rate_limit:
					  unit: minute
					  requests_per_unit: 1`,
		},
	}
	controllerutil.SetControllerReference(m, configMap, r.scheme)
	return configMap
}
