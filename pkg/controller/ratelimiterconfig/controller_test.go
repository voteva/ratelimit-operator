package ratelimiterconfig

import (
	"context"
	"github.com/stretchr/testify/assert"
	"k8s.io/apimachinery/pkg/types"
	v1 "ratelimit-operator/pkg/apis/operators/v1"
	"ratelimit-operator/pkg/utils"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
	"testing"
)

func Test_Reconcile_NotFoundRateLimiterConfig(t *testing.T) {
	t.Parallel()
	a := assert.New(t)

	t.Run("reconcile (NotFoundRateLimiterConfig)", func(t *testing.T) {
		rateLimiter := buildRateLimiter()
		rateLimiterConfig := buildRateLimiterConfig(rateLimiter)
		request := buildReconcileRequest(rateLimiterConfig)
		r := buildEmptyReconciler()

		reconcileResult, err := r.Reconcile(request)

		a.Nil(err)
		a.NotNil(reconcileResult)
		a.False(reconcileResult.Requeue)
	})
}

func Test_Reconcile_NotFoundRateLimiter(t *testing.T) {
	t.Parallel()
	a := assert.New(t)

	t.Run("reconcile (NotFoundRateLimiter)", func(t *testing.T) {
		rateLimiter := buildRateLimiter()
		rateLimiterConfig := buildRateLimiterConfig(rateLimiter)
		rateLimiterConfig.Spec.RateLimiter = utils.BuildRandomString(5)
		r := buildReconciler(rateLimiter)

		errCreate := r.client.Create(context.Background(), rateLimiterConfig)
		a.Nil(errCreate)

		request := buildReconcileRequest(rateLimiterConfig)
		_, err := r.Reconcile(request)

		a.NotNil(err)
	})
}

func Test_Reconcile_NotFoundConfigMap(t *testing.T) {
	t.Parallel()
	a := assert.New(t)

	t.Run("reconcile (CreatedEnvoyFilter)", func(t *testing.T) {
		rateLimiter := buildRateLimiter()
		rateLimiterConfig := buildRateLimiterConfig(rateLimiter)
		request := buildReconcileRequest(rateLimiterConfig)
		r := buildReconciler(rateLimiter)

		errCreate := r.client.Create(context.Background(), rateLimiterConfig)
		a.Nil(errCreate)

		_, err := r.Reconcile(request)

		a.NotNil(err)
	})
}

func Test_Reconcile_CreatedEnvoyFilter(t *testing.T) {
	t.Parallel()
	a := assert.New(t)

	t.Run("reconcile (CreatedEnvoyFilter)", func(t *testing.T) {
		rateLimiter := buildRateLimiter()
		rateLimiterConfig := buildRateLimiterConfig(rateLimiter)
		request := buildReconcileRequest(rateLimiterConfig)
		r := buildReconciler(rateLimiter)

		errCreate := r.client.Create(context.Background(), rateLimiterConfig)
		a.Nil(errCreate)

		errCreateCM := r.client.Create(context.Background(), buildConfigMap(rateLimiter))
		a.Nil(errCreateCM)

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
		rateLimiterConfig := buildRateLimiterConfig(rateLimiter)
		request := buildReconcileRequest(rateLimiterConfig)
		r := buildReconciler(rateLimiter)

		errCreate := r.client.Create(context.Background(), rateLimiterConfig)
		a.Nil(errCreate)

		errCreateCM := r.client.Create(context.Background(), buildConfigMap(rateLimiter))
		a.Nil(errCreateCM)

		errCreateEF := r.client.Create(context.Background(), buildEnvoyFilter(rateLimiterConfig, rateLimiter))
		a.Nil(errCreateEF)

		reconcileResult, err := r.Reconcile(request)

		a.Nil(err)
		a.NotNil(reconcileResult)
		a.False(reconcileResult.Requeue)
	})
}

