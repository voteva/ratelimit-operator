package ratelimiter

import (
	"context"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/equality"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/apimachinery/pkg/util/intstr"
	v1 "ratelimit-operator/pkg/apis/operators/v1"
	"ratelimit-operator/pkg/constants"
	"ratelimit-operator/pkg/utils"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
)

func (r *ReconcileRateLimiter) reconcileServiceForRedis(ctx context.Context, instance *v1.RateLimiter) (reconcile.Result, error) {
	reqLogger := log.WithValues("Instance.Name", instance.Name)

	foundService := &corev1.Service{}
	serviceFromInstance := buildServiceForRedis(instance)
	_ = controllerutil.SetControllerReference(instance, serviceFromInstance, r.scheme)

	err := r.client.Get(ctx, types.NamespacedName{Name: serviceFromInstance.Name, Namespace: instance.Namespace}, foundService)
	if err != nil {
		if errors.IsNotFound(err) {
			reqLogger.Info("Creating a new Service Redis")
			err = r.client.Create(ctx, serviceFromInstance)
			if err != nil {
				reqLogger.Error(err, "Failed to create new Service Redis")
				return reconcile.Result{}, err
			}
			return reconcile.Result{Requeue: true}, nil
		} else {
			reqLogger.Error(err, "Failed to get Service Redis")
			return reconcile.Result{}, err
		}
	}

	if !equality.Semantic.DeepEqual(foundService.Spec, serviceFromInstance.Spec) {
		serviceFromInstance.Spec.ClusterIP = foundService.Spec.ClusterIP
		foundService.Spec = serviceFromInstance.Spec
		r.client.Update(ctx, foundService)
	}

	return reconcile.Result{}, nil
}

func buildServiceForRedis(instance *v1.RateLimiter) *corev1.Service {
	serviceName := buildNameForRedis(instance.Name)
	servicePort := constants.REDIS_PORT

	service := &corev1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name:      serviceName,
			Namespace: instance.Namespace,
		},
		Spec: corev1.ServiceSpec{
			Ports: []corev1.ServicePort{{
				Name:       serviceName,
				Protocol:   corev1.ProtocolTCP,
				Port:       servicePort,
				TargetPort: intstr.IntOrString{Type: intstr.Int, IntVal: servicePort},
			}},
			Selector: utils.SelectorsForApp(serviceName),
		},
	}
	return service
}
