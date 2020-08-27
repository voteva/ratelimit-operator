package ratelimiter

import (
	"context"
	"github.com/stretchr/testify/assert"
	appsv1 "k8s.io/api/apps/v1"
	"ratelimit-operator/pkg/utils"
	"testing"
)

func Test_ReconcileDeploymentForRedis_Create(t *testing.T) {
	t.Parallel()
	a := assert.New(t)

	t.Run("reconcile deployment for Redis (Create)", func(t *testing.T) {
		rateLimiter := buildRateLimiter()
		r := buildReconciler(rateLimiter)

		reconcileResult, err := r.reconcileDeploymentForRedis(context.Background(), rateLimiter)

		foundDeployment := &appsv1.Deployment{}
		namespaceName := buildRedisResourceNamespacedName(rateLimiter)
		errGet := r.client.Get(context.Background(), namespaceName, foundDeployment)

		a.Nil(err)
		a.NotNil(reconcileResult)
		a.True(reconcileResult.Requeue)
		a.Nil(errGet)
		a.NotNil(foundDeployment)
	})
}

func Test_ReconcileDeploymentForRedis_Update(t *testing.T) {
	t.Parallel()
	a := assert.New(t)

	t.Run("reconcile deployment for Redis (Update)", func(t *testing.T) {
		rateLimiter := buildRateLimiter()
		r := buildReconciler(rateLimiter)

		dep := buildDeploymentForRedis(rateLimiter)
		newReplicas := int32(10)
		dep.Spec.Replicas = &newReplicas
		errCreateSrvRL := r.client.Create(context.Background(), dep)
		a.Nil(errCreateSrvRL)

		reconcileResult, err := r.reconcileDeploymentForRedis(context.Background(), rateLimiter)

		foundDeployment := &appsv1.Deployment{}
		namespaceName := buildRedisResourceNamespacedName(rateLimiter)
		errGet := r.client.Get(context.Background(), namespaceName, foundDeployment)

		a.Nil(err)
		a.NotNil(reconcileResult)
		a.False(reconcileResult.Requeue)
		a.Nil(errGet)
		a.NotNil(foundDeployment)
		a.Equal(int32(1), *foundDeployment.Spec.Replicas)
	})
}

func Test_BuildDeploymentForRedis(t *testing.T) {
	t.Parallel()
	a := assert.New(t)

	t.Run("success build Deployment for Redis", func(t *testing.T) {
		rateLimiter := buildRateLimiter()

		expectedDeploymentName := buildNameForRedis(rateLimiter.Name)
		expectedLabels := utils.LabelsForApp(expectedDeploymentName)
		expectedReplicas := int32(1)

		actualResult := buildDeploymentForRedis(rateLimiter)

		a.Equal(expectedDeploymentName, actualResult.Name)
		a.Equal(rateLimiter.Namespace, actualResult.Namespace)
		a.Equal(&expectedReplicas, actualResult.Spec.Replicas)
		a.Equal(expectedLabels, actualResult.Spec.Selector.MatchLabels)
		a.Equal(expectedLabels, actualResult.Spec.Template.ObjectMeta.Labels)
		a.Equal(utils.AnnotationSidecarIstio(), actualResult.Spec.Template.ObjectMeta.Annotations)
		a.Equal(1, len(actualResult.Spec.Template.Spec.Containers))
	})
}
