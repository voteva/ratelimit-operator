package ratelimiter

import (
	"context"
	"github.com/stretchr/testify/assert"
	appsv1 "k8s.io/api/apps/v1"
	"ratelimit-operator/pkg/utils"
	"testing"
)

func Test_ReconcileDeploymentForService_CreateSuccess(t *testing.T) {
	t.Parallel()
	a := assert.New(t)

	t.Run("reconcile deployment for ratelimit-service (CreateSuccess)", func(t *testing.T) {
		rateLimiter := buildRateLimiter()
		r := buildReconciler(rateLimiter)

		reconcileResult, err := r.reconcileDeploymentForService(context.Background(), rateLimiter)

		foundDeployment := &appsv1.Deployment{}
		namespaceName := buildServiceResourceNamespacedName(rateLimiter)
		errGet := r.client.Get(context.Background(), namespaceName, foundDeployment)

		a.Nil(err)
		a.NotNil(reconcileResult)
		a.True(reconcileResult.Requeue)
		a.Nil(errGet)
		a.NotNil(foundDeployment)
	})
}

func Test_ReconcileDeploymentForService_CreateError(t *testing.T) {
	t.Parallel()
	a := assert.New(t)

	t.Run("reconcile deployment for ratelimit-service (CreateError)", func(t *testing.T) {
		rateLimiter := buildRateLimiter()
		rateLimiter.Name = ""
		rateLimiter.Namespace = ""
		r := buildReconciler(rateLimiter)

		_, err := r.reconcileDeploymentForService(context.Background(), rateLimiter)

		a.NotNil(err)
	})
}

func Test_ReconcileDeploymentForService_Update(t *testing.T) {
	t.Parallel()
	a := assert.New(t)

	t.Run("reconcile deployment for ratelimit-service (Update)", func(t *testing.T) {
		rateLimiter := buildRateLimiter()
		r := buildReconciler(rateLimiter)

		dep := buildDeploymentForService(rateLimiter)
		dep.Spec.Selector = nil
		errCreateSrvRL := r.client.Create(context.Background(), dep)
		a.Nil(errCreateSrvRL)

		reconcileResult, err := r.reconcileDeploymentForService(context.Background(), rateLimiter)

		foundDeployment := &appsv1.Deployment{}
		namespaceName := buildServiceResourceNamespacedName(rateLimiter)
		errGet := r.client.Get(context.Background(), namespaceName, foundDeployment)

		a.Nil(err)
		a.NotNil(reconcileResult)
		a.False(reconcileResult.Requeue)
		a.Nil(errGet)
		a.NotNil(foundDeployment)
		a.NotNil(foundDeployment.Spec.Selector)
		a.Equal(utils.LabelsForApp(rateLimiter.Name), foundDeployment.Spec.Selector.MatchLabels)
	})
}

func Test_BuildDeployment(t *testing.T) {
	t.Parallel()
	a := assert.New(t)

	t.Run("success build Deployment for ratelimit-service", func(t *testing.T) {
		rateLimiter := buildRateLimiter()

		actualResult := buildDeploymentForService(rateLimiter)

		a.Equal(rateLimiter.Name, actualResult.Name)
		a.Equal(rateLimiter.Namespace, actualResult.Namespace)
		a.Equal(rateLimiter.Spec.Size, actualResult.Spec.Replicas)
		a.Equal(utils.LabelsForApp(rateLimiter.Name), actualResult.Spec.Selector.MatchLabels)
		a.Equal(utils.LabelsForApp(rateLimiter.Name), actualResult.Spec.Template.ObjectMeta.Labels)
		a.Equal(utils.AnnotationSidecarIstio(), actualResult.Spec.Template.ObjectMeta.Annotations)
		a.Equal(1, len(actualResult.Spec.Template.Spec.Containers))
		a.Equal(1, len(actualResult.Spec.Template.Spec.Volumes))
		a.Equal("config", actualResult.Spec.Template.Spec.Volumes[0].Name)
		a.Equal(rateLimiter.Name, actualResult.Spec.Template.Spec.Volumes[0].VolumeSource.ConfigMap.LocalObjectReference.Name)
	})
}
