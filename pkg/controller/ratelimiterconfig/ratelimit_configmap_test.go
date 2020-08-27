package ratelimiterconfig

import (
	"context"
	"fmt"
	"github.com/stretchr/testify/assert"
	"k8s.io/apimachinery/pkg/types"
	v1 "ratelimit-operator/pkg/apis/operators/v1"
	"ratelimit-operator/pkg/utils"
	"testing"
)

func Test_UpdateConfigMap_Success(t *testing.T) {
	t.Parallel()
	a := assert.New(t)

	t.Run("update ConfigMap (Success)", func(t *testing.T) {
		rateLimiter := buildRateLimiter()
		r := buildReconciler(rateLimiter)

		errCreate := r.client.Create(context.Background(), r.configMap)
		a.Nil(errCreate)

		rateLimiterConfig := buildRateLimiterConfig(rateLimiter)
		reconcileResult, err := r.updateConfigMap(context.Background(), rateLimiterConfig)
		namespacedName := types.NamespacedName{Name: r.configMap.Name, Namespace: r.configMap.Namespace}
		errGet := r.client.Get(context.Background(), namespacedName, r.configMap)

		a.Nil(err)
		a.NotNil(reconcileResult)
		a.Nil(errGet)

		fileName := buildConfigMapDataFileName(rateLimiterConfig.Name)
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

		errCreate := r.client.Create(context.Background(), r.configMap)
		a.Nil(errCreate)

		rateLimiterConfig := buildRateLimiterConfig(rateLimiter)
		fileName := utils.BuildRandomString(3)
		r.configMap.Data[fileName] = buildRateLimitPropertyValue(rateLimiterConfig.Spec.RateLimitProperty)

		reconcileResult, err := r.updateConfigMap(context.Background(), rateLimiterConfig)
		namespacedName := types.NamespacedName{Name: r.configMap.Name, Namespace: r.configMap.Namespace}
		errGet := r.client.Get(context.Background(), namespacedName, r.configMap)

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
		fileName := buildConfigMapDataFileName(rateLimiterConfig.Name)
		r.configMap.Data[fileName] = buildRateLimitPropertyValue(rateLimiterConfig.Spec.RateLimitProperty)

		errCreate := r.client.Create(context.Background(), r.configMap)
		a.Nil(errCreate)

		err := r.deleteFromConfigMap(context.Background(), rateLimiterConfig)
		namespacedName := types.NamespacedName{Name: r.configMap.Name, Namespace: r.configMap.Namespace}
		errGet := r.client.Get(context.Background(), namespacedName, r.configMap)

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

func Test_BuildRateLimitPropertyValue_Success(t *testing.T) {
	t.Parallel()
	a := assert.New(t)

	t.Run("success build RateLimitProperty value", func(t *testing.T) {
		prop := v1.RateLimitProperty{
			Domain: utils.BuildRandomString(3),
		}
		a.NotNil(buildRateLimitPropertyValue(prop))
	})
}

func Test_UnmarshalRateLimitPropertyValue_Success(t *testing.T) {
	t.Parallel()
	a := assert.New(t)

	t.Run("unmarshal RateLimitProperty value (Success)", func(t *testing.T) {
		domain := utils.BuildRandomString(3)
		data := fmt.Sprintf("domain: %s", domain)

		result := unmarshalRateLimitPropertyValue(data)

		a.NotNil(result)
		a.Equal(domain, result.Domain)
	})
}

func Test_BuildConfigMapDataFileName_Success(t *testing.T) {
	t.Parallel()
	a := assert.New(t)

	t.Run("success build ConfigMap.Data file name", func(t *testing.T) {
		fileName := utils.BuildRandomString(3)
		a.Equal(fileName+".yaml", buildConfigMapDataFileName(fileName))
	})
}
