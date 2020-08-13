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
	reqLogger := log.WithValues("Instance.Name", instance.Name)

	foundConfigMap := &corev1.ConfigMap{}

	err := r.Client.Get(ctx, types.NamespacedName{Name: instance.Name, Namespace: instance.Namespace}, foundConfigMap)
	if err != nil {
		if errors.IsNotFound(err) {
			configMapFromInstance := r.buildConfigMap(instance)
			reqLogger.Info("Creating a new ConfigMap")
			err = r.Client.Create(ctx, configMapFromInstance)
			if err != nil {
				reqLogger.Error(err, "Failed to create new ConfigMap")
				return reconcile.Result{}, err
			}
			return reconcile.Result{Requeue: true}, nil
		} else {
			reqLogger.Error(err, "Failed to get ConfigMap")
			return reconcile.Result{}, err
		}
	}
	return reconcile.Result{}, nil
}

func (r *ReconcileRateLimiter) buildConfigMap(instance *v1.RateLimiter) *corev1.ConfigMap {
	configMap := &corev1.ConfigMap{
		ObjectMeta: metav1.ObjectMeta{
			Name:      instance.Name,
			Namespace: instance.Namespace,
		},
		Data: map[string]string{},
	}
	controllerutil.SetControllerReference(instance, configMap, r.Scheme)
	return configMap
}
