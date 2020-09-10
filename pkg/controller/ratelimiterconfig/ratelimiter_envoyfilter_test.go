package ratelimiterconfig

import (
	"context"
	"fmt"
	"github.com/stretchr/testify/assert"
	networking "istio.io/api/networking/v1alpha3"
	"istio.io/client-go/pkg/apis/networking/v1alpha3"
	"k8s.io/apimachinery/pkg/types"
	"ratelimit-operator/pkg/apis/operators/v1"
	"ratelimit-operator/pkg/utils"
	"testing"
)

func Test_ReconcileEnvoyFilter_CreateSuccess(t *testing.T) {
	t.Parallel()
	a := assert.New(t)

	t.Run("reconcile EnvoyFilter (CreateSuccess)", func(t *testing.T) {
		rateLimiter := buildRateLimiter()
		r := buildReconciler(rateLimiter)

		rateLimiterConfig := buildRateLimiterConfig(rateLimiter)
		reconcileResult, err := r.reconcileEnvoyFilter(context.Background(), rateLimiterConfig)

		a.Nil(err)
		a.NotNil(reconcileResult)
		a.True(reconcileResult.Requeue)

		foundEnvoyFilter := &v1alpha3.EnvoyFilter{}
		namespaceName := types.NamespacedName{Name: rateLimiterConfig.Name, Namespace: rateLimiterConfig.Namespace}
		errGet := r.Client.Get(context.Background(), namespaceName, foundEnvoyFilter)

		a.Nil(errGet)
		a.NotNil(foundEnvoyFilter)
	})
}

func Test_ReconcileEnvoyFilter_CreateError(t *testing.T) {
	t.Parallel()
	a := assert.New(t)

	t.Run("reconcile EnvoyFilter (CreateError)", func(t *testing.T) {
		rateLimiter := buildRateLimiter()
		r := buildReconciler(rateLimiter)

		rateLimiterConfig := buildRateLimiterConfig(rateLimiter)
		rateLimiterConfig.Name = ""
		rateLimiterConfig.Namespace = ""
		_, err := r.reconcileEnvoyFilter(context.Background(), rateLimiterConfig)

		a.NotNil(err)
	})
}

func Test_ReconcileEnvoyFilter_Update(t *testing.T) {
	t.Parallel()
	a := assert.New(t)

	t.Run("reconcile EnvoyFilter (Update)", func(t *testing.T) {
		rateLimiter := buildRateLimiter()
		rateLimiterConfig := buildRateLimiterConfig(rateLimiter)
		r := buildReconciler(rateLimiter)

		ef := buildEnvoyFilter(rateLimiterConfig, rateLimiter)
		ef.Spec.WorkloadSelector = nil
		errCreateEF := r.Client.Create(context.Background(), ef)
		a.Nil(errCreateEF)

		reconcileResult, err := r.reconcileEnvoyFilter(context.Background(), rateLimiterConfig)

		a.Nil(err)
		a.NotNil(reconcileResult)
		a.False(reconcileResult.Requeue)

		foundEnvoyFilter := &v1alpha3.EnvoyFilter{}
		namespaceName := types.NamespacedName{Name: rateLimiterConfig.Name, Namespace: rateLimiterConfig.Namespace}
		errGet := r.Client.Get(context.Background(), namespaceName, foundEnvoyFilter)

		a.Nil(errGet)
		a.NotNil(foundEnvoyFilter)
		a.NotNil(foundEnvoyFilter.Spec.WorkloadSelector)
		a.Equal(buildWorkloadSelectorLabels(rateLimiterConfig), foundEnvoyFilter.Spec.WorkloadSelector.Labels)
	})
}

func Test_BuildEnvoyFilter_Success(t *testing.T) {
	t.Parallel()
	a := assert.New(t)

	t.Run("success build envoy filter", func(t *testing.T) {
		rateLimiter := buildRateLimiter()
		rateLimiterConfig := buildRateLimiterConfig(rateLimiter)

		actualPatch := buildEnvoyFilter(rateLimiterConfig, rateLimiter)

		a.Equal(rateLimiterConfig.Name, actualPatch.ObjectMeta.Name)
		a.Equal(rateLimiterConfig.Namespace, actualPatch.ObjectMeta.Namespace)
		a.Equal(buildWorkloadSelectorLabels(rateLimiterConfig), actualPatch.Spec.WorkloadSelector.Labels)
		a.Equal(3, len(actualPatch.Spec.ConfigPatches))
		a.Equal(buildHttpFilterPatch(rateLimiterConfig, rateLimiter), actualPatch.Spec.ConfigPatches[0])
		a.Equal(buildClusterPatch(rateLimiter), actualPatch.Spec.ConfigPatches[1])
		a.Equal(buildVirtualHostPatch(rateLimiterConfig), actualPatch.Spec.ConfigPatches[2])
	})
}

