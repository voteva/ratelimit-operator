package ratelimiter

import (
	"context"

	"istio.io/client-go/pkg/apis/networking/v1alpha3"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/types"
	operatorsv1alpha1 "ratelimit-operator/pkg/apis/operators/v1alpha1"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
)

func (r *ReconcileRateLimiter) reconcileEnvoyFilter(request reconcile.Request, instance *operatorsv1alpha1.RateLimiter) (reconcile.Result, error) {
	reqLogger := log.WithValues("Request.Namespace", request.Namespace, "Request.Name", request.Name)

	foundEnvoyFilter := &v1alpha3.EnvoyFilter{}
	err := r.client.Get(context.TODO(), types.NamespacedName{Name: instance.Name, Namespace: instance.Namespace}, foundEnvoyFilter)
	if err != nil && errors.IsNotFound(err) {
		// Define a new EnvoyFilter
		cm := r.buildEnvoyFilter(instance)
		reqLogger.Info("Creating a new EnvoyFilter", "EnvoyFilter.Namespace", cm.Namespace, "EnvoyFilter.Name", cm.Name)
		err = r.client.Create(context.TODO(), cm)
		if err != nil {
			reqLogger.Error(err, "Failed to create new EnvoyFilter", "EnvoyFilter.Namespace", cm.Namespace, "EnvoyFilter.Name", cm.Name)
			return reconcile.Result{}, err
		}
		// Deployment created successfully - return and requeue
		return reconcile.Result{Requeue: true}, nil
	} else if err != nil {
		reqLogger.Error(err, "Failed to get EnvoyFilter")
		return reconcile.Result{}, err
	}
	return reconcile.Result{}, nil
}

func (r *ReconcileRateLimiter) buildEnvoyFilter(m *operatorsv1alpha1.RateLimiter) *v1alpha3.EnvoyFilter {
	envoyFilter := &v1alpha3.EnvoyFilter{
		// TODO implement
	}
	controllerutil.SetControllerReference(m, envoyFilter, r.scheme)
	return envoyFilter
}
