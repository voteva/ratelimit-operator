package ratelimiter

import (
	"context"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	networking "istio.io/api/networking/v1alpha3"
	"istio.io/client-go/pkg/apis/networking/v1alpha3"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/types"
	operatorsv1alpha1 "ratelimit-operator/pkg/apis/operators/v1alpha1"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
	"github.com/champly/lib4go/encoding"
	proto_types "github.com/gogo/protobuf/types"
)

func (r *ReconcileRateLimiter) reconcileEnvoyFilter(request reconcile.Request, instance *operatorsv1alpha1.RateLimiter) (reconcile.Result, error) {
	reqLogger := log.WithValues("Request.Namespace", request.Namespace, "Request.Name", request.Name)

	foundEnvoyFilter := &v1alpha3.EnvoyFilter{}
	err := r.client.Get(context.TODO(), types.NamespacedName{Name: instance.Name, Namespace: instance.Namespace}, foundEnvoyFilter)
	if err != nil && errors.IsNotFound(err) {
		// Define a new EnvoyFilter
		cm := r.buildEnvoyFilter(instance)
		reqLogger.Info("Creating a new EnvoyFilter", "EnvoyFilter.Namespace", cm.Namespace, "EnvoyFilter.Name", cm.Name)
		err = r.client.Create(context.TODO(), cm)
		if err != nil {
			reqLogger.Error(err, "Failed to create new EnvoyFilter", "EnvoyFilter.Namespace", cm.Namespace, "EnvoyFilter.Name", cm.Name)
			return reconcile.Result{}, err
		}
		// Deployment created successfully - return and requeue
		return reconcile.Result{Requeue: true}, nil
	} else if err != nil {
		reqLogger.Error(err, "Failed to get EnvoyFilter")
		return reconcile.Result{}, err
	}
	return reconcile.Result{}, nil
}

func (r *ReconcileRateLimiter) buildEnvoyFilter(m *operatorsv1alpha1.RateLimiter) *v1alpha3.EnvoyFilter {
	envoyFilter := &v1alpha3.EnvoyFilter{
		ObjectMeta: metav1.ObjectMeta{
			Name:      m.Name,
			Namespace: m.Namespace,
		},
		Spec: networking.EnvoyFilter{
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
						Value:     getPatch(patch1),
					},
				},
				{
					ApplyTo: networking.EnvoyFilter_CLUSTER,
					Match: &networking.EnvoyFilter_EnvoyConfigObjectMatch{
						ObjectTypes: &networking.EnvoyFilter_EnvoyConfigObjectMatch_Cluster{
							Cluster: &networking.EnvoyFilter_ClusterMatch{
								Service: "rate-limit-server.test-project.svc.cluster.local",
							},
						},
					},
					Patch: &networking.EnvoyFilter_Patch{
						Operation: networking.EnvoyFilter_Patch_INSERT_BEFORE,
						Value:     getPatch(patch2),
					},
				},
				{
					ApplyTo: networking.EnvoyFilter_VIRTUAL_HOST,
					Match: &networking.EnvoyFilter_EnvoyConfigObjectMatch{
						Context: networking.EnvoyFilter_GATEWAY,
						ObjectTypes: &networking.EnvoyFilter_EnvoyConfigObjectMatch_RouteConfiguration{
							RouteConfiguration: &networking.EnvoyFilter_RouteConfigurationMatch{
								Vhost: &networking.EnvoyFilter_RouteConfigurationMatch_VirtualHostMatch{
									Name: "host-info-service.org:80",
									Route: &networking.EnvoyFilter_RouteConfigurationMatch_RouteMatch{
										Action: networking.EnvoyFilter_RouteConfigurationMatch_RouteMatch_ANY,
									},
								},
							},
						},
					},
					Patch: &networking.EnvoyFilter_Patch{
						Operation: networking.EnvoyFilter_Patch_INSERT_BEFORE,
						Value:     getPatch(patch3),
					},
				},
			},
		},
	}
	controllerutil.SetControllerReference(m, envoyFilter, r.scheme)
	return envoyFilter
}

func getPatch(str string) *proto_types.Struct {
	res, _ := encoding.YAML2Struct(str)
	return res
}

var patch1 = `
operation: INSERT_BEFORE
value:
  config:
    domain: test
    failure_mode_deny: true
    rate_limit_service:
      grpc_service:
        envoy_grpc:
          cluster_name: rate_limit_service
        timeout: 10s
`

var patch2 = `
operation: ADD
value:
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
                  address: rate-limit-server.test-project.svc.cluster.local
                  port_value: 8081
  name: rate_limit_service
  type: STRICT_DNS
`

var patch3 = `
operation: MERGE
value:
  rate_limits:
    - actions:
        - request_headers:
            descriptor_key: custom-rl-header
            header_name: custom-rl-header
`
