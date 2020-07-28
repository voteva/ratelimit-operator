package ratelimitconfig

import (
	"context"
	"github.com/champly/lib4go/encoding"
	"github.com/ghodss/yaml"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	operatorsv1 "ratelimit-operator/pkg/apis/operators/v1"
	"ratelimit-operator/pkg/controller/ratelimitconfig/envoyfilter_types"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
	proto_types "github.com/gogo/protobuf/types"
	networking "istio.io/api/networking/v1alpha3"
	"istio.io/client-go/pkg/apis/networking/v1alpha3"
)

func (r *ReconcileRateLimitConfig) reconcileEnvoyFilter(ctx context.Context, instance *operatorsv1.RateLimitConfig) (reconcile.Result, error) {
	foundEnvoyFilter := &v1alpha3.EnvoyFilter{}

	err := r.client.Get(ctx, types.NamespacedName{Name: instance.Name, Namespace: instance.Namespace}, foundEnvoyFilter)
	log.Info("reconcileEnvoyFilter !!!!!")

	if err != nil && errors.IsNotFound(err) {
		ef := r.buildEnvoyFilter(instance)
		log.Info("Creating a new EnvoyFilter", "EnvoyFilter.Namespace", ef.Namespace, "EnvoyFilter.Name", ef.Name)
		err = r.client.Create(ctx, ef)
		if err != nil {
			log.Error(err, "Failed to create new EnvoyFilter", "EnvoyFilter.Namespace", ef.Namespace, "EnvoyFilter.Name", ef.Name)
			return reconcile.Result{}, err
		}
		return reconcile.Result{Requeue: true}, nil
	} else if err != nil {
		log.Error(err, "Failed to get EnvoyFilter")
		return reconcile.Result{}, err
	}
	return reconcile.Result{}, nil
}

func (r *ReconcileRateLimitConfig) buildEnvoyFilter(instance *operatorsv1.RateLimitConfig) *v1alpha3.EnvoyFilter {
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
								Service: r.rateLimiter.Name, //instance.Name + "." + instance.Namespace + ".svc.cluster.local",
							},
						},
					},
					Patch: &networking.EnvoyFilter_Patch{
						Operation: networking.EnvoyFilter_Patch_ADD,
						Value:     convertYaml2Struct(r.buildClusterPatch(instance)),
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

func (r *ReconcileRateLimitConfig) buildHttpFilterPatch(instance *operatorsv1.RateLimitConfig) string {
	values := envoyfilter_types.HttpFilterPatchValues{
		Name: "envoy.rate_limit",
		Config: envoyfilter_types.Config{
			Domain:          instance.Spec.RateLimitProperty.Domain,
			FailureModeDeny: instance.Spec.FailureModeDeny,
			RateLimitService: envoyfilter_types.RateLimitService{
				GrpcService: envoyfilter_types.GrpcService{
					Timeout: "10s", // TODO
					EnvoyGrpc: envoyfilter_types.EnvoyGrpc{
						ClusterName: "rate_limit_service",
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

func (r *ReconcileRateLimitConfig) buildClusterPatch(instance *operatorsv1.RateLimitConfig) string {
	values := envoyfilter_types.ClusterPatchValues{
		ConnectTimeout:       "10s", // TODO
		Http2ProtocolOptions: envoyfilter_types.Http2ProtocolOption{},
		LbPolicy:             "ROUND_ROBIN",
		LoadAssignment: envoyfilter_types.LoadAssignment{
			ClusterName: "rate_limit_service",
			Endpoints: []envoyfilter_types.LoadAssignmentEndpoints{{
				LbEndpoints: []envoyfilter_types.LbEndpoint{{
					Endpoint: envoyfilter_types.Endpoint{
						Address: envoyfilter_types.Address{
							SocketAddress: envoyfilter_types.SocketAddress{
								Address:   r.rateLimiter.Name, //"rate-limit.operator-test.svc.cluster.local",
								PortValue: r.rateLimiter.Spec.ServicePort,
							},
						},
					},
				}},
			}},
		},
		Name: "rate_limit_service",
		Type: "STRICT_DNS",
	}

	res, err := yaml.Marshal(&values)
	if err != nil {
		log.Error(err, "Failed to convert object to yaml for cluster patch")
	}
	return string(res)
}

func (r *ReconcileRateLimitConfig) buildVirtualHostPatch(instance *operatorsv1.RateLimitConfig) string {
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
