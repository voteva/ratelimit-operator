package ratelimiter

import (
	corev1 "k8s.io/api/core/v1"
	operatorsv1alpha1 "ratelimit-operator/pkg/apis/operators/v1alpha1"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
)

func (r *ReconcileRateLimiter) ServiceForRateLimiter(m *operatorsv1alpha1.RateLimiter) *corev1.Service {
	service := &corev1.Service{
		// TODO implement
	}
	controllerutil.SetControllerReference(m, service, r.scheme)
	return service
}
