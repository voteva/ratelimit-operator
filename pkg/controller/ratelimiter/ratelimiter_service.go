package ratelimiter

import (
	"context"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
	"ratelimit-operator/pkg/constants"
	"ratelimit-operator/pkg/utils"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/types"
	"ratelimit-operator/pkg/apis/operators/v1"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
)

func (r *ReconcileRateLimiter) reconcileService(ctx context.Context, instance *v1.RateLimiter) (reconcile.Result, error) {
	foundService := &corev1.Service{}

	err := r.client.Get(ctx, types.NamespacedName{Name: instance.Name, Namespace: instance.Namespace}, foundService)
	if err != nil && errors.IsNotFound(err) {
		srv := r.buildService(instance)
		log.Info("Creating a new Service", "Service.Namespace", srv.Namespace, "Service.Name", srv.Name)
		err = r.client.Create(ctx, srv)
		if err != nil {
			log.Error(err, "Failed to create new Service", "Service.Namespace", srv.Namespace, "Service.Name", srv.Name)
			return reconcile.Result{}, err
		}
		return reconcile.Result{Requeue: true}, nil
	} else if err != nil {
		log.Error(err, "Failed to get Service")
		return reconcile.Result{}, err
	}
	return reconcile.Result{}, nil
}

func (r *ReconcileRateLimiter) buildService(instance *v1.RateLimiter) *corev1.Service {
	servicePort := utils.DefaultIfZero(instance.Spec.ServicePort, constants.DEFAULT_PORT)

	service := &corev1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name:      instance.Name,
			Namespace: instance.Namespace,
		},
		Spec: corev1.ServiceSpec{
			Ports: []corev1.ServicePort{{
				Name:       instance.Name,
				Protocol:   corev1.ProtocolTCP,
				Port:       servicePort,
				TargetPort: intstr.IntOrString{Type: intstr.Int, IntVal: servicePort},
			}},
			Selector: SelectorsForRateLimiter(instance.Name),
		},
	}
	controllerutil.SetControllerReference(instance, service, r.scheme)
	return service
}
