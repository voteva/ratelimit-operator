package ratelimiter

import (
	"context"
	"github.com/stretchr/testify/assert"
	corev1 "k8s.io/api/core/v1"
	"testing"
)

func Test_ReconcileConfigMap_CreateSuccess(t *testing.T) {
	t.Parallel()
	a := assert.New(t)

	t.Run("reconcile ConfigMap (CreateSuccess)", func(t *testing.T) {
		rateLimiter := buildRateLimiter()
		r := buildReconciler(rateLimiter)

		reconcileResult, err := r.reconcileConfigMap(context.Background(), rateLimiter)

		foundConfigMap := &corev1.ConfigMap{}
		namespaceName := buildServiceResourceNamespacedName(rateLimiter)
		errGet := r.Client.Get(context.Background(), namespaceName, foundConfigMap)

		a.Nil(err)
		a.NotNil(reconcileResult)
		a.True(reconcileResult.Requeue)
		a.Nil(errGet)
		a.NotNil(foundConfigMap)
	})
}

func Test_ReconcileConfigMap_CreateError(t *testing.T) {
	t.Parallel()
	a := assert.New(t)

	t.Run("reconcile ConfigMap (CreateError)", func(t *testing.T) {
		rateLimiter := buildRateLimiter()
		rateLimiter.Name = ""
		rateLimiter.Namespace = ""
		r := buildReconciler(rateLimiter)

		_, err := r.reconcileConfigMap(context.Background(), rateLimiter)

		a.NotNil(err)
	})
}

func Test_BuildConfigMap(t *testing.T) {
	t.Parallel()
	a := assert.New(t)

	t.Run("success build ConfigMap for ratelimit-service", func(t *testing.T) {
		rateLimiter := buildRateLimiter()

		actualResult := buildConfigMap(rateLimiter)

		a.Equal(rateLimiter.Name, actualResult.Name)
		a.Equal(rateLimiter.Namespace, actualResult.Namespace)
		a.Equal(map[string]string{}, actualResult.Data)
	})
}
