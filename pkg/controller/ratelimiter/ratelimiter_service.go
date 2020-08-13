package ratelimiter

import (
	"context"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/equality"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/apimachinery/pkg/util/intstr"
	"ratelimit-operator/pkg/apis/operators/v1"
	"ratelimit-operator/pkg/utils"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
)

func (r *ReconcileRateLimiter) reconcileServiceForService(ctx context.Context, instance *v1.RateLimiter) (reconcile.Result, error) {
	reqLogger := log.WithValues("Instance.Name", instance.Name)

	foundService := &corev1.Service{}
	serviceFromInstance := r.buildService(instance)

	err := r.Client.Get(ctx, types.NamespacedName{Name: instance.Name, Namespace: instance.Namespace}, foundService)
	if err != nil {
		if errors.IsNotFound(err) {
			reqLogger.Info("Creating a new Service")
			err = r.Client.Create(ctx, serviceFromInstance)
			if err != nil {
				reqLogger.Error(err, "Failed to create new Service")
				return reconcile.Result{}, err
			}
			return reconcile.Result{Requeue: true}, nil
		} else {
			reqLogger.Error(err, "Failed to get Service")
			return reconcile.Result{}, err
		}
	}

	if !equality.Semantic.DeepEqual(foundService.Spec, serviceFromInstance.Spec) {
		foundService.Spec = serviceFromInstance.Spec
		r.Client.Update(ctx, foundService)
	}

	return reconcile.Result{}, nil
}

func (r *ReconcileRateLimiter) buildService(instance *v1.RateLimiter) *corev1.Service {
	port := r.buildRateLimiterServicePort(instance)

	service := &corev1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name:      instance.Name,
			Namespace: instance.Namespace,
		},
		Spec: corev1.ServiceSpec{
			Ports: []corev1.ServicePort{{
				Name:       "grpc-" + instance.Name,
				Protocol:   corev1.ProtocolTCP,
				Port:       port,
				TargetPort: intstr.IntOrString{Type: intstr.Int, IntVal: port},
			}},
			Selector: utils.SelectorsForApp(instance.Name),
		},
	}
	controllerutil.SetControllerReference(instance, service, r.Scheme)
	return service
}
