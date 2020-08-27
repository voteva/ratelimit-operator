package ratelimiter

import (
	"context"
	"github.com/stretchr/testify/assert"
	"k8s.io/apimachinery/pkg/types"
	v1 "ratelimit-operator/pkg/apis/operators/v1"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
	"testing"
)

func Test_Reconcile_NotFoundRateLimiter(t *testing.T) {
	t.Parallel()
	a := assert.New(t)

	t.Run("reconcile (NotFoundRateLimiter)", func(t *testing.T) {
		rateLimiter := buildRateLimiter()
		request := buildReconcileRequest(rateLimiter)
		r := buildEmptyReconciler()

		reconcileResult, err := r.Reconcile(request)

		a.Nil(err)
		a.NotNil(reconcileResult)
		a.False(reconcileResult.Requeue)
	})
}

func Test_Reconcile_NeedUpdateWithDefaults(t *testing.T) {
	t.Parallel()
	a := assert.New(t)

	t.Run("reconcile (NeedUpdateWithDefaults)", func(t *testing.T) {
		rateLimiter := buildRateLimiter()
		rateLimiter.Spec.LogLevel = nil
		rateLimiter.Spec.Port = nil
		rateLimiter.Spec.Size = nil
		request := buildReconcileRequest(rateLimiter)
		r := buildEmptyReconciler()

		errCreate := r.client.Create(context.Background(), rateLimiter)
		a.Nil(errCreate)

		reconcileResult, err := r.Reconcile(request)
		a.Nil(err)
		a.NotNil(reconcileResult)
		a.True(reconcileResult.Requeue)

		errGet := r.client.Get(context.Background(), buildNamespacedName(rateLimiter), rateLimiter)
		a.Nil(errGet)
		a.Equal(v1.WARN, *rateLimiter.Spec.LogLevel)
		a.Equal(int32(8081), *rateLimiter.Spec.Port)
		a.Equal(int32(1), *rateLimiter.Spec.Size)
	})
}

func Test_Reconcile_ReconcileConfigMap(t *testing.T) {
	t.Parallel()
	a := assert.New(t)

	t.Run("reconcile (ReconcileConfigMap)", func(t *testing.T) {
		rateLimiter := buildRateLimiter()
		request := buildReconcileRequest(rateLimiter)
		r := buildEmptyReconciler()

		errCreate := r.client.Create(context.Background(), rateLimiter)
		a.Nil(errCreate)

		reconcileResult, err := r.Reconcile(request)

		a.Nil(err)
		a.NotNil(reconcileResult)
		a.True(reconcileResult.Requeue)
	})
}

func Test_Reconcile_DeploymentForRedis(t *testing.T) {
	t.Parallel()
	a := assert.New(t)

	t.Run("reconcile (DeploymentForRedis)", func(t *testing.T) {
		rateLimiter := buildRateLimiter()
		request := buildReconcileRequest(rateLimiter)
		r := buildEmptyReconciler()

		errCreate := r.client.Create(context.Background(), rateLimiter)
		errCreateCM := r.client.Create(context.Background(), buildConfigMap(rateLimiter))

		a.Nil(errCreate)
		a.Nil(errCreateCM)

		reconcileResult, err := r.Reconcile(request)

		a.Nil(err)
		a.NotNil(reconcileResult)
		a.True(reconcileResult.Requeue)
	})
}

func Test_Reconcile_ServiceForRedis(t *testing.T) {
	t.Parallel()
	a := assert.New(t)

	t.Run("reconcile (ServiceForRedis)", func(t *testing.T) {
		rateLimiter := buildRateLimiter()
		request := buildReconcileRequest(rateLimiter)
		r := buildEmptyReconciler()

		errCreate := r.client.Create(context.Background(), rateLimiter)
		errCreateCM := r.client.Create(context.Background(), buildConfigMap(rateLimiter))
		errCreateDepRedis := r.client.Create(context.Background(), buildDeploymentForRedis(rateLimiter))

		a.Nil(errCreate)
		a.Nil(errCreateCM)
		a.Nil(errCreateDepRedis)

		reconcileResult, err := r.Reconcile(request)

		a.Nil(err)
		a.NotNil(reconcileResult)
		a.True(reconcileResult.Requeue)
	})
}