func Test_BuildHttpFilterPatch_Success(t *testing.T) {
	t.Parallel()
	a := assert.New(t)

	t.Run("success build patch for http filter", func(t *testing.T) {
		rateLimiter := buildRateLimiter()
		rateLimiterConfig := buildRateLimiterConfig(rateLimiter)

		expectedObjectTypes := &networking.EnvoyFilter_EnvoyConfigObjectMatch_Listener{
			Listener: &networking.EnvoyFilter_ListenerMatch{
				FilterChain: &networking.EnvoyFilter_ListenerMatch_FilterChainMatch{
					Filter: &networking.EnvoyFilter_ListenerMatch_FilterMatch{
						Name: "envoy.http_connection_manager",
						SubFilter: &networking.EnvoyFilter_ListenerMatch_SubFilterMatch{
							Name: "envoy.router",
						},
					},
				},
			},
		}

		actualPatch := buildHttpFilterPatch(rateLimiterConfig, rateLimiter)

		a.Equal(networking.EnvoyFilter_HTTP_FILTER, actualPatch.ApplyTo)
		a.IsType(&networking.EnvoyFilter_EnvoyConfigObjectMatch_Listener{}, actualPatch.Match.ObjectTypes)
		a.Equal(expectedObjectTypes, actualPatch.Match.ObjectTypes)
		a.Equal(networking.EnvoyFilter_Patch_INSERT_BEFORE, actualPatch.Patch.Operation)
		a.Equal(convertYaml2Struct(buildHttpFilterPatchValue(rateLimiterConfig, rateLimiter)), actualPatch.Patch.Value)
	})
}

func Test_BuildHttpFilterPatchValue_Success(t *testing.T) {
	t.Parallel()
	a := assert.New(t)

	t.Run("success build patch value for http filter", func(t *testing.T) {
		rateLimiter := buildRateLimiter()
		rateLimiterConfig := buildRateLimiterConfig(rateLimiter)

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
			rateLimiterConfig.Name,
			buildWorkAroundServiceName(rateLimiter))

		expectedPatch := convertYaml2Struct(expectedPatchValue)

		actualPatchValue := buildHttpFilterPatchValue(rateLimiterConfig, rateLimiter)
		actualPatch := convertYaml2Struct(actualPatchValue)

		a.Equal(expectedPatch, actualPatch)
	})
}

func Test_BuildClusterPatch_Success(t *testing.T) {
	t.Parallel()
	a := assert.New(t)

	t.Run("success build patch for cluster", func(t *testing.T) {
		rateLimiter := buildRateLimiter()

		expectedObjectTypes := &networking.EnvoyFilter_EnvoyConfigObjectMatch_Cluster{
			Cluster: &networking.EnvoyFilter_ClusterMatch{
				Service: buildRateLimiterServiceName(rateLimiter),
			},
		}

		actualPatch := buildClusterPatch(rateLimiter)

		a.Equal(networking.EnvoyFilter_CLUSTER, actualPatch.ApplyTo)
		a.IsType(&networking.EnvoyFilter_EnvoyConfigObjectMatch_Cluster{}, actualPatch.Match.ObjectTypes)
		a.Equal(expectedObjectTypes, actualPatch.Match.ObjectTypes)
		a.Equal(networking.EnvoyFilter_Patch_MERGE, actualPatch.Patch.Operation)
		a.Equal(convertYaml2Struct(buildClusterPatchValue(rateLimiter)), actualPatch.Patch.Value)
	})
}

func Test_BuildClusterPatchValue_Success(t *testing.T) {
	t.Parallel()
	a := assert.New(t)

	t.Run("success build patch value for cluster", func(t *testing.T) {
		rateLimiter := buildRateLimiter()

		expectedPatchValue := fmt.Sprintf("name: %s", buildWorkAroundServiceName(rateLimiter))
		expectedPatch := convertYaml2Struct(expectedPatchValue)

		actualPatchValue := buildClusterPatchValue(rateLimiter)
		actualPatch := convertYaml2Struct(actualPatchValue)

		a.Equal(expectedPatch, actualPatch)
	})
}

