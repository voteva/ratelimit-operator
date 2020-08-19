package ratelimiterconfig

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"ratelimit-operator/pkg/apis/operators/v1"
	"ratelimit-operator/pkg/utils"
	"testing"
)

func Test_HttpFilterPatch_Success(t *testing.T) {
	t.Parallel()
	a := assert.New(t)

	t.Run("success build patch for http filter", func(t *testing.T) {
		rateLimiterConfig := &v1.RateLimiterConfig{
			Spec: v1.RateLimiterConfigSpec{
				RateLimitProperty: v1.RateLimitProperty{
					Domain: utils.BuildRandomString(3),
				},
				FailureModeDeny: true,
			},
		}
		rateLimiter := buildRateLimiter()

		expectedPatchValue := fmt.Sprintf(`
          config:
            domain: %s
            failure_mode_deny: true
            rate_limit_service:
              grpc_service:
                envoy_grpc:
                  cluster_name: %s
                timeout: 0.25s
          name: envoy.rate_limit`,
			rateLimiterConfig.Spec.RateLimitProperty.Domain,
			buildWorkAroundServiceName(rateLimiter))

		expectedPatch := convertYaml2Struct(expectedPatchValue)

		actualPatchValue := buildHttpFilterPatchValue(rateLimiterConfig, rateLimiter)
		actualPatch := convertYaml2Struct(actualPatchValue)

		a.Equal(actualPatch, expectedPatch)
	})
}

func Test_ClusterPatch_Success(t *testing.T) {
	t.Parallel()
	a := assert.New(t)

	t.Run("success build patch for cluster", func(t *testing.T) {
		rateLimiter := buildRateLimiter()

		expectedPatchValue := fmt.Sprintf("name: %s", buildWorkAroundServiceName(rateLimiter))
		expectedPatch := convertYaml2Struct(expectedPatchValue)

		actualPatchValue := buildClusterPatchValue(rateLimiter)
		actualPatch := convertYaml2Struct(actualPatchValue)

		a.Equal(actualPatch, expectedPatch)
	})
}

func Test_VirtualHostPatchValue_Success(t *testing.T) {
	t.Parallel()
	a := assert.New(t)

	t.Run("success build patch for virtual host", func(t *testing.T) {
		var strPatch = `
          rate_limits:
            - actions:
                - request_headers:
                    descriptor_key: custom-rl-header
                    header_name: custom-rl-header`

		expectedPatch := convertYaml2Struct(strPatch)

		rateLimiterConfig := &v1.RateLimiterConfig{
			Spec: v1.RateLimiterConfigSpec{
				RateLimitProperty: v1.RateLimitProperty{
					Descriptors: []v1.Descriptor{{
						Key: "custom-rl-header",
					}},
				},
			},
		}

		actualPatchValue := buildVirtualHostPatchValue(rateLimiterConfig)
		actualPatch := convertYaml2Struct(actualPatchValue)

		a.Equal(actualPatch, expectedPatch)
	})
}

func Test_BuildVirtualHostName_Success(t *testing.T) {
	t.Parallel()
	a := assert.New(t)

	t.Run("success build virtual host name", func(t *testing.T) {
		rateLimiterConfig := &v1.RateLimiterConfig{
			Spec: v1.RateLimiterConfigSpec{
				Host: utils.BuildRandomString(3),
				Port: int32(utils.BuildRandomInt(4)),
			},
		}

		expectedResult := fmt.Sprintf("%s:%d", rateLimiterConfig.Spec.Host, rateLimiterConfig.Spec.Port)
		actualResult := buildVirtualHostName(rateLimiterConfig)

		a.Equal(expectedResult, actualResult)
	})
}

func Test_BuildRateLimiterServiceName_Success(t *testing.T) {
	t.Parallel()
	a := assert.New(t)

	t.Run("success build rate limiter service name", func(t *testing.T) {
		rateLimiter := buildRateLimiter()

		expectedResult := fmt.Sprintf("%s.%s.%s", rateLimiter.Name, rateLimiter.Namespace, "svc.cluster.local")
		actualResult := buildRateLimiterServiceName(rateLimiter)

		a.Equal(expectedResult, actualResult)
	})
}

func Test_BuildWorkAroundServiceName_Success(t *testing.T) {
	t.Parallel()
	a := assert.New(t)

	t.Run("success build work around service name", func(t *testing.T) {
		rateLimiter := buildRateLimiter()

		expectedResult := fmt.Sprintf("%s.%s.%s.%s", "patched", rateLimiter.Name, rateLimiter.Namespace, "svc.cluster.local")
		actualResult := buildWorkAroundServiceName(rateLimiter)

		a.Equal(expectedResult, actualResult)
	})
}

func buildRateLimiter() *v1.RateLimiter {
	return &v1.RateLimiter{
		ObjectMeta: metav1.ObjectMeta{
			Name:      utils.BuildRandomString(3),
			Namespace: utils.BuildRandomString(3),
		},
	}
}
