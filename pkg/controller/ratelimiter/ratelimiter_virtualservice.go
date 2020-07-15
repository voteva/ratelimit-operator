package ratelimiter

import (
	"context"
	"istio.io/api/networking/v1beta1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	networkingv1beta1 "istio.io/client-go/pkg/apis/networking/v1beta1"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/types"
	operatorsv1alpha1 "ratelimit-operator/pkg/apis/operators/v1alpha1"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
)

func (r *ReconcileRateLimiter) reconcileVirtualService(request reconcile.Request, instance *operatorsv1alpha1.RateLimiter) (reconcile.Result, error) {
	reqLogger := log.WithValues("Request.Namespace", request.Namespace, "Request.Name", request.Name)

	foundVirtualService := &networkingv1beta1.VirtualService{}
	err := r.client.Get(context.TODO(), types.NamespacedName{Name: instance.Name, Namespace: instance.Namespace}, foundVirtualService)
	if err != nil && errors.IsNotFound(err) {
		// Define a new VirtualService
		cm := r.buildVirtualService(instance)
		reqLogger.Info("Creating a new VirtualService", "VirtualService.Namespace", cm.Namespace, "VirtualService.Name", cm.Name)
		err = r.client.Create(context.TODO(), cm)
		if err != nil {
			reqLogger.Error(err, "Failed to create new VirtualService", "VirtualService.Namespace", cm.Namespace, "VirtualService.Name", cm.Name)
			return reconcile.Result{}, err
		}
		// Deployment created successfully - return and requeue
		return reconcile.Result{Requeue: true}, nil
	} else if err != nil {
		reqLogger.Error(err, "Failed to get VirtualService")
		return reconcile.Result{}, err
	}
	return reconcile.Result{}, nil
}

func (r *ReconcileRateLimiter) buildVirtualService(m *operatorsv1alpha1.RateLimiter) *networkingv1beta1.VirtualService {
	virtualService := &networkingv1beta1.VirtualService{
		ObjectMeta: metav1.ObjectMeta{
			Name:      m.Name,
			Namespace: m.Namespace,
		},
		Spec: v1beta1.VirtualService{
			Gateways: []string{
				"istio-ingressgateway",
			},
			Hosts: []string{
				"rate-limit-server.test-project.svc.cluster.local",
			},
			Http: []*v1beta1.HTTPRoute{{
				Route: []*v1beta1.HTTPRouteDestination{{
					Destination: &v1beta1.Destination{
						Host: "rate-limit-server.test-project.svc.cluster.local",
					},
				}},
			}},
			Tcp: []*v1beta1.TCPRoute{{
				Match: []*v1beta1.L4MatchAttributes{{
					Port: 8081,
				}},
				Route: []*v1beta1.RouteDestination{{
					Destination: &v1beta1.Destination{
						Host: "rate-limit-server.test-project.svc.cluster.local",
						Port: &v1beta1.PortSelector{
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