func Test_BuildVirtualHostPatch_Success(t *testing.T) {
	t.Parallel()
	a := assert.New(t)

	t.Run("success build virtual host patch", func(t *testing.T) {
		host := utils.BuildRandomString(3)
		rateLimiterConfig := &v1.RateLimiterConfig{
			Spec: v1.RateLimiterConfigSpec{
				ApplyTo: v1.GATEWAY,
				Host:    &host,
				Port:    int32(utils.BuildRandomInt(2)),
				Descriptors: []v1.Descriptor{{
					Key: utils.BuildRandomString(3),
				}},
			},
		}

		expectedObjectTypes := &networking.EnvoyFilter_EnvoyConfigObjectMatch_RouteConfiguration{
			RouteConfiguration: &networking.EnvoyFilter_RouteConfigurationMatch{
				Vhost: &networking.EnvoyFilter_RouteConfigurationMatch_VirtualHostMatch{
					Name: buildVirtualHostName(rateLimiterConfig),
					Route: &networking.EnvoyFilter_RouteConfigurationMatch_RouteMatch{
						Action: networking.EnvoyFilter_RouteConfigurationMatch_RouteMatch_ANY,
					},
				},
			},
		}

		actualPatch := buildVirtualHostPatch(rateLimiterConfig)

		a.Equal(networking.EnvoyFilter_VIRTUAL_HOST, actualPatch.ApplyTo)
		a.Equal(buildContext(rateLimiterConfig), actualPatch.Match.Context)
		a.IsType(&networking.EnvoyFilter_EnvoyConfigObjectMatch_RouteConfiguration{}, actualPatch.Match.ObjectTypes)
		a.Equal(expectedObjectTypes, actualPatch.Match.ObjectTypes)
		a.Equal(networking.EnvoyFilter_Patch_MERGE, actualPatch.Patch.Operation)
		a.Equal(convertYaml2Struct(buildVirtualHostPatchValue(rateLimiterConfig)), actualPatch.Patch.Value)
	})
}

func Test_BuildVirtualHostPatchValue_HeaderSuccess(t *testing.T) {
	t.Parallel()
	a := assert.New(t)

	t.Run("success build patch value for virtual host (header)", func(t *testing.T) {
		header := utils.BuildRandomString(3)
		rateLimiterConfig := &v1.RateLimiterConfig{
			Spec: v1.RateLimiterConfigSpec{
				Descriptors: []v1.Descriptor{{
					Key: header,
				}},
				RateLimits: []v1.RateLimits{{
					Actions: []v1.Action{{
						RequestHeaders: &v1.Action_RequestHeaders{
							HeaderName:    header,
							DescriptorKey: header,
						},
					}},
				}},
			},
		}

		var expectedPatchValue = fmt.Sprintf(`
          rate_limits:
            - actions:
                - request_headers:
                    descriptor_key: %s
                    header_name: %s`,
			header, header)

		expectedPatch := convertYaml2Struct(expectedPatchValue)

		actualPatchValue := buildVirtualHostPatchValue(rateLimiterConfig)
		actualPatch := convertYaml2Struct(actualPatchValue)

		a.Equal(expectedPatch, actualPatch)
	})
}

func Test_BuildVirtualHostName_Gateway(t *testing.T) {
	t.Parallel()
	a := assert.New(t)

	t.Run("success build virtual host name (Gateway)", func(t *testing.T) {
		host := utils.BuildRandomString(3)
		rateLimiterConfig := &v1.RateLimiterConfig{
			Spec: v1.RateLimiterConfigSpec{
				ApplyTo: v1.GATEWAY,
				Host:    &host,
				Port:    int32(utils.BuildRandomInt(4)),
			},
		}

		expectedResult := fmt.Sprintf("%s:%d", *rateLimiterConfig.Spec.Host, rateLimiterConfig.Spec.Port)
		actualResult := buildVirtualHostName(rateLimiterConfig)

		a.Equal(expectedResult, actualResult)
	})
}

func Test_BuildVirtualHostName_SidecarOutbound(t *testing.T) {
	t.Parallel()
	a := assert.New(t)

	t.Run("success build virtual host name (SidecarOutbound)", func(t *testing.T) {
		host := utils.BuildRandomString(3)
		rateLimiterConfig := &v1.RateLimiterConfig{
			Spec: v1.RateLimiterConfigSpec{
				ApplyTo: v1.SIDECAR_OUTBOUND,
				Host:    &host,
				Port:    int32(utils.BuildRandomInt(4)),
			},
		}

		expectedResult := fmt.Sprintf("%s:%d", *rateLimiterConfig.Spec.Host, rateLimiterConfig.Spec.Port)
		actualResult := buildVirtualHostName(rateLimiterConfig)

		a.Equal(expectedResult, actualResult)
	})
}

