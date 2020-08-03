package ratelimiterconfig

import (
	"context"
	"github.com/champly/lib4go/encoding"
	"github.com/ghodss/yaml"
	proto_types "github.com/gogo/protobuf/types"
	networking "istio.io/api/networking/v1alpha3"
	"istio.io/client-go/pkg/apis/networking/v1alpha3"
	"k8s.io/apimachinery/pkg/api/equality"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"ratelimit-operator/pkg/apis/operators/v1"
	"ratelimit-operator/pkg/controller/ratelimiterconfig/envoyfilter_types"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
)

func (r *ReconcileRateLimiterConfig) reconcileEnvoyFilter(ctx context.Context, instance *v1.RateLimiterConfig) (reconcile.Result, error) {
	reqLogger := log.WithValues("Instance.Name", instance.Name)

	foundEnvoyFilter := &v1alpha3.EnvoyFilter{}
	envoyFilterFromInstance := r.buildEnvoyFilter(instance)

	err := r.client.Get(ctx, types.NamespacedName{Name: instance.Name, Namespace: instance.Namespace}, foundEnvoyFilter)
	if err != nil {
		if errors.IsNotFound(err) {
			reqLogger.Info("Creating a new EnvoyFilter")
			err = r.client.Create(ctx, envoyFilterFromInstance)
			if err != nil {
				reqLogger.Error(err, "Failed to create new EnvoyFilter")
				return reconcile.Result{}, err
			}
			return reconcile.Result{Requeue: true}, nil
		} else {
			reqLogger.Error(err, "Failed to get EnvoyFilter")
			return reconcile.Result{}, err
		}
	}

	if !equality.Semantic.DeepEqual(foundEnvoyFilter.Spec, envoyFilterFromInstance.Spec) {
		foundEnvoyFilter.Spec = envoyFilterFromInstance.Spec
		r.client.Update(ctx, foundEnvoyFilter)
	}

	return reconcile.Result{}, nil
}

func (r *ReconcileRateLimiterConfig) buildEnvoyFilter(instance *v1.RateLimiterConfig) *v1alpha3.EnvoyFilter {
	envoyFilter := &v1alpha3.EnvoyFilter{
		ObjectMeta: metav1.ObjectMeta{
			Name:      instance.Name,
			Namespace: instance.Namespace,
		},
		Spec: networking.EnvoyFilter{
			WorkloadSelector: &networking.WorkloadSelector{
				Labels: map[string]string{
					"istio": "ingressgateway",
				},
			},
			ConfigPatches: []*networking.EnvoyFilter_EnvoyConfigObjectPatch{
				{
					ApplyTo: networking.EnvoyFilter_HTTP_FILTER,
					Match: &networking.EnvoyFilter_EnvoyConfigObjectMatch{
						Context: networking.EnvoyFilter_GATEWAY,
						ObjectTypes: &networking.EnvoyFilter_EnvoyConfigObjectMatch_Listener{
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
						},
					},
					Patch: &networking.EnvoyFilter_Patch{
						Operation: networking.EnvoyFilter_Patch_INSERT_BEFORE,
						Value:     convertYaml2Struct(r.buildHttpFilterPatch(instance)),
					},
				},
				{
					ApplyTo: networking.EnvoyFilter_CLUSTER,
					Match: &networking.EnvoyFilter_EnvoyConfigObjectMatch{
						ObjectTypes: &networking.EnvoyFilter_EnvoyConfigObjectMatch_Cluster{
							Cluster: &networking.EnvoyFilter_ClusterMatch{
								Service: r.buildServiceName(),
							},
						},
					},
					Patch: &networking.EnvoyFilter_Patch{
						Operation: networking.EnvoyFilter_Patch_MERGE,
						Value:     convertYaml2Struct(r.buildClusterPatch()),
					},
				},
				{
					ApplyTo: networking.EnvoyFilter_VIRTUAL_HOST,
					Match: &networking.EnvoyFilter_EnvoyConfigObjectMatch{
						Context: networking.EnvoyFilter_GATEWAY,
						ObjectTypes: &networking.EnvoyFilter_EnvoyConfigObjectMatch_RouteConfiguration{
							RouteConfiguration: &networking.EnvoyFilter_RouteConfigurationMatch{
								Vhost: &networking.EnvoyFilter_RouteConfigurationMatch_VirtualHostMatch{
									Name: instance.Spec.VirtualHostName,
									Route: &networking.EnvoyFilter_RouteConfigurationMatch_RouteMatch{
										Action: networking.EnvoyFilter_RouteConfigurationMatch_RouteMatch_ANY,
									},
								},
							},
						},
					},
					Patch: &networking.EnvoyFilter_Patch{
						Operation: networking.EnvoyFilter_Patch_MERGE,
						Value:     convertYaml2Struct(r.buildVirtualHostPatch(instance)),
					},
				},
			},
		},
	}
	controllerutil.SetControllerReference(instance, envoyFilter, r.scheme)
	return envoyFilter
}

func convertYaml2Struct(str string) *proto_types.Struct {
	res, _ := encoding.YAML2Struct(str)
	return res
}

func (r *ReconcileRateLimiterConfig) buildHttpFilterPatch(instance *v1.RateLimiterConfig) string {
	values := envoyfilter_types.HttpFilterPatchValues{
		Name: "envoy.rate_limit",
		Config: envoyfilter_types.Config{
			Domain:          instance.Spec.RateLimitProperty.Domain,
			FailureModeDeny: instance.Spec.FailureModeDeny,
			RateLimitService: envoyfilter_types.RateLimitService{
				GrpcService: envoyfilter_types.GrpcService{
					Timeout: "0.25s",
					EnvoyGrpc: envoyfilter_types.EnvoyGrpc{
						ClusterName: r.buildWorkAroundServiceName(),
					},
				},
			},
		},
	}

	res, err := yaml.Marshal(&values)
	if err != nil {
		log.Error(err, "Failed to convert object to yaml for http filter patch")
	}
	return string(res)
}

func (r *ReconcileRateLimiterConfig) buildClusterPatch() string {
	values := envoyfilter_types.ClusterPatchValues{
		Name: r.buildWorkAroundServiceName(),
	}

	res, err := yaml.Marshal(&values)
	if err != nil {
		log.Error(err, "Failed to convert object to yaml for cluster patch")
	}
	return string(res)
}

func (r *ReconcileRateLimiterConfig) buildVirtualHostPatch(instance *v1.RateLimiterConfig) string {
	var actions []envoyfilter_types.Action

	for _, d := range instance.Spec.RateLimitProperty.Descriptors {
		actions = append(actions,
			envoyfilter_types.Action{
				RequestHeaders: envoyfilter_types.RequestHeader{
					DescriptorKey: d.Key,
					HeaderName:    d.Key,
				},
			},
		)
	}

	rateLimits := []envoyfilter_types.RateLimit{{Actions: actions}}
	values := envoyfilter_types.VirtualHostPatchValues{RateLimits: rateLimits}

	res, err := yaml.Marshal(&values)
	if err != nil {
		log.Error(err, "Failed to convert object to yaml for virtual host patch")
	}
	return string(res)
}

func (r *ReconcileRateLimiterConfig) buildServiceName() string {
	return r.rateLimiter.Name + "." + r.rateLimiter.Namespace + ".svc.cluster.local"
}

func (r *ReconcileRateLimiterConfig) buildWorkAroundServiceName() string {
	return "patched." + r.rateLimiter.Name + "." + r.rateLimiter.Namespace + ".svc.cluster.local"
}