func Test_GetRateLimiter_Success(t *testing.T) {
	t.Parallel()
	a := assert.New(t)

	t.Run("get RateLimiter (Success)", func(t *testing.T) {
		rateLimiter := buildRateLimiter()
		rateLimiterConfig := buildRateLimiterConfig(rateLimiter)
		r := buildReconciler(rateLimiter)

		err := r.getRateLimiter(context.Background(), rateLimiterConfig)

		a.Nil(err)
	})
}

func Test_GetRateLimiter_ErrorNotFound(t *testing.T) {
	t.Parallel()
	a := assert.New(t)

	t.Run("get RateLimiter (ErrorNotFound)", func(t *testing.T) {
		rateLimiter := buildRateLimiter()
		rateLimiterConfig := buildRateLimiterConfig(rateLimiter)
		rateLimiterConfig.Spec.RateLimiter = utils.BuildRandomString(5)
		r := buildReconciler(rateLimiter)

		err := r.getRateLimiter(context.Background(), rateLimiterConfig)

		a.NotNil(err)
	})
}

func Test_GetRateLimiterConfigMap_Success(t *testing.T) {
	t.Parallel()
	a := assert.New(t)

	t.Run("get RateLimiter ConfigMap (Success)", func(t *testing.T) {
		rateLimiter := buildRateLimiter()
		rateLimiterConfig := buildRateLimiterConfig(rateLimiter)
		r := buildReconciler(rateLimiter)

		errCreate := r.client.Create(context.Background(), r.configMap)
		a.Nil(errCreate)

		err := r.getRateLimiterConfigMap(context.Background(), rateLimiterConfig)

		a.Nil(err)
	})
}

func Test_GetRateLimiterConfigMap_Error(t *testing.T) {
	t.Parallel()
	a := assert.New(t)

	t.Run("get RateLimiter ConfigMap (Error)", func(t *testing.T) {
		rateLimiter := buildRateLimiter()
		rateLimiterConfig := buildRateLimiterConfig(rateLimiter)
		r := buildReconciler(rateLimiter)

		err := r.getRateLimiterConfigMap(context.Background(), rateLimiterConfig)

		a.NotNil(err)
	})
}

func Test_AddFinalizerIfNotExists(t *testing.T) {
	t.Parallel()

	t.Run("add Finalizer if not exists", func(t *testing.T) {
		rateLimiter := buildRateLimiter()
		rateLimiterConfig := buildRateLimiterConfig(rateLimiter)
		r := buildReconciler(rateLimiter)

		r.addFinalizerIfNotExists(context.Background(), rateLimiterConfig)
	})
}

func Test_ManageCleanUpLogic_Success(t *testing.T) {
	t.Parallel()
	a := assert.New(t)

	t.Run("manage clean up logic (Success)", func(t *testing.T) {
		rateLimiter := buildRateLimiter()
		rateLimiterConfig := buildRateLimiterConfig(rateLimiter)
		r := buildReconciler(rateLimiter)

		errCreate := r.client.Create(context.Background(), r.configMap)
		a.Nil(errCreate)

		err := r.manageCleanUpLogic(context.Background(), rateLimiterConfig)

		a.Nil(err)
	})
}

func Test_ManageCleanUpLogic_Error(t *testing.T) {
	t.Parallel()
	a := assert.New(t)

	t.Run("manage clean up logic (Error)", func(t *testing.T) {
		rateLimiter := buildRateLimiter()
		rateLimiterConfig := buildRateLimiterConfig(rateLimiter)
		r := buildReconciler(rateLimiter)

		err := r.manageCleanUpLogic(context.Background(), rateLimiterConfig)

		a.NotNil(err)
	})
}

func buildReconcileRequest(rateLimiterConfig *v1.RateLimiterConfig) reconcile.Request {
	return reconcile.Request{
		NamespacedName: types.NamespacedName{
			Name:      rateLimiterConfig.Name,
			Namespace: rateLimiterConfig.Namespace,
		},
	}
}