func Test_BuildVirtualHostName_SidecarInbound(t *testing.T) {
	t.Parallel()
	a := assert.New(t)

	t.Run("success build virtual host name (SidecarInbound)", func(t *testing.T) {
		host := utils.BuildRandomString(3)
		rateLimiterConfig := &v1.RateLimiterConfig{
			Spec: v1.RateLimiterConfigSpec{
				ApplyTo: v1.SIDECAR_INBOUND,
				Host:    &host,
				Port:    int32(utils.BuildRandomInt(4)),
			},
		}

		expectedResult := fmt.Sprintf("%s|%d", "inbound|http", rateLimiterConfig.Spec.Port)
		actualResult := buildVirtualHostName(rateLimiterConfig)

		a.Equal(expectedResult, actualResult)
	})
}

func Test_BuildContext_Gateway(t *testing.T) {
	t.Parallel()
	a := assert.New(t)

	t.Run("success build context (Gateway)", func(t *testing.T) {
		rateLimiterConfig := &v1.RateLimiterConfig{
			Spec: v1.RateLimiterConfigSpec{
				ApplyTo: v1.GATEWAY,
			},
		}

		expectedResult := networking.EnvoyFilter_GATEWAY
		actualResult := buildContext(rateLimiterConfig)

		a.Equal(expectedResult, actualResult)
	})
}

func Test_BuildContext_SidecarOutbound(t *testing.T) {
	t.Parallel()
	a := assert.New(t)

	t.Run("success build context (SidecarOutbound)", func(t *testing.T) {
		rateLimiterConfig := &v1.RateLimiterConfig{
			Spec: v1.RateLimiterConfigSpec{
				ApplyTo: v1.SIDECAR_OUTBOUND,
				WorkloadSelector: v1.WorkloadSelector{
					Labels: map[string]string{utils.BuildRandomString(3): utils.BuildRandomString(3)},
				},
			},
		}

		expectedResult := networking.EnvoyFilter_SIDECAR_OUTBOUND
		actualResult := buildContext(rateLimiterConfig)

		a.Equal(expectedResult, actualResult)
	})
}

func Test_BuildContext_SidecarInbound(t *testing.T) {
	t.Parallel()
	a := assert.New(t)

	t.Run("success build context (SidecarInbound)", func(t *testing.T) {
		rateLimiterConfig := &v1.RateLimiterConfig{
			Spec: v1.RateLimiterConfigSpec{
				ApplyTo: v1.SIDECAR_INBOUND,
				WorkloadSelector: v1.WorkloadSelector{
					Labels: map[string]string{utils.BuildRandomString(3): utils.BuildRandomString(3)},
				},
			},
		}

		expectedResult := networking.EnvoyFilter_SIDECAR_INBOUND
		actualResult := buildContext(rateLimiterConfig)

		a.Equal(expectedResult, actualResult)
	})
}

func Test_BuildWorkloadSelectorLabels_EmptyWorkloadSelector(t *testing.T) {
	t.Parallel()
	a := assert.New(t)

	t.Run("success build workload selector labels (EmptyWorkloadSelector)", func(t *testing.T) {
		rateLimiterConfig := &v1.RateLimiterConfig{
			Spec: v1.RateLimiterConfigSpec{
				ApplyTo: v1.GATEWAY,
				WorkloadSelector: v1.WorkloadSelector{
					Labels: map[string]string{},
				},
			},
		}

		expectedResult := map[string]string{}
		actualResult := buildWorkloadSelectorLabels(rateLimiterConfig)

		a.Equal(expectedResult, actualResult)
	})
}

func Test_BuildWorkloadSelectorLabels_ExistsWorkloadSelector(t *testing.T) {
	t.Parallel()
	a := assert.New(t)

	t.Run("success build workload selector labels (ExistsWorkloadSelector)", func(t *testing.T) {
		rateLimiterConfig := &v1.RateLimiterConfig{
			Spec: v1.RateLimiterConfigSpec{
				ApplyTo: v1.SIDECAR_OUTBOUND,
				WorkloadSelector: v1.WorkloadSelector{
					Labels: map[string]string{utils.BuildRandomString(3): utils.BuildRandomString(3)},
				},
			},
		}

		expectedResult := rateLimiterConfig.Spec.WorkloadSelector.Labels
		actualResult := buildWorkloadSelectorLabels(rateLimiterConfig)

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