func Test_Reconcile_DeploymentForService(t *testing.T) {
	t.Parallel()
	a := assert.New(t)

	t.Run("reconcile (DeploymentForService)", func(t *testing.T) {
		rateLimiter := buildRateLimiter()
		request := buildReconcileRequest(rateLimiter)
		r := buildEmptyReconciler()

		errCreate := r.client.Create(context.Background(), rateLimiter)
		errCreateCM := r.client.Create(context.Background(), buildConfigMap(rateLimiter))
		errCreateDepRedis := r.client.Create(context.Background(), buildDeploymentForRedis(rateLimiter))
		errCreateSrvRedis := r.client.Create(context.Background(), buildServiceForRedis(rateLimiter))

		a.Nil(errCreate)
		a.Nil(errCreateCM)
		a.Nil(errCreateDepRedis)
		a.Nil(errCreateSrvRedis)

		reconcileResult, err := r.Reconcile(request)

		a.Nil(err)
		a.NotNil(reconcileResult)
		a.True(reconcileResult.Requeue)
	})
}

func Test_Reconcile_ServiceForService(t *testing.T) {
	t.Parallel()
	a := assert.New(t)

	t.Run("reconcile (ServiceForService)", func(t *testing.T) {
		rateLimiter := buildRateLimiter()
		request := buildReconcileRequest(rateLimiter)
		r := buildEmptyReconciler()

		errCreate := r.client.Create(context.Background(), rateLimiter)
		errCreateCM := r.client.Create(context.Background(), buildConfigMap(rateLimiter))
		errCreateDepRedis := r.client.Create(context.Background(), buildDeploymentForRedis(rateLimiter))
		errCreateSrvRedis := r.client.Create(context.Background(), buildServiceForRedis(rateLimiter))
		errCreateDepRL := r.client.Create(context.Background(), buildDeploymentForService(rateLimiter))

		a.Nil(errCreate)
		a.Nil(errCreateCM)
		a.Nil(errCreateDepRedis)
		a.Nil(errCreateSrvRedis)
		a.Nil(errCreateDepRL)

		reconcileResult, err := r.Reconcile(request)

		a.Nil(err)
		a.NotNil(reconcileResult)
		a.True(reconcileResult.Requeue)
	})
}

func Test_Reconcile_Success(t *testing.T) {
	t.Parallel()
	a := assert.New(t)

	t.Run("reconcile (Success)", func(t *testing.T) {
		rateLimiter := buildRateLimiter()
		request := buildReconcileRequest(rateLimiter)
		r := buildEmptyReconciler()

		errCreate := r.client.Create(context.Background(), rateLimiter)
		errCreateCM := r.client.Create(context.Background(), buildConfigMap(rateLimiter))
		errCreateDepRedis := r.client.Create(context.Background(), buildDeploymentForRedis(rateLimiter))
		errCreateSrvRedis := r.client.Create(context.Background(), buildServiceForRedis(rateLimiter))
		errCreateDepRL := r.client.Create(context.Background(), buildDeploymentForService(rateLimiter))
		errCreateSrvRL := r.client.Create(context.Background(), buildService(rateLimiter))

		a.Nil(errCreate)
		a.Nil(errCreateCM)
		a.Nil(errCreateDepRedis)
		a.Nil(errCreateSrvRedis)
		a.Nil(errCreateDepRL)
		a.Nil(errCreateSrvRL)

		reconcileResult, err := r.Reconcile(request)

		a.Nil(err)
		a.NotNil(reconcileResult)
		a.False(reconcileResult.Requeue)
	})
}

func buildReconcileRequest(rateLimiter *v1.RateLimiter) reconcile.Request {
	return reconcile.Request{
		NamespacedName: types.NamespacedName{
			Name:      rateLimiter.Name,
			Namespace: rateLimiter.Namespace,
		},
	}
}

func buildNamespacedName(rateLimiter *v1.RateLimiter) types.NamespacedName {
	return types.NamespacedName{
		Name:      rateLimiter.Name,
		Namespace: rateLimiter.Namespace,
	}
}
