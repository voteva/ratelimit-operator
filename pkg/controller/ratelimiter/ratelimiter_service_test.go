package ratelimiter

import (
	"context"
	"github.com/stretchr/testify/assert"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/apimachinery/pkg/util/intstr"
	v1 "ratelimit-operator/pkg/apis/operators/v1"
	"ratelimit-operator/pkg/utils"
	"testing"
)

func Test_ReconcileServiceForService_Create(t *testing.T) {
	t.Parallel()
	a := assert.New(t)

	t.Run("reconcile service for ratelimit-service (Create)", func(t *testing.T) {
		rateLimiter := buildRateLimiter()
		r := buildReconciler(rateLimiter)

		reconcileResult, err := r.reconcileServiceForService(context.Background(), rateLimiter)

		foundService := &corev1.Service{}
		namespaceName := buildServiceResourceNamespacedName(rateLimiter)
		errGet := r.client.Get(context.Background(), namespaceName, foundService)

		a.Nil(err)
		a.NotNil(reconcileResult)
		a.True(reconcileResult.Requeue)
		a.Nil(errGet)
		a.NotNil(foundService)
	})
}

func Test_BuildService(t *testing.T) {
	t.Parallel()
	a := assert.New(t)

	t.Run("success build Service for ratelimit-service", func(t *testing.T) {
		rateLimiter := buildRateLimiter()

		actualResult := buildService(rateLimiter)

		a.Equal(rateLimiter.Name, actualResult.Name)
		a.Equal(rateLimiter.Namespace, actualResult.Namespace)
		a.Equal(utils.SelectorsForApp(rateLimiter.Name), actualResult.Spec.Selector)
		a.Equal(1, len(actualResult.Spec.Ports))
		a.Equal("grpc-"+rateLimiter.Name, actualResult.Spec.Ports[0].Name)
		a.Equal(corev1.ProtocolTCP, actualResult.Spec.Ports[0].Protocol)
		a.Equal(*rateLimiter.Spec.Port, actualResult.Spec.Ports[0].Port)
		a.Equal(intstr.IntOrString{Type: intstr.Int, IntVal: *rateLimiter.Spec.Port}, actualResult.Spec.Ports[0].TargetPort)
	})
}

func buildServiceResourceNamespacedName(rateLimiter *v1.RateLimiter) types.NamespacedName {
	return types.NamespacedName{
		Name:      rateLimiter.Name,
		Namespace: rateLimiter.Namespace,
	}
}
