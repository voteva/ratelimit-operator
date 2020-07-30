package ratelimiterconfig
/*
import (
	"github.com/stretchr/testify/assert"
	v1 "ratelimit-operator/pkg/apis/operators/v1"
	"testing"
)

func Test_HttpFilterPatch_Success(t *testing.T) {
	t.Parallel()
	a := assert.New(t)

	t.Run("success build patch for http filter", func(t *testing.T) {
		var strPatch = `
          config:
            domain: test
            failure_mode_deny: true
            rate_limit_service:
              grpc_service:
                envoy_grpc:
                  cluster_name: rate_limit_service
                timeout: 10s
          name: envoy.rate_limit`

		expectedPatch := convertYaml2Struct(strPatch)

		instance := v1.RateLimitConfig{
			Spec: v1.RateLimitConfigSpec{
				RateLimitProperty: v1.RateLimitProperty{
					Domain: "test",
				},
				FailureModeDeny: true,
			},
		}

		result := buildHttpFilterPatch(&instance)
		patch := convertYaml2Struct(result)

		a.Equal(patch, expectedPatch)
	})
}

func Test_ClusterPatch_Success(t *testing.T) {
	t.Parallel()
	a := assert.New(t)

	t.Run("success build patch for cluster", func(t *testing.T) {
		var strPatch = `
          connect_timeout: 10s
          http2_protocol_options: {}
          lb_policy: ROUND_ROBIN
          load_assignment:
            cluster_name: rate_limit_service
            endpoints:
              - lb_endpoints:
                  - endpoint:
                      address:
                        socket_address:
                          address: rate-limit.operator-test.svc.cluster.local
                          port_value: 8081
          name: rate_limit_service
          type: STRICT_DNS`

		expectedPatch := convertYaml2Struct(strPatch)

		instance := v1.RateLimiter{
			Spec: v1.RateLimiterSpec{
				Port: 8081,
			},
		}

		result := buildClusterPatch(&instance)
		patch := convertYaml2Struct(result)

		a.Equal(patch, expectedPatch)
	})
}

func Test_VirtualHostPatch_Success(t *testing.T) {
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

		instance := v1.RateLimiter{
			Spec: v1.RateLimiterSpec{
				RateLimitProperty: v1.RateLimitProperty{
					Descriptors: []v1.Descriptor{{
						Key: "custom-rl-header",
					}},
				},
			},
		}

		result := buildVirtualHostPatch(&instance)
		patch := convertYaml2Struct(result)

		a.Equal(patch, expectedPatch)
	})
}
*/
