package ratelimiter

import (
	"context"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"ratelimit-operator/pkg/apis/operators/v1"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
)

func (r *ReconcileRateLimiter) reconcileConfigMap(ctx context.Context, instance *v1.RateLimiter) (reconcile.Result, error) {
	foundConfigMap := &corev1.ConfigMap{}

	err := r.client.Get(ctx, types.NamespacedName{Name: instance.Name, Namespace: instance.Namespace}, foundConfigMap)
	if err != nil && errors.IsNotFound(err) {
		cm := r.buildConfigMap(instance)
		log.Info("Creating a new ConfigMap", "ConfigMap.Namespace", cm.Namespace, "ConfigMap.Name", cm.Name)
		err = r.client.Create(ctx, cm)
		if err != nil {
			log.Error(err, "Failed to create new ConfigMap", "ConfigMap.Namespace", cm.Namespace, "ConfigMap.Name", cm.Name)
			return reconcile.Result{}, err
		}
		return reconcile.Result{Requeue: true}, nil
	} else if err != nil {
		log.Error(err, "Failed to get ConfigMap")
		return reconcile.Result{}, err
	}
	return reconcile.Result{}, nil
}

func (r *ReconcileRateLimiter) buildConfigMap(m *v1.RateLimiter) *corev1.ConfigMap {
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
