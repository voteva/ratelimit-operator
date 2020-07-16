package ratelimiter

import (
	"context"
	networking "istio.io/api/networking/v1alpha3"
	"istio.io/client-go/pkg/apis/networking/v1alpha3"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"ratelimit-operator/pkg/apis/operators/v1"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
)

func (r *ReconcileRateLimiter) reconcileVirtualService(ctx context.Context, instance *v1.RateLimiter) (reconcile.Result, error) {
	foundVirtualService := &v1alpha3.VirtualService{}

	err := r.client.Get(ctx, types.NamespacedName{Name: instance.Name, Namespace: instance.Namespace}, foundVirtualService)
	if err != nil && errors.IsNotFound(err) {
		vs := r.buildVirtualService(instance)
		log.Info("Creating a new VirtualService", "VirtualService.Namespace", vs.Namespace, "VirtualService.Name", vs.Name)
		err = r.client.Create(ctx, vs)
		if err != nil {
			log.Error(err, "Failed to create new VirtualService", "VirtualService.Namespace", vs.Namespace, "VirtualService.Name", vs.Name)
			return reconcile.Result{}, err
		}
		return reconcile.Result{Requeue: true}, nil
	} else if err != nil {
		log.Error(err, "Failed to get VirtualService")
		return reconcile.Result{}, err
	}
	return reconcile.Result{}, nil
}

func (r *ReconcileRateLimiter) buildVirtualService(m *v1.RateLimiter) *v1alpha3.VirtualService {
	virtualService := &v1alpha3.VirtualService{
		ObjectMeta: metav1.ObjectMeta{
			Name:      m.Name,
			Namespace: m.Namespace,
		},
		Spec: networking.VirtualService{
			Gateways: []string{
				"istio-ingressgateway",
			},
			Hosts: []string{
				"rate-limit-server." + m.Namespace + ".svc.cluster.local",
			},
			Http: []*networking.HTTPRoute{{
				Route: []*networking.HTTPRouteDestination{{
					Destination: &networking.Destination{
						Host: "rate-limit-server." + m.Namespace + ".svc.cluster.local",
					},
				}},
			}},
			Tcp: []*networking.TCPRoute{{
				Match: []*networking.L4MatchAttributes{{
					Port: 8081,
				}},
				Route: []*networking.RouteDestination{{
					Destination: &networking.Destination{
						Host: "rate-limit-server." + m.Namespace + ".svc.cluster.local",
						Port: &networking.PortSelector{
							Number: 8081,
						},
					},
				}},
			}},
		},
	}
	controllerutil.SetControllerReference(m, virtualService, r.scheme)
	return virtualService
}
