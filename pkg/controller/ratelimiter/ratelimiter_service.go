package ratelimiter

import (
	"context"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/types"
	operatorsv1alpha1 "ratelimit-operator/pkg/apis/operators/v1alpha1"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
)

func (r *ReconcileRateLimiter) reconcileService(request reconcile.Request, instance *operatorsv1alpha1.RateLimiter) (reconcile.Result, error) {
	reqLogger := log.WithValues("Request.Namespace", request.Namespace, "Request.Name", request.Name)

	foundService := &corev1.Service{}
	err := r.client.Get(context.TODO(), types.NamespacedName{Name: instance.Name, Namespace: instance.Namespace}, foundService)
	if err != nil && errors.IsNotFound(err) {
		// Define a new Service
		cm := r.buildService(instance)
		reqLogger.Info("Creating a new Service", "Service.Namespace", cm.Namespace, "Service.Name", cm.Name)
		err = r.client.Create(context.TODO(), cm)
		if err != nil {
			reqLogger.Error(err, "Failed to create new Service", "Service.Namespace", cm.Namespace, "Service.Name", cm.Name)
			return reconcile.Result{}, err
		}
		// Deployment created successfully - return and requeue
		return reconcile.Result{Requeue: true}, nil
	} else if err != nil {
		reqLogger.Error(err, "Failed to get Service")
		return reconcile.Result{}, err
	}
	return reconcile.Result{}, nil
}

func (r *ReconcileRateLimiter) buildService(m *operatorsv1alpha1.RateLimiter) *corev1.Service {
	service := &corev1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name:      m.Name,
			Namespace: m.Namespace,
		},
		Spec: corev1.ServiceSpec{
			Ports: []corev1.ServicePort{{
				Name:       "grpc-rate-limiter",
				Protocol:   "TCP",
				Port:       8081,
				TargetPort: intstr.IntOrString{Type: intstr.Int, IntVal: 8081},
			}},
			Selector: map[string]string{
				"app": "rate-limit-server",
			},
		},
	}
	controllerutil.SetControllerReference(m, service, r.scheme)
	return service
}
