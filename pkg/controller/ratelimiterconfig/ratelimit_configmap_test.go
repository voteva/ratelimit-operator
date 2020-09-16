package ratelimiterconfig

import (
	"context"
	"github.com/stretchr/testify/assert"
	"k8s.io/apimachinery/pkg/types"
	"github.com/voteva/ratelimit-operator/pkg/controller/common"
	"github.com/voteva/ratelimit-operator/pkg/utils"
	"testing"
)

func Test_UpdateConfigMap_Success(t *testing.T) {
	t.Parallel()
	a := assert.New(t)

	t.Run("update ConfigMap (Success)", func(t *testing.T) {
		rateLimiter := buildRateLimiter()
		r := buildReconciler(rateLimiter)

		errCreate := r.Client.Create(context.Background(), r.configMap)
		a.Nil(errCreate)

		rateLimiterConfig := buildRateLimiterConfig(rateLimiter)
		reconcileResult, err := r.updateConfigMap(context.Background(), rateLimiterConfig)
		namespacedName := types.NamespacedName{Name: r.configMap.Name, Namespace: r.configMap.Namespace}
		errGet := r.Client.Get(context.Background(), namespacedName, r.configMap)

		a.Nil(err)
		a.NotNil(reconcileResult)
		a.Nil(errGet)

		fileName := common.BuildConfigMapDataFileName(rateLimiterConfig.Name)
		_, found := r.configMap.Data[fileName]
		a.True(found)
	})
}

func Test_UpdateConfigMap_DomainExists(t *testing.T) {
	t.Parallel()
	a := assert.New(t)

	t.Run("update ConfigMap (DomainExists)", func(t *testing.T) {
		rateLimiter := buildRateLimiter()
		r := buildReconciler(rateLimiter)

		errCreate := r.Client.Create(context.Background(), r.configMap)
		a.Nil(errCreate)

		rateLimiterConfig := buildRateLimiterConfig(rateLimiter)
		fileName := utils.BuildRandomString(3)
		r.configMap.Data[fileName] = common.BuildRateLimitPropertyValue(rateLimiterConfig)

		reconcileResult, err := r.updateConfigMap(context.Background(), rateLimiterConfig)
		namespacedName := types.NamespacedName{Name: r.configMap.Name, Namespace: r.configMap.Namespace}
		errGet := r.Client.Get(context.Background(), namespacedName, r.configMap)

		a.Nil(err)
		a.NotNil(reconcileResult)
		a.Nil(errGet)
	})
}

func Test_UpdateConfigMap_ErrorNotFound(t *testing.T) {
	t.Parallel()
	a := assert.New(t)

	t.Run("update ConfigMap (ErrorNotFound)", func(t *testing.T) {
		rateLimiter := buildRateLimiter()
		r := buildReconciler(rateLimiter)
		r.configMap.Data = nil

		rateLimiterConfig := buildRateLimiterConfig(rateLimiter)
		_, err := r.updateConfigMap(context.Background(), rateLimiterConfig)

		a.NotNil(err)
	})
}

func Test_DeleteFromConfigMap_Success(t *testing.T) {
	t.Parallel()
	a := assert.New(t)

	t.Run("delete from ConfigMap (Success)", func(t *testing.T) {
		rateLimiter := buildRateLimiter()
		r := buildReconciler(rateLimiter)

		rateLimiterConfig := buildRateLimiterConfig(rateLimiter)
		fileName := common.BuildConfigMapDataFileName(rateLimiterConfig.Name)
		r.configMap.Data[fileName] = common.BuildRateLimitPropertyValue(rateLimiterConfig)

		errCreate := r.Client.Create(context.Background(), r.configMap)
		a.Nil(errCreate)

		err := r.deleteFromConfigMap(context.Background(), rateLimiterConfig)
		namespacedName := types.NamespacedName{Name: r.configMap.Name, Namespace: r.configMap.Namespace}
		errGet := r.Client.Get(context.Background(), namespacedName, r.configMap)

		a.Nil(err)
		a.Nil(errGet)

		_, found := r.configMap.Data[fileName]
		a.False(found)
	})
}

func Test_DeleteFromConfigMap_SuccessNilData(t *testing.T) {
	t.Parallel()
	a := assert.New(t)

	t.Run("delete from ConfigMap (SuccessNilData)", func(t *testing.T) {
		rateLimiter := buildRateLimiter()
		r := buildReconciler(rateLimiter)
		r.configMap.Data = nil

		rateLimiterConfig := buildRateLimiterConfig(rateLimiter)
		err := r.deleteFromConfigMap(context.Background(), rateLimiterConfig)

		a.Nil(err)
	})
}

func Test_DeleteFromConfigMap_ErrorNotFound(t *testing.T) {
	t.Parallel()
	a := assert.New(t)

	t.Run("delete from ConfigMap (ErrorNotFound)", func(t *testing.T) {
		rateLimiter := buildRateLimiter()
		r := buildReconciler(rateLimiter)

		rateLimiterConfig := buildRateLimiterConfig(rateLimiter)
		err := r.deleteFromConfigMap(context.Background(), rateLimiterConfig)

		a.NotNil(err)
	})
}
