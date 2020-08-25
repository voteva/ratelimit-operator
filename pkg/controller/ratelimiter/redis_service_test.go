package ratelimiter

import (
	"context"
	"github.com/stretchr/testify/assert"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/apimachinery/pkg/util/intstr"
	v1 "ratelimit-operator/pkg/apis/operators/v1"
	"ratelimit-operator/pkg/constants"
	"ratelimit-operator/pkg/utils"
	"testing"
)

func Test_ReconcileServiceForRedis_Create(t *testing.T) {
	t.Parallel()
	a := assert.New(t)

	t.Run("reconcile service for Redis (Create)", func(t *testing.T) {
		rateLimiter := buildRateLimiter()
		r := buildReconciler(rateLimiter)

		reconcileResult, err := r.reconcileServiceForRedis(context.Background(), rateLimiter)

		foundService := &corev1.Service{}
		namespaceName := buildRedisResourceNamespacedName(rateLimiter)
		errGet := r.client.Get(context.Background(), namespaceName, foundService)

		a.Nil(err)
		a.NotNil(reconcileResult)
		a.True(reconcileResult.Requeue)
		a.Nil(errGet)
		a.NotNil(foundService)
	})
}

func Test_BuildServiceForRedis(t *testing.T) {
	t.Parallel()
	a := assert.New(t)

	t.Run("success build Service for Redis", func(t *testing.T) {
		rateLimiter := buildRateLimiter()
		expectedServiceName := buildNameForRedis(rateLimiter.Name)

		actualResult := buildServiceForRedis(rateLimiter)

		a.Equal(expectedServiceName, actualResult.Name)
		a.Equal(rateLimiter.Namespace, actualResult.Namespace)
		a.Equal(1, len(actualResult.Spec.Ports))
		a.Equal(expectedServiceName, actualResult.Spec.Ports[0].Name)
		a.Equal(corev1.ProtocolTCP, actualResult.Spec.Ports[0].Protocol)
		a.Equal(constants.REDIS_PORT, actualResult.Spec.Ports[0].Port)
		a.Equal(intstr.IntOrString{Type: intstr.Int, IntVal: constants.REDIS_PORT}, actualResult.Spec.Ports[0].TargetPort)
		a.Equal(utils.SelectorsForApp(expectedServiceName), actualResult.Spec.Selector)
	})
}

func buildRedisResourceNamespacedName(rateLimiter *v1.RateLimiter) types.NamespacedName {
	return types.NamespacedName{
		Name:      buildNameForRedis(rateLimiter.Name),
		Namespace: rateLimiter.Namespace,
	}
}
